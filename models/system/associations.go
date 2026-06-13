package models

// UserRole 用户角色关联
type UserRole struct {
	UserID int64 `json:"userId" gorm:"primaryKey"`
	RoleID int64 `json:"roleId" gorm:"primaryKey"`
}

// RoleMenu 角色菜单关联
type RoleMenu struct {
	RoleID int64 `json:"roleId" gorm:"primaryKey"`
	MenuID int64 `json:"menuId" gorm:"primaryKey"`
}

// RolePermission 角色权限关联
type RolePermission struct {
	RoleID       int64 `json:"roleId" gorm:"primaryKey"`
	PermissionID int64 `json:"permissionId" gorm:"primaryKey"`
}
