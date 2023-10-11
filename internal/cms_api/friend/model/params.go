package model

import (
	apiFriendModel "im/internal/api/friend/model"
)

type Op struct {
	OperationID string `json:"operation_id" binding:"required"`
}

type ListReq struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type UserFriendListReq struct {
	Op
	ListReq
	UserId   string `json:"user_id" binding:"required"`
	FriendId string `json:"friend_id"`
}

type UserFriendListResp struct {
	Page     int                         `json:"page"`
	PageSize int                         `json:"page_size"`
	Count    int                         `json:"count"`
	List     []apiFriendModel.FriendInfo `json:"list"`
}

type UserAddFriendReq struct {
	Op
	UserId   string `json:"user_id" binding:"required"`
	FriendId string `json:"friend_id" binding:"required"`
}

type UserAddFriendResp struct{}

type UserRemoveFriendReq struct {
	Op
	UserId   string `json:"user_id" binding:"required"`
	FriendId string `json:"friend_id" binding:"required"`
}

type UserRemoveFriendResp struct{}
