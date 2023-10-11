package service

import (
	"im/internal/cms_api/dashboard/model"
	"im/internal/cms_api/dashboard/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

var DashboardService = new(dashboardService)

type dashboardService struct{}

func (s *dashboardService) Info(c *gin.Context) {
	req := new(model.GetDashboardReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	registerNum, err := repo.DashboardRepo.GetRegisterNumToday()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}

	loginNum, err := repo.DashboardRepo.GetLoginNumToday()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}

	registerDailyNum, err := repo.DashboardRepo.GetRegisterNumDaily()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	registerDailyNum = checkMonthData(registerDailyNum)

	loginDailyNum, err := repo.DashboardRepo.GetLoginNumDaily()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	loginDailyNum = checkMonthData(loginDailyNum)

	grouprNum, err := repo.DashboardRepo.GetGroupNum()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}

	sigleMsgCount, err := repo.DashboardRepo.GetTodayMessageCount(1)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}

	groupMsgCount, err := repo.DashboardRepo.GetTodayMessageCount(2)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}

	sigleMsgDaily := getDailyData(sigleMsgCount, 1)

	groupMsgDaily := getDailyData(groupMsgCount, 2)

	onlineMax := repo.DashboardCache.GetOnlineMax()
	ret := new(model.GetDashboardResp)
	ret.RegisterNum = registerNum
	ret.GroupNum = grouprNum
	ret.LoginNum = loginNum
	ret.OnlineMax = onlineMax
	ret.SigleMsgNum = sigleMsgCount
	ret.GroupMsgNum = groupMsgCount

	ret.GroupMsgBar = groupMsgDaily
	ret.LoginBar = loginDailyNum
	ret.RegisterBar = registerDailyNum
	ret.SigleMsgBar = sigleMsgDaily
	http.Success(c, ret)
	return
}

func checkMonthData(data model.DataBarList) model.DataBarList {
	var LOC, _ = time.LoadLocation("Asia/Shanghai")
	now := time.Now()
	currentYear := now.Year()
	currentMonth := now.Month()
	endOfMonth := time.Date(currentYear, currentMonth, 1, 23, 59, 59, 0, LOC)
	lastOfMonth := endOfMonth.AddDate(0, 1, -1)

	for i := 1; i <= lastOfMonth.Day(); i++ {
		day := time.Date(currentYear, currentMonth, i, 0, 0, 0, 0, LOC).Format("2006-01-02")
		var exist = 0
		for _, check := range data {
			if check.Date == day {
				exist = 1
			}
		}
		if exist == 0 {
			data = append(data, model.DataBar{Date: day, Count: 0})
		}
	}
	sort.Sort(data)
	return data
}

func getDailyData(todayCount int64, Type int64) model.DataBarList {
	nowTime := time.Now()
	currentYear := nowTime.Year()
	currentMonth := nowTime.Month()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, repo.LOC)
	endOfMonth := time.Date(currentYear, currentMonth, 1, 23, 59, 59, 0, repo.LOC)
	lastOfMonth := endOfMonth.AddDate(0, 1, -1)

	Dailys, err := repo.DashboardRepo.FindDashboardDailyData(firstOfMonth.Unix(), lastOfMonth.Unix(), Type)
	if err != nil {
		logger.Sugar.Errorw("FindDashboardDailyData", "func", util.GetSelfFuncName(), "error", err)
		return nil
	}
	var DailyDataBar model.DataBarList
	for _, DailyItem := range Dailys {
		DailyDataBar = append(DailyDataBar, model.DataBar{Date: time.Unix(DailyItem.DateTime, 0).Format("2006-01-02"), Count: DailyItem.Count})
	}

	todayStr := nowTime.Format("2006-01-02")
	DailyDataBar = append(DailyDataBar, model.DataBar{Date: todayStr, Count: todayCount})
	Daily := checkMonthData(DailyDataBar)
	return Daily
}

func QueryDailyData(todayCount int64, Type int64, start, end int64) model.DataBarList {
	nowTime := time.Now()

	Dailys, err := repo.DashboardRepo.FindDashboardDailyData(start, end, Type)
	if err != nil {
		logger.Sugar.Errorw("FindDashboardDailyData", "func", util.GetSelfFuncName(), "error", err)
		return nil
	}
	var DailyDataBar model.DataBarList
	for _, DailyItem := range Dailys {
		DailyDataBar = append(DailyDataBar, model.DataBar{Date: time.Unix(DailyItem.DateTime, 0).Format("2006-01-02"), Count: DailyItem.Count})
	}

	todayStr := nowTime.Format("2006-01-02")
	DailyDataBar = append(DailyDataBar, model.DataBar{Date: todayStr, Count: todayCount})
	Daily := checkQueryData(DailyDataBar, start, end)
	return Daily
}

func (s *dashboardService) SingleMessageDaily(c *gin.Context) {
	req := new(model.GetDailyDataReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	sigleMsgCount, err := repo.DashboardRepo.GetTodayMessageCount(1)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}

	sigleMsgDaily := QueryDailyData(sigleMsgCount, 1, req.BeginDate, req.EndDate)

	http.Success(c, sigleMsgDaily)
	return
}

func (s *dashboardService) GroupMessageDaily(c *gin.Context) {
	req := new(model.GetDailyDataReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	groupMsgCount, err := repo.DashboardRepo.GetTodayMessageCount(2)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}

	groupMsgDaily := QueryDailyData(groupMsgCount, 2, req.BeginDate, req.EndDate)

	http.Success(c, groupMsgDaily)
	return
}

func (s *dashboardService) OnlineUserDaily(c *gin.Context) {
	req := new(model.GetDailyDataReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	onlineNum := repo.DashboardCache.GetOnlineMax()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}

	onlineDaily := QueryDailyData(onlineNum, 3, req.BeginDate, req.EndDate)

	http.Success(c, onlineDaily)
	return
}

func checkQueryData(data model.DataBarList, start, end int64) model.DataBarList {
	startTime := time.Unix(start, 0)
	endTime := time.Unix(end, 0)
	endStr := endTime.Format("2006-01-02")
	for {
		day := startTime.Format("2006-01-02")
		var exist = 0
		for _, check := range data {
			if check.Date == day {
				exist = 1
			}
		}
		if exist == 0 {
			data = append(data, model.DataBar{Date: day, Count: 0})
		}
		startTime = startTime.AddDate(0, 0, 1)
		if day == endStr {
			break
		}
	}

	sort.Sort(data)
	return data
}
