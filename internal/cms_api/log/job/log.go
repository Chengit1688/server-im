package job

import (
	"im/config"
	"im/internal/cms_api/log/repo"
	"time"
)

func LogClearJob() {
	cfg := config.Config
	days := cfg.Log.RecordKeepDays
	if days == 0 {
		return
	}
	now := time.Now()
	before := now.Add(time.Hour * time.Duration(24*days*-1))
	repo.LogRepo.OperateLogClear(before.UnixMilli())
}
