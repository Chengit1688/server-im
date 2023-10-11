package repo

import (
	"gorm.io/gorm"
	"im/internal/api/friend/model"
	userModel "im/internal/api/user/model"
	userUse "im/internal/api/user/usecase"
	"im/pkg/db"
	"time"
)

var FriendRequestRepo = new(friendRequestRepo)

type friendRequestRepo struct{}

func (r *friendRequestRepo) GetFriendMemberUserInfo(member *model.FriendRequest, userId string) (model.FriendRequestInfo, error) {
	var (
		err      error
		userInfo *userModel.UserBaseInfo
	)
	if userInfo, err = userUse.UserUseCase.GetBaseInfo(userId); err != nil {
		return model.FriendRequestInfo{}, err
	}
	return model.FriendRequestInfo{
		ID:         member.ID,
		NickName:   userInfo.NickName,
		FaceURL:    userInfo.FaceURL,
		BigFaceURL: userInfo.BigFaceURL,
		Signatures: userInfo.Signatures,
		Gender:     userInfo.Gender,
		UserId:     userInfo.UserId,
		Account:    userInfo.Account,
		Age:        userInfo.Age,
		ReqMsg:     member.ReqMsg,
		CreateTime: member.CreatedAt,
		Status:     int64(member.Status),
		FromUserID: member.FromUserID,
		ToUserID:   member.ToUserID,
	}, nil
}

func (r *friendRequestRepo) GetFriendMemberUserInfoList(members *[]model.FriendRequest, userIdModel string) ([]model.FriendRequestInfo, error) {
	var (
		err    error
		fs     []model.FriendRequestInfo
		f      model.FriendRequestInfo
		userId string
	)

	for _, member := range *members {
		if userIdModel == "from_user_id" {
			userId = member.FromUserID
		} else {
			userId = member.ToUserID
		}
		if f, err = r.GetFriendMemberUserInfo(&member, userId); err != nil {
			return nil, err
		}
		fs = append(fs, f)
	}
	return fs, nil
}

func (r *friendRequestRepo) GetFriendMemberUserInfoAll(members *[]model.FriendRequest, fromUserID string) ([]model.FriendRequestInfo, error) {
	var (
		err error
		fs  []model.FriendRequestInfo
		f   model.FriendRequestInfo
	)
	for _, member := range *members {
		f.IsOwner = 2
		if member.FromUserID == fromUserID {
			if f, err = r.GetFriendMemberUserInfo(&member, member.ToUserID); err != nil {
				return nil, err
			}
			f.IsOwner = 1
		} else {
			if f, err = r.GetFriendMemberUserInfo(&member, member.FromUserID); err != nil {
				return nil, err
			}
		}
		fs = append(fs, f)
	}
	return fs, nil
}

func (r *friendRequestRepo) Get(userID string, friendID string) (friendRequest *model.FriendRequest, err error) {
	friendRequest = new(model.FriendRequest)
	err = db.DB.Model(&model.FriendRequest{}).Where("from_user_id = ? AND to_user_id = ?", userID, friendID).Find(friendRequest).Error
	return
}

func (r *friendRequestRepo) GetByID(id int64) (friendRequest *model.FriendRequest, err error) {
	friendRequest = new(model.FriendRequest)
	err = db.DB.Model(&model.FriendRequest{}).Where("id = ?", id).Find(friendRequest).Error
	return
}

func (r *friendRequestRepo) Create(userID string, friendID string, msg string, remark string) (friendRequest *model.FriendRequest, err error) {
	friendRequest = new(model.FriendRequest)
	friendRequest.CreatedAt = time.Now().Unix()
	friendRequest.FromUserID = userID
	friendRequest.ToUserID = friendID
	friendRequest.ReqMsg = msg
	friendRequest.Remark = remark

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		var count int64

		if err = tx.Model(&model.FriendRequest{}).Where("from_user_id = ? AND to_user_id = ?", userID, friendID).Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			m := map[string]interface{}{
				"req_msg": friendRequest.ReqMsg,
				"remark":  friendRequest.Remark,
				"status":  friendRequest.Status,
			}

			err = tx.Model(&model.FriendRequest{}).Where("from_user_id = ? AND to_user_id = ?", userID, friendID).Updates(m).Error
		} else {
			err = tx.Create(&friendRequest).Error
		}

		if err != nil {
			return err
		}
		return nil
	})
	return
}

func (r *friendRequestRepo) UpdateStatus(userID string, friendID string, status int) (err error) {
	err = db.DB.Model(&model.FriendRequest{}).Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)", userID, friendID, friendID, userID).Updates(map[string]interface{}{
		"status": status,
	}).Error
	return
}
