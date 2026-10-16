package db

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gin/app/facade"
	"gin/database/migrations"
	"gin/pkg/cli"
	"gin/pkg/cli/base"
	"gin/pkg/terminal"

	"gin/pkg/pathx"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Seed struct {
	base.Command
}

// NewSeed 创建命令
func NewSeed() *Seed {
	return &Seed{
		Command: base.NewCommand(base.Config{
			Name:        "db:seed",
			Description: "数据填充",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Short: "i", Long: "id"},
					Desc:     "执行指定Seed ID, 如: 20251212_user_seed",
					Required: false,
				},
				{
					Flag:     base.Flag{Short: "I", Long: "init", Default: "false"},
					Desc:     "初始化数据",
					Required: false,
				},
			},
		}),
	}
}

func (s *Seed) Handle(values map[string]string) error {
	terminal.Infof("开始执行数据填充...")

	if s.StringToBool(values["init"]) {
		s.initDatabase()
		return nil

	}

	db := facade.DB()
	if db == nil {
		terminal.Errorf("数据库未初始化")
		return nil

	}
	db = withoutSQLLog(db)

	id := values["id"]
	total := 0
	for _, seed := range migrations.Seeds() {
		if id != "" && seed.ID() != id {
			continue
		}
		total++
	}

	if id != "" && total == 0 {
		terminal.Errorf("Seed ID: %s 不存在", id)
		return nil

	}

	tracker, progressDone := cli.StartProgress("执行数据填充", total)
	for _, seed := range migrations.Seeds() {
		if id != "" && seed.ID() != id {
			continue
		}
		if err := seed.Run(db); err != nil {
			cli.StopProgress(tracker, progressDone, false)
			terminal.Errorf("Seed %s 执行失败: %v", seed.ID(), err)
			return nil

		}
		tracker.Increment(1)
		terminal.Successf("Seed %s 执行成功", seed.ID())
	}
	cli.StopProgress(tracker, progressDone, true)
	return nil
}

// initDatabase 逐条执行数据库初始化语句
func (s *Seed) initDatabase() {
	file := filepath.Join(pathx.RootPath(), "database", "gin.sql")

	if _, err := os.Stat(file); os.IsNotExist(err) {
		terminal.Errorf("初始化sql文件不存在: %s", file)
		return
	}

	sqlBytes, err := os.ReadFile(file)
	if err != nil {
		terminal.Errorf("读取初始化sql文件失败: %v", err)
		return
	}

	statements, err := parseSQLStatements(string(sqlBytes))
	if err != nil {
		terminal.Errorf("解析初始化sql文件失败: %v", err)
		return
	}
	if len(statements) == 0 {
		terminal.Errorf("初始化sql文件没有可执行语句")
		return
	}

	db := facade.DB()
	if db == nil {
		terminal.Errorf("数据库未初始化")
		return
	}
	db = withoutSQLLog(db)

	tx := db.Begin()
	if tx.Error != nil {
		terminal.Errorf("开启事务失败: %v", tx.Error)
		return
	}

	tracker, progressDone := cli.StartProgress("初始化数据库", len(statements))
	for _, stmt := range statements {
		if err = tx.Exec(stmt).Error; err != nil {
			_ = tx.Rollback().Error
			cli.StopProgress(tracker, progressDone, false)
			terminal.Errorf("执行sql失败: %v\nsql: %s", err, stmt)
			return
		}
		tracker.Increment(1)
	}

	if err = tx.Commit().Error; err != nil {
		cli.StopProgress(tracker, progressDone, false)
		terminal.Errorf("提交事务失败: %v", err)
		return
	}

	cli.StopProgress(tracker, progressDone, true)
	terminal.Successf("数据库初始化成功, 共执行 %d 条sql", len(statements))
}

// withoutSQLLog 临时关闭sql输出
func withoutSQLLog(db *gorm.DB) *gorm.DB {
	return db.Session(&gorm.Session{Logger: logger.Discard})
}

// parseSQLStatements 解析SQL文件提取可执行语句
func parseSQLStatements(sqlContent string) ([]string, error) {
	var statements []string
	var currentStmt strings.Builder
	var quote byte
	var escaped bool
	var lineComment bool
	var blockComment bool

	for i := 0; i < len(sqlContent); i++ {
		ch := sqlContent[i]

		if lineComment {
			if ch == '\n' {
				lineComment = false
				if currentStmt.Len() > 0 {
					currentStmt.WriteByte('\n')
				}
			}
			continue
		}

		if blockComment {
			if ch == '*' && i+1 < len(sqlContent) && sqlContent[i+1] == '/' {
				blockComment = false
				i++
				if currentStmt.Len() > 0 {
					currentStmt.WriteByte(' ')
				}
			}
			continue
		}

		if quote != 0 {
			currentStmt.WriteByte(ch)

			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == quote {
				if i+1 < len(sqlContent) && sqlContent[i+1] == quote {
					currentStmt.WriteByte(sqlContent[i+1])
					i++
					continue
				}
				quote = 0
			}
			continue
		}

		switch {
		case ch == '\'' || ch == '"' || ch == '`':
			quote = ch
			currentStmt.WriteByte(ch)
		case ch == '-' && i+1 < len(sqlContent) && sqlContent[i+1] == '-' && isSQLCommentEnd(sqlContent, i+2):
			lineComment = true
			i++
		case ch == '#':
			lineComment = true
		case ch == '/' && i+1 < len(sqlContent) && sqlContent[i+1] == '*':
			blockComment = true
			i++
		case ch == ';':
			appendSQLStatement(&statements, currentStmt.String())
			currentStmt.Reset()
		default:
			currentStmt.WriteByte(ch)
		}
	}

	if quote != 0 {
		return nil, errors.New("sql字符串未闭合")
	}
	if blockComment {
		return nil, errors.New("sql块注释未闭合")
	}

	lastStmt := strings.TrimSpace(currentStmt.String())
	if lastStmt != "" {
		statements = append(statements, lastStmt)
	}

	return statements, nil
}

// appendSQLStatement 添加非空SQL语句
func appendSQLStatement(statements *[]string, statement string) {
	statement = strings.TrimSpace(statement)
	if statement == "" {
		return
	}
	*statements = append(*statements, statement+";")
}

// isSQLCommentEnd 判断是否为SQL注释结束位置
func isSQLCommentEnd(content string, index int) bool {
	if index >= len(content) {
		return true
	}
	switch content[index] {
	case ' ', '\t', '\r', '\n':
		return true
	default:
		return false
	}
}
