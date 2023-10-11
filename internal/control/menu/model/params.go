package model

import "time"

type MenuListReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Name        string `json:"name" form:"name"`
	Title       string `json:"title" form:"title"`
}
type MenuListResp struct {
	List  []MenuListRespItem `json:"list"`
	Count int64              `json:"count"`
}

type MenuListRespItem struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Icon       string `json:"icon"`
	Path       string `json:"path"`
	Paths      string `json:"paths"`
	Type       int    `json:"type"`
	Action     string `json:"action"`
	Permission string `json:"permission"`
	ParentId   int    `json:"parent_id"`
	NoCache    int    `json:"no_cache"`
	Component  string `json:"component"`
	Sort       int    `json:"sort"`
	Visible    int    `json:"visible"`
	Hidden     int    `json:"hidden"`
	IsFrame    int    `json:"is_frame"`
}

type AddMenuReq struct {
	OperationID string `json:"operation_id"  binding:"required" msg:"操作ID不能为空"`
	Name        string `json:"name" binding:"required" msg:"菜单名不能为空"`
	Title       string `json:"title" binding:"required" msg:"菜单标题不能为空"`
	Icon        string `json:"icon"`
	Path        string `json:"path" binding:"required" msg:"菜单路径不能为空"`
	Paths       string `json:"paths"`
	Type        int    `json:"type" binding:"required" msg:"菜单类型不能为空"`
	Action      string `json:"action"`
	Permission  string `json:"permission"`
	ParentId    *int   `json:"parent_id" binding:"required" msg:"菜单父级不能为空"`
	NoCache     int    `json:"no_cache"`
	Component   string `json:"component"`
	Sort        int    `json:"sort"`
	Visible     int    `json:"visible"`
	Hidden      int    `json:"hidden"`
	IsFrame     int    `json:"is_frame"`
	Apis        []int  `json:"apis"`
}

type AddMenuResp struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Icon       string `json:"icon"`
	Path       string `json:"path"`
	Paths      string `json:"paths"`
	Type       int    `json:"type"`
	Action     string `json:"action"`
	Permission string `json:"permission"`
	ParentId   int    `json:"parent_id"`
	NoCache    int    `json:"no_cache"`
	Component  string `json:"component"`
	Sort       int    `json:"sort"`
	Visible    int    `json:"visible"`
	Hidden     int    `json:"hidden"`
	IsFrame    int    `json:"is_frame"`
	Apis       []int  `json:"apis"`
}

type GetMenuResp struct {
	Name       string `json:"name"`
	Title      string `json:"title"`
	Icon       string `json:"icon"`
	Path       string `json:"path"`
	Paths      string `json:"paths"`
	Type       int    `json:"type"`
	Action     string `json:"action"`
	Permission string `json:"permission"`
	ParentId   int    `json:"parent_id"`
	NoCache    int    `json:"no_cache"`
	Component  string `json:"component"`
	Sort       int    `json:"sort"`
	Visible    int    `json:"visible"`
	Hidden     int    `json:"hidden"`
	IsFrame    int    `json:"is_frame"`
	Apis       []int  `json:"apis"`
}

type UpdateMenuReq struct {
	OperationID string `json:"operation_id"  binding:"required" msg:"操作ID不能为空"`
	Name        string `json:"name" binding:"required" msg:"菜单名不能为空"`
	Title       string `json:"title" binding:"required" msg:"菜单标题不能为空"`
	Icon        string `json:"icon"`
	Path        string `json:"path" binding:"required" msg:"菜单路径不能为空"`
	Paths       string `json:"paths"`
	Type        int    `json:"type" binding:"required" msg:"菜单类型不能为空"`
	Action      string `json:"action"`
	Permission  string `json:"permission"`
	ParentId    *int   `json:"parent_id" binding:"required" msg:"菜单父级不能为空"`
	NoCache     int    `json:"no_cache"`
	Component   string `json:"component"`
	Sort        int    `json:"sort"`
	Visible     int    `json:"visible"`
	Hidden      int    `json:"hidden"`
	IsFrame     int    `json:"is_frame"`
	Apis        []int  `json:"apis" binding:"required" msg:"菜单API关联信息不能为空"`
}

type UpdateMenuResp struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Icon       string `json:"icon"`
	Path       string `json:"path"`
	Paths      string `json:"paths"`
	Type       int    `json:"type"`
	Action     string `json:"action"`
	Permission string `json:"permission"`
	ParentId   int    `json:"parent_id"`
	NoCache    int    `json:"no_cache"`
	Component  string `json:"component"`
	Sort       int    `json:"sort"`
	Visible    int    `json:"visible"`
	Hidden     int    `json:"hidden"`
	IsFrame    int    `json:"is_frame"`
	Apis       []int  `json:"apis"`
}

type GetMenuConfigReq struct {
	Timestamp int64 `json:"timestamp" form:"timestamp" binding:"required" msg:"时间戳不能为空"`
}

type MenuListRespSyncItem struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Title      string    `json:"title"`
	Icon       string    `json:"icon"`
	Path       string    `json:"path"`
	Paths      string    `json:"paths"`
	Type       int       `json:"type"`
	Action     string    `json:"action"`
	Permission string    `json:"permission"`
	ParentId   int       `json:"parent_id"`
	NoCache    int       `json:"no_cache"`
	Component  string    `json:"component"`
	Sort       int       `json:"sort"`
	Visible    int       `json:"visible"`
	Hidden     int       `json:"hidden"`
	IsFrame    int       `json:"is_frame"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}

type GetMenuConfigResp struct {
	Menus     []MenuListRespSyncItem `json:"menus"`
	Timestamp int64                  `json:"timestamp" form:"timestamp"`
}

type DeleteMenuReq struct {
	OperationID string `json:"operation_id"  binding:"required" msg:"操作ID不能为空"`
}
