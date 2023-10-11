package job

import (
	"fmt"
	apiUserUseCase "im/internal/api/user/usecase"
	"im/internal/cms_api/dashboard/model"
	"im/internal/cms_api/dashboard/repo"
	"im/pkg/logger"
	"im/pkg/util"
	"time"
)

func Init() {
	go UpdateOnlineMaxJob()
	go supplementaryData()

}

func ResetUserOnlineMax() {
	repo.DashboardCache.SetOnlineMax(0)
}

func UpdateOnlineMaxJob() {
	ticker := time.NewTicker(time.Minute * 5)

	for {
		select {
		case <-ticker.C:
			UpdateOnlineMax()
		}
	}
}

func UpdateOnlineMax() {
	onlineUserMap, err := apiUserUseCase.UserUseCase.OnlineMap()
	if err != nil {
		logger.Sugar.Errorw("获取在线用户峰值数据失败", "func", util.GetSelfFuncName(), "error", err.Error())
	}
	total := len(onlineUserMap)
	have := repo.DashboardCache.GetOnlineMax()
	if total > int(have) {
		repo.DashboardCache.SetOnlineMax(int64(total))
	}
}

func InsertAnalysisData() {
	now := time.Now()
	yestday := now.AddDate(0, 0, -1)

	singleCount, err := repo.DashboardRepo.GetYestdayMessageCount(1, now)
	if err != nil {
		logger.Sugar.Errorw("GetYestdayMessageCount", "error", err)
	}
	singleAdd := &model.DashboardDailyData{DateTime: yestday.Unix(), Count: singleCount, Type: 1}
	err = repo.DashboardRepo.AddDashboardDailyData(*singleAdd)
	if err != nil {
		logger.Sugar.Errorw("AddDashboardDailyData", "error", err)
	}

	GroupCount, err := repo.DashboardRepo.GetYestdayMessageCount(2, now)
	if err != nil {
		logger.Sugar.Errorw("GetYestdayMessageCount", "error", err)
	}
	groupAdd := &model.DashboardDailyData{DateTime: yestday.Unix(), Count: GroupCount, Type: 2}
	err = repo.DashboardRepo.AddDashboardDailyData(*groupAdd)
	if err != nil {
		logger.Sugar.Errorw("AddDashboardDailyData", "error", err)
	}

	onlineNum := repo.DashboardCache.GetOnlineMax()
	onlineAdd := &model.DashboardDailyData{DateTime: yestday.Unix(), Count: onlineNum, Type: 3}
	err = repo.DashboardRepo.AddDashboardDailyData(*onlineAdd)
	if err != nil {
		logger.Sugar.Errorw("AddDashboardDailyData", "error", err)
	}

	registNum, err := repo.DashboardRepo.GetRegisterNumYesterday()
	if err != nil {
		logger.Sugar.Errorw("GetRegisterNumYesterday", "error", err)
	}
	registAdd := &model.DashboardDailyData{DateTime: yestday.Unix(), Count: registNum, Type: 4}
	err = repo.DashboardRepo.AddDashboardDailyData(*registAdd)
	if err != nil {
		logger.Sugar.Errorw("AddDashboardDailyData", "error", err)
	}
}

func supplementaryData() {
	var LOC, _ = time.LoadLocation("Asia/Shanghai")
	now := time.Now()
	currentYear := now.Year()
	currentMonth := now.Month()

	endOfMonth := time.Date(currentYear, currentMonth, 1, 23, 59, 59, 0, LOC)
	lastOfMonth := endOfMonth.AddDate(0, 1, -1)
	singleDatas, err := repo.DashboardRepo.FindDashboardDailyData(0, lastOfMonth.Unix(), 1)
	if err != nil {
		logger.Sugar.Errorw("FindDashboardDailyData", "error", err)
		return
	}
	if len(singleDatas) == 0 {
		singles, err := repo.DashboardRepo.GetMonthDailyMessage(1)
		if err != nil {
			logger.Sugar.Errorw("GetMonthDailyMessage", "error", err)
			return
		}
		for _, messageItem := range singles {
			exist := 0
			for _, dashboardItem := range singleDatas {
				if messageItem.Date == time.Unix(dashboardItem.DateTime, 0).Format("2006-01-02") {
					exist = 1
				}
			}
			if exist == 0 {
				timeStr := fmt.Sprintf("%s 00:00:01", messageItem.Date)
				date, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, LOC)
				singleAdd := &model.DashboardDailyData{DateTime: date.Unix(), Count: messageItem.Count, Type: 1}
				err = repo.DashboardRepo.AddDashboardDailyData(*singleAdd)
				if err != nil {
					logger.Sugar.Errorw("AddDashboardDailyData", "error", err)
				}
			}
		}
	}
	groupDatas, err := repo.DashboardRepo.FindDashboardDailyData(0, lastOfMonth.Unix(), 2)
	if err != nil {
		logger.Sugar.Errorw("FindDashboardDailyData", "error", err)
		return
	}
	if len(groupDatas) == 0 {

		groups, err := repo.DashboardRepo.GetMonthDailyMessage(2)
		if err != nil {
			logger.Sugar.Errorw("GetMonthDailyMessage", "error", err)
			return
		}
		for _, messageItem := range groups {
			exist := 0
			for _, dashboardItem := range groupDatas {
				if messageItem.Date == time.Unix(dashboardItem.DateTime, 0).Format("2006-01-02") {
					exist = 1
				}
			}
			if exist == 0 {
				timeStr := fmt.Sprintf("%s 00:00:01", messageItem.Date)
				date, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, LOC)
				groupAdd := &model.DashboardDailyData{DateTime: date.Unix(), Count: messageItem.Count, Type: 2}
				err = repo.DashboardRepo.AddDashboardDailyData(*groupAdd)
				if err != nil {
					logger.Sugar.Errorw("AddDashboardDailyData", "error", err)
				}
			}
		}
	}

}
