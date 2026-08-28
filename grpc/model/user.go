package model

const TableNameUser = "user"

// User 用户表
type User struct {
	ID        int32              `gorm:"column:id;primaryKey;autoIncrement;not null;type:int(10) unsigned;comment:ID" json:"id" form:"id"`
	Avatar    string             `gorm:"column:avatar;not null;type:varchar(255);comment:头像" json:"avatar" form:"avatar"`
	Username  string             `gorm:"column:username;not null;type:varchar(10);comment:用户名" json:"username" form:"username"`
	FullName  string             `gorm:"column:full_name;not null;type:varchar(20);comment:姓名" json:"fullName" form:"fullName"`
	Email     string             `gorm:"column:email;not null;type:varchar(50);comment:邮箱" json:"email" form:"email"`
	Password  string             `gorm:"column:password;not null;type:varchar(255);comment:密码" json:"password" form:"password"`
	Nickname  string             `gorm:"column:nickname;not null;type:varchar(50);comment:昵称" json:"nickname" form:"nickname"`
	Gender    int32              `gorm:"column:gender;not null;default:0;type:tinyint(1) unsigned;comment:性别 1=男 2=女" json:"gender" form:"gender"`
	Age       int32              `gorm:"column:age;not null;default:0;type:tinyint(3) unsigned;comment:年龄" json:"age" form:"age"`
	Status    int32              `gorm:"column:status;not null;default:1;type:tinyint(3) unsigned;comment:状态 1=启用 2=停用" json:"status" form:"status"`
	UserRoles []*UserRoles       `gorm:"foreignKey:user_id;references:id;comment:用户角色" json:"userRoles"`
	MainDept  *UserDepartments   `gorm:"foreignKey:user_id;references:id;comment:主部门" json:"mainDept"`
	UserDepts []*UserDepartments `gorm:"foreignKey:user_id;references:id;comment:用户部门" json:"userDepts"`
	CreatedAt *DateTime          `gorm:"column:created_at;type:datetime;comment:创建时间" json:"createdAt" form:"createdAt"`
	UpdatedAt *DateTime          `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt" form:"updatedAt"`
	DeletedAt *DeletedAt         `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt" form:"deletedAt" swaggerignore:"true"`
}

func (*User) TableName() string {
	return TableNameUser
}

// Connection 数据库连接名称
func (m *User) Connection() string {
	return "mysql"
}
