package models

// Menu 菜单模型
type Menu struct {
	ID           int64  `json:"id" gorm:"primaryKey"`
	Label        string `json:"label" gorm:"type:varchar(50);not null"`
	LabelEn      string `json:"labelEn" gorm:"type:varchar(50);not null"`
	Icon         string `json:"icon" gorm:"type:varchar(50)"`
	Type         int    `json:"type" gorm:"type:int;comment:1=目录,2=菜单,3=按钮"`
	Router       string `json:"router" gorm:"type:varchar(255)"`
	Rule         string `json:"rule" gorm:"type:varchar(255)"`
	Order        int    `json:"order" gorm:"type:int;default:0"`
	State        int    `json:"state" gorm:"type:int;default:1;comment:0=隐藏,1=显示"`
	ParentID     *int64     `json:"parentId" gorm:"index"`
	PermissionID *int64     `json:"permissionId" gorm:"index"`
	CreatedAt    CustomTime `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    CustomTime `json:"updatedAt" gorm:"autoUpdateTime"`
	IsDeleted    int        `json:"isDeleted" gorm:"default:0;comment:是否删除"`
	DeletedAt    CustomTime `json:"-" gorm:"comment:删除时间"`

	// 关联关系
	Parent     *Menu       `json:"parent" gorm:"foreignKey:ParentID"`
	Permission *Permission `json:"permission" gorm:"foreignKey:PermissionID"`
	Children   []Menu      `json:"children" gorm:"foreignKey:ParentID"`
	Key        string      `json:"key" gorm:"-"`
	Title      string      `json:"title" gorm:"-"`
	TitleEn    string      `json:"titleEn" gorm:"-"`
	Value      string      `json:"value" gorm:"-"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "sys_menu"
}

// CreateMenuDto 创建菜单请求
type CreateMenuDto struct {
	Label     string   `json:"label" binding:"required"`
	LabelEn   string   `json:"labelEn" binding:"required"`
	Type      int      `json:"type" binding:"required"`
	Icon      string   `json:"icon"`
	Router    string   `json:"router"`
	Rule      string   `json:"rule"`
	Order     int      `json:"order"`
	State     int      `json:"state"`
	ParentID  *int64   `json:"parentId"`
	Actions   []string `json:"actions"`
}

// UpdateMenuDto 更新菜单请求
type UpdateMenuDto struct {
	Label     string `json:"label" binding:"required"`
	LabelEn   string `json:"labelEn" binding:"required"`
	Type      int    `json:"type"`
	Icon      string `json:"icon"`
	Router    string `json:"router"`
	Rule      string `json:"rule"`
	Order     int    `json:"order"`
	State     int    `json:"state"`
	ParentID  *int64 `json:"parentId"`
}

// ChangeMenuStateDto 修改菜单状态请求
type ChangeMenuStateDto struct {
	ID    int64 `json:"id" binding:"required"`
	State int   `json:"state" binding:"required"`
}
