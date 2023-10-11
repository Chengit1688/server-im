package repo

import (
	"im/internal/cms_api/log/model"
	"im/pkg/db"
)

const createBatchSize = 1000

var LogRepo = new(logRepo)

type logRepo struct{}

func (r *logRepo) OperateLogPaging(req model.OperateLogPagingReq) (logs []model.OperateLogs, count int64, err error) {
	req.Pagination.Check()
	tx := db.DB.Model(model.OperateLogs{}).Preload("User")

	if len(req.Search) != 0 {
		tx = tx.Where("service_id = ?", req.Search).Or("operation_id = ?", req.Search).Or("request_url = ?", req.Search)
	}

	if req.CreatedTimeStart != 0 {
		tx = tx.Where("created_at >= ?", req.CreatedTimeStart)
	}
	if req.CreatedTimeEnd != 0 {
		tx = tx.Where("created_at <= ?", req.CreatedTimeEnd)
	}

	err = tx.Offset(req.Offset).Limit(req.Limit).Order("created_at desc").Find(&logs).Limit(-1).Offset(-1).Count(&count).Error
	return
}

func (r *logRepo) OperateLogAdd(log model.OperateLogs) (err error) {
	err = db.DB.Create(&log).Error
	return
}

func (r *logRepo) OperateLogBatchAdd(logs []model.OperateLogs) (err error) {
	err = db.DB.CreateInBatches(&logs, createBatchSize).Error
	return
}

func (r *logRepo) OperateLogClear(before int64) (err error) {
	err = db.DB.Where("created_at <= ?", before).Delete(&model.OperateLogs{}).Error
	return
}
