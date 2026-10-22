package make

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gin/pkg/cli/base"
	"gin/pkg/migration"
	"gin/pkg/terminal"

	"github.com/samber/lo"
)

type MakeSeeder struct {
	base.Command
}

// NewMakeSeeder 创建Seeder生成命令
func NewMakeSeeder() *MakeSeeder {
	return &MakeSeeder{
		Command: base.NewCommand(base.Config{
			Name:        "make:seeder",
			Description: "生成数据库Seeder",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Short: "n", Long: "name"},
					Desc:     "Seeder名称,如:UserSeeder",
					Required: true,
				},
			},
		}),
	}
}

func (s *MakeSeeder) Handle(values map[string]string) error {
	name := strings.TrimSpace(values["name"])
	if !strings.HasSuffix(name, "Seeder") {
		name += "Seeder"
	}

	fileName := lo.SnakeCase(strings.TrimSuffix(name, "Seeder")) + "_seeder.go"
	filePath := filepath.Join("database", "seeders", fileName)
	content := fmt.Sprintf(`package seeders

import (
	"context"

	"gorm.io/gorm"
)

type %s struct{}

// Name Seeder名称
func (s *%s) Name() string {
	return "%s"
}

// Run 执行数据填充
func (s *%s) Run(ctx context.Context, db *gorm.DB) error {
	// todo: 实现填充逻辑
	return nil
}
`, name, name, name, name)

	file := s.CheckDirAndFile(filePath)
	if file == nil {
		return nil
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		terminal.Errorf("写入文件失败: %v", err)
		return err
	}
	if err := migration.RegisterGenerated(
		filepath.Join("database", "seeders", "registry.go"),
		"migration.Seeder",
		"&"+name+"{}",
	); err != nil {
		return fmt.Errorf("注册Seeder失败: %w", err)
	}
	terminal.Successf("Seeder文件生成成功: %s", filePath)
	return nil
}
