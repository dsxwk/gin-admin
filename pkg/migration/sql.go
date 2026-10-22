package migration

import (
	"context"
	"errors"
	"fmt"
	"gin/pkg/cli"
	"gin/pkg/terminal"
	"os"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ExecuteSQLFile 执行初始化SQL文件
func ExecuteSQLFile(ctx context.Context, db *gorm.DB, file string) error {
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	ctx = normalizeContext(ctx)

	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("初始化sql文件不存在: %s", file)
		}
		return fmt.Errorf("检查初始化sql文件失败: %w", err)
	}

	sqlBytes, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("读取初始化sql文件失败: %w", err)
	}

	statements, err := parseSQLStatements(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("解析初始化sql文件失败: %w", err)
	}
	if len(statements) == 0 {
		return fmt.Errorf("初始化sql文件没有可执行语句")
	}

	db = db.WithContext(ctx).Session(&gorm.Session{Logger: logger.Discard})
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("开启事务失败: %w", tx.Error)
	}

	tracker, progressDone := cli.StartProgress("初始化数据库", len(statements))
	for index, statement := range statements {
		if err = tx.Exec(statement).Error; err != nil {
			_ = tx.Rollback().Error
			cli.StopProgress(tracker, progressDone, false)
			return fmt.Errorf("执行第%d条初始化SQL失败: %w\nsql: %s", index+1, err, statement)
		}
		tracker.Increment(1)
	}

	if err = tx.Commit().Error; err != nil {
		cli.StopProgress(tracker, progressDone, false)
		return fmt.Errorf("提交事务失败: %w", err)
	}
	cli.StopProgress(tracker, progressDone, true)
	terminal.Successf("数据库初始化成功,共执行%d条SQL", len(statements))
	return nil
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
