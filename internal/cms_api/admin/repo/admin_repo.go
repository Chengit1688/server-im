package repo

import (
	"fmt"
	"im/internal/cms_api/admin/model"
	"im/pkg/db"
	"im/pkg/util"
	"time"
)

var AdminRepo = new(adminRepo)

type adminRepo struct{}

func (r *adminRepo) GetByUserID(userID string) (user *model.Admin, err error) {
	err = db.DB.Model(model.Admin{}).Where("user_id = ?", userID).First(&user).Error
	return
}

func (r *adminRepo) GetByUsername(username string) (user *model.Admin, err error) {
	err = db.DB.Model(model.Admin{}).Where("username = ?", username).First(&user).Error
	return
}

func (r *adminRepo) Paging(req model.ListReq) (admins []model.ListAdminItem, count int64, err error) {

	req.Pagination.Check()
	tx := db.DB.Table("cms_admins AS t1").
		Joins("JOIN cms_role AS t2 ON t1.role = t2.id")
	if len(req.Username) != 0 {
		tx = tx.Where(fmt.Sprintf("t1.username like %q", ("%" + req.Username + "%")))
	}
	if len(req.Nickname) != 0 {
		tx = tx.Where(fmt.Sprintf("nick_name like %q", ("%" + req.Nickname + "%")))
	}
	if len(req.RoleKey) != 0 {
		tx.Where("t2.role_key = ?", req.RoleKey)
	}
	if req.LoginTimeStart != 0 {
		timeLayout := "2006-01-02 15:04:05"
		timeStart := time.Unix(req.LoginTimeStart, 0).Format(timeLayout)
		tx = tx.Where("last_login_time >= ?", timeStart)
	}
	if req.LoginTimeEnd != 0 {
		timeLayout := "2006-01-02 15:04:05"

		timeEnd := time.Unix(req.LoginTimeEnd, 0).Format(timeLayout)
		tx = tx.Where("last_login_time <= ?", timeEnd)
	}
	tx.Where("t1.deleted_at is null")
	tx.Select("t1.*,t2.role_key")
	err = tx.Offset(req.Pagination.Offset).Limit(req.Pagination.Limit).Find(&admins).Limit(-1).Offset(-1).Count(&count).Error
	return
}

func (r *adminRepo) UpdateLoginInfo(user_id, ip string) (err error) {
	loginTime := time.Now()
	err = db.DB.Model(model.Admin{}).Where("user_id = ?", user_id).Updates(model.Admin{LastloginIp: ip, LastloginTime: &loginTime}).Error
	return
}

func (r *adminRepo) Add(req model.AddReq, createByUserID string) (user model.Admin, err error) {
	util.CopyStructFields(&user, &req)
	userID := util.RandID(db.UserIDSize)
	salt := userID[:8]
	user.CreateUser = createByUserID
	user.Password = util.GetPassword(req.Password, salt)
	user.UserID = util.RandID(db.UserIDSize)
	user.Salt = salt
	user.Role = req.RoleID
	err = db.DB.Model(model.Admin{}).Create(&user).Error
	return
}

func (r *adminRepo) UpdateInfo(req model.UpdateInfoReq, updateByUserID string) (user model.Admin, err error) {
	user.UpdateUser = updateByUserID
	user.Nickname = req.Nickname
	user.Username = req.Username
	user.Role = req.RoleID
	err = db.DB.Model(model.Admin{}).Where("id = ?", req.ID).Updates(&user).Error
	return
}

func (r *adminRepo) UpdatePassword(req model.UpdatePasswordReq, updateByUserID string) (user model.Admin, err error) {
	have := new(model.Admin)
	err = db.DB.Model(model.Admin{}).Where("id = ?", req.ID).First(&have).Error
	if err != nil {
		return
	}
	salt := have.Salt
	user.UpdateUser = updateByUserID
	user.Password = util.GetPassword(req.Password, salt)
	err = db.DB.Model(model.Admin{}).Where("id = ?", req.ID).Updates(&user).Error
	return
}

func (r *adminRepo) Delete(req model.DeleteReq, deleteByUserID string) (user model.Admin, err error) {
	err = db.DB.Model(model.Admin{}).Where("id = ?", req.ID).First(&user).Error
	if err != nil {
		return
	}
	err = db.DB.Model(model.Admin{}).Where("id = ?", req.ID).Update("delete_user", deleteByUserID).Error
	if err != nil {
		return
	}
	err = db.DB.Model(model.Admin{}).Where("id = ?", req.ID).Delete(&user).Error
	return
}

func (r *adminRepo) GetGoogleCodeSecret(UserID string) (user model.Admin, err error) {
	err = db.DB.Model(model.Admin{}).Where("user_id = ?", UserID).First(&user).Error
	if err != nil {
		return
	}
	if user.Google2fSecretKey == "" {
		secret := util.GetGoogleCodeSecret()
		err = db.DB.Model(model.Admin{}).Where("user_id = ?", UserID).Update("google_2f_secret_key", secret).Error
		if err != nil {
			return
		}
		user.Google2fSecretKey = secret
	}
	return
}
