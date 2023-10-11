package model

import (
	"im/internal/api/shopping/model"
	"im/pkg/pagination"
)

type OperatorListReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required"`
	ShopName    string `json:"shop_name" form:"shop_name"`
	ShopID      int64  `json:"shop_id" form:"shop_id"`
	Status      int64  `json:"status" form:"status"`
	pagination.Pagination
}

type OperatorMemberListReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required"`
	ShopID      int64  `json:"shop_id" binding:"required,number"`
	Key         string `json:"key"  binding:"omitempty,min=1"`
	pagination.Pagination
}

type OperatorMemberListResp struct {
	List []TeamInfo `json:"list"`
	pagination.Pagination
}

type TeamInfo struct {
	UserID      string `json:"user_id"`
	Account     string `json:"account"`
	PhoneNumber string `json:"phone_number"`
	CountryCode string `json:"country_code"`
	FaceURL     string `json:"face_url"`
	BigFaceURL  string `json:"big_face_url"`
	Gender      int64  `json:"gender"`
	NickName    string `json:"nick_name"`
	Age         int64  `json:"age"`
}

type OperatorApproveReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required"`
	ShopID      int64  `json:"shop_id" form:"shop_id" binding:"required,number"`
	Status      int64  `json:"status" form:"status" binding:"required,oneof=1 2 3"  msg:"审批状态错误"`
}

type UserInfo struct {
	UserID      string `json:"user_id"`
	Account     string `json:"account"`
	PhoneNumber string `json:"phone_number"`
	CountryCode string `json:"country_code"`
	FaceURL     string `json:"face_url"`
	BigFaceURL  string `json:"big_face_url"`
	Gender      int64  `json:"gender"`
	NickName    string `json:"nick_name"`
}

type ShopListInfo struct {
	UserInfo
	ShopInfo
}

type OperatorListResp struct {
	List []ShopListInfo `json:"list"`
	pagination.Pagination
}

type ShopInfo struct {
	ShopID          int64    `json:"shop_id"`
	Status          int64    `json:"status"`
	ShopName        string   `json:"shop_name"`
	ShopStar        float64  `json:"shop_star"`
	ShopLocation    string   `json:"shop_location"`
	ShopType        string   `json:"shop_type"`
	Description     string   `json:"description"`
	ShopDistance    float64  `json:"shop_distance"`
	DecorationScore float64  `json:"decoration_score"`
	Star            float64  `json:"star"`
	QualityScore    float64  `json:"quality_score"`
	ServiceScore    float64  `json:"service_score"`
	ShopIcon        []string `json:"shop_icon"`
	Longitude       string   `json:"longitude"`
	Latitude        string   `json:"latitude"`
	Address         string   `json:"address"`
	License         string   `json:"license"`
	CreatedAt       int64    `gorm:"column:created_at"`
}

type SearchResp struct {
	List []model.Shop `json:"list"`
	pagination.Pagination
}
