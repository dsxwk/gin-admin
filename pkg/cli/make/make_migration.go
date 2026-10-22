package make

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gin/pkg/cli/base"
	"gin/pkg/migration"
	"gin/pkg/terminal"
)

type MakeMigration struct {
	base.Command
}

// NewMakeMigration 创建命令
func NewMakeMigration() *MakeMigration {
	return &MakeMigration{
		Command: base.NewCommand(base.Config{
			Name:        "make:migration",
			Description: "生成数据库迁移模板",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Short: "n", Long: "name"},
					Desc:     "迁移名称,如:create_orders_table",
					Required: false,
				},
				{
					Flag:     base.Flag{Short: "t", Long: "table"},
					Desc:     "指定表名,未设置name时使用",
					Required: false,
				},
			},
		}),
	}
}

func (s *MakeMigration) Handle(values map[string]string) error {
	timestamp := time.Now().Format("20060102150405")
	plan, err := migration.BuildMigrationPlan(values["name"], values["table"], timestamp)
	if err != nil {
		return err
	}

	dir := filepath.Join("database", "migrations")
	filePath := filepath.Join(dir, plan.FileName)
	content := fmt.Sprintf(`package migrations

import (
	"context"

	"gin/pkg/migration"

	"gorm.io/gorm"
)

type %s struct{}

// Name 迁移名称
func (m *%s) Name() string {
	return "%s"
}

// Up 执行迁移
func (m *%s) Up(ctx context.Context, db *gorm.DB) error {
	return migration.CreateTable(db.WithContext(ctx), "%s", func(t *migration.Table) {
		t.Id()
		t.Timestamps()
	})
}

// Down 回滚迁移
func (m *%s) Down(ctx context.Context, db *gorm.DB) error {
	return migration.DropTable(db.WithContext(ctx), "%s")
}
`, plan.StructName, plan.StructName, plan.Name, plan.StructName, plan.TableName, plan.StructName, plan.TableName)

	// 创建文件
	f := s.CheckDirAndFile(filePath)
	if f == nil {
		return nil

	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		terminal.Errorf("写入文件失败: %v", err)
		return err
	}
	if err := migration.RegisterGenerated(
		filepath.Join(dir, "registry.go"),
		"migration.Migration",
		"&"+plan.StructName+"{}",
	); err != nil {
		return fmt.Errorf("注册迁移失败: %w", err)
	}

	terminal.Successf("迁移文件生成成功: %s", filePath)
	return nil
}
