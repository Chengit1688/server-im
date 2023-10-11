package repo

import (
	"fmt"
	shoppingModel "im/internal/api/operator/model"
	userModel "im/internal/api/user/model"
	"im/internal/cms_api/operator/model"
	"im/pkg/db"
	"time"
)

var OperatorRepo = new(operatorRepo)

type operatorRepo struct{}

func (r *operatorRepo) FetchList(req model.OperatorListReq) (dataShop []shoppingModel.Operator, count int64, err error) {
	tx := db.DB
	tx = tx.Model(&shoppingModel.Operator{}).Where("status != ?", shoppingModel.ShopStatusDeleted).Preload("CreatorUser")
	if req.ShopID != 0 {
		tx = tx.Where("id = ?", req.ShopID)
	}
	if req.Status != 0 {
		tx = tx.Where("status = ?", req.Status)
	}
	if req.ShopName != "" {
		tx = tx.Where("name LIKE ?", "%"+req.ShopName+"%")
	}
	err = tx.Offset(req.Offset).Limit(req.Limit).Order("created_at desc").Find(&dataShop).Limit(-1).Offset(-1).Count(&count).Error

	return
}

func (r *operatorRepo) FetchTeamList(req model.OperatorMemberListReq) (dataTeam []model.TeamInfo, count int64, err error) {
	tx := db.DB
	shopTeamTable := new(shoppingModel.OperatorTeam).TableName()
	userTable := new(userModel.User).TableName()
	tx = tx.Model(&userModel.User{}).Joins(fmt.Sprintf("INNER JOIN %s ON %s.user_id=%s.user_id", shopTeamTable, shopTeamTable, userTable)).Where(fmt.Sprintf("%s.shop_id = ?", shopTeamTable), req.ShopID)
	if req.Key != "" {
		tx = tx.Where(db.DB.Or(fmt.Sprintf("%s.account like ?", userTable), "%"+req.Key+"%").Or(fmt.Sprintf("%s.nick_name like ?", userTable), "%"+req.Key+"%"))
	}
	err = tx.Offset(req.Offset).Limit(req.Limit).Order(fmt.Sprintf("%s.created_at desc", shopTeamTable)).Find(&dataTeam).Limit(-1).Offset(-1).Count(&count).Error

	return
}

func (r *operatorRepo) Approve(req model.OperatorApproveReq) (data shoppingModel.Operator, err error) {
	data = shoppingModel.Operator{
		Status:      req.Status,
		CommonModel: db.CommonModel{UpdatedAt: time.Now().Unix()},
	}
	err = db.DB.Model(&shoppingModel.Operator{}).Where("id = ?", req.ShopID).
		Where("status = ?", shoppingModel.ShopStatusApprove).
		Updates(&data).Error
	return
}

func (r *operatorRepo) FetchTeamListByInviteUserId(req shoppingModel.OperatorAgentLevelListReq) (dataTeam []shoppingModel.OperatorAgentLevelListInfo, count int64, err error) {
	var shopTeam []shoppingModel.OperatorTeam
	tx := db.DB.Preload("Shop").Preload("InviteUser").Preload("User")
	tx = tx.Model(&shoppingModel.OperatorTeam{}).Where("invite_user_id = ?", req.InviteUserId).
		Where("user_id <> ?", req.InviteUserId)
	if req.UserId != "" {
		tx = tx.Where("user_id = ?", req.UserId)
	}
	if req.ShopID != 0 {
		tx = tx.Where("shop_id =?", req.ShopID)
	}
	if req.BeginDate != 0 {
		tx = tx.Where("created_at >=?", req.BeginDate)
	}
	if req.EndDate != 0 {
		tx = tx.Where("created_at <=?", req.EndDate)
	}
	err = tx.Offset(req.Offset).Limit(req.Limit).Order("id desc").Find(&shopTeam).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		return
	}
	for _, team := range shopTeam {
		dataTeam = append(dataTeam, shoppingModel.OperatorAgentLevelListInfo{
			ShopID:   team.ShopID,
			ShopName: team.Shop.Name,
			TeamInfo: shoppingModel.TeamInfo{
				UserID:      team.UserID,
				Account:     team.User.Account,
				PhoneNumber: team.User.PhoneNumber,
				CountryCode: team.User.CountryCode,
				FaceURL:     team.User.FaceURL,
				BigFaceURL:  team.User.BigFaceURL,
				Gender:      team.User.Gender,
				NickName:    team.User.NickName,
				Age:         team.User.Age,
			},
			CreatedAt: team.CreatedAt,
		})
	}
	return
}
