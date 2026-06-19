package services

import (
	"errors"
	"south-admin-gin/database"
	"south-admin-gin/models/system"
	"time"
)

// GetPermissionByID 根据ID查询权限
func GetPermissionByID(id int64) (*models.Permission, error) {
	var permission models.Permission
	err := database.DB.Where("id = ? AND is_deleted = 0", id).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// CreatePermission 创建权限
func CreatePermission(name, description string) (*models.Permission, error) {
	// 检查权限名是否已存在
	var existing models.Permission
	err := database.DB.Where("name = ?", name).First(&existing).Error
	if err == nil {
		return &existing, nil
	}

	permission := models.Permission{
		Name:        name,
		Description: description,
	}

	err = database.DB.Create(&permission).Error
	if err != nil {
		return nil, err
	}

	return &permission, nil
}

// UpdatePermission 更新权限
func UpdatePermission(id int64, name, description string) (*models.Permission, error) {
	permission, err := GetPermissionByID(id)
	if err != nil {
		return nil, errors.New("权限不存在")
	}

	permission.Name = name
	permission.Description = description

	err = database.DB.Save(permission).Error
	if err != nil {
		return nil, err
	}

	return permission, nil
}

// DeletePermission 删除权限
func DeletePermission(id int64) error {
	permission, err := GetPermissionByID(id)
	if err != nil {
		return errors.New("权限不存在")
	}

	return database.DB.Model(&permission).Updates(map[string]interface{}{
		"is_deleted": 1,
		"deleted_at": models.CustomTime(time.Now().Unix()),
	}).Error
}

// GetPermissionList 获取权限列表
func GetPermissionList() ([]models.Permission, error) {
	var permissions []models.Permission
	err := database.DB.Where("is_deleted = 0").
		Select("id, name, description").
		Order("created_at DESC").
		Find(&permissions).Error
	return permissions, err
}

// GetPermissionPage 获取权限分页列表
func GetPermissionPage(page, pageSize int, name string) (*models.PageResult, error) {
	var permissions []models.Permission
	var total int64

	query := database.DB.Model(&models.Permission{}).Where("is_deleted = 0")

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &models.PageResult{
		Items:      permissions,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}
