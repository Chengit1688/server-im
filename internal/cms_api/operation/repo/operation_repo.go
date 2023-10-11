package repo

import (
	"fmt"
	apiUserModel "im/internal/api/user/model"
	configModel "im/internal/cms_api/config/model"
	"im/internal/cms_api/operation/model"
	"im/pkg/db"
)

var OperationRepo = new(operationRepo)

type operationRepo struct{}

type WhereOptionForRegStatistics struct {
	StartDate string
	EndDate   string
}

func (r *operationRepo) GetRegistrationStatistics(req model.RegistrationStatisticsReq) ([]model.RegistrationStatisticsInfo, error) {
	req.Pagination.Check()
	var (
		err error
		m   []model.RegistrationStatisticsInfo
	)
	field := "FROM_UNIXTIME(`created_at`, '%Y-%m-%d') as daily,count(*) as count"
	query := db.DB.Model(apiUserModel.User{}).Select(field)
	if req.BeginDate != 0 {
		query = query.Where("created_at >= ?", req.BeginDate)
	}
	if req.EndDate != 0 {
		query = query.Where("created_at <= ?", req.EndDate)
	}
	err = query.Group("daily").Order("daily ASC").Find(&m).Error

	return m, err
}

func (r *operationRepo) GetInviteCodeStatistics(req model.InviteCodeStatisticsReq) ([]model.InviteCodeStatisticsInfo, int64, error) {
	req.Pagination.Check()
	var (
		err   error
		count int64
		m     []model.InviteCodeStatisticsInfo
	)
	userTable := new(apiUserModel.User).TableName()
	inviteTable := new(configModel.InviteCode).TableName()
	field := fmt.Sprintf("%s.invite_code,count(*) count", inviteTable)
	query := db.DB.Model(configModel.InviteCode{}).Select(field).
		Joins(fmt.Sprintf("join %s as u ON u.invite_code=%s.invite_code", userTable, inviteTable))
	if req.InviteCode != "" {
		query = query.Where(fmt.Sprintf("%s.invite_code like ?", inviteTable), "%"+req.InviteCode+"%")
	}
	query = query.Group(fmt.Sprintf("%s.invite_code", inviteTable))
	query.Count(&count)
	err = query.Offset(req.Offset).Limit(req.Limit).Find(&m).Error
	if len(m) == 1 {
		count = 1
	}
	return m, count, err
}

func (r *operationRepo) GetInviteCodeStatisticsDetails(req model.InviteCodeStatisticsDetailsReq) ([]model.InviteCodeStatisticsDetailsInfo, int64, error) {
	req.Pagination.Check()
	var (
		err   error
		count int64
		m     []model.InviteCodeStatisticsDetailsInfo
	)
	field := "invite_code,user_id,account,phone_number,face_url,nick_name,balance,created_at as registry_time,latest_login_time"
	query := db.DB.Model(apiUserModel.User{}).Select(field)
	if req.InviteCode != "" {
		query = query.Where("invite_code = ?", req.InviteCode)
	}
	if req.Account != "" {
		query = query.Where("account = ?", req.Account)
	}
	if req.NickName != "" {
		query = query.Where("nick_name = ?", req.NickName)
	}
	if req.PhoneNumber != "" {
		query = query.Where("phone_number = ?", req.PhoneNumber)
	}
	if req.BeginDate != 0 {
		query = query.Where("created_at >= ?", req.BeginDate)
	}
	if req.EndDate != 0 {
		query = query.Where("created_at <= ?", req.EndDate)
	}
	err = query.Order("created_at asc").Offset(req.Offset).Limit(req.Limit).Find(&m).Offset(-1).Limit(-1).Count(&count).Error

	return m, count, err
}
