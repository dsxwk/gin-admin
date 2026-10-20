package model

const TableNameFile = "file"

// File 文件表
type File struct {
	ID         int64      `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	GroupId    int64      `gorm:"column:group_id;not null;default:0;type:int(10) unsigned;comment:分组ID 0为未分组" json:"groupId" form:"groupId"`
	Group      *FileGroup `gorm:"foreignKey:group_id;references:id" json:"group"`
	Storage    string     `gorm:"column:storage;not null;default:local;type:varchar(20);comment:存储平台 local,qiniu,aliyun,tencent" json:"storage" form:"storage"`
	OriginName string     `gorm:"column:origin_name;not null;type:varchar(255);comment:原始文件名(如: 商品主图.jpg)" json:"originName" form:"originName"`
	ObjectKey  string     `gorm:"column:object_key;not null;type:varchar(255);comment:存储路径/Key(如: 2026/10/20/xxx.jpg)" json:"objectKey" form:"objectKey"`
	Url        string     `gorm:"column:url;not null;type:varchar(255);comment:完整访问URL" json:"url" form:"url"`
	MimeType   string     `gorm:"column:mime_type;type:varchar(100);comment:MIME类型(如: application/octet-stream)" json:"mimeType" form:"mimeType"`
	Size       int64      `gorm:"column:size;not null;default:0;type:int(10) unsigned;comment:文件大小(字节)" json:"size" form:"size"`
	Hash       string     `gorm:"column:hash;type:varchar(64);comment:文件哈希值（MD5/SHA1，用于秒传校验）" json:"hash" form:"hash"`
	Ext        *JsonValue `gorm:"column:ext;type:json;comment:扩展信息（如图片宽高、视频时长等）" json:"ext" form:"ext"`
	UploaderId int64      `gorm:"column:uploader_id;default:0;type:int(10) unsigned;comment:上传者ID（关联用户表）" json:"uploaderId" form:"uploaderId"`
	CreatedAt  *DateTime  `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt  *DateTime  `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
}

func (*File) TableName() string {
	return TableNameFile
}

// Connection 数据库连接名称
func (*File) Connection() string {
	return "mysql"
}
