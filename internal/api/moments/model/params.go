package model

import "im/pkg/pagination"

type IssueReq struct {
	OperationID    string   `json:"operation_id" binding:"required"`
	Content        string   `json:"content"`
	ShareTagID     []int64  `json:"share_tag_id"`
	ShareFriendID  []string `json:"share_friend_id"`
	InviteFriendID []string `json:"invite_friend_id"`
	Location       string   `json:"location"`
	NoComment      int      `json:"no_comment"`
	CanSee         int      `json:"can_see"`
	Images         []string `json:"images"`
	Videos         []string `json:"videos"`
	VideoImg       []string `json:"video_img"`
}

type DelIssueReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	MomentsID   int64  `json:"moments_id" binding:"required"`
}

type CommentReq struct {
	OperationID string   `json:"operation_id" binding:"required"`
	Content     string   `json:"content"`
	MomentsID   int64    `json:"moments_id" binding:"required"`
	Images      []string `json:"images"`
	ReplyToId   string   `json:"reply_to_id"`
}

type DelCommentReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	CommentsID  int64  `json:"coments_id" binding:"required"`
}

type MomentsReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	MomentsID   int64  `json:"moments_id" binding:"required,number"`
}

type MomentsCommentsLikeReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	MomentsID   int64  `json:"moments_id" binding:"required"`
}

type IssueList struct {
	pagination.Pagination
	List []IssueInfo `json:"list"`
}

type IssueListReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	IsOwner     int64  `json:"is_owner"`
	UserID      string `json:"user_id"`
	pagination.Pagination
}

type IssueInfo struct {
	MomentsID     int64 `json:"moments_id"`
	User          UserInfo
	MomentsInfo   *MomentsMessage
	CommentsCount int64 `json:"comments_count"`
	Comments      []CommentsInfo
	LikesCount    int64 `json:"likes_count"`
	Likes         []MomentsCommentsLikeInfo
}

type MomentsDetailResp struct {
	Data IssueInfo `json:"data"`
	pagination.Pagination
}

type MomentsCommentsReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	MomentsID   int64  `json:"moments_id" binding:"required"`
	pagination.Pagination
}

type MomentsCommentsList struct {
	pagination.Pagination
	List []CommentsInfo `json:"list"`
}

type CommentsInfo struct {
	ID           int64 `json:"id"`
	MomentsID    int64 `json:"moments_id"`
	User         UserInfo
	Content      string   `json:"content"`
	IsOwnComment int64    `json:"is_own_comment"`
	ReplyToUser  UserInfo `json:"reply_user"`
	Images       []string `json:"images"`
}

type MomentsCommentsLikeList struct {
	pagination.Pagination
	List []MomentsCommentsLikeInfo `json:"list"`
}

type MomentsCommentsLikeInfo struct {
	ID        int64 `json:"id"`
	MomentsID int64 `json:"moments_id"`
	User      UserInfo
}

type UserInfo struct {
	UserId   string `json:"user_id"`
	FaceURL  string `json:"face_url"`
	NickName string `json:"nick_name"`
}

type TagReq struct {
	OperationID string   `json:"operation_id" binding:"required"`
	Title       string   `json:"title" binding:"required"`
	UserID      []string `json:"user_id" binding:"required"`
	TagID       int64    `json:"tag_id"`
}

type TagAddFriendReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
	TagID       int64  `json:"tag_id"  binding:"required"`
}

type FetchFriendTag struct {
	OperationID string `json:"operation_id" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
}

type TagAddFriendResp struct {
	List []TagListInfo
}

type CommonIDReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	ID          int64  `json:"id" binding:"required"`
}

type TagListReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	pagination.Pagination
}

type TagListResp struct {
	List []TagListInfo `json:"list"`
	pagination.Pagination
}

type TagListInfo struct {
	TagID     int64  `json:"tag_id"`
	Title     string `json:"title"`
	UserTotal int    `json:"user_total"`
}

type TagDetailResp struct {
	Total int           `json:"total" `
	List  []TagUserInfo `json:"list" `
}

type TagUserInfo struct {
	UserID     string `json:"user_id"`
	Nickname   string `json:"nickname"`
	FaceUrl    string `json:"face-url"`
	BigFaceURL string `json:"big_face_url"`
}

type CommonResp struct {
	Code int      `json:"code"`
	Data struct{} `json:"data"`
	Msg  string   `json:"msg"`
}

type MomentsInviteFriendPush struct {
	ConversationID    string `json:"conversation_id"`
	Timestamp         int64  `json:"timestamp"`
	PublisherUserID   string `json:"publisher_user_id"`
	FriendUserID      string `json:"friend_user_id"`
	MomentsID         int64  `json:"moments_id"`
	Type              int64  `json:"type"`
	FaceURL           string `json:"face_url"`
	PublisherNickname string `json:"publisher_nickname"`
	LikeStatus        int64  `json:"like_status"`
}

type GroupAtPush struct {
	ConversationID    string `json:"conversation_id"`
	Timestamp         int64  `json:"timestamp"`
	PublisherUserID   string `json:"publisher_user_id"`
	PublisherNickname string `json:"publisher_nickname"`
	FaceURL           string `json:"face_url"`
	FriendUserID      string `json:"friend_user_id"`
	Type              int64  `json:"type"`
}
