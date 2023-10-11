package control

import (
	"bufio"
	"im/config"
	domainModel "im/internal/control/domain/model"
	errorModel "im/internal/control/error/model"
	menuModel "im/internal/control/menu/model"
	"im/pkg/db"
	"im/pkg/logger"
	"os"
	"path"
)

func Init() {
	_ = db.DB.AutoMigrate(new(errorModel.ErrLog),
		new(menuModel.CmsMenu),
		new(menuModel.CmsMenuApi),
		new(menuModel.Config),
		new(domainModel.DomainSite),
		new(domainModel.DomainWarning),
	)
	var needInitTableData int64
	err := db.DB.Model(menuModel.CmsMenu{}).Count(&needInitTableData).Error
	if err != nil {
		logger.Sugar.Errorw("Init Count cms_menus error.")
		return
	}
	if needInitTableData == 0 {
		InitTableData()
	}
}

func InitTableData() {
	cfg := config.Config
	sqlPath := "resources/initsql/control.sql"
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
