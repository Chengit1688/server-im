package repo

import (
	configModel "im/internal/cms_api/config/model"
	"im/pkg/code"
	"im/pkg/common/constant"
	"im/pkg/db"
)

var AppVersion = new(appVersion)

type appVersion struct{}

type WhereOptionForVersion struct {
	Id           int64
	Platform     int64
	Status       int64
	DeleteStatus int64
	BeginDate    int64
	EndDate      int64
	Ext          string
	Version      string
	PageSize     int64
	Size         int64
}

func (r *appVersion) Create(data *configModel.AppVersion) (*configModel.AppVersion, error) {
	err := db.DB.Create(&data).Error
	return data, err
}

func (r *appVersion) GetInfo(opts ...WhereOptionForVersion) (*configModel.AppVersion, error) {
	var m configModel.AppVersion
	var opt WhereOptionForVersion
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		return nil, code.ErrSettingNotExist
	}
	query := db.DB
	if opt.Id != 0 {
		query = query.Where("id = ? ", opt.Id)
	}
	if opt.Platform != 0 {
		query = query.Where("platform = ?", opt.Platform)
	}
	if opt.Version != "" {
		query = query.Where("version = ?", opt.Version)
	}
	if opt.Ext != "" {
		query = query.Where(opt.Ext)
	}
	err := query.First(&m).Error
	return &m, err
}

func (r *appVersion) Exists(opts ...WhereOptionForVersion) (int64, error) {
	var (
		opt   WhereOptionForVersion
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
	if opt.Platform != 0 {
		query = query.Where("platform = ?", opt.Platform)
	}
	if opt.Version != "" {
		query = query.Where("version = ?", opt.Version)
	}
	if opt.Status != 0 {
		query = query.Where("status = ?", opt.Status)
	}
	if opt.Ext != "" {
		query = query.Where(opt.Ext)
	}
	query = query.Where("delete_status != ?", constant.SwitchOn)

	err = query.First(&configModel.AppVersion{}).Count(&total).Error

	return total, err
}

func (r *appVersion) UpdateById(opt WhereOptionForVersion, data *configModel.AppVersion) error {
	var err error
	if opt.Id == 0 {
		return code.ErrBadRequest
	}
	query := db.DB.Model(&configModel.AppVersion{})
	if opt.Id != 0 {
		query = query.Where("id = ?", opt.Id)
	}

	if err = query.UpdateColumns(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r *appVersion) GetList(req configModel.VersionListReq) (list []configModel.AppVersion, count int64, err error) {
	req.Pagination.Check()
	query := db.DB.Model(configModel.AppVersion{})
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}
	if req.Platform != 0 {
		query = query.Where("platform = ?", req.Platform)
	}
	if req.BeginDate != 0 {
		query = query.Where("created_at >= ?", req.BeginDate)
	}
	if req.EndDate != 0 {
		query = query.Where("created_at <= ?", req.EndDate)
	}
	if req.Limit == 0 {
		req.Limit = 5
	}
	query = query.Where("delete_status != ?", constant.SwitchOn)

	err = query.Order("created_at asc").Offset(req.Offset).Limit(req.Limit).Find(&list).Error
	_ = query.Offset(-1).Limit(-1).Count(&count).Error
	return
}
