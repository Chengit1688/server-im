package repo

import (
	"fmt"
	chatModel "im/internal/api/chat/model"
	groupModel "im/internal/api/group/model"
	userModel "im/internal/api/user/model"
	"im/internal/cms_api/dashboard/model"
	"im/pkg/db"
	"im/pkg/util"
	"time"
)

var LOC, _ = time.LoadLocation("Asia/Shanghai")

var DashboardRepo = new(dashboardRepo)

type dashboardRepo struct{}

func (r *dashboardRepo) GetRegisterNumToday() (count int64, err error) {
	today := time.Now().Format("2006-01-02")
	todayStartStr := fmt.Sprintf("%s 00:00:00", today)
	todayEndStr := fmt.Sprintf("%s 23:59:59", today)
	todayStart, _ := time.ParseInLocation("2006-01-02 15:04:05", todayStartStr, LOC)
	todayEnd, _ := time.ParseInLocation("2006-01-02 15:04:05", todayEndStr, LOC)
	tx := db.DB.Model(userModel.User{})
	tx = tx.Where("created_at >= ?", todayStart.Unix())
	tx = tx.Where("created_at <= ?", todayEnd.Unix())
	err = tx.Count(&count).Error
	return
}

func (r *dashboardRepo) GetRegisterNumYesterday() (count int64, err error) {
	yestday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	yestdayStartStr := fmt.Sprintf("%s 00:00:00", yestday)
	yestdayEndStr := fmt.Sprintf("%s 23:59:59", yestday)
	todayStart, _ := time.ParseInLocation("2006-01-02 15:04:05", yestdayStartStr, LOC)
	todayEnd, _ := time.ParseInLocation("2006-01-02 15:04:05", yestdayEndStr, LOC)
	tx := db.DB.Model(userModel.User{})
	tx = tx.Where("created_at >= ?", todayStart.Unix())
	tx = tx.Where("created_at <= ?", todayEnd.Unix())
	err = tx.Count(&count).Error
	return
}

func (r *dashboardRepo) GetLoginNumToday() (count int64, err error) {
	today := time.Now().Format("2006-01-02")
	todayStartStr := fmt.Sprintf("%s 00:00:00", today)
	todayEndStr := fmt.Sprintf("%s 23:59:59", today)
	todayStart, _ := time.ParseInLocation("2006-01-02 15:04:05", todayStartStr, LOC)
	todayEnd, _ := time.ParseInLocation("2006-01-02 15:04:05", todayEndStr, LOC)
	tx := db.DB.Model(userModel.User{})
	tx = tx.Where("latest_login_time >= ?", todayStart.Unix())
	tx = tx.Where("latest_login_time <= ?", todayEnd.Unix())
	err = tx.Count(&count).Error
	return
}

func (r *dashboardRepo) GetGroupNum() (count int64, err error) {
	tx := db.DB.Model(groupModel.Group{})
	tx = tx.Where("status = ?", 1)
	err = tx.Count(&count).Error
	return
}

func (r *dashboardRepo) GetRegisterNumDaily() (data model.DataBarList, err error) {
	now := time.Now()
	currentYear := now.Year()
	currentMonth := now.Month()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, LOC)
	endOfMonth := time.Date(currentYear, currentMonth, 1, 23, 59, 59, 0, LOC)
	lastOfMonth := endOfMonth.AddDate(0, 1, -1)

	field := "FROM_UNIXTIME(`created_at`, '%Y-%m-%d') as date,count(*) as count"
	query := db.DB.Model(userModel.User{}).Select(field)
	query = query.Where("created_at >= ?", firstOfMonth.Unix())
	query = query.Where("created_at <= ?", lastOfMonth.Unix())
	err = query.Group("date").Order("date ASC").Find(&data).Error
	return
}

func (r *dashboardRepo) GetLoginNumDaily() (data model.DataBarList, err error) {
	now := time.Now()
	currentYear := now.Year()
	currentMonth := now.Month()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, LOC)
	endOfMonth := time.Date(currentYear, currentMonth, 1, 23, 59, 59, 0, LOC)
	lastOfMonth := endOfMonth.AddDate(0, 1, -1)

	field := "FROM_UNIXTIME(`latest_login_time`, '%Y-%m-%d') as date,count(*) as count"
	query := db.DB.Model(userModel.User{}).Select(field)
	query = query.Where("latest_login_time >= ?", firstOfMonth.Unix())
	query = query.Where("latest_login_time <= ?", lastOfMonth.Unix())
	err = query.Group("date").Order("date ASC").Find(&data).Error
	return
}

func (c *dashboardRepo) GetTodayMessageCount(cType chatModel.ConversationType) (count int64, err error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, LOC)
	err = db.DB.Model(chatModel.Message{}).Where("conversation_type = ?", cType).Where("send_time >= ?", util.UnixMilliTime(today)).Count(&count).Error
	return
}

func (c *dashboardRepo) GetMonthDailyMessage(cType chatModel.ConversationType) (data model.DataBarList, err error) {
	now := time.Now()
	currentYear := now.Year()
	currentMonth := now.Month()
	start := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, LOC)
	end := time.Date(currentYear, currentMonth, now.Day(), 0, 0, 0, 0, LOC)

	field := "FROM_UNIXTIME(left(`send_time`,10), '%Y-%m-%d') as date,count(*) as count"
	query := db.DB.Model(chatModel.Message{}).Select(field)
	query = query.Where("conversation_type = ?", cType)
	query = query.Where("send_time >= ?", util.UnixMilliTime(start))
	query = query.Where("send_time < ?", util.UnixMilliTime(end))
	err = query.Group("date").Order("date ASC").Find(&data).Error
	return
}

func (c *dashboardRepo) GetYestdayMessageCount(cType chatModel.ConversationType, now time.Time) (count int64, err error) {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, LOC)
	yestday := today.AddDate(0, 0, -1)
	err = db.DB.Model(chatModel.Message{}).Where("conversation_type = ?", cType).Where("send_time >= ?", util.UnixMilliTime(yestday)).Where("send_time < ?", util.UnixMilliTime(today)).Count(&count).Error
	return
}

func (c *dashboardRepo) AddDashboardDailyData(data model.DashboardDailyData) (err error) {
	err = db.DB.Model(model.DashboardDailyData{}).Create(&data).Error
	return
}

func (c *dashboardRepo) FindDashboardDailyData(start, end, dType int64) (datas []model.DashboardDailyData, err error) {
	tx := db.DB.Model(model.DashboardDailyData{})
	tx.Where("datetime >= ?", start)
	tx.Where("datetime <= ?", end)
	tx.Where("type = ?", dType)
	err = tx.Find(&datas).Error
	return
}
