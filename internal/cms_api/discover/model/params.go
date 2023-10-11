package model

import "im/pkg/pagination"

type GetDiscoverInfoReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
}

type GetDiscoverInfoResp struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Url       string `json:"url"`
	Icon      string `json:"icon"`
	Sort      int    `json:"sort"`
	CreatedAt int64  `json:"created_at"`
}

type AddDiscoverReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Name        string `json:"name" binding:"required,min=1,max=20" msg:"名称长度1～20"`
	Url         string `json:"url" binding:"required,min=1,max=200" msg:"链接长度1～20"`
	Icon        string `json:"icon"`
	Sort        int    `json:"sort"`
}

type GetDiscoverOpenReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
}

type SetDiscoverOpenReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Status      int    `json:"status" binding:"required,min=1,max=2" msg:"1开2关"`
}

type AddNewsReq struct {
	OperationID string   `json:"operation_id" binding:"required" msg:"操作ID不能为空"`
	ID          int64    `json:"id" binding:"omitempty,number"`
	Title       string   `json:"title" binding:"required,min=1,max=100"`
	Content     string   `json:"content" binding:"required,min=1"`
	CategoryID  int64    `json:"category_id"`
	Status      int64    `json:"status"`
	Image       []string `json:"image"`
	Video       string   `json:"video"`
}

type DeleteNewsReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	ID          int64  `json:"id" binding:"required,min=1"`
}

type ListNewsReq struct {
	OperationID  string `json:"operation_id" binding:"required"`
	CreateUserID string `json:"create_user_id"`
	Title        string `json:"title"`
	OrderBy      string `json:"order_by"`
	pagination.Pagination
}

type NewsInfo struct {
	ID             int64    `json:"id"`
	CreateUserID   string   `json:"create_user_id"`
	CreateNickname string   `json:"create_nickname"`
	CreateAccount  string   `json:"create_account"`
	Title          string   `json:"title"`
	Content        string   `json:"content"`
	ViewTotal      int64    `json:"view_total"`
	CategoryID     int64    `json:"category_id"`
	Image          []string `json:"image"`
	Video          string   `json:"video"`
	CreatedAt      int64    `json:"created_at"`
	UpdatedAt      int64    `json:"updated_at"`
}

type ListNewsResp struct {
	List []NewsInfo `json:"list"`
	pagination.Pagination
}

type AddPrizeReq struct {
	OperationID string      `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	ID          int64       `json:"id"`
	List        []PrizeList `json:"list" form:"list" binding:"required" msg:"兑换奖品"`
}

type UpdatePrizeReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	ID          int64  `json:"id" binding:"required"`
	PrizeList
}

type DeletePrizeReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	ID          int64  `json:"id" binding:"required"`
}

type RedeemPrizeLogReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	pagination.Pagination
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
	RedeemPrizeLog
}

type PrizeListReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	pagination.Pagination
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
	PrizeList
}

type PrizeListResp struct {
	pagination.Pagination
	List []PrizeList `json:"list"`
}
type RedeemPrizeLogResp struct {
	pagination.Pagination
	List []RedeemPrizeLog `json:"list"`
}

type SetRedeemPrizeReq struct {
	OperationID   string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	UserName      string `json:"user_name"`
	Address       string `json:"address"`
	Mobile        string `json:"mobile"`
	ExpressNumber string `json:"express_number"`
	Status        int64  `json:"status"`
	ID            int64  `json:"id" binding:"required"`
}
