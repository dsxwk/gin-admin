package migration

import (
	"context"
	"time"

	"gorm.io/gorm"
)

const recordTableName = "migrations"

// Record 迁移执行记录
type Record struct {
	ID        uint       `gorm:"primaryKey;comment:ID"`
	Migration string     `gorm:"size:191;not null;uniqueIndex;comment:迁移名称"`
	Batch     int        `gorm:"not null;default:0;index;comment:迁移批次"`
	CreatedAt *time.Time `gorm:"comment:创建时间"`
}

// TableName 表名
func (Record) TableName() string {
	return recordTableName
}

type repository struct {
	db *gorm.DB
}

// newRepository 创建迁移仓库
func newRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

// install 初始化迁移记录表
func (r *repository) install(ctx context.Context) error {
	return r.db.WithContext(ctx).AutoMigrate(&Record{})
}

// all 获取全部迁移记录
func (r *repository) all(ctx context.Context) ([]Record, error) {
	var records []Record
	if err := r.db.WithContext(ctx).
		Order("batch ASC, migration ASC").
		Find(&records).
		Error; err != nil {
		return nil, err
	}
	return records, nil
}

// maxBatch 获取最大迁移批次
func (r *repository) maxBatch(ctx context.Context) (int, error) {
	var batch int
	if err := r.db.WithContext(ctx).
		Model(&Record{}).
		Select("COALESCE(MAX(batch), 0)").
		Scan(&batch).
		Error; err != nil {
		return 0, err
	}
	return batch, nil
}

// create 保存迁移记录
func (r *repository) create(ctx context.Context, name string, batch int) error {
	now := time.Now()
	return r.db.WithContext(ctx).Create(&Record{
		Migration: name,
		Batch:     batch,
		CreatedAt: &now,
	}).Error
}

// delete 删除迁移记录
func (r *repository) delete(ctx context.Context, name string) error {
	return r.db.WithContext(ctx).
		Where("migration = ?", name).
		Delete(&Record{}).
		Error
}
