package model

const TableNameArticle = "article"

// Article 文章表
type Article struct {
	ID         int64      `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	Uid        int64      `gorm:"column:uid;not null;default:0;type:int(11);comment:用户id" json:"uid" form:"uid"`
	Title      string     `gorm:"column:title;not null;type:varchar(50);comment:标题" json:"title" form:"title"`
	Content    string     `gorm:"column:content;not null;type:longtext;comment:内容" json:"content" form:"content"`
	CategoryId int64      `gorm:"column:category_id;not null;default:0;type:int(11);comment:分类id" json:"categoryId" form:"categoryId"`
	DataSource int64      `gorm:"column:data_source;not null;default:2;type:tinyint(3) unsigned;comment:数据来源 1=文章库 2=自建" json:"dataSource" form:"dataSource"`
	IsPublish  int64      `gorm:"column:is_publish;not null;default:1;type:tinyint(3) unsigned;comment:是否发布 0=待发布 1=已发布 2=已下架" json:"isPublish" form:"isPublish"`
	Tag        *JsonValue `gorm:"column:tag;type:json;comment:标签" json:"tag" form:"tag"`
	User       *User      `gorm:"foreignKey:uid;references:id" json:"user"`
	CreatedAt  *DateTime  `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt  *DateTime  `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt  *DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*Article) TableName() string {
	return TableNameArticle
}

// Connection 数据库连接名称
func (m *Article) Connection() string {
	return "mysql"
}
