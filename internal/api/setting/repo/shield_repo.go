package repo

import (
	"im/internal/api/setting/model"
	configModel "im/internal/cms_api/config/model"
	"im/pkg/code"
	"im/pkg/common/constant"
	"im/pkg/db"
)

var ShieldRepo = new(shieldRepo)

type shieldRepo struct{}

type WhereOptionForShield struct {
	Id           int64
	Status       int64
	DeleteStatus int64
	BeginDate    int64
	EndDate      int64
	Ext          string
	ShieldWords  string
}

func (r *shieldRepo) GetInfo(opts ...WhereOptionForShield) (*configModel.ShieldWords, error) {
	var m configModel.ShieldWords
	var opt WhereOptionForShield
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		return nil, code.ErrSettingNotExist
	}
	query := db.DB
	if opt.Id != 0 {
		query = query.Where("id = ? ", opt.Id)
	}
	if opt.ShieldWords != "" {
		query = query.Where("shield_words = ? ", opt.ShieldWords)
	}
	if opt.Ext != "" {
		query = query.Where(opt.Ext)
	}
	err := query.First(&m).Error
	return &m, err
}
func (r *shieldRepo) GetList(req model.ShieldListReq) (list []configModel.ShieldWords, count int64, err error) {
	req.Pagination.Check()
	query := db.DB.Model(configModel.ShieldWords{}).Select("shield_words")
	query = query.Where("status = ?", constant.SwitchOn)
	query = query.Where("delete_status = ?", constant.SwitchOff)
	err = query.Offset(req.Offset).Limit(req.Limit).Find(&list).Offset(-1).Limit(-1).Count(&count).Error
	return
}
