package model

import (
	"im/pkg/pagination"
	"time"
)

type GetinfoResp struct {
	Username    string            `json:"username"`
	UserID      string            `json:"user_id"`
	NickName    string            `json:"nick_name"`
	PhoneNumber string            `json:"phone_number"`
	Menus       []GetinfoMenuResp `json:"menus"`
}

type GetinfoMenuResp struct {
	MenuID     int               `json:"menu_id"`
	MenuName   string            `json:"menu_name"`
	Title      string            `json:"title"`
	Icon       string            `json:"icon"`
	Path       string            `json:"path"`
	MenuType   int               `json:"menu_type"`
	Action     string            `json:"action"`
	Permission string            `json:"permission"`
	ParentId   int               `json:"parent_id"`
	NoCache    int               `json:"no_cache"`
	Component  string            `json:"component"`
	Sort       int               `json:"sort"`
	Visible    int               `json:"visible"`
	Hidden     int               `json:"hidden"`
	IsFrame    int               `json:"is_frame"`
	Children   []GetinfoMenuResp `json:"children,omitempty"`
}

type LoginReq struct {
	OperationID string `json:"operation_id"  binding:"required" msg:"操作ID不能为空"`
	Username    string `json:"username" binding:"required" msg:"账号不能为空"`
	Password    string `json:"password" binding:"required" msg:"密码不能为空"`
	GoogleCode  int32  `json:"google_code"`
}

type LoginResp struct {
	Token  string `json:"token"`
	Expire int64  `json:"expire"`
}

type RefreshTokenReq struct {
	OperationID string `json:"operation_id"  binding:"required" msg:"操作ID不能为空"`
}

type ListReq struct {
	OperationID    string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Username       string `json:"username" form:"username"`
	Nickname       string `json:"nick_name" form:"nick_name"`
	RoleKey        string `json:"role_key" form:"role_key"`
	LoginTimeStart int64  `json:"login_time_start" form:"login_time_start"`
	LoginTimeEnd   int64  `json:"login_time_end" form:"login_time_end"`
	pagination.Pagination
}

type ListRespItem struct {
	ID            int    `json:"id"`
	Nickname      string `json:"nick_name" gorm:"column:nick_name"`
	Username      string `json:"username"`
	RoleKey       string `json:"role_key" gorm:"column:role_key"`
	LastloginIp   string `json:"last_login_ip" gorm:"column:last_login_ip"`
	LastloginTime int64  `json:"last_login_time" gorm:"column:last_login_time"`
}

type ListAdminItem struct {
	ID            int       `json:"id"`
	Nickname      string    `json:"nick_name" gorm:"column:nick_name"`
	Username      string    `json:"username"`
	RoleKey       string    `json:"role_key" gorm:"column:role_key"`
	LastloginIp   string    `json:"last_login_ip" gorm:"column:last_login_ip"`
	LastloginTime time.Time `json:"last_login_time" gorm:"column:last_login_time"`
}

type ListResp struct {
	List     []ListRespItem `json:"list"`
	Count    int64          `json:"count"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

type AddReq struct {
	OperationID string `json:"operation_id"  binding:"required" msg:"操作ID不能为空"`
	Nickname    string `json:"nick_name"  binding:"required,min=1,max=16" msg:"管理员名称不能为空"`
	Username    string `json:"username" binding:"required,min=4,max=16" msg:"管理员账号不能为空"`
	Password    string `json:"password" binding:"required,min=6,max=16" msg:"管理员密码不能空"`
	RoleID      int    `json:"role_id" binding:"required" msg:"管理员角色不能空"`
}

type AddResp struct {
	ID       int    `json:"id"`
	Nickname string `json:"nick_name"`
	Username string `json:"username"`
	RoleID   int    `json:"role_id"`
}

type UpdateInfoReq struct {
	OperationID string `json:"operation_id"  binding:"required" msg:"操作ID不能为空"`
	ID          int    `json:"id"  binding:"required" msg:"ID不能为空"`
	Nickname    string `json:"nick_name"  binding:"required,min=1,max=16" msg:"管理员名称不能为空"`
	Username    string `json:"username" binding:"required,min=4,max=16" msg:"管理员账号不能为空"`
	RoleID      int    `json:"role_id" binding:"required" msg:"管理员角色不能空"`
}

type UpdatePasswordReq struct {
	OperationID string `json:"operation_id"  binding:"required" msg:"操作ID不能为空"`
	ID          int    `json:"id"  binding:"required" msg:"ID不能为空"`
	Password    string `json:"password" binding:"required,min=6,max=16" msg:"管理员密码不能空"`
}

type DeleteReq struct {
	OperationID string `json:"operation_id"  binding:"required" msg:"操作ID不能为空"`
	ID          int    `json:"id"  binding:"required" msg:"ID不能为空"`
}

type GetGoogleCodeReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
}

type GetGoogleCodeResp struct {
	Username string `json:"username"`
	Secret   string `json:"secret"`
	Image    string `json:"image"`
}
