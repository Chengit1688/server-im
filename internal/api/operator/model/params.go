package model

import "im/pkg/pagination"

type SearchOperatorReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required"`
	Longitude   string `json:"longitude" form:"longitude"`
	Latitude    string `json:"latitude" form:"latitude"`
	Key         string `json:"key" form:"key"`
	CityCode    int    `json:"city_code" form:"city_code"`
	pagination.Pagination
}

type OperatorTeamMemberInfoReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required"`
	ShopID      int64  `json:"shop_id" binding:"required,number"`
	UserID      string `json:"user_id"`
}
type OperatorTeamMemberInfoResp struct {
	ShopID   int64    `json:"shop_id" binding:"required,number"`
	UserID   string   `json:"user_id"`
	Role     int64    `json:"role"`
	UserInfo TeamInfo `json:"user_info"`
}

type OperatorTeamLeaderInfoReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	UserID      string `json:"user_id" binding:"omitempty,min=1"`
}

type OperatorTeamLeaderInfoResp struct {
	ShopID     int64    `json:"shop_id"`
	ShopName   string   `json:"shop_name"`
	UserID     string   `json:"user_id"`
	ShopStatus int64    `json:"shop_status"`
	Role       int64    `json:"role"`
	HasShop    int64    `json:"has_shop"`
	UserInfo   TeamInfo `json:"user_info"`
}

type OperatorJoinTeamReq struct {
	OperationID  string `json:"operation_id" form:"operation_id" binding:"required"`
	InviteCode   string `json:"invite_code"  binding:"omitempty,min=1"`
	ShopID       int64  `json:"shop_id" binding:"omitempty,number"`
	UserID       string `json:"user_id"`
	InviteUserId string `json:"invite_user_id" binding:"omitempty,min=1"`
}

type OperatorJoinTeamResp struct {
	ShopID   int64  `json:"shop_id"`
	UserID   string `json:"user_id"`
	TeamName string `json:"team_name"`
}

type OperatorRemoveTeamReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required"`
	ShopID      int64  `json:"shop_id" binding:"required,number"`
	UserID      string `json:"user_id"`
}

type OperatorIDCommonReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required"`
	ShopID      int64  `json:"shop_id" binding:"required,number"`
}

type OperatorTeamListReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required"`
	ShopID      int64  `json:"shop_id" binding:"required,number"`
	Key         string `json:"key"  binding:"omitempty,min=1"`
	pagination.Pagination
}

type SearchInfo struct {
	ShopID          int64    `json:"shop_id"`
	ShopName        string   `json:"shop_name"`
	ShopStar        float64  `json:"shop_star"`
	ShopLocation    string   `json:"shop_location"`
	ShopType        string   `json:"shop_type"`
	ShopDistance    float64  `json:"shop_distance"`
	DecorationScore float64  `json:"decoration_score"`
	Star            float64  `json:"star"`
	QualityScore    float64  `json:"quality_score"`
	ServiceScore    float64  `json:"service_score"`
	ShopIcon        []string `json:"shop_icon"`
	Longitude       string   `json:"longitude"`
	Latitude        string   `json:"latitude"`
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
type OperatorTeamResp struct {
	List []TeamInfo `json:"list"`
	pagination.Pagination
}

type SearchOperatorResp struct {
	List []SearchInfo `json:"list"`
	pagination.Pagination
}

type OperatorDetailResp struct {
	ID              int64    `json:"id"`
	Name            string   `json:"name"`
	Longitude       string   `json:"longitude"`
	Latitude        string   `json:"latitude"`
	Address         string   `json:"address"`
	Image           []string `json:"image"`
	Description     string   `json:"description"`
	DecorationScore float64  `json:"decoration_score"`
	QualityScore    float64  `json:"quality_score"`
	ServiceScore    float64  `json:"service_score"`
	CreatedAt       int64    `json:"created_at"`
	UpdatedAt       int64    `json:"updated_at"`
	ShopType        string   `json:"shop_type"`
	License         string   `json:"license"`
	CreatorId       string   `json:"creator_id"`
	InviteCode      string   `json:"invite_code"`
	Status          int64    `json:"status"`
}

type OperatorApplyForReq struct {
	OperationID string   `json:"operation_id" binding:"required,gte=1"`
	ShopID      int64    `json:"shop_id" binding:"omitempty,number"`
	Name        string   `json:"name" binding:"required"`
	Longitude   string   `json:"longitude" binding:"omitempty"`
	Latitude    string   `json:"latitude" binding:"omitempty"`
	Address     string   `json:"address" binding:"omitempty,min=1"`
	License     string   `json:"license" binding:"required"`
	Image       []string `json:"image" binding:"required"`
	Description string   `json:"description" binding:"omitempty,min=1"`
	CityCode    string   `json:"city_code" binding:"omitempty"`
	ShopType    string   `json:"shop_type" binding:"omitempty"`
	Star        float64  `json:"-"`
}

type SearchDTO struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	Longitude       string  `json:"longitude"`
	Latitude        string  `json:"latitude"`
	Address         string  `json:"address"`
	License         string  `json:"license"`
	Image           string  `json:"image"`
	Description     string  `json:"description"`
	DecorationScore float64 `json:"decoration_score"`
	QualityScore    float64 `json:"quality_score"`
	ServiceScore    float64 `json:"service_score"`
	Star            float64 `json:"star"`
	ShopType        string  `json:"shop_type"`
	CityCode        int     `json:"city_code"`
	Distance        float64 `json:"distance"`
	CreatedAt       int64   `json:"created_at"`
}

type UpdateOperatorReq struct {
	OperationID     string   `json:"operation_id" binding:"required,gte=1"`
	ShopID          int64    `json:"shop_id" binding:"omitempty,number"`
	Name            string   `json:"name" binding:"required"`
	Longitude       string   `json:"longitude" binding:"omitempty"`
	Latitude        string   `json:"latitude" binding:"omitempty"`
	Address         string   `json:"address" binding:"omitempty,min=1"`
	License         string   `json:"license" binding:"required"`
	Image           []string `json:"image" binding:"required"`
	Description     string   `json:"description" binding:"required"`
	CityCode        string   `json:"city_code" binding:"omitempty"`
	ShopType        string   `json:"shop_type" binding:"omitempty"`
	DecorationScore float64  `json:"decoration_score"`
	QualityScore    float64  `json:"quality_score"`
	ServiceScore    float64  `json:"service_score"`
	Star            float64  `json:"star"`
	Status          int64    `json:"status"`
}

type OperatorAgentLevelListReq struct {
	OperationID  string `json:"operation_id"  binding:"required,min=1" msg:"日志id必须传"`
	ShopID       int64  `json:"shop_id" binding:"omitempty,number"`
	UserId       string `json:"user_id" binding:"omitempty,number,min=1" msg:"请输入user_id"`
	InviteUserId string `json:"invite_user_id" binding:"omitempty,number,min=1" binding:"required" msg:"请输入邀请人ID"`
	BeginDate    int64  `json:"begin_date" binding:"omitempty,number" msg:"请输入开始时间"`
	EndDate      int64  `json:"end_date"  binding:"omitempty,number" msg:"请输入结束时间"`
	pagination.Pagination
}

type OperatorAgentLevelListResp struct {
	List []OperatorAgentLevelListInfo `json:"list"`
	pagination.Pagination
}

type OperatorAgentLevelListInfo struct {
	ShopID    int64  `json:"shop_id"`
	ShopName  string `json:"shop_name"`
	CreatedAt int64  `json:"created_at"`
	TeamInfo
}
