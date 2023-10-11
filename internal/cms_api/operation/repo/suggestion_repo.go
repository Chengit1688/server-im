package repo

import (
	"im/internal/cms_api/operation/model"
	"im/pkg/code"
	"im/pkg/db"
	"im/pkg/pagination"
)

var SuggestionRepo = new(suggestionRepo)

type suggestionRepo struct{}

type WhereOptionForSuggestion struct {
	Account   string
	NickName  string
	UserId    string
	Content   string
	Brand     string
	Platform  int64
	BeginDate int64
	EndDate   int64
	pagination.Pagination
}

func (r *suggestionRepo) GetList(opts ...WhereOptionForSuggestion) (list []model.SuggestionInfo, count int64, err error) {
	var opt WhereOptionForSuggestion
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		return []model.SuggestionInfo{}, 0, code.ErrSuggestion
	}
	query := db.DB.Model(model.Suggestion{})
	if opt.Account != "" {
		query = query.Where("account = ?", opt.Account)
	}
	if opt.UserId != "" {
		query = query.Where("user_id = ?", opt.UserId)
	}
	if opt.NickName != "" {
		query = query.Where("nick_name like ?", "%"+opt.NickName+"%")
	}
	if opt.Content != "" {
		query = query.Where("content like ?", "%"+opt.Content+"%")
	}
	if opt.Brand != "" {
		query = query.Where("brand = ?", opt.Brand)
	}
	if opt.Platform != 0 {
		query = query.Where("platform = ?", opt.Platform)
	}
	if opt.BeginDate != 0 {
		query = query.Where("created_at >= ?", opt.BeginDate)
	}
	if opt.EndDate != 0 {
		query = query.Where("created_at <= ?", opt.EndDate)
	}
	err = query.Order("created_at DESC").Offset(opt.Offset).Limit(opt.Limit).Find(&list).Offset(-1).Limit(-1).Count(&count).Error
	if len(list) == 1 && opt.Page == 1 {
		count = 1
	}

	return
}

func (r *suggestionRepo) Create(data model.Suggestion) (model.Suggestion, error) {
	err := db.DB.Create(&data).Error
	return data, err
}
