package model

import (
	"im/pkg/pagination"
)

type OperateLogPagingReq struct {
	OperationID      string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Search           string `json:"search" form:"search" binding:"omitempty,gte=1"`
	CreatedTimeStart int64  `json:"created_time_start" form:"created_time_start"`
	CreatedTimeEnd   int64  `json:"created_time_end" form:"created_time_end"`
	pagination.Pagination
}

type OperateLogPagingItemResp struct {
	ID               int64  `json:"id"`
	CreatedAt        int64  `json:"created_at"`
	UserID           string `json:"user_id"`
	Username         string `json:"username"`
	NickName         string `json:"nick_name"`
	ServiceID        string `json:"service_id"`
	Env              string `json:"env"`
	LogLevel         string `json:"log_level"`
	RequestUrl       string `json:"request_url"`
	RequestParams    string `json:"request_params"`
	RequestMethod    string `json:"request_method"`
	RequestUserAgent string `json:"request_user_agent"`
	ServiceIp        string `json:"service_ip"`
	ServiceHost      string `json:"service_host"`
	OperationID      string `json:"operation_id"`
	LogRemark        string `json:"log_remark"`
}

type OperateLogPagingResp struct {
	List     []OperateLogPagingItemResp `json:"list"`
	Count    int64                      `json:"count"`
	Page     int                        `json:"page"`
	PageSize int                        `json:"page_size"`
}
