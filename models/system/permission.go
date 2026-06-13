package models

// Permission 权限模型
type Permission struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"type:varchar(100);uniqueIndex;not null"`
	Description string `json:"description" gorm:"type:text"`
	CreatedAt   CustomTime `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   CustomTime `json:"updatedAt" gorm:"autoUpdateTime"`
	IsDeleted   int        `json:"isDeleted" gorm:"default:0;comment:是否删除"`
	DeletedAt   CustomTime `json:"-" gorm:"comment:删除时间"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "sys_permission"
}

// CreatePermissionDto 创建权限请求
type CreatePermissionDto struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdatePermissionDto 更新权限请求
type UpdatePermissionDto struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
