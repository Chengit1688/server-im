package repo

import (
	configModel "im/internal/cms_api/config/model"
	"im/pkg/code"
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

func (r *shieldRepo) Create(data *configModel.ShieldWords) (*configModel.ShieldWords, error) {
	err := db.DB.Create(&data).Error
	return data, err
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

func (r *shieldRepo) Exists(opts ...WhereOptionForShield) (int64, error) {
	var (
		opt   WhereOptionForShield
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
	if opt.ShieldWords != "" {
		query = query.Where("shield_words = ?", opt.ShieldWords)
	}
	if opt.Status != 0 {
		query = query.Where("status = ?", opt.Status)
	}
	if opt.DeleteStatus != 0 {
		query = query.Where("delete_status = ?", opt.DeleteStatus)
	}
	if opt.Ext != "" {
		query = query.Where(opt.Ext)
	}
	err = query.First(&configModel.ShieldWords{}).Count(&total).Error

	return total, err
}

func (r *shieldRepo) UpdateById(opt WhereOptionForShield, data *configModel.ShieldWords) error {
	var err error
	if opt.Id == 0 {
		return code.ErrBadRequest
	}
	query := db.DB.Model(&configModel.ShieldWords{})
	if opt.Id != 0 {
		query = query.Where("id = ?", opt.Id)
	}

	if err = query.UpdateColumns(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r *shieldRepo) DeleteById(opt WhereOptionForShield) error {
	if opt.Id == 0 {
		return code.ErrBadRequest
	}
	return db.DB.Where("id = ?", opt.Id).Delete(&configModel.ShieldWords{}).Error
}

func (r *shieldRepo) GetList(req configModel.ShieldListReq) (list []configModel.ShieldListInfo, count int64, err error) {
	req.Pagination.Check()
	query := db.DB.Model(configModel.ShieldWords{})
	if req.ShieldWords != "" {
		query = query.Where("shield_words like ?", "%"+req.ShieldWords+"%")
	}
	if req.OperationUser != "" {
		query = query.Where("operation_user = ?", req.OperationUser)
	}
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}
	if req.BeginDate != 0 {
		query = query.Where("created_at >= ?", req.BeginDate)
	}
	if req.EndDate != 0 {
		query = query.Where("created_at <= ?", req.EndDate)
	}
	if req.DeleteStatus != 0 {
		query = query.Where("delete_status = ?", req.DeleteStatus)
	}
	if req.Limit == 0 {
		req.Limit = 5
	}
	err = query.Order("created_at DESC").Offset(req.Offset).Limit(req.Limit).Find(&list).Offset(-1).Limit(-1).Count(&count).Error
	return
}
