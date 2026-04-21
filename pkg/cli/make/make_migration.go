package make

import (
	"fmt"
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg/cli"
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"time"
)

type MakeMigration struct {
	base.BaseCommand
}

func (s *MakeMigration) Name() string {
	return "make:migration"
}

func (s *MakeMigration) Description() string {
	return "生成数据库迁移模板"
}

func (s *MakeMigration) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			base.Flag{Short: "t", Long: "table"},
			"指定表名，留空生成所有表模板",
			false,
		},
	}
}

func (s *MakeMigration) Execute(args []string) {
	values := s.ParseFlags(s.Name(), args, s.Help())
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
		return
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		flag.Errorf("写入文件失败: %v", err)
		return
	}

	flag.Successf("迁移文件生成成功: %s", filePath)
}

func init() {
	cli.Register(&MakeMigration{})
}
