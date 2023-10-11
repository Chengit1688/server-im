package repo

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	friendusecase "im/internal/api/friend/usecase"
	groupUse "im/internal/api/group/usecase"
	"im/internal/api/moments/model"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/util"
	"strings"
	"time"
)

var MomentsRepo = new(momentsRepo)

type momentsRepo struct{}

func (r *momentsRepo) MomentsAdd(add model.MomentsMessage) (momentsMessage model.MomentsMessage, err error) {
	tx := db.DB.Begin()
	if err = tx.Model(&model.MomentsMessage{}).Create(&add).Error; err != nil {
		tx.Rollback()
		return momentsMessage, err
	}

	if err = r.PushMomentsToFriends(tx, add); err != nil {
		tx.Rollback()
		return momentsMessage, err
	}
	tx.Commit()
	return add, nil
}

func (r *momentsRepo) PushMomentsToFriends(tx *gorm.DB, add model.MomentsMessage) (err error) {
	var (
		pushMoments []model.MomentsInbox
		idList      []string
	)
	switch add.CanSee {
	case model.CanSeePublic:
		idList = friendusecase.FriendUseCase.GetUserFriendIdList(add.UserId)
	case model.CanSeeFriend:
		if add.ShareTagID != "" {
			for _, groupID := range strings.Split(add.ShareTagID, ",") {

				if !groupUse.GroupMemberUseCase.CheckMember(groupID, add.UserId) {
					continue
				}
				gInfo, err1 := groupUse.GroupUseCase.GroupInfo(groupID)
				if err1 != nil {
					continue
				}

				if gInfo.Group.CreateUserId != add.UserId {
					continue
				}
				idList = append(idList, groupUse.GroupUseCase.GroupMemberIdList(groupID)...)
			}
		}

		if add.ShareFriendID != "" {
			friendIDs := strings.Split(add.ShareFriendID, ",")
			for _, friendID := range friendIDs {

				if !friendusecase.FriendUseCase.CheckFriend(add.UserId, friendID) {
					continue
				}
				idList = append(idList, friendID)
			}
		}
	}

	idList = append(idList, add.UserId)

	idList = util.RemoveDuplicateElement(idList)
	for _, userID := range idList {
		pushMoments = append(pushMoments, model.MomentsInbox{
			CreatedAt:       time.Now().Unix(),
			UpdatedAt:       time.Now().Unix(),
			MomentsID:       add.ID,
			FriendUserID:    userID,
			PublisherUserId: add.UserId,
			CommonModel: db.CommonModel{
				CreatedAt: time.Now().Unix(),
			},
		})
	}
	if err = tx.Model(&model.MomentsInbox{}).CreateInBatches(&pushMoments, len(pushMoments)).Error; err != nil {
		logger.Sugar.Errorw("PushMomentsToFriends", "func", util.GetSelfFuncName(), "error", err)
		return
	}
	return
}

func (r *momentsRepo) MomentsMessageInfo(id int64) (momentsMessage model.MomentsMessage, err error) {
	tx := db.DB.Model(&model.MomentsMessage{})
	if err = tx.Where("status = ? and id = ?", 1, id).First(&momentsMessage).Error; err != nil {
		return momentsMessage, err
	}
	return momentsMessage, nil
}

func (r *momentsRepo) MomentsList(ids []string, offset int, limit int, count *int64) (moments []model.MomentsMessage, err error) {
	var listDB *gorm.DB
	listDB = db.DB.Model(&model.MomentsMessage{})

	listDB.Where("status = ?", 1)
	if len(ids) > 0 {
		listDB.Where("user_id IN (?) OR can_see = 2", ids)
	}
	if err = listDB.Count(count).Error; err != nil {
		return
	}
	listDB = listDB.Offset(offset).Limit(limit)
	listDB = listDB.Order("created_at desc")
	if err = listDB.Find(&moments).Error; err != nil {
		return
	}

	return
}

func (r *momentsRepo) GetFriendMomentsList(req model.IssueListReq, count *int64) (moments []model.MomentsMessage, err error) {
	var listDB *gorm.DB
	momentsTable := new(model.MomentsMessage).TableName()
	inboxTable := new(model.MomentsInbox).TableName()
	listDB = db.DB.Table(fmt.Sprintf("%s AS mt", momentsTable)).
		Joins(fmt.Sprintf("JOIN %s AS mi ON mt.id = mi.moments_id", inboxTable)).
		Where("mi.friend_user_id = ?", req.UserID)
	err = listDB.Offset(req.Offset).Limit(req.Limit).Order("mt.created_at desc").Find(&moments).Limit(-1).Offset(-1).Count(count).Error
	return
}

func (r *momentsRepo) GetSelfMomentsList(userID string, req model.IssueListReq, count *int64) (moments []model.MomentsMessage, err error) {
	var listDB *gorm.DB
	listDB = db.DB.Model(&model.MomentsMessage{}).
		Where("status = ?", 1).
		Where("user_id = ?", req.UserID)

	if req.UserID != userID {
		listDB = listDB.Where("can_see <> ?", model.CanSeePrivate)
	}
	err = listDB.Offset(req.Offset).Limit(req.Limit).Order("created_at desc").Find(&moments).Limit(-1).Offset(-1).Count(count).Error
	return
}

func (r *momentsRepo) GetMomentsDetail(ID int64) (moments model.MomentsMessage, err error) {
	err = db.DB.Model(&model.MomentsMessage{}).
		Where("status = ?", 1).
		Where("id = ?", ID).Find(&moments).Error
	return
}

func (r *momentsRepo) MomentsDel(del model.DelIssueReq, userId string, isPrivilege int64) (moments model.MomentsMessage, err error) {

	tx := db.DB.Begin()
	info := new(model.MomentsMessage)
	if err = tx.Model(&model.MomentsMessage{}).Where("id = ?", del.MomentsID).First(&info).Error; err != nil {
		return moments, err
	}
	canDel := false
	if isPrivilege == 1 {
		canDel = true
	} else if userId == info.UserId {
		canDel = true
	} else {
		err = errors.New("非特权用户或者非自己发布的消息没有删除权限")
	}
	if canDel {
		err = tx.Model(&model.MomentsMessage{}).Where("id = ?", del.MomentsID).UpdateColumn("status", 2).Error
	}
	if err != nil {
		tx.Rollback()
		return *info, err
	}

	if err = tx.Where("publisher_user_id = ?", userId).Delete(&model.MomentsInbox{}).Error; err != nil {
		tx.Rollback()
		return *info, err
	}
	tx.Commit()
	return *info, err
}

func (r *momentsRepo) MomentsCommentsAdd(add model.MomentsComments) (err error) {
	tx := db.DB.Model(&model.MomentsComments{})
	if err = tx.Create(&add).Error; err != nil {
		return err
	}
	return
}

func (r *momentsRepo) MomentsInfo(id int64, moments *model.MomentsMessage) (err error) {
	tx := db.DB.Model(moments)
	if err = tx.Where("id = ?", id).First(moments).Error; err != nil {
		return err
	}
	return nil
}

func (r *momentsRepo) MomentsCommentsDel(del model.DelCommentReq, userId string) (moments model.MomentsComments, err error) {

	tx := db.DB.Model(&moments)
	if err = tx.Where("id = ?", del.CommentsID).Where("status = ?", 1).First(&moments).Error; err != nil {
		return moments, err
	}
	momentsMessage, err := r.MomentsMessageInfo(moments.MomentsID)
	if err != nil {
		return moments, err
	}
	if momentsMessage.UserId != userId {
		return moments, errors.New("不是消息发布人,没有删除权限")
	}
	err = tx.Where("id = ?", del.CommentsID).Where("status = ?", 1).UpdateColumn("status", 2).Error
	return moments, err
}

func (r *momentsRepo) MomentsCommentsList(momentsID int64, count *int64, offset, limit int) (moments []model.MomentsComments, err error) {
	tx := db.DB.Model(&model.MomentsComments{}).Preload("ReplyUser")
	tx.Where("status = ? and moments_id = ?", 1, momentsID)
	if err = tx.Count(count).Error; err != nil {
		return
	}

	tx = tx.Order("created_at desc")
	if err = tx.Find(&moments).Error; err != nil {
		return
	}
	return moments, nil
}

func (r *momentsRepo) MomentsCommentsLikeGet(id int64, userId string) (info model.MomentsCommentsLike, err error) {

	tx := db.DB.Model(model.MomentsCommentsLike{})
	err = tx.Where("moments_id = ? and friend_user_id = ?", id, userId).First(&info).Error
	return
}

func (r *momentsRepo) MomentsCommentsLikeAdd(req model.MomentsCommentsLikeReq, userId string) (info model.MomentsCommentsLike, err error) {

	tx := db.DB.Model(model.MomentsCommentsLike{})

	info.MomentsID = req.MomentsID
	info.FriendUserID = userId
	info.CreatedAt = time.Now().Unix()
	info.Status = 2
	err = tx.Create(&info).Error
	if err != nil {
		logger.Sugar.Errorw("db error", "func", util.GetSelfFuncName(), "error", err)
		return info, err
	}
	return info, nil
}

func (r *momentsRepo) MomentsCommentsLikeUpdate(Status int64, userId string, momentsID int64) (status int64, err error) {
	tx := db.DB.Model(model.MomentsCommentsLike{})
	if Status == 1 {
		Status = 2
	} else {
		Status = 1
	}

	err = tx.Where("moments_id = ? and friend_user_id = ?", momentsID, userId).UpdateColumn("status", Status).Error
	return Status, err
}

func (r *momentsRepo) MomentsCommentsLikeList(MomentsID int64, count *int64) (list []model.MomentsCommentsLike, err error) {
	tx := db.DB.Model(&model.MomentsCommentsLike{})
	if err = tx.Where("moments_id = ? and status = 2", MomentsID).Count(count).Error; err != nil {
		return list, err
	}
	if err = tx.Where("moments_id = ? and status = 2", MomentsID).Find(&list).Error; err != nil {
		return list, err
	}
	return list, nil
}
