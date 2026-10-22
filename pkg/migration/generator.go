package migration

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/samber/lo"
)

// MigrationPlan 迁移生成计划
type MigrationPlan struct {
	Name       string
	StructName string
	TableName  string
	FileName   string
}

// BuildMigrationPlan 根据参数生成迁移命名
func BuildMigrationPlan(name, table, timestamp string) (MigrationPlan, error) {
	name = strings.TrimSpace(name)
	table = strings.TrimSpace(table)
	if name == "" && table == "" {
		return MigrationPlan{}, fmt.Errorf("迁移名称或表名不能为空")
	}

	if name == "" {
		name = "create_" + table + "_table"
	}
	name = strings.TrimPrefix(name, timestamp+"_")

	tableName := table
	if tableName == "" {
		tableName = strings.TrimSuffix(strings.TrimPrefix(name, "create_"), "_table")
	}

	structName := lo.PascalCase(name)
	if table != "" &&
		strings.HasPrefix(name, "create_") &&
		strings.HasSuffix(name, "_table") {
		structName = "Create" + lo.PascalCase(pluralizeTableName(table)) + "Table"
	}

	return MigrationPlan{
		Name:       timestamp + "_" + name,
		StructName: structName,
		TableName:  tableName,
		FileName:   timestamp + "_" + name + ".go",
	}, nil
}

// RegisterGenerated 注册生成项到清单
func RegisterGenerated(file, listType, item string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	if bytes.Contains(content, []byte(item)) {
		return nil
	}

	newline := "\n"
	if bytes.Contains(content, []byte("\r\n")) {
		newline = "\r\n"
	}

	marker := strings.ReplaceAll(
		"func All() []"+listType+" {\n\treturn []"+listType+"{",
		"\n",
		newline,
	)
	start := bytes.Index(content, []byte(marker))
	if start < 0 {
		return fmt.Errorf("未找到注册清单:%s", file)
	}

	bodyStart := start + len(marker)
	closeMarker := []byte(newline + "\t}")
	end := bytes.Index(content[bodyStart:], closeMarker)
	if end < 0 {
		return fmt.Errorf("注册清单格式不正确:%s", file)
	}

	insertAt := bodyStart + end
	itemLine := []byte(newline + "\t\t" + item + ",")
	result := make([]byte, 0, len(content)+len(itemLine))
	result = append(result, content[:insertAt]...)
	result = append(result, itemLine...)
	result = append(result, content[insertAt:]...)
	return os.WriteFile(file, result, 0644)
}

// pluralizeTableName 生成表名复数形式
func pluralizeTableName(name string) string {
	lower := strings.ToLower(name)
	if strings.HasSuffix(lower, "s") ||
		strings.HasSuffix(lower, "x") ||
		strings.HasSuffix(lower, "z") ||
		strings.HasSuffix(lower, "ch") ||
		strings.HasSuffix(lower, "sh") {
		return name
	}
	if strings.HasSuffix(lower, "y") && len(name) > 1 {
		return name[:len(name)-1] + "ies"
	}
	return name + "s"
}
