package usecase

import (
	"fmt"
	chatRepo "im/internal/api/chat/repo"
	configUseCase "im/internal/cms_api/config/usecase"
	"im/pkg/common"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/util"
	"time"
)

var MessageUseCase = new(messageUseCase)

type messageUseCase struct{}

func (c *messageUseCase) CronClear() {
	l := util.NewLock(db.RedisCli, common.LockMessageExpireClear)
	if l.IsLock() {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("cron task is running"))
		return
	}
	if err := l.Lock(); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("lock error, error: %v", err))
		return
	}
	defer l.Unlock()

	cfg, err := configUseCase.ConfigUseCase.GetParameterConfig()
	if err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get parameter config error, error: %v", err))
		return
	}

	if cfg.HistoryTime == 0 {
		return
	}

	clearTime := util.UnixMilliTime(time.Now()) - cfg.HistoryTime*24*60*60*int64(1000)
	for {
		var count int64
		if count, err = chatRepo.MessageRepo.Clear(clearTime); err != nil {
			break
		}

		if count < 100 {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}
