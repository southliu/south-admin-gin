package services

import (
	"errors"
	"serve-wechat-gin/database"
	"serve-wechat-gin/models/system"
	"strconv"
	"time"
)

// GetRoleByName 根据角色名查询角色
func GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	err := database.DB.Where("name = ? AND is_deleted = 0", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRoleByID 根据ID查询角色
func GetRoleByID(id int) (*models.Role, error) {
	var role models.Role
	err := database.DB.Where("id = ? AND is_deleted = 0", id).
		Preload("Menus").
		Preload("Permissions").
		First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRolesByIDs 根据ID列表查询角色
func GetRolesByIDs(ids []int) ([]models.Role, error) {
	var roles []models.Role
	err := database.DB.Where("id IN ? AND is_deleted = 0", ids).
		Preload("Permissions").
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// CreateRole 创建角色
func CreateRole(dto models.CreateRoleDto) (*models.Role, error) {
	// 检查角色名是否已存在
	var count int64
	database.DB.Model(&models.Role{}).Where("name = ? AND is_deleted = 0", dto.Name).Count(&count)
	if count > 0 {
		return nil, errors.New("角色名已存在")
	}

	role := models.Role{
		Name:        dto.Name,
		Description: dto.Description,
		IsDeleted:   0,
	}

	err := database.DB.Create(&role).Error
	if err != nil {
		return nil, err
	}

	// 关联菜单和权限
	if len(dto.Authorize) > 0 {
		if err := updateRoleAuthorize(role.ID, dto.Authorize); err != nil {
			return nil, err
		}
	}

	return &role, nil
}

// UpdateRole 更新角色
func UpdateRole(id int, dto models.UpdateRoleDto) (*models.Role, error) {
	role, err := GetRoleByID(id)
	if err != nil {
		return nil, errors.New("角色不存在")
	}

	// admin角色不允许修改
	if role.Name == "admin" {
		return nil, errors.New("不允许修改admin角色")
	}

	// 检查角色名是否已存在
	if dto.Name != role.Name {
		var count int64
		database.DB.Model(&models.Role{}).Where("name = ? AND id != ? AND is_deleted = 0", dto.Name, id).Count(&count)
		if count > 0 {
			return nil, errors.New("角色名已存在")
		}
	}

	role.Name = dto.Name
	role.Description = dto.Description

	err = database.DB.Save(role).Error
	if err != nil {
		return nil, err
	}

	// 更新菜单和权限关联
	if dto.Authorize != nil {
		if err := updateRoleAuthorize(int64(id), dto.Authorize); err != nil {
			return nil, err
		}
	}

	return role, nil
}

// DeleteRole 删除角色
func DeleteRole(id int) error {
	role, err := GetRoleByID(id)
	if err != nil {
		return errors.New("角色不存在")
	}

	// admin角色不允许删除
	if role.Name == "admin" {
		return errors.New("不允许删除admin角色")
	}

	return database.DB.Model(&role).Updates(map[string]interface{}{
		"is_deleted": 1,
		"deleted_at": models.CustomTime(time.Now().Unix()),
	}).Error
}

// GetRoleList 获取角色列表
func GetRoleList() ([]models.Role, error) {
	var roles []models.Role
	err := database.DB.Where("is_deleted = 0").
		Select("id, name, description, created_at, updated_at").
		Order("created_at DESC").
		Find(&roles).Error
	return roles, err
}

// GetRolePage 获取角色分页列表
func GetRolePage(page, pageSize int, name string) (*models.PageResult, error) {
	var roles []models.Role
	var total int64

	query := database.DB.Model(&models.Role{}).Where("is_deleted = 0")

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).
		Preload("Permissions").
		Preload("Menus").
		Order("created_at DESC").
		Find(&roles).Error
	if err != nil {
		return nil, err
	}

	// 添加统计信息
	type RoleWithCounts struct {
		models.Role
		PermissionCount int `json:"permissionCount"`
		MenuCount       int `json:"menuCount"`
	}

	var items []RoleWithCounts
	for _, role := range roles {
		items = append(items, RoleWithCounts{
			Role:            role,
			PermissionCount: len(role.Permissions),
			MenuCount:       len(role.Menus),
		})
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &models.PageResult{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// GetRoleAuthorize 获取角色授权信息
func GetRoleAuthorize(roleId int) (map[string]interface{}, error) {
	role, err := GetRoleByID(roleId)
	if err != nil {
		return nil, errors.New("角色不存在")
	}

	menuIds := make([]string, 0)
	for _, menu := range role.Menus {
		menuIds = append(menuIds, strconv.FormatInt(menu.ID, 10))
	}

	// 获取菜单树
	treeData, err := GetMenuTree()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"defaultCheckedKeys": menuIds,
		"treeData":           treeData,
	}, nil
}

// SaveRoleAuthorize 保存角色授权
func SaveRoleAuthorize(roleId int, menuIds []int) error {
	role, err := GetRoleByID(roleId)
	if err != nil {
		return errors.New("角色不存在")
	}

	// admin角色不允许修改授权
	if role.Name == "admin" {
		return errors.New("不允许修改admin角色授权")
	}

	return updateRoleAuthorize(int64(roleId), menuIds)
}

// updateRoleAuthorize 更新角色的菜单和权限关联
func updateRoleAuthorize(roleId int64, menuIds []int) error {
	// 清除原有关联
	database.DB.Table("role_menu").Where("role_id = ?", roleId).Delete(nil)
	database.DB.Table("role_permission").Where("role_id = ?", roleId).Delete(nil)

	if len(menuIds) == 0 {
		return nil
	}

	// 查询菜单及其关联的权限
	var menus []models.Menu
	err := database.DB.Where("id IN ? AND is_deleted = 0", menuIds).
		Preload("Permission").
		Find(&menus).Error
	if err != nil {
		return err
	}

	// 关联菜单
	for _, menu := range menus {
		database.DB.Table("role_menu").Create(map[string]interface{}{
			"role_id": roleId,
			"menu_id": menu.ID,
		})

		// 关联权限
		if menu.PermissionID != nil {
			database.DB.Table("role_permission").Create(map[string]interface{}{
				"role_id":       roleId,
				"permission_id": *menu.PermissionID,
			})
		}
	}

	return nil
}
