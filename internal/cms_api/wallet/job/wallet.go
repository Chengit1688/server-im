package job

import (
	"fmt"
	"im/internal/cms_api/wallet/model"
	"im/internal/cms_api/wallet/repo"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/util"
)

func RedpackSingleReturn() {
	redpacks, err := repo.WalletRepo.RedpackSingleNeedReturns()
	if err != nil {
		logger.Sugar.Debugw("个人红包退回协程", "error", err)
		return
	}
	for _, redpack := range redpacks {
		key := fmt.Sprintf(repo.UserWalletKey, redpack.SenderID)
		lock := util.NewLock(db.RedisCli, key)
		if err = lock.Lock(); err != nil {
			return
		}
		err = repo.WalletRepo.RedpackSingleReturn(redpack)
		if err != nil {
			logger.Sugar.Debugw("个人红包退回失败", "error", err)
		}
		lock.Unlock()
	}
}

func RedpackGroupReturn() {
	redpacks, err := repo.WalletRepo.RedpackGroupNeedReturns()
	if err != nil {
		logger.Sugar.Debugw("群红包退回协程", "error", err)
		return
	}

	for _, redpack := range redpacks {
		if redpack.RemainderAmount == 0 {
			logger.Sugar.Debugw(fmt.Sprintf("群ID：%s 群红包ID：%d ,剩余金额为0，无需执行回退逻辑", redpack.GroupID, redpack.GroupRecordsID))
			continue
		}
		logger.Sugar.Infow(fmt.Sprintf("群ID：%s 群红包ID：%d ,开始执行回退逻辑", redpack.GroupID, redpack.GroupRecordsID))
		key := fmt.Sprintf(repo.UserWalletKey, redpack.SenderID)
		lock := util.NewLock(db.RedisCli, key)
		if err = lock.Lock(); err != nil {
			return
		}
		r := model.RedpackGroupRecvs{
			SenderID:       redpack.SenderID,
			RedpackGroupID: redpack.GroupRecordsID,
			Amount:         redpack.RemainderAmount,
			RedpackGroup:   model.RedpackGroupRecords{ID: redpack.GroupRecordsID},
		}
		err = repo.WalletRepo.RedpackGroupReturn(r)
		if err != nil {
			logger.Sugar.Debugw("群红包退回失败", "error", err)
		}
		lock.Unlock()
		logger.Sugar.Infow(fmt.Sprintf("群ID：%s 群红包ID：%d ,结束执行回退逻辑", redpack.GroupID, redpack.GroupRecordsID))
	}

}
