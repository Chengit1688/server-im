package repo

import (
	"fmt"
	"gorm.io/gorm/clause"
	friendModel "im/internal/api/friend/model"
	configModel "im/internal/cms_api/config/model"
	"im/pkg/code"
	"im/pkg/db"
	"im/pkg/logger"
)

var DefaultFriendRepo = new(defaultFriendRepo)

type defaultFriendRepo struct{}

type WhereOptionForDefaultFriend struct {
	Id     int64
	UserId string
	Ext    string
}

func (r *defaultFriendRepo) CreateOrUpdate(users *configModel.DefaultFriend) (*configModel.DefaultFriend, error) {
	err := db.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"greet_msg", "remarks", "operation_user_id", "user_primary_id"}),
	}).Create(&users).Error
	return users, err
}

func (r *defaultFriendRepo) Create(data *configModel.DefaultFriend) (*configModel.DefaultFriend, error) {
	err := db.DB.Create(&data).Error
	return data, err
}

func (r *defaultFriendRepo) DeleteById(id int64) error {
	return db.DB.Where("id = ?", id).Delete(&configModel.DefaultFriend{}).Error
}

func (r *defaultFriendRepo) DeleteByUserId(userId string) error {
	return db.DB.Where("user_id = ?", userId).Delete(&configModel.DefaultFriend{}).Error
}

func (r *defaultFriendRepo) GetInfo(opts ...WhereOptionForDefaultFriend) (*configModel.DefaultFriend, error) {
	var m configModel.DefaultFriend
	var opt WhereOptionForDefaultFriend
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		return nil, code.ErrSettingNotExist
	}
	query := db.DB
	if opt.Id != 0 {
		query = query.Where("id = ? ", opt.Id)
	}
	if opt.UserId != "" {
		query = query.Where("user_id = ?", opt.UserId)
	}
	if opt.Ext != "" {
		query = query.Where(opt.Ext)
	}
	err := query.First(&m).Error
	return &m, err
}

func (r *defaultFriendRepo) Exists(opts ...WhereOptionForDefaultFriend) (int64, error) {
	var (
		opt   WhereOptionForDefaultFriend
		total int64
		err   error
	)
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		return 0, code.ErrSettingNotExist
	}
	query := db.DB
	if opt.Id != 0 {
		query = query.Where("id = ? ", opt.Id)
	}
	if opt.UserId != "" {
		query = query.Where("user_id = ?", opt.UserId)
	}
	if opt.Ext != "" {
		query = query.Where(opt.Ext)
	}
	err = query.First(&configModel.DefaultFriend{}).Count(&total).Error

	return total, err
}

func (r *defaultFriendRepo) UpdateById(opt WhereOptionForDefaultFriend, data *configModel.DefaultFriend) error {
	var err error
	if opt.Id == 0 {
		return code.ErrBadRequest
	}
	query := db.DB.Model(&configModel.DefaultFriend{})
	if opt.Id != 0 {
		query = query.Where("id = ?", opt.Id)
	}

	if err = query.UpdateColumns(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r *defaultFriendRepo) GetList(req configModel.DefaultFriendListReq) ([]configModel.DefaultFriendListInfo, int64, error) {
	var (
		err   error
		count int64
		list  []configModel.DefaultFriendListInfo
	)
	req.Pagination.Check()
	d := new(configModel.DefaultFriend).TableName()
	f := new(friendModel.Friend).TableName()
	dsql := fmt.Sprintf("any_value(%s.id) id ,%s.user_id,any_value(u.account) account,any_value(u.nick_name) nick_name,any_value(%s.greet_msg) greet_msg,any_value(%s.remarks) remarks,count(f.friend_user_id) as friend_total", d, d, d, d)
	query := db.DB.Model(configModel.DefaultFriend{}).Select(dsql)
	query = query.Joins(fmt.Sprintf("join users as u ON u.user_id = %s.user_id", d))
	query = query.Joins(fmt.Sprintf("left join %s as f ON f.owner_user_id = %s.user_id", f, d))

	if req.Account != "" {
		query = query.Where("u.account = ?", req.Account)
	}
	if req.NickName != "" {
		query = query.Where("u.nick_name = ?", req.NickName)
	}
	if req.Remarks != "" {
		query = query.Where(fmt.Sprintf("any_value(%s.remarks) like ?", d), "%"+req.Remarks+"%")
	}
	if req.UserId != "" {
		query = query.Where("u.user_id = ?", req.UserId)
	}
	query = query.Group(fmt.Sprintf("%s.user_id", d))
	if req.Limit == 0 {
		req.Limit = 5
	}
	if err = query.Order(fmt.Sprintf("any_value(%s.created_at) asc", d)).Offset(req.Offset).Limit(req.Limit).Find(&list).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if len(list) == 1 && req.Page == 1 {
		count = 1
	}

	return list, count, nil
}

func (r *defaultFriendRepo) RandFriend() (*configModel.DefaultFriend, error) {
	var (
		err error

		cnt  int64
		resp []configModel.DefaultFriend
	)

	if err = db.DB.Model(&configModel.DefaultFriend{}).Count(&cnt).Error; err != nil {
		return nil, err
	}
	version, _ := DefaultFriendCahce.GetVersion()
	if err = db.DB.Model(&configModel.DefaultFriend{}).Find(&resp).Error; err != nil || len(resp) <= 0 {
		return nil, code.ErrUserIdNotExist
	}
	index := 0
	if version < cnt {
		index = int(version)
	} else {
		index = int(version % cnt)
	}
	res := resp[index]
	logger.Sugar.Infof("设置初始值为:%d", index)

	DefaultFriendCahce.setVersionIncr()

	return &res, err
}

func (r *defaultFriendRepo) GetAllDefalutFriendIdsAndGreetMsg() ([]configModel.DefaultFriend, error) {
	var result []configModel.DefaultFriend
	err := db.DB.Select("user_id, greet_msg").Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
