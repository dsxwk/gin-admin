package make

import (
	"fmt"
	"gin/pkg/cli/base"
	"gin/pkg/terminal"
	"os"
	"path/filepath"
	"time"

	"github.com/samber/lo"
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
					Flag:     base.Flag{Short: "t", Long: "table"},
					Desc:     "指定表名，留空生成所有表模板",
					Required: false,
				},
			},
		}),
	}
}

func (s *MakeMigration) Handle(values map[string]string) error {
	table := values["table"]

	timestamp := time.Now().Format("20060102")
	name := table
	filename := fmt.Sprintf("%s_%s.go", timestamp, "create_"+name+"_table")
	id := timestamp + "_create_" + name + "_table"
	if name == "" {
		name = "all_tables"
		filename = fmt.Sprintf("%s_%s.go", timestamp, "create_"+name)
		id = timestamp + "_create_" + name
	}

	dir := filepath.Join("database", "migrations")
	filePath := filepath.Join(dir, filename)
	content := fmt.Sprintf(`package migrations

import (
	"gorm.io/gorm"
)

type %s struct{}

func (m *%s) ID() string {
	return "%s"
}

func (m *%s) Migrate(db *gorm.DB) error {
	// db.AutoMigrate(&model.User{})
	// todo: 实现迁移逻辑
	return nil
}

func (m *%s) Rollback(db *gorm.DB) error {
	// db.Migrator().DropTable(&model.UserRoles{})
	// todo: 实现回滚逻辑
	return nil
}
`, "Create"+lo.PascalCase(name)+timestamp, "Create"+lo.PascalCase(name)+timestamp, id, "Create"+lo.PascalCase(name)+timestamp, "Create"+lo.PascalCase(name)+timestamp)

	// 创建文件
	f := s.CheckDirAndFile(filePath)
	if f == nil {
		return nil

	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		terminal.Errorf("写入文件失败: %v", err)
		return nil

	}

	terminal.Successf("迁移文件生成成功: %s", filePath)
	return nil
}
