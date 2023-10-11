package repo

import (
	"gorm.io/gorm"
	"im/internal/api/user/model"
	"im/pkg/code"
	"im/pkg/db"
)

var UserConfigRepo = new(userConfigRepo)

type userConfigRepo struct{}

type WhereOptionForUserConfig struct {
	Id      int64
	UserId  string
	Version int64
}

func (r *userConfigRepo) UpdateBy(opt WhereOptionForUserConfig, data model.UserConfig) (model.UserConfig, error) {
	var err error
	query := db.DB.Model(&model.UserConfig{})
	if opt.Id > 0 {
		query = query.Where("id = ?", opt.Id)
	}
	if opt.UserId != "" {
		query = query.Where("user_id = ? ", opt.UserId)
	}
	if opt.Version != 0 {
		query = query.Where("version = ? ", opt.Version)
	}

	if err = query.UpdateColumns(&data).Error; err != nil {
		return model.UserConfig{}, err
	}
	return data, nil
}

func (r *userConfigRepo) CreateOrUpdate(userConfig model.UserConfig) (model.UserConfig, error) {
	m := model.UserConfig{}
	var err, err2 error
	if err = db.DB.Where("user_id = ?", userConfig.UserId).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err2 = db.DB.Create(&userConfig).Error; err2 != nil {
				return model.UserConfig{}, err2
			}
			return userConfig, nil
		}
		return model.UserConfig{}, err
	}
	values := map[string]interface{}{"content": userConfig.Content, "version": gorm.Expr("version + ?", 1)}
	if err = db.DB.Model(model.UserConfig{}).Where("user_id = ?", userConfig.UserId).Updates(values).Error; err != nil {
		return model.UserConfig{}, err
	}
	return r.GetUserConfig(WhereOptionForUserConfig{UserId: userConfig.UserId})

}

func (r *userConfigRepo) GetUserConfig(opts ...WhereOptionForUserConfig) (model.UserConfig, error) {
	var opt WhereOptionForUserConfig
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		return model.UserConfig{}, code.ErrUserConfig
	}
	query := db.DB
	if opt.Id != 0 {
		query = query.Where("id = ?", opt.Id)
	}
	if opt.UserId != "" {
		query = query.Where("user_id = ?", opt.UserId)
	}
	if opt.Version != 0 {
		query = query.Where("version = ?", opt.Version)
	}
	m := model.UserConfig{}
	if err := query.First(&m).Error; err != nil {
		return model.UserConfig{}, err
	}

	return m, nil
}
