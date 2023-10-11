package repo

import (
	"im/internal/api/moments/model"
	userModel "im/internal/api/user/model"
	"im/pkg/code"
	"im/pkg/db"
	"strings"
	"time"
)

var TagRepo = new(tagRepo)

type tagRepo struct{}

func (r *tagRepo) AddTag(creatorId string, add model.TagReq) (err error) {
	var count int64
	tx := db.DB.Model(&model.ContactsTag{})
	if err = tx.Where("title = ?", add.Title).Where("creator_id = ?", creatorId).Count(&count).Error; err != nil {
		return err
	}
	if count >= 1 {
		return code.ErrTagExists
	}
	insertData := model.ContactsTag{
		Title:        add.Title,
		CreatorId:    creatorId,
		FriendUserID: strings.Join(add.UserID, ","),
		FriendLength: len(add.UserID),
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}

	return tx.Create(&insertData).Error
}

func (r *tagRepo) AddFriendTag(creatorId, friendID string, tagID int64) (err error) {
	tx := db.DB.Model(&model.ContactsTag{})
	rowData := &model.ContactsTag{}
	if err = tx.Where("creator_id = ?", creatorId).Where("id = ?", tagID).Find(&rowData).Error; err != nil {
		return
	}
	friendIds := strings.Split(rowData.FriendUserID, ",")
	for _, id := range friendIds {
		if id == friendID {
			return code.ErrTagFriendExists
		}
	}
	friendIds = append(friendIds, friendID)
	upData := model.ContactsTag{
		FriendUserID: strings.Join(friendIds, ","),
		FriendLength: rowData.FriendLength + 1,
		UpdatedAt:    time.Now().Unix(),
	}
	if len(friendIds) == 1 {
		upData.FriendUserID = friendID
	}
	err = tx.Where("creator_id = ?", creatorId).Where("id = ?", tagID).Updates(&upData).Error

	return
}

func (r *tagRepo) CheckFriendTag(creatorId, friendID string, tagID int64) (isExists bool, err error) {
	tx := db.DB.Model(&model.ContactsTag{})
	rowData := &model.ContactsTag{}
	if err = tx.Where("creator_id = ?", creatorId).Where("id = ?", tagID).Find(&rowData).Error; err != nil {
		return
	}
	if rowData.FriendLength == 0 {
		return
	}
	friendIds := strings.Split(rowData.FriendUserID, ",")
	for _, id := range friendIds {
		if id == friendID {
			return true, nil
		}
	}
	return
}

func (r *tagRepo) GetFriendTag(creatorId, friendID string) (rowData []model.ContactsTag, err error) {
	if err = db.DB.Model(&model.ContactsTag{}).Where("creator_id = ?", creatorId).
		Where("find_in_set(?,user_id)", friendID).Find(&rowData).Error; err != nil {
		return
	}
	return
}

func (r *tagRepo) GetTag(tagID int64) (rowData model.ContactsTag, err error) {
	if err = db.DB.Model(&model.ContactsTag{}).Where("id = ?", tagID).Find(&rowData).Error; err != nil {
		return
	}
	return
}

func (r *tagRepo) EditeTag(tagID int64, creatorId string, add model.TagReq) (tag model.ContactsTag, err error) {
	userIDJoin := strings.Join(add.UserID, ",")
	data := model.ContactsTag{
		Title:        add.Title,
		FriendUserID: userIDJoin,
		FriendLength: len(add.UserID),
		UpdatedAt:    time.Now().Unix(),
	}
	if len(add.UserID) == 1 {
		data.FriendUserID = add.UserID[0]
		data.FriendLength = 1
	}
	if len(add.UserID) == 0 {
		u := map[string]interface{}{"user_id": "", "friend_length": 0, "title": add.Title, "updated_at": time.Now().Unix()}
		err = db.DB.Model(&model.ContactsTag{}).Where("id = ?", tagID).Where("creator_id = ?", creatorId).Updates(u).Error
		_ = db.DB.Model(&model.ContactsTag{}).Where("id = ?", tagID).Find(&data).Error
		return data, err
	}
	err = db.DB.Model(&model.ContactsTag{}).Where("id = ?", tagID).Where("creator_id = ?", creatorId).Updates(&data).Error
	return data, err
}

func (r *tagRepo) DeleteTagByID(tagID int64, creatorId string) (err error) {
	err = db.DB.Model(&model.ContactsTag{}).Where("id = ?", tagID).Where("creator_id = ?", creatorId).Delete(&model.ContactsTag{}).Error
	return err
}

func (r *tagRepo) ListTagByCreatorID(creatorId string, req model.TagListReq) (tags []model.ContactsTag, count int64, err error) {
	tx := db.DB.Where("creator_id = ?", creatorId)
	err = tx.Offset(req.Offset).Limit(req.Limit).Order("created_at desc").Find(&tags).Limit(-1).Offset(-1).Count(&count).Error
	return
}

func (r *tagRepo) TagCountByCreatorID(creatorId string, tagID int64) (count int64, err error) {
	data := model.ContactsTag{}
	err = db.DB.Model(&model.ContactsTag{}).Where("creator_id = ?", creatorId).Where("id = ?", tagID).Find(&data).Error
	if err != nil {
		return
	}
	return int64(data.FriendLength), nil
}

func (r *tagRepo) TagDetailByID(creatorId string, tagID int64) (tagInfos model.TagDetailResp, err error) {
	var (
		tag   model.ContactsTag
		users []userModel.User
	)
	err = db.DB.Where("creator_id = ?", creatorId).Where("id = ?", tagID).Find(&tag).Error
	if err != nil {
		return
	}
	if err = db.DB.Model(&userModel.User{}).Where("user_id IN ?", strings.Split(tag.FriendUserID, ",")).Find(&users).Error; err != nil {
		return
	}
	tagInfos.Total = tag.FriendLength
	for _, user := range users {
		tagInfos.List = append(tagInfos.List, model.TagUserInfo{
			UserID:     user.UserID,
			Nickname:   user.NickName,
			FaceUrl:    user.FaceURL,
			BigFaceURL: user.BigFaceURL,
		})
	}
	return
}
