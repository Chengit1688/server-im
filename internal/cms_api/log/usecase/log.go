package usecase

import (
	"im/internal/cms_api/log/model"
	"im/internal/cms_api/log/repo"
)

var LogUseCase = new(logUseCase)

type logUseCase struct{}

func (u *logUseCase) AddOperateLog(log model.OperateLogs) (err error) {
	return repo.LogRepo.OperateLogAdd(log)
}

func (u *logUseCase) OperateLogBatchAdd(logs []model.OperateLogs) (err error) {
	return repo.LogRepo.OperateLogBatchAdd(logs)
}
