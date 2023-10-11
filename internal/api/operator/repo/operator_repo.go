package repo

import (
	"fmt"
	"im/internal/api/operator/model"
	userModel "im/internal/api/user/model"
	"im/pkg/code"
	"im/pkg/db"
	"im/pkg/util"
	"strings"
	"time"
)

var OperatorRepo = new(operatorRepo)

type operatorRepo struct{}

func (r *operatorRepo) AddShop(creatorId string, req model.OperatorApplyForReq) (data model.Operator, err error) {
	var count int64
	tx := db.DB.Model(&model.Operator{})
	if err = tx.Where("creator_id = ?", creatorId).Where("status != ?", model.ShopStatusDeleted).Count(&count).Error; err != nil {
		return
	}
	if count >= 1 {
		return data, code.ErrShopExists
	}
	InviteCode := util.RandStringInt(8)
	if err = tx.Where("invite_code = ?", InviteCode).Count(&count).Error; err != nil {
		return
	}
	if count > 1 {
		return data, code.ErrInviteCodeExists
	}
	if req.Longitude == "0" {
		req.Longitude = ""
	}
	if req.Latitude == "0" {
		req.Latitude = ""
	}
	tx = db.DB.Begin()
	data = model.Operator{
		Name:        req.Name,
		Longitude:   req.Longitude,
		Latitude:    req.Latitude,
		Address:     req.Address,
		License:     req.License,
		Image:       strings.Join(req.Image, ","),
		Description: req.Description,
		CreatorId:   creatorId,
		InviteCode:  InviteCode,
		Status:      model.ShopStatusApprove,
		CityCode:    req.CityCode,
		Star:        req.Star,
		CommonModel: db.CommonModel{CreatedAt: time.Now().Unix()},
		ShopType:    req.ShopType,
	}
	if err = tx.Model(&model.Operator{}).Create(&data).Error; err != nil {
		tx.Rollback()
		return
	}

	teamData := model.OperatorTeam{
		ShopID:       data.ID,
		Role:         model.TeamRoleLeader,
		InviteUserId: creatorId,
		UserID:       creatorId,
		Status:       model.ShopStatusPass,
		CommonModel:  db.CommonModel{CreatedAt: time.Now().Unix()},
	}
	if err = tx.Model(&model.OperatorTeam{}).Create(&teamData).Error; err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func (r *operatorRepo) UpdateShop(creatorId string, req model.OperatorApplyForReq) (data model.Operator, err error) {
	var (
		shop model.Operator
	)
	tx := db.DB.Model(&model.Operator{})
	if err = tx.Where("creator_id = ?", creatorId).Where("id = ?", req.ShopID).Find(&shop).Error; err != nil {
		return
	}
	if shop.ID == 0 {
		return data, code.ErrShopNotExists
	}
	if shop.Status == model.ShopStatusDeleted {
		return data, code.ErrShopNotExists
	}
	data = model.Operator{
		Name:        req.Name,
		Longitude:   req.Longitude,
		Latitude:    req.Latitude,
		Address:     req.Address,
		License:     req.License,
		Image:       strings.Join(req.Image, ","),
		Description: req.Description,
		CreatorId:   creatorId,
		CommonModel: db.CommonModel{UpdatedAt: time.Now().Unix()},
	}
	if err = tx.Updates(&data).Error; err != nil {
		return
	}

	return
}

func (r *operatorRepo) UpdateShopByID(req model.UpdateOperatorReq) (data model.Operator, err error) {
	var (
		shop model.Operator
	)
	tx := db.DB.Model(&model.Operator{})
	if err = tx.Where("id = ?", req.ShopID).Find(&shop).Error; err != nil {
		return
	}
	if shop.ID == 0 {
		return data, code.ErrShopNotExists
	}
	if shop.Status == model.ShopStatusDeleted {
		return data, code.ErrShopNotExists
	}
	data = model.Operator{
		Name:            req.Name,
		Longitude:       req.Longitude,
		Latitude:        req.Latitude,
		Address:         req.Address,
		License:         req.License,
		Image:           strings.Join(req.Image, ","),
		Description:     req.Description,
		CommonModel:     db.CommonModel{UpdatedAt: time.Now().Unix()},
		DecorationScore: req.DecorationScore,
		QualityScore:    req.QualityScore,
		ServiceScore:    req.ServiceScore,
		Star:            req.Star,
		ShopType:        req.ShopType,
		Status:          req.Status,
	}
	if err = tx.Where("id = ?", req.ShopID).Updates(&data).Error; err != nil {
		return
	}

	return
}

func (r *operatorRepo) FetchShop(req model.OperatorIDCommonReq) (data model.OperatorDetailResp, err error) {
	var shop model.Operator
	tx := db.DB.Model(&model.Operator{})
	if err = tx.Where("id = ?", req.ShopID).
		Where("status = ?", model.ShopStatusPass).Find(&shop).Error; err != nil {
		return
	}
	if shop.ID == 0 {
		return data, code.ErrShopNotExists
	}
	data = model.OperatorDetailResp{
		ID:              shop.ID,
		Name:            shop.Name,
		Longitude:       shop.Longitude,
		Latitude:        shop.Latitude,
		Address:         shop.Address,
		Image:           strings.Split(shop.Image, ","),
		Description:     shop.Description,
		DecorationScore: shop.DecorationScore,
		QualityScore:    shop.QualityScore,
		ServiceScore:    shop.ServiceScore,
		CreatedAt:       shop.CreatedAt,
		UpdatedAt:       shop.UpdatedAt,
		ShopType:        shop.ShopType,
		License:         shop.License,
		CreatorId:       shop.CreatorId,
		InviteCode:      shop.InviteCode,
		Status:          shop.Status,
	}
	return
}

func (r *operatorRepo) FetchTeamList(creatorId string, req model.OperatorTeamListReq) (dataTeam []model.TeamInfo, count int64, err error) {
	var shop model.Operator
	tx := db.DB
	if err = tx.Model(&model.Operator{}).Where("id = ?", req.ShopID).Where("status = ?", model.ShopStatusPass).Find(&shop).Error; err != nil {
		return
	}

	if shop.CreatorId != creatorId {
		return nil, 0, code.ErrDB
	}
	shopTeamTable := new(model.OperatorTeam).TableName()
	userTable := new(userModel.User).TableName()
	tx = tx.Model(&userModel.User{}).Joins(fmt.Sprintf("INNER JOIN %s ON %s.user_id=%s.user_id", shopTeamTable, shopTeamTable, userTable)).Where(fmt.Sprintf("%s.shop_id = ?", shopTeamTable), req.ShopID)
	if req.Key != "" {
		tx = tx.Where(db.DB.Or(fmt.Sprintf("%s.account like ?", userTable), "%"+req.Key+"%").Or(fmt.Sprintf("%s.nick_name like ?", userTable), "%"+req.Key+"%"))
	}
	err = tx.Offset(req.Offset).Limit(req.Limit).Order(fmt.Sprintf("%s.created_at desc", shopTeamTable)).Find(&dataTeam).Limit(-1).Offset(-1).Count(&count).Error

	return
}

func (r *operatorRepo) CheckShop(shopID int64) (data model.OperatorTeam, err error) {
	var shop model.Operator
	tx := db.DB.Model(&model.Operator{})
	if err = tx.Where("id = ?", shopID).
		Where("status = ?", model.ShopStatusPass).Find(&shop).Error; err != nil {
		return
	}
	if shop.ID == 0 {
		return data, code.ErrShopNotExists
	}
	return
}

func (r *operatorRepo) JoinShopTeam(req model.OperatorJoinTeamReq) (data model.OperatorTeam, shop model.Operator, err error) {
	var (
		count int64
	)
	tx := db.DB.Model(&model.Operator{}).Where("status = ?", model.ShopStatusPass)
	if req.ShopID != 0 {
		tx = tx.Where("id = ?", req.ShopID)
	}
	if req.InviteCode != "" {
		tx = tx.Where("invite_code = ?", req.InviteCode)
	}
	if err = tx.Find(&shop).Error; err != nil {
		return
	}
	if shop.ID == 0 {
		return data, shop, code.ErrShopNotExists
	}
	if req.InviteCode != "" && shop.InviteCode != req.InviteCode {
		return data, shop, code.ErrInviteCodeNotExists
	}
	if err = db.DB.Model(&model.OperatorTeam{}).Where("user_id = ?", req.UserID).Count(&count).Error; err != nil {
		return data, shop, err
	}
	if count >= 1 {
		return data, shop, code.ErrShopTeamMemberExists
	}

	if shop.CreatorId == req.UserID {
		return data, shop, code.ErrShopExSelf
	}

	data = model.OperatorTeam{
		ShopID:       shop.ID,
		UserID:       req.UserID,
		Role:         model.TeamRoleNobody,
		InviteUserId: shop.CreatorId,
		Status:       model.ShopStatusPass,
		CommonModel:  db.CommonModel{CreatedAt: time.Now().Unix()},
	}
	if req.InviteUserId != "" {
		data.InviteUserId = req.InviteUserId
	}
	if err = db.DB.Model(&model.OperatorTeam{}).Create(&data).Error; err != nil {
		return
	}

	return
}

func (r *operatorRepo) RemoveShopTeam(req model.OperatorRemoveTeamReq) (err error) {
	err = db.DB.Where("user_id = ?", req.UserID).Where("shop_id = ?", req.ShopID).Delete(&model.OperatorTeam{}).Error
	return
}

func (r *operatorRepo) FetchShopTeamUser(req model.OperatorTeamMemberInfoReq) (info model.OperatorTeamMemberInfoResp, err error) {
	var shop model.Operator
	tx := db.DB.Model(&model.Operator{})
	if err = tx.Where("id = ?", req.ShopID).Where("status = ?", model.ShopStatusPass).Find(&shop).Error; err != nil {
		return
	}
	info.ShopID = shop.ID
	info.UserID = req.UserID
	info.Role = model.TeamRoleNobody
	if req.UserID == shop.CreatorId {
		info.Role = model.TeamRoleLeader
	}
	var userInfo userModel.User
	if err = db.DB.Model(&userModel.User{}).Where("user_id = ?", req.UserID).Find(&userInfo).Error; err != nil {
		return model.OperatorTeamMemberInfoResp{}, err
	}
	if userInfo.UserID == "" {
		return model.OperatorTeamMemberInfoResp{}, code.ErrUserNotFound
	}
	info.UserInfo = model.TeamInfo{
		UserID:      userInfo.UserID,
		Account:     userInfo.Account,
		PhoneNumber: userInfo.PhoneNumber,
		CountryCode: userInfo.CountryCode,
		FaceURL:     userInfo.FaceURL,
		BigFaceURL:  userInfo.BigFaceURL,
		Gender:      userInfo.Gender,
		NickName:    userInfo.NickName,
		Age:         userInfo.Age,
	}
	return
}

func (r *operatorRepo) FetchShopTeamLeaderUser(req model.OperatorTeamLeaderInfoReq) (info model.OperatorTeamLeaderInfoResp, err error) {
	var shopTeam model.OperatorTeam
	tx := db.DB.Model(&model.OperatorTeam{})
	if err = tx.Where("user_id = ?", req.UserID).Preload("Shop").Find(&shopTeam).Error; err != nil {
		return
	}
	info.ShopID = shopTeam.Shop.ID
	info.UserID = req.UserID
	info.HasShop = 1
	info.ShopName = shopTeam.Shop.Name
	info.Role = shopTeam.Role
	if shopTeam.Shop.ID == 0 {
		info.HasShop = 2
	}
	info.ShopStatus = shopTeam.Shop.Status
	var userInfo userModel.User
	if err = db.DB.Model(&userModel.User{}).Where("user_id = ?", req.UserID).Find(&userInfo).Error; err != nil {
		return model.OperatorTeamLeaderInfoResp{}, err
	}
	if userInfo.UserID == "" {
		return model.OperatorTeamLeaderInfoResp{}, code.ErrUserNotFound
	}
	info.UserInfo = model.TeamInfo{
		UserID:      userInfo.UserID,
		Account:     userInfo.Account,
		PhoneNumber: userInfo.PhoneNumber,
		CountryCode: userInfo.CountryCode,
		FaceURL:     userInfo.FaceURL,
		BigFaceURL:  userInfo.BigFaceURL,
		Gender:      userInfo.Gender,
		NickName:    userInfo.NickName,
		Age:         userInfo.Age,
	}
	return
}

func (r *operatorRepo) SearchShop(req model.SearchOperatorReq) (shops []model.SearchDTO, count int64, err error) {
	var fields string
	tx := db.DB.Model(&model.Operator{})
	if req.Latitude == "" || req.Longitude == "" || req.Longitude == "0" || req.Latitude == "0" {
		fields = "id,address,license,latitude,shop_type,description,address,longitude,image,decoration_score,quality_score,service_score,name,created_at,0 AS distance"
		tx = tx.Order("created_at desc")
	} else {
		acos := fmt.Sprintf("6371 * acos (  cos ( radians(%s) ) * cos( radians( latitude ) ) * cos( radians( longitude ) - radians(%s) )  + sin ( radians(%s) ) * sin( radians( latitude ) )  ) ", req.Latitude, req.Longitude, req.Latitude)

		fields = fmt.Sprintf("id,latitude,longitude,shop_type,description,address,license,image,decoration_score,quality_score,service_score,name,created_at,(%s) AS distance", acos)
		tx = tx.Order("distance asc")
	}
	tx = tx.Select(fields).Where("status = ?", model.ShopStatusPass)
	if req.CityCode != 0 {
		tx = tx.Where("city_code=?", req.CityCode)
	}
	if req.Key != "" {
		tx = tx.Where("name like ?", "%"+req.Key+"%")
	}
	err = tx.Offset(req.Offset).Limit(req.Limit).Find(&shops).Limit(-1).Offset(-1).Count(&count).Error

	return
}
