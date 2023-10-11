package model

import (
	"im/pkg/pagination"
)

type PageReq struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}
type FriendInfo struct {
	ID               int64  `json:"id"`
	UserId           string `json:"user_id"`
	FaceURL          string `json:"face_url"`
	BigFaceURL       string `json:"big_face_url"`
	NickName         string `json:"nick_name"`
	Signatures       string `json:"signatures"`
	Gender           int64  `json:"gender"`
	Account          string `json:"account"`
	PhoneNumber      string `json:"phone_number"`
	Age              int64  `json:"age"`
	Remark           string `json:"remark"`
	LoginIp          string `json:"login_ip"`
	LoginIpLocaltion string `json:"login_ip_localtion"`
	OnlineStatus     int    `json:"online_status"`
	CreateTime       int64  `json:"create_time"`
	BlackStatus      int    `json:"black_status"`
	Status           int64  `json:"status"`
	Version          int64  `json:"version"`
	FriendLabel      string `json:"friend_label"`
	Online           bool   `json:"online"`
	OfflineInfo      string `json:"offline_info"`
}

type FriendRequestInfo struct {
	ID         int64  `json:"id"`
	UserId     string `json:"user_id"`
	FaceURL    string `json:"face_url"`
	BigFaceURL string `json:"big_face_url"`
	NickName   string `json:"nick_name"`
	Account    string `json:"account"`
	ReqMsg     string `json:"req_msg"`
	CreateTime int64  `json:"create_time"`
	Status     int64  `json:"status"`
	Signatures string `json:"signatures"`
	Gender     int64  `json:"gender"`
	Age        int64  `json:"age"`
	FromUserID string `json:"owner_user_id"`
	ToUserID   string `json:"to_user_id"`
	IsOwner    int64  `json:"is_owner"`
}

type ParamsCommFriend struct {
	OperationID string `json:"operation_id" binding:"required"`
	UserId      string `json:"user_id" binding:"required"`
}

type AddFriendReq struct {
	ParamsCommFriend
	ReqMsg string `json:"req_msg"`
	Remark string `json:"remark"`
}

type AddFriendResp FriendRequestInfo

type AddFriendAckReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	ReqId       int64  `json:"req_id" binding:"required"`
	Status      int32  `json:"status" binding:"required,oneof=1 2"`
}

type AddFriendAckResp FriendRequestInfo

type DeleteFriendReq struct {
	ParamsCommFriend
}

type DeleteFriendResp FriendInfo

type AddBlackReq ParamsCommFriend

type AddBlackResp FriendInfo

type GetBlackListReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	pagination.Pagination
}

type GetBlackListResp struct {
	pagination.Pagination
	List []FriendInfo `json:"list"`
}

type SetFriendRemarkReq struct {
	ParamsCommFriend
	Remark string `json:"remark"`
}

type GetFriendRemarkReq struct {
	OperationID string `json:"operation_id" binding:"required"`
}

type CheckFriendRemarkReq struct {
	ParamsCommFriend
	Remark string `json:"remark"  binding:"required"`
}
type SetFriendRemarkResp FriendInfo

type GetFriendRemarkResp struct {
	Remark []string
}

type RemoveBlackReq struct {
	ParamsCommFriend
}

type RemoveBlackResp FriendInfo

type GetFriendInfoReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	UserId      string `json:"user_id" binding:"required"`
}

type GetFriendInfoResp FriendInfo

type GetFriendListResp struct {
	pagination.Pagination
	List []FriendInfo `json:"list"`
}
type GetFriendListMaxSeqResp struct {
	pagination.Pagination
	List []MaxFriendsSeq `json:"list"`
}
type GetFriendListMaxSeqReq struct {
	pagination.Pagination
	OperationID string `json:"operation_id" binding:"required"`
}

type GetFriendListReq struct {
	pagination.Pagination
	OperationID string `json:"operation_id" binding:"required"`
	SearchKey   string `json:"search_key"`
	FriendLabel string `json:"friend_label"`
	Version     int32  `json:"version"`
	BlackStatus int    `json:"black_status"`
}
type GetFriendMaxReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	FromUserID  string `json:"from_user_id"`
}

type IsFriendReq struct {
	OperationID  string `json:"operation_id" binding:"required"`
	FriendUserID string `json:"friend_user_id" binding:"required"`
}

type MaxSeq struct {
	Version       int32 `json:"version"`
	MaxUpdateTime int64 `json:"max_update_time`
}

type MaxFriendsSeq struct {
	ConversationID string `json:"conversation_id"`
	MaxFriendSeq   int64  `json:"max_friend_seq`
}

type Message struct {
	FromUserId string `json:"from_user_id"`
	Msg        string `json:"msg"`
	Status     int32  `json:"status"`
}

type GetFriendApplyListReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	PageReq
}

type GetFriendApplyListRes struct {
	pagination.Pagination
	List []FriendRequestInfo `json:"list"`
}

type FriendSearchReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	pagination.Pagination
	Keyword string `json:"keyword" binding:"required"`
}

type FriendSearchResp struct {
	pagination.Pagination
	List []FriendInfo `json:"list"`
}

type SearchFriendReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	UserId      string `json:"user_id"  binding:"required"`
}

type SearchFriendRes FriendInfo

type FriendListSyncResp struct {
	pagination.Pagination
	List []FriendInfo `json:"list"`
}

type FriendListSyncReq struct {
	pagination.Pagination
	OperationID string `json:"operation_id" binding:"required"`
	Version     int    `json:"version"`
}

type FriendLabelInfo struct {
	LabelId   string `json:"label_id"`
	LabelName string `json:"label_name"`
}

type CreateFriendLabelReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	LabelName   string `json:"label_name"  binding:"required"`
}

type CreateFriendLabelResp FriendLabelInfo

type DeleteFriendLabelReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	LabelId     string `json:"label_id"  binding:"required"`
}

type DeleteFriendLabelResp struct {
}

type UpdateFriendLabelReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	LabelId     string `json:"label_id" binding:"required"`
	LabelName   string `json:"label_name" binding:"required"`
}

type UpdateFriendLabelResp FriendLabelInfo

type ChangeFriendLabelReq struct {
	OperationID string   `json:"operation_id" binding:"required"`
	FriendList  []string `json:"friend_list" binding:"required"`
	LabelId     string   `json:"label_id" binding:"required"`
}

type ChangeFriendLabelResp struct {
}

type GetFriendLabelReq struct {
	OperationID string `json:"operation_id" binding:"required"`
}
type GetFriendLabelResp struct {
	List []FriendLabelInfo `json:"list"`
}

type BlackFriendInfo struct {
	FromUserId string `json:"from_user_id"`
	ToUserId   string `json:"to_user_id"`
	Msg        string `json:"msg"`
}
