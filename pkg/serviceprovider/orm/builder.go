package orm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Builder 查询构建器
type Builder struct {
	db     *gorm.DB
	schema *schema.Schema
}

// New 创建查询构建器
func New(db *gorm.DB, model any) (*Builder, error) {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return nil, err
	}

	return &Builder{
		db:     db.Model(model),
		schema: stmt.Schema,
	}, nil
}

// DB 返回数据库实例
func (b *Builder) DB() *gorm.DB {
	return b.db
}

// Schema 返回模型Schema
func (b *Builder) Schema() *schema.Schema {
	return b.schema
}

// Build 构建查询SQL
func (b *Builder) Build(search map[string]any) (string, []any, error) {
	if len(search) == 0 {
		return "", nil, nil
	}

	return parse(b.schema, search)
}

// Apply 应用查询条件
func (b *Builder) Apply(search map[string]any) (*gorm.DB, error) {
	sql, args, err := b.Build(search)
	if err != nil {
		return nil, err
	}

	if sql != "" {
		b.db = b.db.Where(sql, args...)
	}

	return b.db, nil
}

/**
Example:
// 普通字段、逻辑查询、EXISTS
Request:
GET /api/v1/user?page=1&pageSize=10&__search={
    "or":[
        {
            "and":[
                {
                    "createdAt":[">","2025-01-01"]
                },
                {
                    "createdAt":["<","2026-01-01"]
                },
                {
                    "not exist":{
                        "UserRoles":{
                            "name":""
                        }
                    }
                }
            ]
        },
        {
            "username":"admin"
        }
    ]
}

Code:
	db := facade.DB().Model(&model.User{})

	whereSql, args, err := orm.BuildCondition(
		db,
		&model.User{},
		search,
	)
	if err != nil {
		return err
	}

	if whereSql != "" {
		db = db.Where(whereSql, args...)
	}

SQL:
SELECT *
FROM user
WHERE (
    (
        (
            user.created_at > '2025-01-01'
        )
        AND (
            user.created_at < '2026-01-01'
        )
        AND (
            NOT EXISTS (
                SELECT 1
                FROM user_roles
                WHERE user_roles.user_id = user.id
                AND user_roles.name = ''
            )
        )
    )
    OR (
        user.username = 'admin'
    )
)
AND user.deleted_at IS NULL
ORDER BY id DESC
LIMIT 10;

// JSON字段查询
Request:
GET /api/v1/menu?page=1&pageSize=10&__search={
    "and":[
        {
            "createdAt":[">","2025-01-01"]
        },
        {
            "createdAt":["<","2026-01-01"]
        },
        {
            "name":""
        },
        {
            "$.meta.icon":"ele-Collection"
        }
    ]
}

Code:
	db := facade.DB().Model(&model.Menu{})

	whereSql, args, err := orm.BuildCondition(
		db,
		&model.Menu{},
		search,
	)
	if err != nil {
		return err
	}

	if whereSql != "" {
		db = db.Where(whereSql, args...)
	}

SQL:
SELECT *
FROM menu
WHERE (
    menu.created_at > '2025-01-01'
)
AND (
    menu.created_at < '2026-01-01'
)
AND (
    menu.name = ''
)
AND (
    JSON_EXTRACT(menu.meta,'$.icon') = 'ele-Collection'
)
AND menu.deleted_at IS NULL
ORDER BY id DESC
LIMIT 10;

// 关联字段查询
Request:
GET /api/v1/user?page=1&pageSize=10&__search={
    "UserRoles.name":"admin"
}

SQL:
SELECT *
FROM user
WHERE EXISTS (
    SELECT 1
    FROM user_roles
    WHERE user_roles.user_id = user.id
    AND user_roles.name = 'admin'
);

// IN查询
Request:
{
    "id":["in",[1,2,3]]
}

SQL:
WHERE user.id IN (1,2,3)

// BETWEEN查询
Request:
{
    "createdAt":[
        "between",
        [
            "2025-01-01",
            "2025-12-31"
        ]
    ]
}

SQL:
WHERE user.created_at BETWEEN '2025-01-01' AND '2025-12-31'


// LIKE查询
Request:
{
    "username":["like","admin"]
}

SQL:
WHERE user.username LIKE '%admin%'


// NULL查询
Request:
{
    "deletedAt":["is null"]
}

SQL:
WHERE user.deleted_at IS NULL

// 多级关联查询
Request:
{
    "exist":{
        "UserRoles.Menu":{
            "name":"系统管理"
        }
    }
}

SQL:
SELECT *
FROM user
WHERE EXISTS (
    SELECT 1
    FROM user_roles
    ...
);

Supported Operator:
=, !=, <>, >, >=, <, <=,
like, left like, right like,
in, not in,
between, not between,
is null, is not null

Supported Logic:
and
or
exist
not exist

Supported Field:
name
createdAt
UserRoles.name
$.meta.icon
*/

// BuildCondition 构建查询条件
func BuildCondition(db *gorm.DB, model any, search map[string]any) (string, []any, error) {
	builder, err := New(db, model)
	if err != nil {
		return "", nil, err
	}

	return builder.Build(search)
}

// ApplyCondition 应用查询条件
func ApplyCondition(db *gorm.DB, model any, search map[string]any) (*gorm.DB, error) {
	builder, err := New(db, model)
	if err != nil {
		return nil, err
	}

	return builder.Apply(search)
}
