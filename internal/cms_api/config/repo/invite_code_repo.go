package repo

import (
	configModel "im/internal/cms_api/config/model"
	"im/pkg/code"
	"im/pkg/common/constant"
	"im/pkg/db"
	"im/pkg/util"
	"strings"
)

var InviteCode = new(inviteCode)

type inviteCode struct{}

type WhereOptionForInvite struct {
	Id            int64
	Ids           []int64
	InviteCode    string
	Status        int64
	DeleteStatus  int64
	BeginDate     int64
	EndDate       int64
	Ext           string
	Remarks       string
	OperationUser string
}

func (r *inviteCode) Create(data *configModel.InviteCode) (*configModel.InviteCode, error) {
	err := db.DB.Create(&data).Error
	return data, err
}

func (r *inviteCode) GetInfo(opts ...WhereOptionForInvite) (*configModel.InviteCode, error) {
	var m configModel.InviteCode
	var opt WhereOptionForInvite
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		return nil, code.ErrSettingNotExist
	}
	query := db.DB
	if opt.Id != 0 {
		query = query.Where("id = ? ", opt.Id)
	}
	if opt.InviteCode != "" {
		query = query.Where(" invite_code = ?", opt.InviteCode)
	}
	if opt.Status != 0 {
		query = query.Where("status = ?", opt.Status)
	}
	if opt.Ext != "" {
		query = query.Where(opt.Ext)
	}
	if err := query.First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *inviteCode) Exists(opts ...WhereOptionForInvite) (int64, error) {
	var (
		opt   WhereOptionForInvite
		total int64
		err   error
	)
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		return 0, code.ErrSettingNotExist
	}
	query := db.DB
	if opt.Id != 0 {
		query = query.Where("id = ? ", opt.Id)
	}
	if opt.InviteCode != "" {
		query = query.Where("invite_code = ?", opt.InviteCode)
	}

	if opt.Status != 0 {
		query = query.Where("status = ?", opt.Status)
	}
	if opt.Ext != "" {
		query = query.Where(opt.Ext)
	}

	err = query.First(&configModel.InviteCode{}).Count(&total).Error

	return total, err
}

func (r *inviteCode) UpdateById(opt WhereOptionForInvite, data *configModel.InviteCode) error {
	var err error
	if opt.Id == 0 && len(opt.Ids) == 0 {
		return code.ErrBadRequest
	}
	query := db.DB.Model(&configModel.InviteCode{})
	if opt.Id != 0 {
		query = query.Where("id = ?", opt.Id)
	}
	if len(opt.Ids) != 0 {
		query = query.Where("id IN ?", opt.Ids)
	}
	if err = query.UpdateColumns(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r *inviteCode) UpdateInviteById(opt WhereOptionForInvite, req configModel.InviteUpdateReq) error {
	var err error
	if opt.Id == 0 && len(opt.Ids) == 0 {
		return code.ErrBadRequest
	}

	updates, err := util.StructToMapWithoutNil(req, "json")
	delete(updates, "operation_id")
	delete(updates, "id")

	query := db.DB.Model(&configModel.InviteCode{})
	if opt.Id != 0 {
		query = query.Where("id = ?", opt.Id)
	}
	if len(opt.Ids) != 0 {
		query = query.Where("id IN ?", opt.Ids)
	}
	if err = query.Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

func (r *inviteCode) GetList(req configModel.InviteListReq) (list []configModel.InviteCode, count int64, err error) {
	req.Pagination.Check()
	query := db.DB.Model(configModel.InviteCode{})
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}
	if req.InviteCode != "" {
		query = query.Where("invite_code = ?", req.InviteCode)
	}
	if req.DefaultFriends != "" {
		var fs []string
		var fv []interface{}
		for _, v := range strings.Split(req.DefaultFriends, ",") {
			fs = append(fs, "find_in_set(?,`default_friends`)")
			fv = append(fv, v)
		}
		query = query.Where(strings.Join(fs, " OR "), fv...)
	}
	if req.DefaultGroups != "" {
		var ds []string
		var dv []interface{}
		for _, v := range strings.Split(req.DefaultGroups, ",") {
			ds = append(ds, "find_in_set(?,`default_groups`)")
			dv = append(dv, v)
		}
		query = query.Where(strings.Join(ds, " OR "), dv...)
	}
	if req.Remarks != "" {
		query = query.Where("remarks like ?", "%"+req.Remarks+"%")
	}
	if req.OperationUser != "" {
		query = query.Where("operation_user = ?", req.OperationUser)
	}
	if req.BeginDate != 0 {
		query = query.Where("created_at >= ?", req.BeginDate)
	}
	if req.EndDate != 0 {
		query = query.Where("created_at <= ?", req.EndDate)
	}
	if req.Limit == 0 {
		req.Limit = 5
	}
	query = query.Where("delete_status != ?", constant.SwitchOn)

	err = query.Order("created_at asc").Offset(req.Offset).Limit(req.Limit).Find(&list).Error
	_ = query.Offset(-1).Limit(-1).Count(&count).Error
	return
}

func (r *inviteCode) UpdateMapInfo(opt WhereOptionForInvite, data map[string]interface{}) error {
	var err error
	if opt.Id == 0 && len(opt.Ids) == 0 {
		return code.ErrBadRequest
	}
	query := db.DB.Model(&configModel.InviteCode{})
	if opt.Id != 0 {
		query = query.Where("id = ?", opt.Id)
	}
	if len(opt.Ids) != 0 {
		query = query.Where("id IN ?", opt.Ids)
	}
	if err = query.UpdateColumns(data).Error; err != nil {
		return err
	}
	return nil
}
