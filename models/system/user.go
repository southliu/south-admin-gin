package models

// User 用户模型
type User struct {
	ID        int64  `json:"id" gorm:"primaryKey"`
	Username  string `json:"username" gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string `json:"-" gorm:"type:varchar(255);not null"`
	Name      string `json:"name" gorm:"type:varchar(100)"`
	Phone     string `json:"phone" gorm:"type:varchar(20)"`
	Email     string `json:"email" gorm:"type:varchar(100)"`
	Status    int    `json:"status" gorm:"default:1;comment:状态 1=启用 0=禁用"`
	CreatedAt CustomTime `json:"createdAt" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt CustomTime `json:"updatedAt" gorm:"autoUpdateTime;comment:更新时间"`
	IsDeleted int        `json:"isDeleted" gorm:"default:0;comment:是否删除"`
	DeletedAt CustomTime `json:"-" gorm:"comment:删除时间"`

	// 关联关系
	Roles []Role `json:"roles" gorm:"many2many:user_role;"`
}

// TableName 指定表名
func (User) TableName() string {
	return "sys_user"
}

// LoginData 登录请求数据
type LoginData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResult 登录响应数据
type LoginResult struct {
	Token       string   `json:"token"`
	User        UserInfo  `json:"user"`
	Roles       []int    `json:"roles"`
	Permissions []string `json:"permissions"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Status   int    `json:"status"`
	Roles    []int  `json:"roles"`
}

// CreateUserDto 创建用户请求
type CreateUserDto struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   int    `json:"status"`
	RoleIds  []int  `json:"roleIds"`
}

// UpdateUserDto 更新用户请求
type UpdateUserDto struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   int    `json:"status"`
	RoleIds  []int  `json:"roleIds"`
}

// UpdatePasswordDto 更新密码请求
type UpdatePasswordDto struct {
	OldPassword    string `json:"oldPassword" binding:"required"`
	NewPassword    string `json:"newPassword" binding:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

// CommonResponse 通用响应
type CommonResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// PageResult 分页结果
type PageResult struct {
	Items      interface{} `json:"items"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"totalPages"`
}

// PaginationDto 分页请求
type PaginationDto struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"page_size"`
}

// GetPage 获取页码，默认1
func (p *PaginationDto) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetPageSize 获取每页数量，默认10
func (p *PaginationDto) GetPageSize() int {
	if p.PageSize <= 0 {
		return 10
	}
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}
