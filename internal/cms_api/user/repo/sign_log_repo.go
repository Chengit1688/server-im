package repo

import (
	"fmt"
	userApiModel "im/internal/api/user/model"
	configModel "im/internal/cms_api/config/model"
	userModel "im/internal/cms_api/user/model"
	"im/internal/cms_api/wallet/model"
	"im/pkg/code"
	"im/pkg/common/constant"
	"im/pkg/db"
	"im/pkg/util"
	"time"
)

var SignLogRepo = new(signLogRepo)

type signLogRepo struct{}

func (r *signLogRepo) GetMonthSignLogByUserId(uid string) (int, []int) {
	dbConn := db.DB.Model(userModel.SignLog{})
	dbConn = dbConn.Where("user_id =?", uid)

	var total int64
	dbConn.Count(&total)

	var days []int

	startTime := time.Now().Format("2006-01") + "-01 00:00:00"
	endTime := time.Now().AddDate(0, 1, 0).Format("2006-01") + "-01 00:00:00"

	startTimeUnix, _ := time.Parse("2006-01-02 15:04:05", startTime)
	endTimeUnix, _ := time.Parse("2006-01-02 15:04:05", endTime)

	dbConn = dbConn.Where("created_at >= ? and created_at <?", startTimeUnix.Unix(), endTimeUnix.Unix())
	var h []*userModel.SignLog
	err := dbConn.Find(&h).Error
	if err != nil {
		return 0, days
	}
	for _, v := range h {
		days = append(days, v.Day)
	}
	return int(total), days
}

func (r *signLogRepo) GetMonthSignLogByUserIdV2(uid string) (total int64, days []string) {
	var (
		h []*userModel.SignLog
	)
	dbConn := db.DB.Model(userModel.SignLog{})
	dbConn = dbConn.Where("user_id =?", uid).Where("created_at >= ? and created_at < ?", util.GetFirstDateOfWeek(time.Now()).Unix(), util.GetLastDateOfWeek(time.Now()).Unix())
	dbConn.Count(&total)
	err := dbConn.Find(&h).Error
	if err != nil {
		return 0, days
	}
	for _, v := range h {
		days = append(days, fmt.Sprintf("%02d-%02d-%02d", v.Year, v.Month, v.Day))
	}
	return total, days
}

func (r *signLogRepo) CreateSignLog(uid string, balance, award int64) error {
	timeInt64 := time.Now().Unix()
	var err error
	dbConn := db.DB.Model(&userModel.SignLog{})
	n := time.Now()
	data := &userModel.SignLog{
		UserId: uid,
		Month:  int(n.Month()),
		Reward: award,
		Day:    n.Day(),
		Year:   n.Year(),
		Time:   n.Format("2006-01-02"),
		CommonModel: db.CommonModel{
			CreatedAt: time.Now().Unix(),
		},
	}
	var had int64
	dbConn.Where("user_id =? and `time`=?", data.UserId, data.Time).Count(&had)
	if had > 0 {
		return code.ErrSignLog
	}
	err = dbConn.Create(data).Error
	if err != nil {
		return err
	}

	if award > 0 {

		err = db.DB.Model(&model.BillingRecords{}).Create(map[string]interface{}{
			"type":          model.TypeBillingRecordWithdrawSignAward,
			"sender_id":     uid,
			"receiver_id":   uid,
			"amount":        award,
			"change_before": balance,
			"change_after":  balance + award,
			"note":          "签到奖励",
			"created_at":    timeInt64,
			"updated_at":    timeInt64,
		}).Error

		if err != nil {
			return err
		}
		return nil
	}

	return err
}

func (r *signLogRepo) GetList(req configModel.SignLogListReq) ([]configModel.SignLogInfo, int64, error) {
	req.Pagination.Check()
	var (
		list  []configModel.SignLogInfo
		count int64
		err   error
	)
	userTable := new(userApiModel.User).TableName()
	signLogTable := new(userModel.SignLog).TableName()
	field := fmt.Sprintf("%s.id,u.nick_name,u.user_id,%s.created_at,%s.reward", signLogTable, signLogTable, signLogTable)
	query := db.DB.Model(userModel.SignLog{}).Select(field).Joins(fmt.Sprintf("join %s as u ON u.user_id = %s.user_id", userTable, signLogTable))
	if req.Id != 0 {
		query = query.Where(fmt.Sprintf("%s.id = ?", signLogTable), req.Id)
	}
	if req.UserId != "" {
		query = query.Where("u.user_id = ?", req.UserId)
	}
	if req.NickName != "" {
		query = query.Where("u.nick_name like ?", "%"+req.NickName+"%")
	}
	if req.BeginDate != 0 {
		query = query.Where(fmt.Sprintf("%s.created_at >= ?", signLogTable), req.BeginDate)
	}
	if req.EndDate != 0 {
		query = query.Where(fmt.Sprintf("%s.created_at <= ?", signLogTable), req.EndDate)
	}
	query = query.Where(fmt.Sprintf("%s.status = ?", signLogTable), constant.SwitchOn)
	err = query.Order(fmt.Sprintf("%s.created_at DESC", signLogTable)).Offset(req.Offset).Limit(req.Limit).Find(&list).Offset(-1).Limit(-1).Count(&count).Error

	return list, count, err
}
