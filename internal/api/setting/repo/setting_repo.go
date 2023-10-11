package repo

import (
	"im/internal/api/setting/model"
	"im/pkg/code"
	"im/pkg/db"
)

var SettingRepo = new(settingRepo)

type settingRepo struct{}

type WhereOption struct {
	Id         int64
	ConfigType []string
	Status     int64
}

func (r *settingRepo) GetInfo(opts ...WhereOption) (*model.SettingConfig, error) {
	m := &model.SettingConfig{}
	var opt WhereOption
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		return nil, code.ErrSettingNotExist
	}
	query := db.DB
	if opt.Id != 0 {
		query = query.Where("id = ? ", opt.Id)
	}
	if len(opt.ConfigType) != 0 {
		query = query.Where("config_type IN ?", opt.ConfigType)
	}
	if err := query.First(m).Limit(1).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r *settingRepo) GetConfigInfo(opts ...WhereOption) ([]*model.SettingConfig, error) {
	var m []*model.SettingConfig
	var opt WhereOption
	if len(opts) > 0 {
		opt = opts[0]
	}
	query := db.DB
	if opt.Id != 0 {
		query = query.Where("id = ? ", opt.Id)
	}
	if len(opt.ConfigType) != 0 {
		query = query.Where("config_type IN ?", opt.ConfigType)
	}
	if err := query.Find(&m).Error; err != nil {
		return nil, err
	}
	return m, nil
}
