package repo

import (
	"fmt"
	"im/internal/api/user/model"
	"im/pkg/code"
	"im/pkg/db"
	"im/pkg/util"
	"time"

	"gorm.io/gorm"
)

var UserRepo = new(userRepo)

type userRepo struct{}

type WhereOption struct {
	Id               int64
	UserId           string
	DeviceId         string
	RegisterDeviceId string
	RegisterIp       string
	Account          string
	NickName         string
	Password         string
	PhoneNumber      string
	LoginType        int64
	Platform         int64
	Status           int64
	CountryCode      string
	IP               string
	UserModel        int64
	IsRealAuth       int
}

func (r *userRepo) ListUserID() (userIDList []string, err error) {
	err = db.DB.Model(&model.User{}).Select("user_id").Find(&userIDList).Error
	return
}

func (r *userRepo) Search(keyword string, offset int, limit int) (list []model.User, count int64, err error) {
	needPaging := !(offset == limit && limit == 0)
	var listDB *gorm.DB
	listDB = db.DB.Model(&model.User{}).Where("(user_id = ? OR account = ? OR phone_number = ? OR nick_name = ?)", keyword, keyword, keyword, keyword)

	if needPaging {
		listDB = listDB.Offset(offset).Limit(limit)
	}

	if err = listDB.Find(&list).Error; err != nil {
		return
	}

	if needPaging {
		if err = listDB.Count(&count).Error; err != nil {
			return
		}
	}
	return
}

func (r *userRepo) GetByPassword(opts ...WhereOption) (*model.User, error) {
	m := &model.User{}
	var opt WhereOption
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		return nil, code.ErrUserIdNotExist
	}
	query := db.DB
	if opt.Account != "" && opt.Password != "" {
		query = query.Where("account = ?", opt.Account).Where("password = ?", opt.Password)
		if err := query.First(m).Error; err != nil {
			return nil, err
		}
		return m, nil
	}
	if opt.PhoneNumber != "" && opt.Password != "" {
		query = query.Where("phone_number = ?", opt.PhoneNumber).Where("password = ?", opt.Password)
		if err := query.First(m).Error; err != nil {
			return nil, err
		}
		return m, nil
	}

	return nil, code.ErrUserIdNotExist
}

func (r *userRepo) GetByUserID(opts ...WhereOption) (*model.User, error) {
	m := &model.User{}
	var opt WhereOption
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		return nil, code.ErrBadRequest
	}
	if opt.Id == 0 && opt.UserId == "" && opt.Account == "" && opt.PhoneNumber == "" {
		return nil, code.ErrUserIdNotExist
	}
	query := db.DB
	if opt.Id != 0 {
		query = query.Where("id = ?", opt.Id)
	}
	if opt.Status != 0 {
		query = query.Where("status = ?", opt.Status)
	}
	if opt.NickName != "" {
		query = query.Where("nick_name = ?", opt.NickName)
	}
	orCondition := db.DB
	if opt.UserId != "" {
		orCondition = orCondition.Or("user_id = ?", opt.UserId)
	}
	if opt.Account != "" {
		orCondition = orCondition.Or("account = ?", opt.Account)
	}
	if opt.PhoneNumber != "" {
		orCondition = orCondition.Or("phone_number = ?", fmt.Sprintf("%s%s", opt.CountryCode, opt.PhoneNumber))
	}
	if opt.PhoneNumber != "" && opt.CountryCode != "" {
		orCondition = orCondition.Or("phone_number = ? and country_code=?", opt.PhoneNumber, opt.CountryCode)
	}
	query = query.Where(orCondition)
	if err := query.First(m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, code.ErrUserIdNotExist
		}
		return nil, err
	}

	return m, nil
}

func (r *userRepo) GetByFriendID(opts ...WhereOption) (*model.User, error) {
	m := &model.User{}
	var opt WhereOption
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		return nil, code.ErrBadRequest
	}
	query := db.DB
	if opt.Id != 0 {
		query = query.Where("id = ?", opt.Id)
	}
	if opt.Status != 0 {
		query = query.Where("status = ?", opt.Status)
	}
	orCondition := db.DB
	if opt.UserId != "" {
		orCondition = orCondition.Or("user_id like ?", "%"+opt.UserId+"%")
	}
	if opt.Account != "" {
		orCondition = orCondition.Or("account like ?", "%"+opt.Account+"%")
	}
	if opt.NickName != "" {
		orCondition = orCondition.Or("nick_name like ?", "%"+opt.NickName+"%")
	}
	if opt.PhoneNumber != "" {
		orCondition = orCondition.Or("phone_number like ?", "%"+opt.PhoneNumber+"%")
	}
	query = query.Where(orCondition)
	if err := query.First(m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, code.ErrUserIdNotExist
		}
		return nil, err
	}

	return m, nil
}

func (r *userRepo) OrExists(opt WhereOption) bool {
	query := db.DB
	if opt.UserId != "" {
		query = query.Where("user_id = ?", opt.UserId)
	}
	if opt.Account != "" {
		query = query.Or("account = ?", opt.Account)
	}
	if opt.PhoneNumber != "" && opt.CountryCode != "" {
		query = query.Or("phone_number = ?", fmt.Sprintf("%s%s", opt.CountryCode, opt.PhoneNumber))
	}
	pnCond := db.DB
	if opt.PhoneNumber != "" {
		pnCond = pnCond.Where("phone_number = ?", opt.PhoneNumber)
	}
	if opt.CountryCode != "" {
		pnCond.Where("country_code = ?", opt.CountryCode)
	}
	query = query.Or(pnCond)
	var c int64
	if query.First(&model.User{}).Count(&c); c == 0 {
		return false
	}

	return true
}
func (r *userRepo) UpdateRealName(userID string, req model.RealNameReq) (user model.User, err error) {
	user = model.User{
		RealName:    req.RealName,
		IDNo:        req.IDNo,
		IDFrontImg:  req.IDFrontImg,
		IDBackImg:   req.IDBackImg,
		IsRealAuth:  2,
		RealAuthMsg: "审核中请耐心等待",
	}
	err = db.DB.Model(&model.User{}).Where("user_id = ?", userID).
		Where("is_real_auth IN ?", []int{1, 4}).Updates(&user).Error
	if err != nil {
		return
	}
	err = db.DB.Model(&model.User{}).Where("user_id  = ?", userID).Find(&user).Error
	return
}

func (r *userRepo) GetRealName(userID string) (user model.User, err error) {
	err = db.DB.Model(&model.User{}).Where("user_id  = ?", userID).Find(&user).Error
	return
}

func (r *userRepo) UpdateBy(opt WhereOption, data *model.User) (*model.User, error) {
	var err error
	if opt.Id == 0 && opt.PhoneNumber == "" && opt.Account == "" && opt.UserId == "" {
		return nil, code.ErrBadRequest
	}
	query := db.DB.Model(&model.User{})
	if opt.Id > 0 {
		query = query.Where("id = ?", opt.Id)
	}
	if opt.UserId != "" {
		query = query.Where("user_id = ? ", opt.UserId)
	}
	if opt.Account != "" {
		query = query.Where("account = ? ", opt.Account)
	}
	if opt.PhoneNumber != "" {
		query = query.Where("phone_number = ? ", opt.PhoneNumber)
	}
	if opt.CountryCode != "" {
		query = query.Where("country_code = ? ", opt.CountryCode)
	}
	if err = query.UpdateColumns(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *userRepo) UserInfoUpdate(user_id string, req model.UserInfoUpdateReq) (err error) {
	updates, err := util.StructToMapWithoutNil(req, "json")
	delete(updates, "operation_id")
	if err = db.DB.Model(&model.User{}).Where("user_id = ?", user_id).Updates(updates).Error; err != nil {
		return err
	}
	return err
}

func (r *userRepo) Create(data *model.User) (*model.User, error) {
	if data.PhoneNumber == "" && data.UserID == "" {
		return nil, code.ErrBadRequest
	}
	err := db.DB.Omit("deleted_at").Create(&data).Error
	return data, err
}

func (r *userRepo) GetBaseInfoByUserId(userId string) (*model.UserBaseInfo, error) {
	u, err := r.GetByUserID(WhereOption{UserId: userId})
	if err != nil {
		return nil, err
	}
	return &model.UserBaseInfo{
		UserId:      u.UserID,
		Account:     u.Account,
		FaceURL:     u.FaceURL,
		BigFaceURL:  u.BigFaceURL,
		Gender:      u.Gender,
		NickName:    u.NickName,
		Signatures:  u.Signatures,
		Age:         u.Age,
		IsPrivilege: u.IsPrivilege,
		PhoneNumber: u.PhoneNumber,
		CountryCode: u.CountryCode,
	}, nil
}
func (r *userRepo) GetBaseInfoByPhoneNUmber(phoneNumber string) (*model.UserBaseInfo, error) {
	u, err := r.GetByUserID(WhereOption{PhoneNumber: phoneNumber})
	if err != nil {
		return nil, err
	}
	return &model.UserBaseInfo{
		UserId:      u.UserID,
		Account:     u.Account,
		FaceURL:     u.FaceURL,
		BigFaceURL:  u.BigFaceURL,
		Gender:      u.Gender,
		NickName:    u.NickName,
		Signatures:  u.Signatures,
		Age:         u.Age,
		IsPrivilege: u.IsPrivilege,
		PhoneNumber: u.PhoneNumber,
		CountryCode: u.CountryCode,
	}, nil
}
func (r *userRepo) GetByDeviceID(opts ...WhereOption) (*model.User, error) {
	m := &model.User{}
	var opt WhereOption
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		return nil, code.ErrBadRequest
	}
	if opt.DeviceId == "" {
		return nil, code.ErrUserIdNotExist
	}
	query := db.DB
	if opt.DeviceId != "" {
		query = query.Where("device_id = ?", opt.DeviceId)
	}
	if opt.Platform != 0 {
		query = query.Where("platform = ?", opt.Platform)
	}
	if opt.Status != 0 {
		query = query.Where("status = ?", opt.Status)
	}
	if opt.UserModel != 0 {
		query = query.Where("user_model = ?", opt.UserModel)
	}
	if err := query.First(m).Error; err != nil {
		return nil, err
	}

	return m, nil
}

func (r *userRepo) GetInfo(opts ...WhereOption) (*model.User, error) {
	m := &model.User{}
	var opt WhereOption
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		return nil, code.ErrBadRequest
	}
	query := db.DB
	if opt.DeviceId != "" {
		query = query.Where("device_id = ?", opt.DeviceId)
	}
	if opt.Platform != 0 {
		query = query.Where("platform = ?", opt.Platform)
	}
	if opt.Status != 0 {
		query = query.Where("status = ?", opt.Status)
	}
	if opt.UserId != "" {
		query = query.Where("user_id = ?", opt.UserId)
	}
	if opt.Account != "" {
		query = query.Where("account = ?", opt.Account)
	}
	if opt.NickName != "" {
		query = query.Where("nick_name = ?", opt.NickName)
	}
	if err := query.First(m).Error; err != nil {
		return nil, err
	}

	return m, nil
}

func (r *userRepo) CountRegisterInfo(opts ...WhereOption) (int64, error) {
	var (
		count int64
		opt   WhereOption
	)
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		return 0, code.ErrBadRequest
	}
	tableName := new(model.User).TableName()
	query := db.DB.Table(tableName)
	if opt.IP != "" {
		query = query.Where("ip = ?", opt.IP)
	}
	if opt.RegisterIp != "" {
		query = query.Where("register_ip = ?", opt.RegisterIp)
	}
	if opt.DeviceId != "" {
		query = query.Where("device_id = ?", opt.DeviceId)
	}
	if opt.RegisterDeviceId != "" {
		query = query.Where("register_device_id = ?", opt.RegisterDeviceId)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *userRepo) RecordUserDeviceAndIP(userID, deviceID, ip string, platform int64) (err error) {
	var count int64
	err = db.DB.Model(model.UserDevice{}).Where("user_id = ?", userID).Where("device_id = ?", deviceID).Count(&count).Error
	if count == 0 {
		err = db.DB.Create(&model.UserDevice{UserID: userID, DeviceID: deviceID, Platform: platform, Status: 1}).Error
	}
	err = db.DB.Model(model.UserIp{}).Where("user_id = ?", userID).Where("ip = ?", ip).Count(&count).Error
	if count == 0 {
		err = db.DB.Create(&model.UserIp{UserID: userID, Ip: ip, Status: 1}).Error
	}
	return
}

func (r *userRepo) UpdateWallet(userID, flag string, amount int64) (err error) {
	err = db.DB.Model(&model.User{}).Where("user_id = ? ", userID).UpdateColumn("balance", gorm.Expr(fmt.Sprintf("balance %s ?", flag), amount)).Error
	return err
}

func (r *userRepo) RecordUserLoginHistory(userID, deviceID, ip, brand string, platform int64) (err error) {
	add := model.LoginHistory{
		CreatedAt: time.Now().Unix(),
		UserID:    userID,
		Ip:        ip,
		Platform:  platform,
		DeviceID:  deviceID,
		Brand:     brand,
	}
	err = db.DB.Create(&add).Error
	return
}

func (r *userRepo) FavoriteImagePaging(userID string, offset, limit int) (records []model.FavoriteImage, count int64, err error) {
	err = db.DB.Model(&model.FavoriteImage{}).Where("user_id = ? ", userID).Offset(offset).Limit(limit).Find(&records).Offset(-1).Limit(-1).Count(&count).Error
	return
}

func (r *userRepo) FavoriteImageAdd(add model.FavoriteImage) error {
	err := db.DB.Create(&add).Error
	return err
}

func (r *userRepo) FavoriteImageDel(del model.FavoriteImage) error {
	err := db.DB.Where("user_id = ?", del.UserID).Where("uuid = ?", del.UUID).Delete(&model.FavoriteImage{}).Error
	return err
}

func (r *userRepo) FavoriteImageExist(img model.FavoriteImage) bool {
	var count int64
	err := db.DB.Model(model.FavoriteImage{}).Where("user_id = ?", img.UserID).Where("uuid = ?", img.UUID).Count(&count).Error
	if err != nil {
		return false
	}
	return !(count == 0)
}
func (r *userRepo) FavoriteImageCountLimit(img model.FavoriteImage) bool {
	var count int64
	err := db.DB.Model(model.FavoriteImage{}).Where("user_id = ?", img.UserID).Count(&count).Error
	if err != nil {
		return false
	}
	return count >= 100
}

func (r *userRepo) CheckUserOnline(userId string) bool {
	var count int64
	err := db.DB.Model(&model.User{}).
		Where("user_id = ? AND online_status = ?", userId, model.OnlineStatusTypeOnline).
		Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func (r *userRepo) GetUserInfoInIDs(userId []string) (users []model.User, err error) {
	err = db.DB.Model(&model.User{}).
		Where("user_id  IN ?", userId).Find(&users).Error
	return
}
