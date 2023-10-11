package repo

import (
	"errors"
	"im/internal/control/menu/model"
	"im/pkg/db"
	"im/pkg/util"
	"time"

	"gorm.io/gorm"
)

var MenuRepo = new(menuRepo)

const CONFIG_MENU_TIMESTAMP = "CONFIG_MENU_TIMESTAMP"

type menuRepo struct{}

func (r *menuRepo) MenuList(name, title string) (menus []model.CmsMenu, count int64, err error) {

	tx := db.DB.Model(model.CmsMenu{})
	if len(name) != 0 {
		tx = tx.Where("name LIKE  ?", name+"%")
	}
	if len(title) != 0 {
		tx = tx.Where("title LIKE ?", title+"%")
	}
	err = tx.Find(&menus).Count(&count).Error
	return
}

func (r *menuRepo) MenuAdd(req model.AddMenuReq) (id int, err error) {

	add := new(model.CmsMenu)
	util.CopyStructFields(add, req)
	db.DB.Transaction(func(tx *gorm.DB) error {

		if err = tx.Create(&add).Find(&add).Error; err != nil {
			return err
		}
		id = int(add.ID)
		for _, item := range req.Apis {

			addRule := new(model.CmsMenuApi)
			addRule.MenuID = id
			addRule.ApiID = item
			if err := tx.Create(&addRule).Error; err != nil {
				return err
			}
		}
		return nil
	})

	r.MenuUpdateConfigTime()
	return
}

func (r *menuRepo) MenuUpdate(id string, req model.UpdateMenuReq) (err error) {

	update := new(model.CmsMenu)
	have := new(model.CmsMenu)
	util.CopyStructFields(update, req)
	updates, err := util.StructToMap(update, "json")
	delete(updates, "apis")

	if err != nil {
		return
	}
	db.DB.Transaction(func(tx *gorm.DB) error {

		if err = tx.Model(&have).Where("id = ?", id).First(&have).Error; err != nil {
			return err
		}

		if err = tx.Model(&have).Updates(updates).Error; err != nil {
			return err
		}

		menuApis, err := r.GetMenuApisByMenuID(id)
		if err != nil {
			return err
		}
		var haveApi []int

		for _, item := range menuApis {
			haveApi = append(haveApi, item.ApiID)
		}
		needAdd := util.Difference(req.Apis, haveApi)
		needDelete := util.Difference(haveApi, req.Apis)

		for _, addID := range needAdd {
			add := new(model.CmsMenuApi)
			add.MenuID = util.StringToInt(id)
			add.ApiID = addID
			if err := tx.Create(&add).Error; err != nil {
				return err
			}
		}

		for _, deleteID := range needDelete {
			if err := tx.Where("menu_id = ?", id).Where("api_id = ?", deleteID).Delete(&model.CmsMenuApi{}).Error; err != nil {
				return err
			}
		}
		return nil
	})

	r.MenuUpdateConfigTime()
	return
}

func (r *menuRepo) GetMenuApisByMenuID(id string) (menuApis []model.CmsMenuApi, err error) {
	err = db.DB.Model(model.CmsMenuApi{}).Where("menu_id = ?", id).Find(&menuApis).Error
	return
}

func (r *menuRepo) MenuDelete(id string) (err error) {

	have := new(model.CmsMenu)
	db.DB.Transaction(func(tx *gorm.DB) error {

		if err = tx.Model(&have).Where("id = ?", id).First(&have).Error; err != nil {

			return err
		}

		if err = tx.Model(&have).Delete(&have).Error; err != nil {

			return err
		}

		if err = tx.Where("menu_id = ?", id).Delete(&model.CmsMenuApi{}).Error; err != nil {
			return err
		}
		return nil
	})

	r.MenuUpdateConfigTime()
	return
}

func (r *menuRepo) MenuGet(id string) (menu model.CmsMenu, err error) {

	err = db.DB.Model(model.CmsMenu{}).Where("id = ?", id).Preload("Apis").First(&menu).Error
	return
}

func (r *menuRepo) MenuGetConfigTime() (timeString int64, err error) {

	config := new(model.Config)
	err = db.DB.Model(model.Config{}).Where("name = ?", CONFIG_MENU_TIMESTAMP).First(&config).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		timeString = time.Now().UnixNano()
		config.Name = CONFIG_MENU_TIMESTAMP
		config.Value = util.Int64ToString(timeString)
		db.DB.Save(&config)
	}
	timeString = util.StringToInt64(config.Value)
	return
}
func (r *menuRepo) MenuUpdateConfigTime() (err error) {

	timeString := time.Now().UnixNano()
	err = db.DB.Model(model.Config{}).Where("name = ?", CONFIG_MENU_TIMESTAMP).Update("value", timeString).Error
	return
}

func (r *menuRepo) GetAllMenu() (menus []model.CmsMenu, err error) {
	err = db.DB.Raw("select * from cms_menus").Scan(&menus).Error
	return
}
