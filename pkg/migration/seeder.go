package migration

import (
	"context"
	"fmt"
	"gin/pkg/cli"
	"strings"

	"gorm.io/gorm"
)

// Seeder 数据填充接口
type Seeder interface {
	Name() string
	Run(ctx context.Context, db *gorm.DB) error
}

// SeedResult Seed执行结果
type SeedResult struct {
	Name string
}

// RunSeeders 执行全部或指定Seeder
func RunSeeders(ctx context.Context, db *gorm.DB, items []Seeder, name string) ([]SeedResult, error) {
	if db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	ctx = normalizeContext(ctx)

	validated, err := validatedSeeders(items)
	if err != nil {
		return nil, err
	}

	name = strings.TrimSpace(name)
	selected := validated
	if name != "" {
		target, exists := findSeeder(validated, name)
		if !exists {
			return nil, fmt.Errorf("Seeder不存在: %s", name)
		}
		selected = []Seeder{target}
	}
	if len(selected) == 0 {
		return nil, nil
	}

	result := make([]SeedResult, 0, len(selected))
	tracker, progressDone := cli.StartProgress("执行数据填充", len(selected))
	for _, item := range selected {
		seeder := item
		seederName := seeder.Name()
		if err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			return seeder.Run(ctx, tx)
		}); err != nil {
			cli.StopProgress(tracker, progressDone, false)
			return result, &RunError{Name: seederName, Err: err}
		}
		result = append(result, SeedResult{Name: seederName})
		tracker.Increment(1)
	}
	cli.StopProgress(tracker, progressDone, true)
	return result, nil
}

// validatedSeeders 校验Seeder
func validatedSeeders(items []Seeder) ([]Seeder, error) {
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		if item == nil {
			return nil, fmt.Errorf("Seeder列表中存在空项")
		}
		name := strings.TrimSpace(item.Name())
		if name == "" {
			return nil, fmt.Errorf("Seeder名称不能为空")
		}
		if _, exists := seen[name]; exists {
			return nil, fmt.Errorf("Seeder名称重复: %s", name)
		}
		seen[name] = struct{}{}
	}
	return append([]Seeder(nil), items...), nil
}

// findSeeder 查找Seeder
func findSeeder(items []Seeder, name string) (Seeder, bool) {
	for _, item := range items {
		if strings.EqualFold(item.Name(), name) {
			return item, true
		}
	}
	return nil, false
}
