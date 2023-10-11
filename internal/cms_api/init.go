package cms_api

import (
	"bufio"
	"im/config"
	adminModel "im/internal/cms_api/admin/model"
	chatModel "im/internal/cms_api/chat/model"
	configModel "im/internal/cms_api/config/model"
	"im/internal/cms_api/cron"
	cornJob "im/internal/cms_api/cron/job"
	dashboardjob "im/internal/cms_api/dashboard/job"
	dashboardModel "im/internal/cms_api/dashboard/model"
	discoverModel "im/internal/cms_api/discover/model"
	ipblacklistModel "im/internal/cms_api/ipblacklist/model"
	ipblacklistRepo "im/internal/cms_api/ipblacklist/repo"
	ipwhitelistModel "im/internal/cms_api/ipwhitelist/model"
	ipwhitelistRepo "im/internal/cms_api/ipwhitelist/repo"
	logModel "im/internal/cms_api/log/model"
	"im/internal/cms_api/operation/model"
	roleModel "im/internal/cms_api/role/model"
	userjob "im/internal/cms_api/user/job"
	walletModel "im/internal/cms_api/wallet/model"
	"im/pkg/db"
	"im/pkg/logger"
	"os"
	"path"
)

func Init() {

	if err := db.DB.AutoMigrate(new(adminModel.Admin)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(roleModel.AdminApi)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(roleModel.CmsMenu)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(roleModel.CmsMenuApi)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(roleModel.CmsMenuRole)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(roleModel.CmsRole)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(configModel.Config)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(discoverModel.Discover)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(ipblacklistModel.IPBlackList)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(configModel.ShieldWords)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(chatModel.MultiSendRecord)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(chatModel.MultiSendUser)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(model.Suggestion)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(dashboardModel.DashboardDailyData)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(walletModel.BillingRecords)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(walletModel.RedpackSingleRecords)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(walletModel.WithdrawRecords)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(walletModel.RedpackGroupRecords)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(walletModel.RedpackGroupRecvs)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(logModel.OperateLogs)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(ipwhitelistModel.IPWhiteList)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(discoverModel.PrizeList)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(discoverModel.RedeemPrizeLog)); err != nil {
		panic(err)
	}

	go cron.Start()

	dashboardjob.Init()

	ipblacklistRepo.IPBlackListRepo.SyncCache()

	ipwhitelistRepo.IPWhiteListRepo.SyncCache()

	cornJob.Init()

	userjob.Init()

	var needInitTableData int64
	err := db.DB.Model(adminModel.Admin{}).Count(&needInitTableData).Error
	if err != nil {
		logger.Sugar.Errorw("Init Count cms_admins error.")
		return
	}
	if needInitTableData == 0 {
		InitTableData()
	}

}

func InitTableData() {
	cfg := config.Config
	sqlPath := "resources/initsql/cms.sql"
	sqlFile := path.Join(cfg.Captcha.DefaultResourceRoot, sqlPath)
	_, err := os.Stat(sqlFile)
	if err != nil {
		logger.Sugar.Errorw("InitTableData not found initsql.", "cfg.Captcha.DefaultResourceRoot:", cfg.Captcha.DefaultResourceRoot, "sqlfile:", sqlPath)
		return
	}
	file, _ := os.Open(sqlFile)
	defer file.Close()
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		sql := sc.Text()
		if len(sql) < 6 {
			continue
		}
		err := db.DB.Exec(sql).Error
		if err != nil {
			logger.Sugar.Errorw("InitTableData Exec sql error.", "sql:", sql, "error:", err)
			return
		}
	}
	logger.Sugar.Infow("InitTableData Success.")
}
