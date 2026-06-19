package services

import (
	"errors"
	"south-admin-gin/database"
	"south-admin-gin/models/system"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// GetUserByUsername 根据用户名查询用户
func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("username = ? AND is_deleted = 0", username).
		Preload("Roles").
		Preload("Roles.Menus.Permission").
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据ID查询用户
func GetUserByID(id int64) (*models.User, error) {
	var user models.User
	err := database.DB.Where("id = ? AND is_deleted = 0", id).
		Preload("Roles").
		Preload("Roles.Menus.Permission").
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CheckPassword 校验密码是否正确（bcrypt）
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// HashPassword 对明文密码进行 bcrypt 哈希
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CreateUser 创建用户
func CreateUser(dto models.CreateUserDto) (*models.User, error) {
	// 检查用户名是否已存在
	var count int64
	database.DB.Model(&models.User{}).Where("username = ? AND is_deleted = 0", dto.Username).Count(&count)
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username:  dto.Username,
		Password:  hashedPassword,
		Name:      dto.Name,
		Email:     dto.Email,
		Phone:     dto.Phone,
		Status:    dto.Status,
		IsDeleted: 0,
	}

	// 设置默认状态
	if user.Status == 0 {
		user.Status = 1
	}

	err = database.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}

	// 关联角色
	if len(dto.RoleIds) > 0 {
		if err := updateUserRoles(&user, dto.RoleIds); err != nil {
			return nil, err
		}
	}

	return &user, nil
}

// UpdateUser 更新用户
func UpdateUser(id int64, dto models.UpdateUserDto) (*models.User, error) {
	user, err := GetUserByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查用户名是否已存在
	if dto.Username != user.Username {
		var count int64
		database.DB.Model(&models.User{}).Where("username = ? AND id != ? AND is_deleted = 0", dto.Username, id).Count(&count)
		if count > 0 {
			return nil, errors.New("用户名已存在")
		}
	}

	user.Username = dto.Username

	// 更新密码
	if dto.Password != "" {
		hashedPassword, err := HashPassword(dto.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	if dto.Name != "" {
		user.Name = dto.Name
	}
	if dto.Email != "" {
		user.Email = dto.Email
	}
	if dto.Phone != "" {
		user.Phone = dto.Phone
	}
	if dto.Status != 0 {
		user.Status = dto.Status
	}

	err = database.DB.Save(user).Error
	if err != nil {
		return nil, err
	}

	// 更新角色关联
	if dto.RoleIds != nil {
		if err := updateUserRoles(user, dto.RoleIds); err != nil {
			return nil, err
		}
	}

	return user, nil
}

// DeleteUser 删除用户
func DeleteUser(id int64) error {
	user, err := GetUserByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	return database.DB.Model(&user).Updates(map[string]interface{}{
		"is_deleted": 1,
		"deleted_at": models.CustomTime(time.Now().Unix()),
	}).Error
}

// GetUserList 获取用户列表
func GetUserList() ([]models.User, error) {
	var users []models.User
	err := database.DB.Where("is_deleted = 0 AND status = 1").
		Select("id, username, name, email, phone").
		Find(&users).Error
	return users, err
}

// GetUserPage 获取用户分页列表
func GetUserPage(page, pageSize int, username string) (*models.PageResult, error) {
	var users []models.User
	var total int64

	query := database.DB.Model(&models.User{}).Where("is_deleted = 0")

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).
		Preload("Roles").
		Order("created_at DESC").
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	// 添加角色名称
	type UserWithRolesName struct {
		models.User
		RolesName string `json:"rolesName"`
		RoleCount int    `json:"roleCount"`
	}

	var items []UserWithRolesName
	for _, user := range users {
		roleNames := make([]string, 0)
		for _, role := range user.Roles {
			roleNames = append(roleNames, role.Name)
		}
		rolesName := ""
		if len(roleNames) > 0 {
			for i, name := range roleNames {
				if i > 0 {
					rolesName += ","
				}
				rolesName += name
			}
		}
		items = append(items, UserWithRolesName{
			User:      user,
			RolesName: rolesName,
			RoleCount: len(user.Roles),
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

// UpdatePassword 更新密码
func UpdatePassword(userId int64, dto models.UpdatePasswordDto) error {
	if dto.NewPassword != dto.ConfirmPassword {
		return errors.New("新旧密码不一致")
	}

	user, err := GetUserByID(userId)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if !CheckPassword(user.Password, dto.OldPassword) {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := HashPassword(dto.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return database.DB.Save(user).Error
}

// RefreshPermissions 刷新用户权限
func RefreshPermissions(userId int64) (map[string]interface{}, error) {
	user, err := GetUserByID(userId)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	permissions := getUserPermissions(user)
	roleIds := getUserRoleIds(user)

	return map[string]interface{}{
		"user": models.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
			Phone:    user.Phone,
			Email:    user.Email,
			Status:   user.Status,
			Roles:    roleIds,
		},
		"permissions": permissions,
		"roles":       roleIds,
	}, nil
}

// updateUserRoles 更新用户角色关联
func updateUserRoles(user *models.User, roleIds []int) error {
	// 清除原有关联
	database.DB.Table("user_role").Where("user_id = ?", user.ID).Delete(nil)

	if len(roleIds) == 0 {
		return nil
	}

	// 获取角色
	roles, err := GetRolesByIDs(roleIds)
	if err != nil {
		return err
	}

	// 关联角色
	for _, role := range roles {
		database.DB.Table("user_role").Create(map[string]interface{}{
			"user_id": user.ID,
			"role_id": role.ID,
		})
	}

	return nil
}

// GetUserPermissions 获取用户权限列表（公开方法）
func GetUserPermissions(user *models.User) []string {
	return getUserPermissions(user)
}

// getUserPermissions 获取用户权限列表（通过角色的菜单获取权限）
func getUserPermissions(user *models.User) []string {
	permissionSet := make(map[string]bool)

	for _, role := range user.Roles {
		for _, menu := range role.Menus {
			if menu.Permission != nil {
				permissionSet[menu.Permission.Name] = true
			}
		}
	}

	permissions := make([]string, 0)
	for name := range permissionSet {
		permissions = append(permissions, name)
	}
	return permissions
}

// getUserRoleIds 获取用户角色ID列表
func getUserRoleIds(user *models.User) []int {
	roleIds := make([]int, 0)
	for _, role := range user.Roles {
		roleIds = append(roleIds, int(role.ID))
	}
	return roleIds
}
