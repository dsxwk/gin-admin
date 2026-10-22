package migration

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"gorm.io/gorm"
)

const (
	mysqlMigrationLock = "gin_admin_migrations"
	mysqlLockTimeout   = 10
)

// Runner 数据库迁移执行器
type Runner struct {
	db         *gorm.DB
	migrations []Migration
}

// NewRunner 创建迁移执行器
func NewRunner(db *gorm.DB, items []Migration) (*Runner, error) {
	if db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	migrations, err := sortedMigrations(items)
	if err != nil {
		return nil, err
	}
	return &Runner{db: db, migrations: migrations}, nil
}

// Migrate 执行全部待处理迁移
func (r *Runner) Migrate(ctx context.Context) ([]string, error) {
	ctx = normalizeContext(ctx)
	var executed []string

	err := r.withLock(ctx, func() error {
		repo := newRepository(r.db)
		if err := repo.install(ctx); err != nil {
			return fmt.Errorf("初始化迁移记录表失败: %w", err)
		}

		records, err := repo.all(ctx)
		if err != nil {
			return fmt.Errorf("读取迁移记录失败: %w", err)
		}
		applied := make(map[string]struct{}, len(records))
		for _, record := range records {
			applied[record.Migration] = struct{}{}
		}

		pending := make([]Migration, 0, len(r.migrations))
		for _, item := range r.migrations {
			if _, exists := applied[item.Name()]; !exists {
				pending = append(pending, item)
			}
		}
		if len(pending) == 0 {
			return nil
		}

		batch, err := repo.maxBatch(ctx)
		if err != nil {
			return fmt.Errorf("读取迁移批次失败: %w", err)
		}
		batch++

		for _, item := range pending {
			name := item.Name()
			if err = item.Up(ctx, r.db.WithContext(ctx)); err != nil {
				return &RunError{Name: name, Err: err}
			}
			if err = repo.create(ctx, name, batch); err != nil {
				return fmt.Errorf("记录迁移%s失败: %w", name, err)
			}
			executed = append(executed, name)
		}
		return nil
	})
	return executed, err
}

// MarkMigrated 标记全部迁移已执行
func (r *Runner) MarkMigrated(ctx context.Context) ([]string, error) {
	ctx = normalizeContext(ctx)
	var marked []string

	err := r.withLock(ctx, func() error {
		repo := newRepository(r.db)
		if err := repo.install(ctx); err != nil {
			return fmt.Errorf("初始化迁移记录表失败: %w", err)
		}

		records, err := repo.all(ctx)
		if err != nil {
			return fmt.Errorf("读取迁移记录失败: %w", err)
		}
		applied := make(map[string]struct{}, len(records))
		for _, record := range records {
			applied[record.Migration] = struct{}{}
		}

		pending := make([]Migration, 0, len(r.migrations))
		for _, item := range r.migrations {
			if _, exists := applied[item.Name()]; !exists {
				pending = append(pending, item)
			}
		}
		if len(pending) == 0 {
			return nil
		}

		batch, err := repo.maxBatch(ctx)
		if err != nil {
			return fmt.Errorf("读取迁移批次失败: %w", err)
		}
		batch++
		for _, item := range pending {
			name := item.Name()
			if err = repo.create(ctx, name, batch); err != nil {
				return fmt.Errorf("记录迁移%s失败: %w", name, err)
			}
			marked = append(marked, name)
		}
		return nil
	})
	return marked, err
}

// Rollback 回滚最近N个批次
func (r *Runner) Rollback(ctx context.Context, steps int) ([]string, error) {
	if steps <= 0 {
		return nil, fmt.Errorf("回滚批次必须大于0")
	}
	ctx = normalizeContext(ctx)

	var executed []string
	err := r.withLock(ctx, func() error {
		repo := newRepository(r.db)
		if err := repo.install(ctx); err != nil {
			return fmt.Errorf("初始化迁移记录表失败: %w", err)
		}

		records, err := repo.all(ctx)
		if err != nil {
			return fmt.Errorf("读取迁移记录失败: %w", err)
		}
		if len(records) == 0 {
			return nil
		}

		batches := latestBatches(records, steps)
		targets := make([]Record, 0, len(records))
		for _, record := range records {
			if _, exists := batches[record.Batch]; exists {
				targets = append(targets, record)
			}
		}
		sort.Slice(targets, func(i, j int) bool {
			if targets[i].Batch != targets[j].Batch {
				return targets[i].Batch > targets[j].Batch
			}
			return targets[i].Migration > targets[j].Migration
		})

		items := migrationMap(r.migrations)
		for _, record := range targets {
			item, exists := items[record.Migration]
			if !exists {
				return fmt.Errorf("迁移文件不存在,无法回滚: %s", record.Migration)
			}
			if err = item.Down(ctx, r.db.WithContext(ctx)); err != nil {
				return &RunError{Name: record.Migration, Err: err}
			}
			if err = repo.delete(ctx, record.Migration); err != nil {
				return fmt.Errorf("删除迁移记录%s失败: %w", record.Migration, err)
			}
			executed = append(executed, record.Migration)
		}
		return nil
	})
	return executed, err
}

// RollbackAll 回滚全部迁移
func (r *Runner) RollbackAll(ctx context.Context) ([]string, error) {
	ctx = normalizeContext(ctx)

	var executed []string
	err := r.withLock(ctx, func() error {
		repo := newRepository(r.db)
		if err := repo.install(ctx); err != nil {
			return fmt.Errorf("初始化迁移记录表失败: %w", err)
		}
		records, err := repo.all(ctx)
		if err != nil {
			return fmt.Errorf("读取迁移记录失败: %w", err)
		}
		sort.Slice(records, func(i, j int) bool {
			if records[i].Batch != records[j].Batch {
				return records[i].Batch > records[j].Batch
			}
			return records[i].Migration > records[j].Migration
		})

		items := migrationMap(r.migrations)
		for _, record := range records {
			item, exists := items[record.Migration]
			if !exists {
				return fmt.Errorf("迁移文件不存在,无法回滚: %s", record.Migration)
			}
			if err = item.Down(ctx, r.db.WithContext(ctx)); err != nil {
				return &RunError{Name: record.Migration, Err: err}
			}
			if err = repo.delete(ctx, record.Migration); err != nil {
				return fmt.Errorf("删除迁移记录%s失败: %w", record.Migration, err)
			}
			executed = append(executed, record.Migration)
		}
		return nil
	})
	return executed, err
}

// Status 获取全部迁移状态
func (r *Runner) Status(ctx context.Context) ([]Status, error) {
	ctx = normalizeContext(ctx)
	repo := newRepository(r.db)
	if err := repo.install(ctx); err != nil {
		return nil, fmt.Errorf("初始化迁移记录表失败: %w", err)
	}
	records, err := repo.all(ctx)
	if err != nil {
		return nil, fmt.Errorf("读取迁移记录失败: %w", err)
	}

	applied := make(map[string]Record, len(records))
	for _, record := range records {
		applied[record.Migration] = record
	}

	result := make([]Status, 0, len(r.migrations)+len(records))
	for _, item := range r.migrations {
		name := item.Name()
		record, exists := applied[name]
		if !exists {
			result = append(result, Status{Name: name})
			continue
		}
		result = append(result, Status{
			Name:      name,
			Batch:     record.Batch,
			Ran:       true,
			CreatedAt: record.CreatedAt,
		})
	}

	known := make(map[string]struct{}, len(r.migrations))
	for _, item := range r.migrations {
		known[item.Name()] = struct{}{}
	}
	for _, record := range records {
		if _, exists := known[record.Migration]; exists {
			continue
		}
		result = append(result, Status{
			Name:      record.Migration,
			Batch:     record.Batch,
			Ran:       true,
			CreatedAt: record.CreatedAt,
		})
	}
	return result, nil
}

// withLock 执行迁移时获取数据库锁
func (r *Runner) withLock(ctx context.Context, run func() error) error {
	if r.db.Dialector.Name() != "mysql" {
		return run()
	}

	sqlDB, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}
	conn, err := sqlDB.Conn(ctx)
	if err != nil {
		return fmt.Errorf("获取迁移锁连接失败: %w", err)
	}
	defer func() { _ = conn.Close() }()

	var acquired int
	if err = conn.QueryRowContext(ctx,
		"SELECT GET_LOCK(?, ?)",
		mysqlMigrationLock,
		mysqlLockTimeout,
	).Scan(&acquired); err != nil {
		return fmt.Errorf("获取迁移锁失败: %w", err)
	}
	if acquired != 1 {
		return fmt.Errorf("获取迁移锁超时")
	}
	defer func() {
		_, _ = conn.ExecContext(context.Background(), "SELECT RELEASE_LOCK(?)", mysqlMigrationLock)
	}()

	return run()
}

// sortedMigrations 校验并排序迁移
func sortedMigrations(items []Migration) ([]Migration, error) {
	result := append([]Migration(nil), items...)
	seen := make(map[string]struct{}, len(result))
	for _, item := range result {
		if item == nil {
			return nil, fmt.Errorf("迁移列表中存在空项")
		}
		name := strings.TrimSpace(item.Name())
		if name == "" {
			return nil, fmt.Errorf("迁移名称不能为空")
		}
		if _, exists := seen[name]; exists {
			return nil, fmt.Errorf("迁移名称重复: %s", name)
		}
		seen[name] = struct{}{}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name() < result[j].Name()
	})
	return result, nil
}

// migrationMap 创建迁移映射
func migrationMap(items []Migration) map[string]Migration {
	result := make(map[string]Migration, len(items))
	for _, item := range items {
		result[item.Name()] = item
	}
	return result
}

// latestBatches 获取最近N个批次
func latestBatches(records []Record, steps int) map[int]struct{} {
	batchSet := make(map[int]struct{})
	for _, record := range records {
		if record.Batch > 0 {
			batchSet[record.Batch] = struct{}{}
		}
	}
	batches := make([]int, 0, len(batchSet))
	for batch := range batchSet {
		batches = append(batches, batch)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(batches)))
	if len(batches) > steps {
		batches = batches[:steps]
	}
	result := make(map[int]struct{}, len(batches))
	for _, batch := range batches {
		result[batch] = struct{}{}
	}
	return result
}

// normalizeContext 标准化上下文
func normalizeContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}
