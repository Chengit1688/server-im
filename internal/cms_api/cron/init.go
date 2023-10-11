package cron

import (
	"github.com/robfig/cron"
	configjob "im/internal/cms_api/config/job"
)

const (
	per30Second = "*/30 * * * * ?"
)

func Start() {
	c := cron.New()
	c.AddFunc(per30Second, configjob.GetCtrlMenu)
	c.Start()
}
