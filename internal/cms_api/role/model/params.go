package model

import "im/pkg/pagination"

type GetRoleReq struct {
	OperationID string `json:"operation_id"  binding:"required" msg:"操作ID不能为空"`
	RoleName    string `json:"role_name"`
	RoleKey     string `json:"role_key"`
	Page        string `json:"page"`
	PageSize    string `json:"page_size"`
}

type GetRoleResp struct {
	List  []GetRoleItem `json:"list"`
	Count int           `json:"count"`
}

type GetRoleItem struct {
	ID       int    `json:"id"`
	RoleName string `json:"role_name"`
	RoleKey  string `json:"role_key"`
}

type AddRoleResp struct {
	OperationID string `json:"operation_id"  binding:"required" msg:"操作ID不能为空"`
	RoleName    string `json:"role_name" binding:"required" msg:"角色描述不能为空"`
	RoleKey     string `json:"role_key" binding:"required" msg:"角色名称不能为空"`
	Menus       []int  `json:"menus"`
}

type UpdateRoleReq struct {
	OperationID string `json:"operation_id"  binding:"required" msg:"操作ID不能为空"`
	RoleName    string `json:"role_name" binding:"required,min=1,max=20" msg:"角色描述不能为空"`
	RoleKey     string `json:"role_key" binding:"required,min=1,max=16" msg:"角色名不能为空"`
	Menus       []int  `json:"menus" binding:"required" msg:"角色菜单不能为空"`
}

type UpdateRoleResp struct {
	ID       int    `json:"id"`
	RoleName string `json:"role_name"`
	RoleKey  string `json:"role_key"`
	Menus    []int  `json:"menus"`
}

type RoleListReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	RoleName    string `json:"role_name" form:"role_name"`
	RoleKey     string `json:"role_key" form:"role_key"`
	pagination.Pagination
}

type RoleListResp struct {
	List     []RoleListRespList `json:"list"`
	Count    int64              `json:"count"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}
type RoleListRespList struct {
	ID       int    `json:"id"`
	RoleName string `json:"role_name"`
	RoleKey  string `json:"role_key"`
}

type RoleAddReq struct {
	OperationID string `json:"operation_id"  binding:"required" msg:"操作ID不能为空"`
	RoleName    string `json:"role_name" binding:"required,min=1,max=20" msg:"角色描述不能为空"`
	RoleKey     string `json:"role_key" binding:"required,min=1,max=16" msg:"角色名不能为空"`
	Menus       []int  `json:"menus"`
}

type RoleAddResp struct {
	ID       int    `json:"id"`
	RoleName string `json:"role_name"`
	RoleKey  string `json:"role_key"`
}

type DeleteRoleReq struct {
	OperationID string `json:"operation_id"  binding:"required" msg:"操作ID不能为空"`
}

type GetRoleByIDReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
}

type GetRoleByIDResp struct {
	ID       int                        `json:"id"`
	RoleName string                     `json:"role_name"`
	RoleKey  string                     `json:"role_key"`
	CmsMenu  []GetRoleByIDMenusItemResp `json:"menus"`
}

type GetRoleByIDMenusItemResp struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Path       string `json:"path"`
	Type       int    `json:"type"`
	Permission string `json:"permission"`
	ParentId   int    `json:"parent_id"`
	NoCache    int    `json:"no_cache"`
	Component  string `json:"component"`
	Sort       int    `json:"sort"`
	Visible    int    `json:"visible"`
	Hidden     int    `json:"hidden"`
	IsFrame    int    `json:"is_frame"`
}
