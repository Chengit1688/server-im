package model

import "im/pkg/pagination"

type UploadInfo struct {
	OperationID string `json:"operation_id"  binding:"required" msg:"操作ID不能为空"`
	AppName     string `json:"app_name" binding:"required"`
	UserId      string `json:"user_id" binding:"required"`
	MacType     string `json:"mac_type" binding:"required"`
	PhoneType   string `json:"phone_type" binding:"required"`
	Info        string `json:"info" binding:"required"`
	Extra       string `json:"extra"`
}

type SearchInfo struct {
	OperationID string `json:"operation_id"  binding:"required" msg:"操作ID不能为空"`
	UserId      string `json:"user_id"`
	Stime       int64  `json:"stime" binding:"required"`
	Etime       int64  `json:"etime" binding:"required"`
	MacType     string `json:"mac_type"`
	pagination.Pagination
}
type ErrlogInfo struct {
	UploadInfo
	CreateTime int64
}

type SearchInfoResp struct {
	List []ErrlogInfo `json:"list"`
	pagination.Pagination
}
