package repo

import (
	"fmt"
	"im/internal/cms_api/ipwhitelist/model"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/util"
)

var IPWhiteListRepo = new(ipWhitelistRepo)

type ipWhitelistRepo struct{}

func (r *ipWhitelistRepo) Paging(req model.GetIPListReq) (ips []model.IPWhiteList, count int64, err error) {

	req.Pagination.Check()
	tx := db.DB.Model(model.IPWhiteList{})
	if len(req.IP) != 0 {
		tx = tx.Where(fmt.Sprintf("ip like %q", ("%" + req.IP + "%")))
	}
	if len(req.Note) != 0 {
		tx = tx.Where(fmt.Sprintf("note like %q", ("%" + req.Note + "%")))
	}
	err = tx.Offset(req.Offset).Limit(req.Limit).Find(&ips).Limit(-1).Offset(-1).Count(&count).Error
	return
}

func (r *ipWhitelistRepo) Get(id string) (ip model.IPWhiteList, err error) {

	tx := db.DB.Model(model.IPWhiteList{})
	err = tx.Where("id = ?", id).First(&ip).Error
	return
}

func (r *ipWhitelistRepo) GetByIP(ip string) (item model.IPWhiteList, err error) {

	tx := db.DB.Model(model.IPWhiteList{})
	err = tx.Where("ip = ?", ip).First(&item).Error
	return
}

func (r *ipWhitelistRepo) Add(req model.IPWhiteList) (ip model.IPWhiteList, err error) {

	tx := db.DB.Model(model.IPWhiteList{})
	err = tx.Create(&req).Error
	if err == nil {

		err = IpWhiteListCache.Add(req.IP)
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

func (r *ipWhitelistRepo) Update(id string, req model.IPWhiteList) (ip model.IPWhiteList, err error) {

	old, err := r.Get(id)
	updates, err := util.StructToMap(req, "json")
	delete(updates, "id")
	if err = db.DB.Model(&model.IPWhiteList{}).Where("id = ?", id).Updates(&updates).Error; err != nil {
		return ip, err
	}
	if old.IP != req.IP {

		err = IpWhiteListCache.Del(old.IP)
		if err != nil {
			return
		}
		err = IpWhiteListCache.Add(req.IP)
		if err != nil {
			return
		}
	}
	return r.Get(id)
}

func (r *ipWhitelistRepo) Delete(id string) (err error) {

	have, err := r.Get(id)
	tx := db.DB.Model(model.IPWhiteList{})
	err = tx.Where("id = ?", id).Delete(&model.IPWhiteList{}).Error
	if err == nil {

		err = IpWhiteListCache.Del(have.IP)
	}
	return
}

func (r *ipWhitelistRepo) SyncCache() {

	dbIpInfos := new([]model.IPWhiteList)
	tx := db.DB.Model(model.IPWhiteList{})
	err := tx.Find(&dbIpInfos).Error
	if err != nil {
		logger.Sugar.Errorw("db error", "func", util.GetSelfFuncName(), "error", err)
		return
	}
	cacheIpInfos := IpWhiteListCache.GetAllInfo()

	for _, item := range *dbIpInfos {
		if _, ok := cacheIpInfos[item.IP]; !ok {
			err = IpWhiteListCache.Add(item.IP)
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
			err = IpWhiteListCache.Del(k)
			if err != nil {
				logger.Sugar.Errorw("redis error", "func", util.GetSelfFuncName(), "error", err)
				return
			}
		}
	}
	var root_id string
	err = db.DB.Table("cms_admins").Where("username = ?", "root").Select("user_id").Scan(&root_id).Error
	if err != nil {
		logger.Sugar.Errorw("cms_admins get root id error", "func", util.GetSelfFuncName(), "error", err)
		return
	}
	IpWhiteListCache.SetRootID(root_id)
}
