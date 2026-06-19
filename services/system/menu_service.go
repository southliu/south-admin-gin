package services

import (
	"errors"
	"south-admin-gin/database"
	"south-admin-gin/models/system"
	"strconv"
	"time"
)

// GetMenuByID 根据ID查询菜单
func GetMenuByID(id int64) (*models.Menu, error) {
	var menu models.Menu
	err := database.DB.Where("id = ? AND is_deleted = 0", id).
		Preload("Parent").
		Preload("Permission").
		First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// CreateMenu 创建菜单
func CreateMenu(dto models.CreateMenuDto, userId int64) (*models.Menu, error) {
	// 验证父级菜单
	if dto.ParentID != nil {
		parent, err := GetMenuByID(*dto.ParentID)
		if err != nil || parent.IsDeleted == 1 {
			return nil, errors.New("父级菜单不存在")
		}
	}

	menu := models.Menu{
		Label:    dto.Label,
		LabelEn:  dto.LabelEn,
		Type:     dto.Type,
		Icon:     dto.Icon,
		Router:   dto.Router,
		Rule:     dto.Rule,
		Order:    dto.Order,
		State:    dto.State,
		ParentID: dto.ParentID,
	}

	// 设置默认状态
	if menu.State == 0 {
		menu.State = 1
	}

	err := database.DB.Create(&menu).Error
	if err != nil {
		return nil, err
	}

	// 处理权限
	if dto.Rule != "" {
		permission, err := getOrCreatePermission(dto.Rule, dto.Label)
		if err != nil {
			return nil, err
		}
		menu.PermissionID = &permission.ID
		database.DB.Save(&menu)
	}

	// 将菜单关联到用户的角色
	if err := associateMenuToUserRoles(userId, menu.ID); err != nil {
		return nil, err
	}

	// 处理按钮菜单
	if len(dto.Actions) > 0 {
		if err := createButtonMenus(&menu, dto.Actions, dto.Rule, userId); err != nil {
			return nil, err
		}
	}

	return &menu, nil
}

// UpdateMenu 更新菜单
func UpdateMenu(id int64, dto models.UpdateMenuDto) (*models.Menu, error) {
	menu, err := GetMenuByID(id)
	if err != nil {
		return nil, errors.New("菜单不存在")
	}

	// 验证父级菜单
	if dto.ParentID != nil {
		if *dto.ParentID == id {
			return nil, errors.New("菜单不能是它自己的父菜单")
		}
		if *dto.ParentID != 0 {
			parent, err := GetMenuByID(*dto.ParentID)
			if err != nil || parent.IsDeleted == 1 {
				return nil, errors.New("父级菜单不存在")
			}
		}
	}

	menu.Label = dto.Label
	menu.LabelEn = dto.LabelEn
	if dto.Type != 0 {
		menu.Type = dto.Type
	}
	menu.Icon = dto.Icon
	menu.Router = dto.Router
	menu.Rule = dto.Rule
	if dto.Order != 0 {
		menu.Order = dto.Order
	}
	if dto.State != 0 {
		menu.State = dto.State
	}
	menu.ParentID = dto.ParentID

	// 处理权限
	if dto.Rule != "" {
		// 如果菜单已有权限，更新现有权限的name和description
		if menu.PermissionID != nil {
			var existingPermission models.Permission
			err := database.DB.Where("id = ?", *menu.PermissionID).First(&existingPermission).Error
			if err == nil {
				// 检查新名称是否已被其他权限使用
				if existingPermission.Name != dto.Rule {
					var duplicatePermission models.Permission
					checkErr := database.DB.Where("name = ? AND id != ?", dto.Rule, *menu.PermissionID).First(&duplicatePermission).Error
					if checkErr == nil {
						return nil, errors.New("权限名称 " + dto.Rule + " 已被使用")
					}
				}

				// 更新现有权限的name和description
				existingPermission.Name = dto.Rule
				existingPermission.Description = dto.Label
				err = database.DB.Save(&existingPermission).Error
				if err != nil {
					return nil, err
				}
			} else {
				// 如果权限不存在，创建新权限
				permission, err := getOrCreatePermission(dto.Rule, dto.Label)
				if err != nil {
					return nil, err
				}
				menu.PermissionID = &permission.ID
			}
		} else {
			// 如果菜单没有权限，创建新权限
			permission, err := getOrCreatePermission(dto.Rule, dto.Label)
			if err != nil {
				return nil, err
			}
			menu.PermissionID = &permission.ID
		}
	}

	return menu, database.DB.Save(menu).Error
}

// DeleteMenu 删除菜单
func DeleteMenu(id int64) error {
	menu, err := GetMenuByID(id)
	if err != nil {
		return errors.New("菜单不存在")
	}

	// 检查是否有子菜单
	var childCount int64
	database.DB.Model(&models.Menu{}).Where("parent_id = ? AND is_deleted = 0", id).Count(&childCount)
	if childCount > 0 {
		return errors.New("不能删除有子菜单的菜单")
	}

	// 删除菜单关联关系
	database.DB.Table("role_menu").Where("menu_id = ?", id).Delete(nil)

	// 查找并删除权限关联关系
	// 优先通过 PermissionID 查找，兜底通过 Rule(name) 查找
	var permission models.Permission
	permissionFound := false
	if menu.PermissionID != nil {
		err := database.DB.Where("id = ?", *menu.PermissionID).First(&permission).Error
		if err == nil {
			permissionFound = true
		}
	}
	if !permissionFound && menu.Rule != "" {
		err := database.DB.Where("name = ? AND is_deleted = 0", menu.Rule).First(&permission).Error
		if err == nil {
			permissionFound = true
		}
	}
	if permissionFound {
		// 软删除权限
		database.DB.Model(&permission).Updates(map[string]interface{}{
			"is_deleted": 1,
			"deleted_at": models.CustomTime(time.Now().Unix()),
		})
	}

	// 软删除菜单
	return database.DB.Model(&menu).Updates(map[string]interface{}{
		"is_deleted": 1,
		"deleted_at": models.CustomTime(time.Now().Unix()),
	}).Error
}

// BatchDeleteMenu 批量删除菜单
func BatchDeleteMenu(ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的菜单")
	}

	// 检查是否有子菜单
	var childCount int64
	database.DB.Model(&models.Menu{}).Where("parent_id IN ? AND is_deleted = 0", ids).Count(&childCount)
	if childCount > 0 {
		return errors.New("不能删除有子菜单的菜单")
	}

	// 获取要删除的菜单及其rule字段
	var menus []models.Menu
	err := database.DB.Where("id IN ? AND is_deleted = 0", ids).Find(&menus).Error
	if err != nil {
		return err
	}

	// 收集所有rule值
	var rules []string
	for _, menu := range menus {
		if menu.Rule != "" {
			rules = append(rules, menu.Rule)
		}
	}

	// 删除菜单关联关系
	database.DB.Table("role_menu").Where("menu_id IN ?", ids).Delete(nil)

	// 查找并删除权限关联关系
	// 优先通过 PermissionID 查找，兜底通过 Rule(name) 查找
	permissionIdSet := make(map[int64]bool)
	var permissionIds []int64
	for _, menu := range menus {
		if menu.PermissionID != nil && !permissionIdSet[*menu.PermissionID] {
			permissionIdSet[*menu.PermissionID] = true
			permissionIds = append(permissionIds, *menu.PermissionID)
		}
	}
	// 兜底：通过 Rule 查找尚未收集的权限
	if len(rules) > 0 {
		var permissions []models.Permission
		database.DB.Where("name IN ? AND is_deleted = 0", rules).Find(&permissions)
		for _, p := range permissions {
			if !permissionIdSet[p.ID] {
				permissionIdSet[p.ID] = true
				permissionIds = append(permissionIds, p.ID)
			}
		}
	}

	if len(permissionIds) > 0 {
		// 批量软删除权限
		database.DB.Model(&models.Permission{}).Where("id IN ?", permissionIds).Updates(map[string]interface{}{
			"is_deleted": 1,
			"deleted_at": models.CustomTime(time.Now().Unix()),
		})
	}

	// 批量软删除菜单
	return database.DB.Model(&models.Menu{}).Where("id IN ?", ids).Updates(map[string]interface{}{
		"is_deleted": 1,
		"deleted_at": models.CustomTime(time.Now().Unix()),
	}).Error
}

// GetMenuDetail 获取菜单详情
func GetMenuDetail(id int64) (*models.Menu, error) {
	menu, err := GetMenuByID(id)
	if err != nil {
		return nil, err
	}

	// 转换permission为rule字段
	if menu.Permission != nil {
		menu.Rule = menu.Permission.Name
		menu.Permission = nil
	}

	return menu, nil
}

// GetMenuList 获取用户菜单列表
func GetMenuList(userId int64) ([]models.Menu, error) {
	// 获取用户信息及角色
	var user models.User
	err := database.DB.Where("id = ?", userId).
		Preload("Roles").
		First(&user).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if len(user.Roles) == 0 {
		return []models.Menu{}, nil
	}

	// 获取角色关联的菜单ID
	var roleMenuIds []int64
	for _, role := range user.Roles {
		var menuIds []int64
		database.DB.Table("role_menu").Where("role_id = ?", role.ID).Pluck("menu_id", &menuIds)
		roleMenuIds = append(roleMenuIds, menuIds...)
	}

	// 去重
	uniqueMenuIds := make(map[int64]bool)
	var uniqueIds []int64
	for _, id := range roleMenuIds {
		if !uniqueMenuIds[id] {
			uniqueMenuIds[id] = true
			uniqueIds = append(uniqueIds, id)
		}
	}

	if len(uniqueIds) == 0 {
		return []models.Menu{}, nil
	}

	// 获取所有有效菜单
	var allMenus []models.Menu
	err = database.DB.Where("id IN ? AND is_deleted = 0 AND state = 1 AND type != 3", uniqueIds).
		Preload("Parent").
		Preload("Permission").
		Order("`order` ASC").
		Find(&allMenus).Error
	if err != nil {
		return nil, err
	}

	// 获取所有父级菜单ID
	parentIds := make(map[int64]bool)
	for _, menu := range allMenus {
		if menu.ParentID != nil {
			parentIds[*menu.ParentID] = true
		}
	}

	// 确保显示的菜单ID列表：包含关联的菜单ID + 它们的父级菜单ID
	menuIdsToShow := make(map[int64]bool)
	for _, id := range uniqueIds {
		menuIdsToShow[id] = true
	}
	for id := range parentIds {
		menuIdsToShow[id] = true
	}

	// 过滤出需要显示的菜单
	var filteredMenus []models.Menu
	for _, menu := range allMenus {
		if menuIdsToShow[menu.ID] {
			// 转换permission为rule字段
			if menu.Permission != nil {
				menu.Rule = menu.Permission.Name
				menu.Permission = nil
			}
			// 设置前端需要的字段
			menu.Key = menu.Router
			if menu.Key == "" {
				menu.Key = strconv.FormatInt(menu.ID, 10)
			}
			menu.Title = menu.Label
			menu.TitleEn = menu.LabelEn
			menu.Value = strconv.FormatInt(menu.ID, 10)
			filteredMenus = append(filteredMenus, menu)
		}
	}

	return buildMenuTree(filteredMenus, nil), nil
}

// GetMenuPage 获取菜单分页列表
func GetMenuPage(page, pageSize int, label, labelEn string, state *int, rule string) (map[string]interface{}, error) {
	var menus []models.Menu

	query := database.DB.Model(&models.Menu{}).Where("is_deleted = 0")

	if label != "" {
		query = query.Where("label LIKE ?", "%"+label+"%")
	}
	if labelEn != "" {
		query = query.Where("label_en LIKE ?", "%"+labelEn+"%")
	}
	if state != nil {
		query = query.Where("state = ?", *state)
	}
	if rule != "" {
		query = query.Where("rule LIKE ?", "%"+rule+"%")
	}

	err := query.Preload("Parent").
		Preload("Permission").
		Order("`order` ASC").
		Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// 转换permission为rule字段
	for i := range menus {
		if menus[i].Permission != nil {
			menus[i].Rule = menus[i].Permission.Name
			menus[i].Permission = nil
		}
	}

	tree := buildMenuTree(menus, nil)

	// 计算总数
	var total int64
	database.DB.Model(&models.Menu{}).
		Where("is_deleted = 0 AND ((type = 1 AND parent_id IS NULL) OR (type = 2 AND parent_id IS NULL))").
		Count(&total)

	return map[string]interface{}{
		"items": tree,
		"total": total,
	}, nil
}

// ChangeMenuState 修改菜单状态
func ChangeMenuState(dto models.ChangeMenuStateDto) error {
	menu, err := GetMenuByID(dto.ID)
	if err != nil {
		return errors.New("菜单不存在")
	}

	menu.State = dto.State
	return database.DB.Save(menu).Error
}

// GetMenuTree 获取菜单树
func GetMenuTree() ([]models.Menu, error) {
	var menus []models.Menu
	err := database.DB.Where("is_deleted = 0 AND state = 1").
		Preload("Parent").
		Preload("Permission").
		Order("`order` ASC").
		Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// 转换permission为rule字段，并设置前端需要的字段
	for i := range menus {
		if menus[i].Permission != nil {
			menus[i].Rule = menus[i].Permission.Name
			menus[i].Permission = nil
		}
		// 设置前端需要的字段
		menus[i].Key = strconv.FormatInt(menus[i].ID, 10)
		menus[i].Title = menus[i].Label
		menus[i].TitleEn = menus[i].LabelEn
		menus[i].Value = strconv.FormatInt(menus[i].ID, 10)
	}

	return buildMenuTree(menus, nil), nil
}

// buildMenuTree 构建菜单树
func buildMenuTree(menus []models.Menu, parentId *int64) []models.Menu {
	var tree []models.Menu
	for _, menu := range menus {
		var menuParentId *int64
		if menu.ParentID != nil {
			menuParentId = menu.ParentID
		}

		if (parentId == nil && menuParentId == nil) ||
			(parentId != nil && menuParentId != nil && *parentId == *menuParentId) {
			children := buildMenuTree(menus, &menu.ID)
			if len(children) > 0 {
				menu.Children = children
			}
			tree = append(tree, menu)
		}
	}
	return tree
}

// getOrCreatePermission 获取或创建权限
func getOrCreatePermission(name, description string) (*models.Permission, error) {
	var permission models.Permission
	err := database.DB.Where("name = ?", name).First(&permission).Error
	if err == nil {
		return &permission, nil
	}

	permission = models.Permission{
		Name:        name,
		Description: description,
	}
	err = database.DB.Create(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// associateMenuToUserRoles 将菜单关联到用户的角色
func associateMenuToUserRoles(userId int64, menuId int64) error {
	var user models.User
	err := database.DB.Where("id = ?", userId).Preload("Roles").First(&user).Error
	if err != nil {
		return err
	}

	for _, role := range user.Roles {
		// 检查是否已存在关联
		var count int64
		database.DB.Table("role_menu").
			Where("role_id = ? AND menu_id = ?", role.ID, menuId).
			Count(&count)
		if count == 0 {
			database.DB.Table("role_menu").Create(map[string]interface{}{
				"role_id": role.ID,
				"menu_id": menuId,
			})
		}
	}
	return nil
}

// createButtonMenus 创建按钮菜单
func createButtonMenus(parentMenu *models.Menu, actions []string, rule string, userId int64) error {
	buttonType := 3
	actionLabels := map[string][2]string{
		"create": {"新增", "Create"},
		"update": {"修改", "Update"},
		"delete": {"删除", "Delete"},
		"detail": {"详情", "Detail"},
		"export": {"导出权限", "Export"},
		"status": {"状态权限", "Status"},
	}

	for _, action := range actions {
		labels, ok := actionLabels[action]
		if !ok {
			continue
		}

		buttonRule := rule + "/" + action
		buttonMenu := models.Menu{
			Label:    labels[0],
			LabelEn:  labels[1],
			Type:     buttonType,
			Rule:     buttonRule,
			Router:   "",
			Order:    0,
			State:    1,
			ParentID: &parentMenu.ID,
		}

		// 创建权限
		permission, err := getOrCreatePermission(buttonRule, labels[0])
		if err != nil {
			return err
		}

		// 关联权限到菜单
		buttonMenu.PermissionID = &permission.ID

		err = database.DB.Create(&buttonMenu).Error
		if err != nil {
			return err
		}

		// 关联菜单到用户的角色
		if err := associateMenuToUserRoles(userId, buttonMenu.ID); err != nil {
			return err
		}
	}

	return nil
}
