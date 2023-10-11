package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"im/internal/cms_api/role/model"
	"im/pkg/db"
	"im/pkg/util"

	"gorm.io/gorm"
)

var RoleRepo = new(roleRepo)

type roleRepo struct{}

func (r *roleRepo) GetRoleKeyByRoleID(id int) (role *model.CmsRole, err error) {
	err = db.DB.Model(model.CmsRole{}).Where("id = ?", id).First(&role).Error
	return
}

func (r *roleRepo) GetRoleMenu(roleKey string) (role model.CmsRole, err error) {

	err = db.DB.Model(model.CmsRole{}).Where("role_key = ?", roleKey).Preload("CmsMenu", db.DB.Where(&model.CmsMenu{Visible: 1})).First(&role).Error
	return
}

func (r *roleRepo) GetAdminMenu() (menus []model.CmsMenu, err error) {

	err = db.DB.Model(model.CmsMenu{}).Where("visible = ?", 1).Find(&menus).Error
	return
}

func (r *roleRepo) RolePaging(req model.RoleListReq) (roles []model.CmsRole, count int64, err error) {

	req.Pagination.Check()
	tx := db.DB.Debug().Model(model.CmsRole{})
	if len(req.RoleName) != 0 {
		tx = tx.Where(fmt.Sprintf("role_name like %q", ("%" + req.RoleName + "%")))
	}
	if len(req.RoleKey) != 0 {
		tx = tx.Where(fmt.Sprintf("role_key like %q", ("%" + req.RoleKey + "%")))
	}
	err = tx.Offset(req.Offset).Limit(req.Limit).Find(&roles).Limit(-1).Offset(-1).Count(&count).Error
	return
}

func (r *roleRepo) RoleAdd(req model.RoleAddReq) (id int, err error) {

	add := new(model.CmsRole)
	util.CopyStructFields(add, req)
	db.DB.Transaction(func(tx *gorm.DB) error {

		if err = tx.Create(&add).Find(&add).Error; err != nil {
			return err
		}
		id = int(add.ID)
		for _, item := range req.Menus {

			addRule := new(model.CmsMenuRole)
			addRule.MenuID = item
			addRule.RoleID = int(add.ID)
			if err = tx.Create(&addRule).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return
}

func (r *roleRepo) RoleUpdate(id string, req model.UpdateRoleReq) (err error) {

	update := new(model.CmsRole)
	have := new(model.CmsRole)
	util.CopyStructFields(update, req)
	updates, err := util.StructToMap(update, "json")
	delete(updates, "cms_menu")
	delete(updates, "menu_ids")
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

		menuRoles, err := r.GetRoleMenusByRoleID(id)
		if err != nil {
			return err
		}
		var haveMenu []int

		for _, item := range menuRoles {
			haveMenu = append(haveMenu, item.MenuID)
		}
		needAdd := util.Difference(req.Menus, haveMenu)
		needDelete := util.Difference(haveMenu, req.Menus)

		for _, addID := range needAdd {
			add := new(model.CmsMenuRole)
			add.RoleID = util.StringToInt(id)
			add.MenuID = addID
			if err := tx.Create(&add).Error; err != nil {
				return err
			}

		}

		for _, deleteID := range needDelete {
			if err = tx.Where("menu_id = ?", deleteID).Where("role_id = ?", id).Delete(&model.CmsMenuRole{}).Error; err != nil {
				return err
			}

		}
		return nil
	})
	return
}

func (r *roleRepo) GetRoleMenusByRoleID(id string) (menuRoles []model.CmsMenuRole, err error) {
	err = db.DB.Model(model.CmsMenuRole{}).Where("role_id = ?", id).Find(&menuRoles).Error
	return
}

func (r *roleRepo) RoleDelete(id string) (err error) {

	have := new(model.CmsRole)
	db.DB.Transaction(func(tx *gorm.DB) error {

		if err = tx.Model(&have).Where("id = ?", id).First(&have).Error; err != nil {

			return err
		}
		roleKey := have.RoleKey

		if err = tx.Where("role_id = ?", id).Delete(&model.CmsMenuRole{}).Error; err != nil {
			return err
		}

		if err = tx.Model(&have).Delete(&have).Error; err != nil {

			return err
		}

		if err = tx.Exec("DELETE FROM `cms_casbin_rule` where v0 = @role_key", sql.Named("role_key", roleKey)).Error; err != nil {
			return err
		}
		return nil
	})
	return
}

func (r *roleRepo) RoleGet(id string) (role model.CmsRole, err error) {

	err = db.DB.Model(model.CmsRole{}).Where("id = ?", id).Preload("CmsMenu", db.DB.Where(&model.CmsMenu{Visible: 1})).First(&role).Error
	return
}

func (r *roleRepo) RoleGetByName(name string, id string) (role model.CmsRole, err error) {

	tx := db.DB.Model(model.CmsRole{})
	if len(id) != 0 {
		tx.Where("id != ?", id)
	}
	err = tx.Where("role_name = ?", name).First(&role).Error
	return
}
func (r *roleRepo) RoleGetByKey(key string, id string) (role model.CmsRole, err error) {

	tx := db.DB.Model(model.CmsRole{})
	if len(id) != 0 {
		tx.Where("id != ?", id)
	}
	err = tx.Where("role_key = ?", key).First(&role).Error
	return
}

func (r *roleRepo) SyncMenu(menus []model.CmsMenu) (err error) {
	db.DB.Transaction(func(tx *gorm.DB) error {

		for _, menu := range menus {
			have := new(model.CmsMenu)
			err = tx.Raw("select * from cms_menus where id = ?", menu.ID).Scan(&have).Error
			if errors.Is(err, gorm.ErrRecordNotFound) || have.ID == 0 {

				if err = tx.Model(&have).Create(&menu).Error; err != nil {
					return err
				}
			} else {
				if err != nil {

					return err
				} else {

					updates, err := util.StructToMap(menu, "json")
					delete(updates, "admin_api")
					delete(updates, "apis")
					delete(updates, "params")
					delete(updates, "children")
					delete(updates, "is_select")
					updates["updated_at"] = menu.UpdatedAt
					updates["created_at"] = menu.CreatedAt
					updates["deleted_at"] = menu.DeletedAt
					if err = tx.Model(&have).Where("id = ?", have.ID).Updates(updates).Error; err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
	return
}
