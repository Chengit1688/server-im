package model

import (
	"im/internal/api/user/model"
	"im/pkg/db"
)

const CanSeeFriend = 1
const CanSeePublic = 2
const CanSeePrivate = 3

type MomentsMessage struct {
	ID                int64  `gorm:"column:id;primary_key" json:"id"`
	CreatedAt         int64  `gorm:"column:created_at"`
	UpdatedAt         int64  `gorm:"column:updated_at"`
	Description       string `gorm:"column:description;size:2000" json:"description"`
	Content           string `gorm:"column:content;type:varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci" json:"content"`
	Location          string `gorm:"column:location;size:200" json:"location"`
	Image             string `gorm:"column:image;size:2000;" json:"image"`
	Video             string `gorm:"column:videos;size:1000;" json:"video"`
	VideoImg          string `gorm:"column:video_img;size:1000;" json:"video_img"`
	NoComment         int    `gorm:"column:no_comment;default:1;" json:"no_comment"`
	CanSee            int    `gorm:"column:can_see;default:1;" json:"can_see"`
	InviteFriendID    string `gorm:"column:invite_friend_id;size:1000;default:'';" json:"invite_friend_id"`
	ShareTagID        string `gorm:"column:share_tag_id;size:1000;default:'';" json:"share_tag_id"`
	ShareFriendID     string `gorm:"column:share_friend_id;size:1000;default:'';" json:"share_friend_id"`
	DontShareTagID    string `gorm:"column:dont_share_tag_id;size:1000;default:'';" json:"dont_share_tag_id"`
	DontShareFriendID string `gorm:"column:dont_share_friend_id;size:1000;default:'';" json:"dont_share_friend_id"`
	UserId            string `gorm:"column:user_id;index;size:64" json:"user_id"`
	Year              int    `gorm:"column:year;default:0;" json:"year"`
	Month             int    `gorm:"column:month;default:0;" json:"month"`
	Day               int    `gorm:"column:day;default:0;" json:"day"`
	Status            int64  `gorm:"column:status;default:1;" json:"status"`
}

func (d *MomentsMessage) TableName() string {
	return "moments_messages"
}

type MomentsComments struct {
	ID           int64      `gorm:"column:id;primary_key" json:"id"`
	MomentsID    int64      `gorm:"column:moments_id;index"`
	CreatedAt    int64      `gorm:"column:created_at"`
	UpdatedAt    int64      `gorm:"column:updated_at"`
	Description  string     `gorm:"column:description;size:2000" json:"description"`
	Content      string     `gorm:"column:content;type:varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci" json:"content"`
	Location     string     `gorm:"column:location;size:200" json:"location"`
	Image        string     `gorm:"column:image;size:2000" json:"image"`
	UserID       string     `gorm:"column:user_id;index;size:64;not null;default:''" json:"user_id"`
	IsOwnComment int64      `gorm:"column:is_own_comment;default:1;" json:"is_own_comment"`
	Status       int64      `gorm:"column:status;default:1;" json:"status"`
	ReplyToId    string     `gorm:"column:reply_to_id;index;size:64;" json:"reply_to_id"`
	ReplyUser    model.User `gorm:"foreignKey:UserID;references:ReplyToId"`
}

func (d *MomentsComments) TableName() string {
	return "moments_comments"
}

type MomentsCommentsLike struct {
	ID           int64  `gorm:"column:id;primary_key" json:"id"`
	MomentsID    int64  `gorm:"column:moments_id"`
	CreatedAt    int64  `gorm:"column:created_at"`
	UpdatedAt    int64  `gorm:"column:updated_at"`
	FriendUserID string `gorm:"column:friend_user_id;index;size:64" json:"friend_user_id"`
	Status       int64  `gorm:"column:status;default:2;" json:"status"`
}

func (d *MomentsCommentsLike) TableName() string {
	return "moments_comments_likes"
}

type ContactsTag struct {
	ID           int64  `gorm:"column:id;primary_key" json:"id"`
	Title        string `gorm:"column:title;default:'';" json:"title"`
	CreatorId    string `gorm:"column:creator_id;default:0;" json:"creator_id"`
	FriendUserID string `gorm:"column:user_id;default:0;size:1000;" json:"friend_user_id"`
	FriendLength int    `gorm:"column:friend_length;default:0;" json:"friend_length"`
	CreatedAt    int64  `gorm:"column:created_at"`
	UpdatedAt    int64  `gorm:"column:updated_at"`
}

func (d *ContactsTag) TableName() string {
	return "contacts_tags"
}

type MomentsInbox struct {
	db.CommonModel
	CreatedAt       int64  `gorm:"column:created_at"`
	UpdatedAt       int64  `gorm:"column:updated_at"`
	MomentsID       int64  `gorm:"column:moments_id;default:0;index;" json:"moments_id"`
	FriendUserID    string `gorm:"column:friend_user_id;default:0;size:80;index;" json:"friend_user_id"`
	PublisherUserId string `gorm:"column:publisher_user_id;size:80;default:0;index;" json:"publisher_user_id"`
}

func (d *MomentsInbox) TableName() string {
	return "moments_inbox"
}
