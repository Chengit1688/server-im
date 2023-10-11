package repo

import (
	"im/internal/api/friend/model"
	"im/pkg/code"
	"im/pkg/db"
	"im/pkg/util"
	"time"

	"gorm.io/gorm"
)

var FriendRepo = new(friendRepo)

type friendRepo struct{}

type WhereOption struct {
	FriendUserID  string
	UserId        string
	BlackUserId   string
	Remark        string
	Version       int64
	MaxUpdateTime int64
	Offset        int
	Limit         int
}

func (r *friendRepo) Count(userID string) (count int64, err error) {
	err = db.DB.Model(&model.Friend{}).Where("owner_user_id = ? AND status = 1", userID).Count(&count).Error
	return
}

func (r *friendRepo) UpdateFriendRemark(opt WhereOption, data *model.Friend) error {
	if opt.FriendUserID == "" || opt.UserId == "" {
		return code.ErrFriendNotExist
	}

	updates := map[string]interface{}{}
	updates["updated_at"] = time.Now().Unix()
	updates["version"] = data.Version
	updates["remark"] = data.Remark

	query := db.DB.Model(&model.Friend{})
	return query.Where("owner_user_id = ? and friend_user_id = ?", opt.UserId, opt.FriendUserID).Updates(updates).Error
}

func (r *friendRepo) GetFriendRemark(userID string) (data []model.Friend, err error) {
	query := db.DB.Model(&model.Friend{})
	err = query.Where("owner_user_id = ?", userID).Find(&data).Error
	return
}

func (r *friendRepo) UpdateBlackStatus(opt WhereOption, status int) error {
	if opt.FriendUserID == "" || opt.UserId == "" {
		return code.ErrFriendNotExist
	}

	updates := map[string]interface{}{}
	updates["updated_at"] = time.Now().Unix()
	updates["black_status"] = status
	query := db.DB.Model(&model.Friend{})
	return query.Where("owner_user_id = ? and friend_user_id = ?", opt.UserId, opt.FriendUserID).Updates(updates).Error
}

func (r *friendRepo) FriendMaxSeq(opt WhereOption) (int32, int64) {
	query := db.DB
	if opt.UserId != "" {
		query = query.Where("owner_user_id= ?", opt.UserId)
	}
	var c int64
	m := &model.Friend{}
	if query.First(m).Count(&c); c == 0 {
		return 0, 0
	}

	query.First(&model.Friend{}).Select("updated_at").Order("updated_at DESC").First(m)

	return int32(c), m.UpdatedAt
}

func (r *friendRepo) GetFriend(userID string, friendID string) (friend *model.Friend, err error) {
	friend = new(model.Friend)
	err = db.DB.Model(&model.Friend{}).Where("owner_user_id = ? AND friend_user_id = ?", userID, friendID).Find(friend).Error
	return
}

func (r *friendRepo) AddFriend(userID string, friendID string, operatorUserID string, remark string, version int64) (err error) {
	err = db.DB.Transaction(func(tx *gorm.DB) error {
		var (
			count  int64
			friend model.Friend
		)

		if err = tx.Model(&model.Friend{}).Where("owner_user_id = ? AND friend_user_id = ?", userID, friendID).Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			m := map[string]interface{}{
				"operator_user_id": operatorUserID,
				"remark":           remark,
				"version":          version,
				"status":           1,
				"black_status":     model.NotBlack,
			}

			err = tx.Model(&model.Friend{}).Where("owner_user_id = ? AND friend_user_id = ?", userID, friendID).Updates(m).Error
		} else {
			friend.OwnerUserID = userID
			friend.FriendUserID = friendID
			friend.OperatorUserID = operatorUserID
			friend.Remark = remark
			friend.Version = version
			friend.Status = 1
			friend.BlackStatus = model.NotBlack
			err = tx.Create(&friend).Error
		}

		if err != nil {
			return err
		}
		return nil
	})
	return
}

func (r *friendRepo) DeleteFriend(userID string, friendID string, version int64) (err error) {
	m := map[string]interface{}{
		"status":  2,
		"version": version,
	}

	err = db.DB.Model(&model.Friend{}).Where("owner_user_id = ? AND friend_user_id = ?", userID, friendID).Updates(m).Error
	return
}

func (r *friendRepo) CreateFriendLabel(userID, labelID, labelName string) (int64, error) {
	friendLabel := &model.FriendLabel{
		UserId:     userID,
		LabelId:    labelID,
		LabelName:  labelName,
		CreateTime: time.Now().Unix(),
	}
	if err := db.DB.Create(friendLabel).Error; err != nil {
		return 0, err
	}
	return friendLabel.ID, nil
}

func (r *friendRepo) IsLabelNameExist(userID, labelName string) (bool, error) {
	var count int64
	if err := db.DB.Model(&model.FriendLabel{}).Where("user_id = ? AND label_name = ?", userID, labelName).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *friendRepo) GetAllFriendLabels(userID string) ([]*model.FriendLabel, error) {
	friendLabels := make([]*model.FriendLabel, 0)
	err := db.DB.Model(&model.FriendLabel{}).
		Where("user_id = ?", userID).
		Order("create_time ASC").
		Find(&friendLabels).Error
	if err != nil {
		return nil, err
	}
	if len(friendLabels) == 0 {
		return []*model.FriendLabel{}, nil
	}
	return friendLabels, nil
}

func (r *friendRepo) UpdateFriendLabel(userID, labelID, labelName string) error {
	friendLabel := &model.FriendLabel{}
	err := db.DB.Where("user_id = ? AND label_id = ?", userID, labelID).First(friendLabel).Error
	if err != nil {
		return err
	}
	friendLabel.LabelName = labelName
	err = db.DB.Save(friendLabel).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *friendRepo) DeleteFriendLabel(userID, labelID string) error {
	friendLabel := &model.FriendLabel{}
	err := db.DB.Where("user_id = ? AND label_id = ?", userID, labelID).First(friendLabel).Error
	if err != nil {
		return err
	}
	err = db.DB.Delete(friendLabel).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *friendRepo) IsFriendHasLabel(userID, labelID string) (bool, error) {
	var count int64
	if err := db.DB.Model(&model.Friend{}).Where("owner_user_id = ? AND friend_label = ?", userID, labelID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *friendRepo) IsFriendHasRemark(userID, friendUserID, remark string) (bool, error) {
	var count int64
	if err := db.DB.Model(&model.Friend{}).Where("owner_user_id = ? AND friend_user_id = ? AND remark = ?", userID, friendUserID, remark).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *friendRepo) IsFriendLabelExist(userID, labelID string) (bool, error) {
	var count int64
	if err := db.DB.Model(&model.FriendLabel{}).Where("user_id = ? AND label_id = ?", userID, labelID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *friendRepo) ChangeFriendLabel(userID, labelID string, friendID string) error {
	updateFields := map[string]interface{}{
		"friend_label": labelID,
		"updated_at":   time.Now().Unix(),
	}

	query := db.DB.Model(&model.Friend{}).Where("owner_user_id = ? and friend_user_id = ?", userID, friendID)

	return query.Updates(updateFields).Error
}

func (r *friendRepo) ChangeFriendVersion(userID, friendID string) error {
	updateFields := map[string]interface{}{
		"updated_at": time.Now().Unix(),
	}
	query := db.DB.Model(&model.Friend{}).Where("owner_user_id = ? and friend_user_id = ?", userID, friendID)

	return query.Updates(updateFields).Error
}

func (r *friendRepo) GetFriendLabel(userId, labelId string) (*model.FriendLabel, error) {
	var friendLabel model.FriendLabel
	err := db.DB.Where("user_id = ? AND label_id = ?", userId, labelId).First(&friendLabel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &friendLabel, nil
}

func (r *friendRepo) FetchBlackUserList(userID string, req model.GetBlackListReq) (data []model.FriendInfo, count int64, err error) {
	var blackData []model.Black
	tx := db.DB.Model(&model.Black{}).Preload("BlockUser").Where("owner_user_id = ?", userID).Where("status = ?", model.InBlack)
	err = tx.Offset(req.Offset).Limit(req.Limit).Find(&blackData).Limit(-1).Offset(-1).Count(&count).Error
	var fD model.FriendInfo
	for _, datum := range blackData {
		_ = util.Copy(datum.BlockUser, &fD)
		data = append(data, fD)
	}
	return
}
