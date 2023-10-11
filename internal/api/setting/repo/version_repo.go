package repo

import (
	model "im/internal/cms_api/config/model"
	"im/pkg/code"
	"im/pkg/db"
)

var VersionRepo = new(versionRepo)

type versionRepo struct{}

type WhereOptionByVersion struct {
	Id       int64
	Platform int64
	Status   int64
}

func (r *versionRepo) GetVersionInfo(opts ...WhereOptionByVersion) (*model.AppVersion, error) {
	m := &model.AppVersion{}
	var opt WhereOptionByVersion
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
		query = query.Where("platform = ? ", opt.Platform)
	}
	if opt.Status != 0 {
		query = query.Where("status = ? ", opt.Status)
	}
	if err := query.First(m).Order("created_at desc, updated_at desc").Limit(1).Error; err != nil {
		return nil, err
	}
	return m, nil
}
