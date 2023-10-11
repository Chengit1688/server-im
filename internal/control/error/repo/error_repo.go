package repo

import (
	"im/internal/control/error/model"
	"im/pkg/db"

	"gorm.io/gorm"
)

var ErrorRepo = new(errorRepo)

type errorRepo struct{}

func (r *errorRepo) InsertIntoErrLog(toInsertInfo model.ErrLog) error {
	err := db.DB.Table("errlog").Create(&toInsertInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *errorRepo) QueryErrlogInfo(offset, limit int, stime, etime int64, userId string, macType string) (list []model.ErrLog, count int64, err error) {
	needPaging := !(offset == limit && limit == 0)

	var listDB *gorm.DB
	listDB = db.DB.Model(model.ErrLog{}).Where("1=1")

	if macType != "" {
		listDB = listDB.Where("mac_type=?", macType)
	}
	if userId != "" {
		listDB = listDB.Where("user_id = ?", userId)
	}
	if stime > 0 {
		listDB = listDB.Where("create_time >= ?", stime)
	}
	if etime > 0 {
		listDB = listDB.Where("create_time <= ?", etime)
	}
	if needPaging {
		if err = listDB.Count(&count).Error; err != nil {
			return
		}
		listDB = listDB.Offset(offset).Limit(limit)
	}

	if err = listDB.Order("create_time desc").Find(&list).Error; err != nil {
		return
	}

	return
}
