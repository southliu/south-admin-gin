package models

// Role 角色模型
type Role struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"type:varchar(50);uniqueIndex;not null"`
	Description string `json:"description" gorm:"type:text"`
	CreatedAt   CustomTime `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   CustomTime `json:"updatedAt" gorm:"autoUpdateTime"`
	IsDeleted   int        `json:"isDeleted" gorm:"default:0;comment:是否删除"`
	DeletedAt   CustomTime `json:"-" gorm:"comment:删除时间"`

	// 关联关系
	Menus       []Menu       `json:"menus" gorm:"many2many:role_menu;"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permission;"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "sys_role"
}

// CreateRoleDto 创建角色请求
type CreateRoleDto struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Authorize   []int  `json:"authorize"`
}

// UpdateRoleDto 更新角色请求
type UpdateRoleDto struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Authorize   []int  `json:"authorize"`
}

// AuthorizeRoleDto 角色授权请求
type AuthorizeRoleDto struct {
	RoleId  int   `json:"roleId" binding:"required"`
	MenuIds []int `json:"menuIds" binding:"required"`
}
