package db

import (
	"gin/app/facade"
	"gin/common/base"
	"gin/common/flag"
	"gin/database/migrations"
	"gin/pkg"
	"gin/pkg/cli"
	"os"
	"strings"
)

type Seed struct {
	base.BaseCommand
}

func (s *Seed) Name() string {
	return "db:seed"
}

func (s *Seed) Description() string {
	return "数据填充"
}

func (s *Seed) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			base.Flag{Short: "i", Long: "id"},
			"执行指定Seed ID, 如: 20251212_user_seed",
			false,
		},
		{
			base.Flag{Short: "I", Long: "init", Default: "false"},
			"初始化数据",
			false,
		},
	}
}

func (s *Seed) Execute(args []string) {
	values := s.ParseFlags(s.Name(), args, s.Help())
	flag.Infof("开始执行数据填充...")

	db := facade.DB()
	if values["init"] == "true" || values["init"] == "1" {
		file := pkg.GetRootPath() + "/database/gin.sql"

		// 检查文件是否存在
		if _, err := os.Stat(file); os.IsNotExist(err) {
			flag.Errorf("SQL文件不存在: %s", file)
			return
		}

		// 读取文件内容
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			flag.Errorf("读取SQL文件失败: %v", err)
			return
		}

		sqlContent := string(sqlBytes)

		// 开启事务
		tx := db.Begin()

		// 分割并执行SQL语句
		statements := parseSQLStatements(sqlContent)
		successCount := 0

		for _, stmt := range statements {
			if strings.TrimSpace(stmt) == "" {
				continue
			}
			if err = tx.Exec(stmt).Error; err != nil {
				tx.Rollback()
				flag.Errorf("执行SQL失败: %v\nSQL: %s", err, stmt)
				return
			}
			successCount++
		}

		// 提交事务
		if err = tx.Commit().Error; err != nil {
			flag.Errorf("提交事务失败: %v", err)
			return
		}

		flag.Successf("数据库初始化成功,共执行 %d 条SQL语句", successCount)
	} else {
		id := values["id"]
		for _, seed := range migrations.AllSeeds() {
			if id != "" && seed.ID() != id {
				continue
			}

			if err := seed.Run(db); err != nil {
				flag.Errorf("Seed %s 执行失败: %v", seed.ID(), err)
				return
			}
			flag.Successf("Seed %s 执行成功", seed.ID())
		}
	}
}

// parseSQLStatements 解析SQL文件,提取可执行的SQL语句
func parseSQLStatements(sqlContent string) []string {
	var statements []string
	var currentStmt strings.Builder
	inString := false
	inComment := false

	lines := strings.Split(sqlContent, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 跳过空行
		if trimmed == "" {
			continue
		}

		// 跳过单行注释（以 -- 开头）
		if strings.HasPrefix(trimmed, "--") {
			continue
		}

		// 处理多行注释
		if strings.Contains(trimmed, "/*") && !inComment {
			inComment = true
		}
		if inComment {
			if strings.Contains(trimmed, "*/") {
				inComment = false
			}
			continue
		}

		for i := 0; i < len(line); i++ {
			ch := line[i]

			// 处理字符串内的分号
			if ch == '\'' || ch == '"' {
				inString = !inString
			}

			// 如果不在字符串内,且遇到分号,则结束当前语句
			if !inString && ch == ';' {
				currentStmt.WriteByte(ch)
				stmt := strings.TrimSpace(currentStmt.String())
				if stmt != "" {
					statements = append(statements, stmt)
				}
				currentStmt.Reset()
				continue
			}

			currentStmt.WriteByte(ch)
		}

		// 每行结束后添加换行符
		if currentStmt.Len() > 0 {
			currentStmt.WriteByte('\n')
		}
	}

	// 处理最后没有分号的语句
	lastStmt := strings.TrimSpace(currentStmt.String())
	if lastStmt != "" {
		statements = append(statements, lastStmt)
	}

	return statements
}

func init() {
	cli.Register(&Seed{})
}
