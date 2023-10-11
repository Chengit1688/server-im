package repo

import (
	"im/internal/api/user/model"
	"im/pkg/db"
	"time"
)

var SignInV2 = new(signInV2)

type signInV2 struct{}

func (s *signInV2) FetchOne(userID string) (data model.SignInV2, err error) {
	err = db.DB.Model(&model.SignInV2{}).Where("user_id = ?", userID).Find(&data).Error
	return
}

func (s *signInV2) HasOne(userID string) (count int64, err error) {
	err = db.DB.Model(&model.SignInV2{}).Where("user_id = ?", userID).Count(&count).Error
	return
}

func (s *signInV2) Add(userID string, ip string) (data model.SignInV2, err error) {
	data = model.SignInV2{
		CommonModel: db.CommonModel{
			CreatedAt: time.Now().Unix(),
		},
		UserID:       userID,
		Ip:           ip,
		LastTime:     time.Now().Unix(),
		ContinueDays: 1,
	}
	err = db.DB.Model(&model.SignInV2{}).Create(&data).Error
	return
}

func (s *signInV2) SetContinueDays(userID string, days int64) (data model.SignInV2, err error) {
	data = model.SignInV2{
		CommonModel: db.CommonModel{
			UpdatedAt: time.Now().Unix(),
		},
		ContinueDays: days,
		LastTime:     time.Now().Unix(),
	}
	err = db.DB.Model(&model.SignInV2{}).Where("user_id = ?", userID).Updates(&data).Error
	return
}

func (s *signInV2) UpdateLastTime(userID string) (data model.SignInV2, err error) {
	data = model.SignInV2{
		CommonModel: db.CommonModel{
			UpdatedAt: time.Now().Unix(),
		},
		LastTime: time.Now().Unix(),
	}
	err = db.DB.Model(&model.SignInV2{}).Where("user_id = ?", userID).Updates(&data).Error
	return
}
