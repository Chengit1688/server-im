package model

import "gorm.io/gorm"

type AdminApi struct {
	gorm.Model
	Handle string `gorm:"column:handle;size:128;comment:接口函数"`
	Title  string `gorm:"column:title;size:128;comment:接口标题"`
	Path   string `gorm:"column:path;size:128;comment:接口地址"`
	Action string `gorm:"column:action;size:16;comment:请求类型"`
}

func (AdminApi) TableName() string {
	return "cms_admin_apis"
}

type CmsMenu struct {
	gorm.Model
	Name       string     `json:"name" gorm:"size:128;"`
	Title      string     `json:"title" gorm:"size:128;"`
	Icon       string     `json:"icon" gorm:"size:128;"`
	Path       string     `json:"path" gorm:"size:128;"`
	Paths      string     `json:"paths" gorm:"size:128;"`
	Type       int        `json:"type" gorm:"size:1;"`
	Action     string     `json:"action" gorm:"size:16;"`
	Permission string     `json:"permission" gorm:"size:255;"`
	ParentId   int        `json:"parent_id" gorm:"size:11;"`
	NoCache    int        `json:"no_cache" gorm:"size:1;"`
	Component  string     `json:"component" gorm:"size:255;"`
	Sort       int        `json:"sort" gorm:"size:4;"`
	Visible    int        `gorm:"size:1;" json:"visible"`
	Hidden     int        `gorm:"size:1;" json:"hidden"`
	IsFrame    int        `json:"is_frame" gorm:"size:1;DEFAULT:1;"`
	AdminApi   []AdminApi `json:"admin_api" gorm:"many2many:cms_menu_api_rule"`
	Apis       []int      `json:"apis" gorm:"-"`
	Params     string     `json:"params" gorm:"-"`
	Children   []CmsMenu  `json:"children,omitempty" gorm:"-"`
	IsSelect   bool       `json:"is_select" gorm:"-"`
}

type CmsMenuSlice []CmsMenu

func (x CmsMenuSlice) Len() int           { return len(x) }
func (x CmsMenuSlice) Less(i, j int) bool { return x[i].Sort < x[j].Sort }
func (x CmsMenuSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func (CmsMenu) TableName() string {
	return "cms_menus"
}

func (e *CmsMenu) GetId() interface{} {
	return e.ID
}

type CmsMenuApi struct {
	ID     int `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT"`
	MenuID int `gorm:"column:menu_id;type:int(11)"`
	ApiID  int `gorm:"column:api_id;type:int(11)"`
}

func (m *CmsMenuApi) TableName() string {
	return "cms_menu_api"
}

type CmsRole struct {
	ID       int        `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT"`
	RoleName string     `gorm:"size:128;Index:role_name,unique" json:"role_name"`
	Status   string     `gorm:"size:4;" json:"status"`
	RoleKey  string     `gorm:"size:128;Index:role_key,unique" json:"role_key"`
	RoleSort int        `json:"role_sort" gorm:""`
	Remark   string     `json:"remark" gorm:"size:255;"`
	Admin    bool       `json:"admin" gorm:"size:4;"`
	MenuIds  []int      `json:"menu_ids" gorm:"-"`
	CmsMenu  *[]CmsMenu `json:"cms_menu" gorm:"many2many:cms_menu_role;foreignKey:ID;joinForeignKey:role_id;references:ID;joinReferences:menu_id;"`
}

func (CmsRole) TableName() string {
	return "cms_role"
}

type CmsMenuRole struct {
	ID     int `gorm:"primaryKey;column:id;type:int(11);autoIncrement"`
	MenuID int `gorm:"column:menu_id;type:int(11);uniqueIndex:M_R"`
	RoleID int `gorm:"column:role_id;type:int(11);uniqueIndex:M_R"`
}

func (m *CmsMenuRole) TableName() string {
	return "cms_menu_role"
}
