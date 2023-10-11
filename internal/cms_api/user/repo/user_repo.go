package repo

import (
	"fmt"
	friendUseCase "im/internal/api/friend/usecase"
	apiUserModel "im/internal/api/user/model"
	apiUserUseCase "im/internal/api/user/usecase"
	ipRepo "im/internal/cms_api/ipblacklist/repo"
	"im/internal/cms_api/user/model"
	cmsapiUserModel "im/internal/cms_api/user/model"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/util"
	"strings"
	"time"

	"gorm.io/gorm"
)

var UserRepo = new(userRepo)

type userRepo struct{}

func (r *userRepo) UserPaging(req cmsapiUserModel.UserListReq) (users []apiUserModel.User, count int64, err error) {
	req.Pagination.Check()
	tx := db.DB.Model(apiUserModel.User{})
	if len(req.NickName) != 0 {
		tx = tx.Where(fmt.Sprintf("nick_name like %q", ("%" + req.NickName + "%")))
	}
	if len(req.UserID) != 0 {
		tx = tx.Where(fmt.Sprintf("user_id like %q", ("%" + req.UserID + "%")))
	}
	if len(req.Account) != 0 {
		tx = tx.Where(fmt.Sprintf("account like %q", ("%" + req.Account + "%")))
	}
	if len(req.PhoneNumber) != 0 {
		tx = tx.Where("phone_number = ?", req.PhoneNumber)
	}
	if req.RegisterTimeStart != 0 {
		tx = tx.Where("created_at >= ?", req.RegisterTimeStart)
	}
	if req.RegisterTimeEnd != 0 {
		tx = tx.Where("created_at <= ?", req.RegisterTimeEnd)
	}
	if req.Status != 0 {
		tx = tx.Where("status = ?", req.Status)
	}
	if len(req.LoginIp) != 0 {
		tx = tx.Where("login_ip = ?", req.LoginIp)
	}
	if req.Gender != 0 {
		tx = tx.Where("gender = ?", req.Gender)
	}
	if req.InviteCode != nil {
		tx = tx.Where("invite_code = ?", req.InviteCode)
	}
	if req.IsPrivilege != 0 {
		tx = tx.Where("is_privilege = ?", req.IsPrivilege)
	}

	if req.IsCustomer != 0 {
		tx = tx.Where("is_customer = ?", req.IsCustomer)
	}
	err = tx.Offset(req.Offset).Limit(req.Limit).Order("created_at desc").Find(&users).Limit(-1).Offset(-1).Count(&count).Error
	return
}

func (r *userRepo) RealNameListPaging(req cmsapiUserModel.RealNameListReq) (users []apiUserModel.User, count int64, err error) {
	req.Pagination.Check()
	tx := db.DB.Model(apiUserModel.User{})
	if len(req.UserID) != 0 {
		tx = tx.Where(fmt.Sprintf("user_id like %q", "%"+req.UserID+"%"))
	}
	if req.IsRealAuth != 0 {
		tx = tx.Where("is_real_auth = ?", req.IsRealAuth)
	} else {
		tx = tx.Where("is_real_auth >= 2")
	}
	err = tx.Offset(req.Offset).Limit(req.Limit).Order("created_at desc").Find(&users).Limit(-1).Offset(-1).Count(&count).Error
	return
}

func (r *userRepo) UserGroupPaging(req cmsapiUserModel.UserListReq) (users []apiUserModel.User, count int64, err error) {
	req.Pagination.Check()
	tx := db.DB.Model(apiUserModel.User{})

	if req.NickName != "" && req.PhoneNumber != "" && req.UserID != "" {
		tx = tx.Where(fmt.Sprintf("nick_name like %q OR user_id like %q OR phone_number like %q", "%"+req.NickName+"%", "%"+req.UserID+"%", "%"+req.PhoneNumber+"%"))
	}

	if req.Status != 0 {
		tx = tx.Where("status = ?", req.Status)
	}

	if err = tx.Order("created_at desc").Offset(req.Offset).Limit(req.Limit).Find(&users).Error; err != nil {
		return
	}

	err = tx.Offset(-1).Limit(-1).Count(&count).Error
	logger.Sugar.Debugw(fmt.Sprintf("%v", req.Pagination))
	return
}

func (r *userRepo) UserBatchAdd(req cmsapiUserModel.UserAddBatchReq) (users []apiUserModel.User, err error) {
	db.DB.Transaction(func(tx *gorm.DB) error {

		for _, user := range req.Users {
			add := new(apiUserModel.User)
			util.CopyStructFields(&add, &user)
			add.UserID = util.RandID(10)
			add.Salt = util.RandString(6)
			add.Password = util.GetPassword(user.Password, add.Salt)
			add.LatestLoginTime = time.Now().Unix()
			if err = tx.Create(&add).Error; err != nil {
				return err
			}
			err = apiUserUseCase.AuthUseCase.CreateAuthUsername(add.UserID)
			if err != nil {
				return err
			}

			friendUseCase.FriendUseCase.CreateFriendLabel(add.UserID, add.UserID, "我的好友")
			users = append(users, *add)
		}
		return nil
	})
	return
}

func (r *userRepo) UserInfoUpdate(req cmsapiUserModel.UpdateUserInfoReq) (err error) {
	updates, err := util.StructToMap(req, "json")
	delete(updates, "operation_id")
	delete(updates, "user_id")
	if err = db.DB.Model(&apiUserModel.User{}).Where("user_id = ?", req.UserID).Updates(updates).Error; err != nil {
		return err
	}
	return err
}

func (r *userRepo) FreezeUser(userID string) (err error) {
	if err = db.DB.Model(&apiUserModel.User{}).Where("user_id = ?", userID).Update("status", 2).Error; err != nil {
		return err
	}
	return err
}

func (r *userRepo) UnFreezeUser(userID string) (err error) {
	if err = db.DB.Model(&apiUserModel.User{}).Where("user_id = ?", userID).Update("status", 1).Error; err != nil {
		return err
	}
	return err
}

func (r *userRepo) SetUserPassword(userID string, password string) (err error) {
	user := new(apiUserModel.User)
	if err = db.DB.Model(&apiUserModel.User{}).Where("user_id = ?", userID).First(&user).Error; err != nil {
		return err
	}
	passwordHash := util.GetPassword(password, user.Salt)
	if err = db.DB.Model(&apiUserModel.User{}).Where("user_id = ?", userID).Update("password", passwordHash).Error; err != nil {
		return err
	}
	return err
}

func (r *userRepo) DisabledManagermentUser(req model.DMUserListReq) (users []model.DMUserListItemResp, count int64, err error) {
	req.Pagination.Check()
	var tx, countTx *gorm.DB

	if len(req.Search) == 0 {
		tx = db.DB.Table("users AS t1").
			Joins("JOIN user_devices AS t2 ON t1.user_id = t2.user_id").
			Joins("JOIN user_ips AS t3 ON t1.user_id = t3.user_id").
			Joins("JOIN (SELECT user_id FROM users where status = ? ORDER BY created_at DESC  LIMIT ? OFFSET ?) as o on o.user_id = t1.user_id ", req.Status, req.Pagination.Limit, req.Pagination.Offset).Group("t1.id,t1.user_id,t1.nick_name,t1.account,t1.phone_number")
		countTx = db.DB.Table("users AS t1").
			Joins("JOIN user_devices AS t2 ON t1.user_id = t2.user_id").
			Joins("JOIN user_ips AS t3 ON t1.user_id = t3.user_id").
			Joins("JOIN (SELECT user_id FROM users where status = ? ORDER BY created_at DESC) as o on o.user_id = t1.user_id ", req.Status).Group("t1.id,t1.user_id,t1.nick_name,t1.account,t1.phone_number")
	} else {
		tx = db.DB.Table("users AS t1").
			Joins("JOIN user_devices AS t2 ON t1.user_id = t2.user_id").
			Joins("JOIN user_ips AS t3 ON t1.user_id = t3.user_id ").
			Joins(fmt.Sprintf("JOIN (SELECT user_id FROM users where status = %d and (user_id is not null or user_id like %q or account like %q or nick_name like %q or phone_number like %q ) ORDER BY created_at DESC) as o on o.user_id = t1.user_id ",
				req.Status, ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"))).Where(fmt.Sprintf("t1.user_id like %q or t1.account like %q or t1.nick_name like %q or t1.phone_number like %q or t2.device_id like %q or t3.ip like %q", ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"))).Group("t1.id,t1.user_id,t1.nick_name,t1.account,t1.phone_number")
		countTx = tx
	}
	tx.Select("t1.id,t1.user_id,t1.nick_name,t1.account,t1.phone_number,COUNT(DISTINCT t2.device_id) as count_device,COUNT(DISTINCT t3.ip) as count_ip").Order("t1.id DESC")
	err = tx.Find(&users).Error
	if err != nil {
		return
	}
	countTx = countTx.Select("COUNT(DISTINCT t1.id)")
	err = countTx.Count(&count).Error
	return
}

func (r *userRepo) DisabledManagermentDevice(req model.DMDeviceListReq) (devices []model.DMDeviceListItemResp, count int64, err error) {
	req.Pagination.Check()
	var tx, countTx *gorm.DB

	if len(req.Search) == 0 {
		tx = db.DB.Table("users AS t1").
			Joins("JOIN user_devices AS t2 ON t1.user_id = t2.user_id").
			Joins("JOIN user_ips AS t3 ON t1.user_id = t3.user_id").
			Joins("JOIN (SELECT DISTINCT device_id FROM user_devices where status = ?  LIMIT ? OFFSET ?) as o on o.device_id = t2.device_id ", req.Status, req.Pagination.Limit, req.Pagination.Offset).Group("t2.device_id,t2.platform")
		countTx = db.DB.Table("users AS t1").
			Joins("JOIN user_devices AS t2 ON t1.user_id = t2.user_id").
			Joins("JOIN user_ips AS t3 ON t1.user_id = t3.user_id").
			Joins("JOIN (SELECT DISTINCT device_id FROM user_devices where status = ? ) as o on o.device_id = t2.device_id ", req.Status).Group("t2.device_id,t2.platform")
	} else {
		var baseCount int64
		err = db.DB.Table("user_devices").Where("status = ?", req.Status).Where("device_id like ?", req.Search).Select("DISTINCT device_id").Count(&baseCount).Error
		if err != nil {
			logger.Sugar.Errorw("db error", "func", util.GetSelfFuncName(), "error", err)
		}
		if baseCount == 0 {
			tx = db.DB.Table("users AS t1").
				Joins("JOIN user_devices AS t2 ON t1.user_id = t2.user_id").
				Joins("JOIN user_ips AS t3 ON t1.user_id = t3.user_id").
				Where("t2.status = ?", req.Status).
				Where(fmt.Sprintf("t1.user_id like %q or t1.account like %q or t1.nick_name like %q or t1.phone_number like %q or t3.ip like %q ", ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"))).Limit(req.Pagination.Limit).Offset(req.Pagination.Offset).Group("t2.device_id,t2.platform")
			countTx = db.DB.Table("users AS t1").
				Joins("JOIN user_devices AS t2 ON t1.user_id = t2.user_id").
				Joins("JOIN user_ips AS t3 ON t1.user_id = t3.user_id").
				Where("t2.status = ?", req.Status).
				Where(fmt.Sprintf("t1.user_id like %q or t1.account like %q or t1.nick_name like %q or t1.phone_number like %q or t3.ip like %q ", ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"))).Group("t2.device_id,t2.platform")
		} else {
			var deviceIDlist []string
			err = db.DB.Table("user_devices").Where("status = ?", req.Status).Where("device_id like ?", req.Search).Select("DISTINCT device_id").Find(&deviceIDlist).Error
			if err != nil {
				logger.Sugar.Errorw("db error", "func", util.GetSelfFuncName(), "error", err)
			}
			tx = db.DB.Table("users AS t1").
				Joins("JOIN user_devices AS t2 ON t1.user_id = t2.user_id").
				Joins("JOIN user_ips AS t3 ON t1.user_id = t3.user_id").
				Where("t2.status = ? and t2.device_id in ?", req.Status, deviceIDlist).
				Or(fmt.Sprintf("t1.user_id like %q or t1.account like %q or t1.nick_name like %q or t1.phone_number like %q or t3.ip like %q ", ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"))).Limit(req.Pagination.Limit).Offset(req.Pagination.Offset).Group("t2.device_id,t2.platform")
			countTx = db.DB.Table("users AS t1").
				Joins("JOIN user_devices AS t2 ON t1.user_id = t2.user_id").
				Joins("JOIN user_ips AS t3 ON t1.user_id = t3.user_id").
				Where("t2.status = ? and t2.device_id in ?", req.Status, deviceIDlist).
				Or(fmt.Sprintf("t1.user_id like %q or t1.account like %q or t1.nick_name like %q or t1.phone_number like %q or t3.ip like %q ", ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"))).Group("t2.device_id,t2.platform")
		}

	}
	tx.Select("max(t2.id) as id,t2.device_id,t2.platform,COUNT(DISTINCT t1.user_id) as count_user,COUNT(DISTINCT t3.ip) as count_ip").Order("max(t2.id) DESC")
	err = tx.Find(&devices).Error
	if err != nil {
		return
	}
	countTx = countTx.Select("COUNT(DISTINCT t2.device_id)")
	err = countTx.Count(&count).Error
	return
}

func (r *userRepo) DisabledManagermentIP(req model.DMIPListReq) (ips []model.DMIPListItemResp, count int64, err error) {
	req.Pagination.Check()
	var tx, countTx *gorm.DB

	if len(req.Search) == 0 {
		tx = db.DB.Table("users AS t1").
			Joins("JOIN user_devices AS t2 ON t1.user_id = t2.user_id").
			Joins("JOIN user_ips AS t3 ON t1.user_id = t3.user_id").
			Joins("inner join (SELECT DISTINCT ip FROM user_ips where status = ?  LIMIT ? OFFSET ?) as o on o.ip = t3.ip ", req.Status, req.Pagination.Limit, req.Pagination.Offset).Group("t3.ip")
		countTx = db.DB.Table("user_ips").Where("status = ?", req.Status)
	} else {
		var baseCount int64
		err = db.DB.Table("user_ips").Where("status = ?", req.Status).Where("ip like ?", req.Search).Select("DISTINCT ip").Count(&baseCount).Error
		if err != nil {
			logger.Sugar.Errorw("db error", "func", util.GetSelfFuncName(), "error", err)
		}
		if baseCount == 0 {
			tx = db.DB.Table("users AS t1").
				Joins("JOIN user_devices AS t2 ON t1.user_id = t2.user_id").
				Joins("JOIN user_ips AS t3 ON t1.user_id = t3.user_id").
				Where("t3.status = ?", req.Status).
				Where(fmt.Sprintf("t1.user_id like %q or t1.account like %q or t1.nick_name like %q or t1.phone_number like %q or t2.device_id like %q ", ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"))).Limit(req.Pagination.Limit).Offset(req.Pagination.Offset).Group("t3.ip")
			countTx = db.DB.Table("users AS t1").
				Joins("JOIN user_devices AS t2 ON t1.user_id = t2.user_id").
				Joins("JOIN user_ips AS t3 ON t1.user_id = t3.user_id").
				Where("t3.status = ?", req.Status).
				Where(fmt.Sprintf("t1.user_id like %q or t1.account like %q or t1.nick_name like %q or t1.phone_number like %q or t2.device_id like %q ", ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"))).Group("t3.ip")
		} else {
			var iplist []string
			err = db.DB.Table("user_ips").Where("status = ?", req.Status).Where("ip like ?", req.Search).Select("DISTINCT ip").Find(&iplist).Error
			if err != nil {
				logger.Sugar.Errorw("db error", "func", util.GetSelfFuncName(), "error", err)
			}
			tx = db.DB.Table("users AS t1").
				Joins("JOIN user_devices AS t2 ON t1.user_id = t2.user_id").
				Joins("JOIN user_ips AS t3 ON t1.user_id = t3.user_id").
				Where("t3.status = ? and t3.ip in ?", req.Status, iplist).
				Or(fmt.Sprintf("t1.user_id like %q or t1.account like %q or t1.nick_name like %q or t1.phone_number like %q or t2.device_id like %q ", ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"))).Limit(req.Pagination.Limit).Offset(req.Pagination.Offset).Group("t3.ip")
			countTx = db.DB.Table("users AS t1").
				Joins("JOIN user_devices AS t2 ON t1.user_id = t2.user_id").
				Joins("JOIN user_ips AS t3 ON t1.user_id = t3.user_id").
				Where("t3.status = ? and t3.ip in ?", req.Status, iplist).
				Or(fmt.Sprintf("t1.user_id like %q or t1.account like %q or t1.nick_name like %q or t1.phone_number like %q or t2.device_id like %q ", ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"), ("%" + req.Search + "%"))).Group("t3.ip")
		}
	}

	tx.Select("max(t3.id) as id,t3.ip,COUNT(DISTINCT t1.user_id) as count_user,COUNT(DISTINCT t2.device_id) as count_device").Order("max(t3.id) DESC")

	err = tx.Find(&ips).Error
	if err != nil {
		return
	}
	countTx = countTx.Select("COUNT(DISTINCT ip)")
	err = countTx.Count(&count).Error

	if req.Status == 2 {

		var ipList []string
		err = db.DB.Table("user_ips").Select("DISTINCT ip").Find(&ipList).Error

		var otherList []string
		err = db.DB.Table("cms_ipblacklist").Select("DISTINCT ip").Where("ip not in ?", ipList).Find(&otherList).Error

		var pool []model.DMIPListItemResp
		for _, ip := range otherList {
			if len(req.Search) != 0 {
				if strings.Contains(ip, req.Search) {
					pool = append(pool, model.DMIPListItemResp{ID: util.RandInt(999999, 9999999), IP: ip})
				}
			} else {
				pool = append(pool, model.DMIPListItemResp{ID: util.RandInt(999999, 9999999), IP: ip})
			}
		}

		ips, count = r.splitArrayData(req, count, pool, ips)
	}
	return
}

func (r *userRepo) splitArrayData(req model.DMIPListReq, count int64, pool []model.DMIPListItemResp, ips []model.DMIPListItemResp) (arr []model.DMIPListItemResp, arrCount int64) {
	poolSize := len(pool)

	if (req.Page * req.PageSize) > int(count) {
		needAddNum := (req.Page * req.PageSize) - int(count)

		if needAddNum > req.PageSize {

			sub := needAddNum / req.PageSize
			yu := needAddNum % req.PageSize

			if poolSize >= yu {
				result := splitArrays(pool[yu:], int64(req.PageSize))
				if len(result) >= sub {
					ips = append(ips, result[sub-1]...)
				}
			}
		} else {

			if needAddNum <= poolSize {
				ips = append(ips, pool[:needAddNum]...)
			} else {
				ips = append(ips, pool...)
			}
		}
	}
	arr = ips
	arrCount = count + int64(poolSize)
	return
}

func splitArrays(arr []model.DMIPListItemResp, num int64) [][]model.DMIPListItemResp {
	max := int64(len(arr))
	if max <= num {
		return [][]model.DMIPListItemResp{arr}
	}
	var quantity int64
	if max%num == 0 {
		quantity = max / num
	} else {
		quantity = (max / num) + 1
	}
	var segments = make([][]model.DMIPListItemResp, 0)
	var start, end, i int64
	for i = 1; i <= quantity; i++ {
		end = i * num
		if i != quantity {
			segments = append(segments, arr[start:end])
		} else {
			segments = append(segments, arr[start:])
		}
		start = i * num
	}
	return segments
}

func (r *userRepo) CountUser() (count int64, err error) {
	err = db.DB.Model(apiUserModel.User{}).Count(&count).Error
	return
}

func (r *userRepo) CountUserDevice() (count int64, err error) {
	err = db.DB.Model(apiUserModel.UserDevice{}).Count(&count).Error
	return
}

func (r *userRepo) CountUserIP() (count int64, err error) {
	err = db.DB.Model(apiUserModel.UserIp{}).Count(&count).Error
	return
}

func (r *userRepo) AddUserDeviceBatch(adds []apiUserModel.UserDevice) (err error) {
	err = db.DB.CreateInBatches(&adds, 500).Error
	return
}

func (r *userRepo) AddUserIPBatch(adds []apiUserModel.UserIp) (err error) {
	err = db.DB.CreateInBatches(&adds, 500).Error
	return
}

func (r *userRepo) GetUserDeviceIP() (users []apiUserModel.User, err error) {
	err = db.DB.Model(apiUserModel.User{}).Select("user_id", "platform", "device_id", "login_ip").Find(&users).Error
	return
}

func (r *userRepo) DMDeviceSyncCache() {

	dbDeviceInfos := new([]apiUserModel.UserDevice)
	tx := db.DB.Model(apiUserModel.UserDevice{})
	err := tx.Where("status = ?", 2).Find(&dbDeviceInfos).Error
	if err != nil {
		logger.Sugar.Errorw("db error", "func", util.GetSelfFuncName(), "error", err)
		return
	}
	cacheDeviceInfos := DeviceListCache.GetAllInfo()

	for _, item := range *dbDeviceInfos {
		if _, ok := cacheDeviceInfos[item.DeviceID]; !ok {
			err = DeviceListCache.Add(item.DeviceID)
			if err != nil {
				logger.Sugar.Errorw("redis error", "func", util.GetSelfFuncName(), "error", err)
				return
			}
		}
	}

	for k := range cacheDeviceInfos {
		var status int = 0
		for _, item := range *dbDeviceInfos {
			if item.DeviceID == k {
				status = 1
			}
		}
		if status != 1 {
			err = DeviceListCache.Del(k)
			if err != nil {
				logger.Sugar.Errorw("redis error", "func", util.GetSelfFuncName(), "error", err)
				return
			}
		}
	}
}

func (r *userRepo) DMDeviceLock(deviceIDs []string) (err error) {
	err = db.DB.Model(apiUserModel.UserDevice{}).Where("device_id IN ?", deviceIDs).UpdateColumn("status", 2).Error
	for _, device_id := range deviceIDs {
		err = DeviceListCache.Add(device_id)
		if err != nil {
			logger.Sugar.Errorw("redis error", "func", util.GetSelfFuncName(), "error", err)
			return
		}
	}
	return
}

func (r *userRepo) DMDeviceUnLock(deviceIDs []string) (err error) {
	err = db.DB.Model(apiUserModel.UserDevice{}).Where("device_id IN ?", deviceIDs).UpdateColumn("status", 1).Error
	for _, device_id := range deviceIDs {
		err = DeviceListCache.Del(device_id)
		if err != nil {
			logger.Sugar.Errorw("redis error", "func", util.GetSelfFuncName(), "error", err)
			return
		}
	}
	return
}

func (r *userRepo) DMIPLock(ips []string) (err error) {
	err = db.DB.Model(apiUserModel.UserIp{}).Where("ip IN ?", ips).UpdateColumn("status", 2).Error
	for _, ip := range ips {
		item, err := ipRepo.IPBlackListRepo.GetByIP(ip)
		if err != nil {
			item.IP = ip
			item.Note = "封禁管理-IP封禁"
			_, err = ipRepo.IPBlackListRepo.Add(item)
			if err != nil {
				logger.Sugar.Errorw("IPBlackListRepo Add error", "func", util.GetSelfFuncName(), "error", err)
				return err
			}
		}
	}
	return
}

func (r *userRepo) DMIPUnLock(ips []string) (err error) {
	err = db.DB.Model(apiUserModel.UserIp{}).Where("ip IN ?", ips).UpdateColumn("status", 1).Error
	for _, ip := range ips {
		item, err := ipRepo.IPBlackListRepo.GetByIP(ip)
		if err != nil {
			continue
		}
		err = ipRepo.IPBlackListRepo.Delete(fmt.Sprintf("%d", item.ID))
		if err != nil {
			logger.Sugar.Errorw("IPBlackListRepo Delete error", "func", util.GetSelfFuncName(), "error", err)
			return err
		}
	}
	return
}

func (r *userRepo) UserExport(req cmsapiUserModel.UserListReq) (users []apiUserModel.User, err error) {
	tx := db.DB.Model(apiUserModel.User{})
	if len(req.NickName) != 0 {
		tx = tx.Where(fmt.Sprintf("nick_name like %q", ("%" + req.NickName + "%")))
	}
	if len(req.UserID) != 0 {
		tx = tx.Where(fmt.Sprintf("user_id like %q", ("%" + req.UserID + "%")))
	}
	if len(req.Account) != 0 {
		tx = tx.Where(fmt.Sprintf("account like %q", ("%" + req.Account + "%")))
	}
	if len(req.PhoneNumber) != 0 {
		tx = tx.Where("phone_number = ?", req.PhoneNumber)
	}
	if req.RegisterTimeStart != 0 {
		tx = tx.Where("created_at >= ?", req.RegisterTimeStart)
	}
	if req.RegisterTimeEnd != 0 {
		tx = tx.Where("created_at <= ?", req.RegisterTimeEnd)
	}
	if req.Status != 0 {
		tx = tx.Where("status = ?", req.Status)
	}
	if len(req.LoginIp) != 0 {
		tx = tx.Where("login_ip = ?", req.LoginIp)
	}
	if req.Gender != 0 {
		tx = tx.Where("gender = ?", req.Gender)
	}
	if req.InviteCode != nil {
		tx = tx.Where("invite_code = ?", req.InviteCode)
	}
	if req.IsPrivilege != 0 {
		tx = tx.Where("is_privilege = ?", req.IsPrivilege)
	}
	err = tx.Order("created_at desc").Find(&users).Error
	return
}

func (r *userRepo) UserSearch(req cmsapiUserModel.UserSearchReq) (users []apiUserModel.User, count int64, err error) {
	tx := db.DB.Model(apiUserModel.User{})

	tx = tx.Where(fmt.Sprintf("nick_name like %q", ("%" + req.Search + "%"))).
		Or(fmt.Sprintf("user_id like %q", ("%" + req.Search + "%"))).
		Or(fmt.Sprintf("account like %q", ("%" + req.Search + "%"))).
		Or(fmt.Sprintf("phone_number like %q", ("%" + req.Search + "%")))

	err = tx.Order("created_at desc").Find(&users).Count(&count).Error
	return
}

func (r *userRepo) GetOnlineUsers(keyword string, offset int, limit int) (list []apiUserModel.User, count int64, err error) {
	needPaging := !(offset == limit && limit == 0)
	var query string
	var args []interface{}

	query = "online_status = ?"
	args = append(args, apiUserModel.OnlineStatusTypeOnline)

	if keyword != "" {
		query += " AND (user_id LIKE ? OR nick_name LIKE ?)"
		args = append(args, fmt.Sprintf("%%%s%%", keyword), fmt.Sprintf("%%%s%%", keyword))
	}

	var listDB *gorm.DB
	listDB = db.DB.Model(&apiUserModel.User{})

	if needPaging {
		if err = listDB.Where(query, args...).Count(&count).Error; err != nil {
			return
		}
	}

	if needPaging {
		listDB = listDB.Offset(offset).Limit(limit)
	}

	if err = listDB.Where(query, args...).Find(&list).Error; err != nil {
		return
	}

	return
}

func (r *userRepo) GetOnlineUsersCount() (count int64, err error) {
	err = db.DB.Model(&apiUserModel.User{}).Where("online_status = ?", apiUserModel.OnlineStatusTypeOnline).Count(&count).Error
	return
}

func (r *userRepo) GetOneDayOnlineUsers() (users []apiUserModel.User, err error) {
	err = db.DB.Model(&apiUserModel.User{}).Where("online_status = ? AND latest_login_time < ?", apiUserModel.OnlineStatusTypeOnline, time.Now().Unix()-24*60*60).Find(&users).Error
	return
}

func (r *userRepo) UpdateOnlineStatus(userID string, status apiUserModel.OnlineStatusType, latestLoginTime int64) (err error) {
	m := map[string]interface{}{
		"online_status": status,
	}

	if latestLoginTime != 0 {
		m["latest_login_time"] = latestLoginTime
	}

	err = db.DB.Model(&apiUserModel.User{}).Where("user_id = ?", userID).Updates(m).Error
	return
}

func (r *userRepo) UpdateRealName(req model.RealNameAuthReq) (err error) {
	m := apiUserModel.User{
		RealAuthMsg: req.RealAuthMsg,
		IsRealAuth:  req.IsRealAuth,
	}
	err = db.DB.Model(&apiUserModel.User{}).Where("user_id = ?", req.UserID).Updates(&m).Error
	return
}

func (r *userRepo) SetPrivilegeUserStatus(status int) (users []apiUserModel.User, err error) {

	if err = db.DB.Model(&apiUserModel.User{}).Where("is_privilege = ?", 1).Find(&users).Update("status", status).Error; err != nil {
		return
	}
	return
}

func (r *userRepo) LoginHistoryPaging(req cmsapiUserModel.LoginHistoryReq) (records []apiUserModel.LoginHistory, count int64, err error) {
	req.Pagination.Check()
	tx := db.DB.Model(apiUserModel.LoginHistory{})

	tx = tx.Where("user_id = ?", req.UserID)

	err = tx.Offset(req.Offset).Limit(req.Limit).Order("created_at desc").Find(&records).Limit(-1).Offset(-1).Count(&count).Error
	return
}

func (r *userRepo) LoginHistoryClear(days int) (err error) {
	now := time.Now()
	now = now.Add(time.Duration(-1*days*24) * time.Hour)
	tx := db.DB.Model(apiUserModel.LoginHistory{})

	tx = tx.Where("created_at <= ?", now.Unix())

	err = tx.Delete(&apiUserModel.LoginHistory{}).Error
	return
}
