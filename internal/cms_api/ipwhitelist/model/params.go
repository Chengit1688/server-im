package model

import "im/pkg/pagination"

type GetIPListReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	IP          string `json:"ip" form:"ip"`
	Note        string `json:"note" form:"note"`
	pagination.Pagination
}

type GetIPListResp struct {
	List []GetIPListItem `json:"list"`
	pagination.Pagination
}

type GetIPListItem struct {
	ID        int    `json:"id"`
	IP        string `json:"ip"`
	CreatedAt int64  `json:"created_at"`
	Note      string `json:"note"`
}

type AddIPInfoReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	IP          string `json:"ip" binding:"required,ip" msg:"IP不能为空"`
	Note        string `json:"note" binding:"required" msg:"备注不能为空"`
	UserID      string `json:"user_id"`
}

type AddIPInfoResp struct {
	ID   int    `json:"id"`
	IP   string `json:"ip"`
	Note string `json:"note"`
}

type UpdateIPInfoReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	IP          string `json:"ip" binding:"required,ip" msg:"IP不能为空"`
	Note        string `json:"note" binding:"required" msg:"备注不能为空"`
	UserID      string `json:"user_id"`
}

type UpdateIPInfoResp struct {
	ID   int    `json:"id"`
	IP   string `json:"ip"`
	Note string `json:"note"`
}

type DeleteIPInfoReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
}
