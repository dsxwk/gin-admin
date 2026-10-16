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

type MakeSeed struct {
	base.Command
}

// NewMakeSeed 创建命令
func NewMakeSeed() *MakeSeed {
	return &MakeSeed{
		Command: base.NewCommand(base.Config{
			Name:        "make:seed",
			Description: "生成数据库seeder模板",
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

func (s *MakeSeed) Handle(values map[string]string) error {
	table := values["table"]

	timestamp := time.Now().Format("20060102")
	name := table
	if name == "" {
		name = "all_tables"
	}

	filename := fmt.Sprintf("%s_%s.go", timestamp, name)
	dir := filepath.Join("database", "seeds")
	filePath := filepath.Join(dir, filename)
	content := fmt.Sprintf(`package seeds

import (
	"gorm.io/gorm"
)

type %s struct{}

func (s *%s) ID() string {
	return "%s"
}

func (s *%s) Run(db *gorm.DB) error {
	/*users := []model.User{
		{
			Username: "admin",
			Password: "$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6",
			Status:   1,
		},
		{
			Username: "test",
			Password: "$2a$10$OcSkSCBe8D5tGL2ulmJhTe0Xboy/fzwS1H7AdmkJjpQZfeGUHr5S6",
			Status:   1,
		},
	}

	if err := db.Create(&users).Error; err != nil {
		return err
	}*/
	// todo: 实现填充逻辑
	return nil
}
`, lo.PascalCase(name), lo.PascalCase(name), name+"_"+timestamp, lo.PascalCase(name))

	// 创建文件
	f := s.CheckDirAndFile(filePath)
	if f == nil {
		return nil

	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		terminal.Errorf("写入文件失败: %v", err)
		return nil

	}

	terminal.Successf("seed文件生成成功: %s", filePath)
	return nil
}
