package model

import "im/pkg/pagination"

type NewsReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	Title       string `json:"title"`
	OrderBy     string `json:"order_by"`
	pagination.Pagination
}

type NewsDetailInfo struct {
	ID           int64    `json:"id"`
	CreateUserID string   `json:"create_user_id"`
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	Image        []string `json:"image"`
	Video        string   `json:"video"`
	ViewTotal    int64    `json:"view_total"`
	CategoryID   int64    `json:"category_id"`
	CreatedAt    int64    `json:"created_at"`
	UpdatedAt    int64    `json:"updated_at"`
}

type NewsResp struct {
	List []NewsDetailInfo `json:"list"`
	pagination.Pagination
}

type NewsDetailReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	ID          int64  `json:"id"`
}
