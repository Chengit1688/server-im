package job

import (
	"im/internal/api/user/model"
	"im/internal/cms_api/user/repo"
	"im/pkg/logger"
	"im/pkg/util"
)

func Init() {

	DisabledManagermentDataMigration()

	repo.UserRepo.DMDeviceSyncCache()
}

func DisabledManagermentDataMigration() {
	countUser, err := repo.UserRepo.CountUser()
	if err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		return
	}
	countUserDevice, err := repo.UserRepo.CountUserDevice()
	if err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		return
	}
	countUserIP, err := repo.UserRepo.CountUserIP()
	if err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		return
	}
	if countUser/2 > countUserDevice && countUser/2 > countUserIP {
		users, err := repo.UserRepo.GetUserDeviceIP()
		if err != nil {
			logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
			return
		}
		var addUserDevice []model.UserDevice
		var addUserIP []model.UserIp
		for _, user := range users {
			if len(user.DeviceId) != 0 && user.Platform != 0 {
				addUserDevice = append(addUserDevice, model.UserDevice{UserID: user.UserID, DeviceID: user.DeviceId, Platform: user.Platform, Status: 1})
			}
			if len(user.LoginIp) != 0 {
				addUserIP = append(addUserIP, model.UserIp{UserID: user.UserID, Ip: user.LoginIp, Status: 1})
			} else if len(user.RegisterIp) != 0 {
				addUserIP = append(addUserIP, model.UserIp{UserID: user.UserID, Ip: user.RegisterIp, Status: 1})
			}
		}
		err = repo.UserRepo.AddUserDeviceBatch(addUserDevice)
		if err != nil {
			logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
			return
		}
		err = repo.UserRepo.AddUserIPBatch(addUserIP)
		if err != nil {
			logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
			return
		}
	}
}

func ClearUserLoginHistory() {
	keepDays := 30
	err := repo.UserRepo.LoginHistoryClear(keepDays)
	if err != nil {
		logger.Sugar.Errorw("ClearUserLoginHistory", "func", util.GetSelfFuncName(), "error", err)
	}
}
