package repo

import (
	"im/internal/api/user/model"
	"im/pkg/db"
)

var SignGifts = new(signGifts)

type signGifts struct{}

func (s *signGifts) FetchOne(day int64) (data model.SignGifts, err error) {
	err = db.DB.Model(&model.SignGifts{}).Where("day = ?", day).Find(&data).Error
	return
}
