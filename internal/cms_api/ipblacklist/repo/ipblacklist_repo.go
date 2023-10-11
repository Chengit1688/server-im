package repo

import (
	"fmt"
	"im/internal/cms_api/ipblacklist/model"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/util"
)

var IPBlackListRepo = new(ipblacklistRepo)

type ipblacklistRepo struct{}

func (r *ipblacklistRepo) Paging(req model.GetIPListReq) (ips []model.IPBlackList, count int64, err error) {

	req.Pagination.Check()
	tx := db.DB.Model(model.IPBlackList{})
	if len(req.IP) != 0 {
		tx = tx.Where(fmt.Sprintf("ip like %q", ("%" + req.IP + "%")))
	}
	if len(req.Note) != 0 {
		tx = tx.Where(fmt.Sprintf("note like %q", ("%" + req.Note + "%")))
	}
	err = tx.Offset(req.Offset).Limit(req.Limit).Find(&ips).Limit(-1).Offset(-1).Count(&count).Error
	return
}

func (r *ipblacklistRepo) Get(id string) (ip model.IPBlackList, err error) {

	tx := db.DB.Model(model.IPBlackList{})
	err = tx.Where("id = ?", id).First(&ip).Error
	return
}

func (r *ipblacklistRepo) GetByIP(ip string) (item model.IPBlackList, err error) {

	tx := db.DB.Model(model.IPBlackList{})
	err = tx.Where("ip = ?", ip).First(&item).Error
	return
}

func (r *ipblacklistRepo) Add(req model.IPBlackList) (ip model.IPBlackList, err error) {

	tx := db.DB.Model(model.IPBlackList{})
	err = tx.Create(&req).Error
	if err == nil {

		err = db.DB.Table("user_ips").Where("ip = ?", req.IP).Where("status = ?", 1).UpdateColumn("status", 2).Error
		if err != nil {
			logger.Sugar.Errorw("db error", "func", util.GetSelfFuncName(), "error", err)
			return ip, err
		}
		err = IpBlackListCache.Add(req.IP)
		if err != nil {
			logger.Sugar.Errorw("redis error", "func", util.GetSelfFuncName(), "error", err)
			return
		}
	}
	if err != nil {
		logger.Sugar.Errorw("db error", "func", util.GetSelfFuncName(), "error", err)
		return ip, err
	}
	return req, nil
}

func (r *ipblacklistRepo) Update(id string, req model.IPBlackList) (ip model.IPBlackList, err error) {

	old, err := r.Get(id)
	updates, err := util.StructToMap(req, "json")
	delete(updates, "id")
	if err = db.DB.Model(&model.IPBlackList{}).Where("id = ?", id).Updates(&updates).Error; err != nil {
		return ip, err
	}
	if old.IP != req.IP {

		err = IpBlackListCache.Del(old.IP)
		if err != nil {
			return
		}
		err = db.DB.Table("user_ips").Where("ip = ?", old.IP).Where("status = ?", 2).UpdateColumn("status", 1).Error
		if err != nil {
			return
		}
		err = IpBlackListCache.Add(req.IP)
		if err != nil {
			return
		}
		err = db.DB.Table("user_ips").Where("ip = ?", req.IP).Where("status = ?", 1).UpdateColumn("status", 2).Error
		if err != nil {
			return
		}
	}
	return r.Get(id)
}

func (r *ipblacklistRepo) Delete(id string) (err error) {

	have, err := r.Get(id)
	tx := db.DB.Model(model.IPBlackList{})
	err = tx.Where("id = ?", id).Delete(&model.IPBlackList{}).Error
	if err == nil {

		err = IpBlackListCache.Del(have.IP)
		err = db.DB.Table("user_ips").Where("ip = ?", have.IP).Where("status = ?", 2).UpdateColumn("status", 1).Error
	}
	return
}

func (r *ipblacklistRepo) DeleteInBatch(ips []string) (err error) {

	tx := db.DB.Model(model.IPBlackList{})
	err = tx.Where("ip IN ?", ips).Delete(&model.IPBlackList{}).Error
	if err == nil {

		for _, ip := range ips {
			err = IpBlackListCache.Del(ip)
			if err != nil {
				return
			}
		}
		err = db.DB.Table("user_ips").Where("ip IN ?", ips).Where("status = ?", 2).UpdateColumn("status", 1).Error
	}
	return
}

func (r *ipblacklistRepo) SyncCache() {

	dbIpInfos := new([]model.IPBlackList)
	tx := db.DB.Model(model.IPBlackList{})
	err := tx.Find(&dbIpInfos).Error
	if err != nil {
		logger.Sugar.Errorw("db error", "func", util.GetSelfFuncName(), "error", err)
		return
	}
	cacheIpInfos := IpBlackListCache.GetAllInfo()

	for _, item := range *dbIpInfos {
		if _, ok := cacheIpInfos[item.IP]; !ok {
			err = IpBlackListCache.Add(item.IP)
			if err != nil {
				logger.Sugar.Errorw("redis error", "func", util.GetSelfFuncName(), "error", err)
				return
			}
		}
	}

	for k, _ := range cacheIpInfos {
		var status int = 0
		for _, item := range *dbIpInfos {
			if item.IP == k {
				status = 1
			}
		}
		if status != 1 {
			err = IpBlackListCache.Del(k)
			if err != nil {
				logger.Sugar.Errorw("redis error", "func", util.GetSelfFuncName(), "error", err)
				return
			}
		}
	}
}

func (r *ipblacklistRepo) AddInBatch(ips []string) (err error) {
	var count int64
	for _, ip := range ips {
		err = db.DB.Model(model.IPBlackList{}).Where("ip = ?", ip).Count(&count).Error
		if count == 0 {
			add := &model.IPBlackList{IP: ip, Note: "IP封禁"}
			err = db.DB.Create(&add).Error
			if err != nil {
				logger.Sugar.Errorw("db error", "func", util.GetSelfFuncName(), "error", err)
				return err
			}
			err = db.DB.Table("user_ips").Where("ip = ?", ip).Where("status = ?", 1).UpdateColumn("status", 2).Error
			if err != nil {
				logger.Sugar.Errorw("db error", "func", util.GetSelfFuncName(), "error", err)
				return err
			}
			err = IpBlackListCache.Add(ip)
			if err != nil {
				logger.Sugar.Errorw("redis error", "func", util.GetSelfFuncName(), "error", err)
				return
			}
		} else {
			err = db.DB.Table("user_ips").Where("ip = ?", ip).Where("status = ?", 1).UpdateColumn("status", 2).Error
			if err != nil {
				logger.Sugar.Errorw("redis error", "func", util.GetSelfFuncName(), "error", err)
				return
			}
		}
	}
	return
}
