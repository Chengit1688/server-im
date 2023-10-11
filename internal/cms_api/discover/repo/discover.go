package repo

import (
	"fmt"
	"gorm.io/gorm"
	apiUserModel "im/internal/api/user/model"
	"im/internal/cms_api/discover/model"
	"im/pkg/code"
	"im/pkg/db"
	"im/pkg/util"
)

var DiscoverRepo = new(discoverRepo)

type discoverRepo struct{}

func (r *discoverRepo) GetDiscovers() (discovers []model.Discover, err error) {
	err = db.DB.Model(model.Discover{}).Order("sort desc").Find(&discovers).Error
	return
}

func (r *discoverRepo) AddDiscover(add model.Discover) (model.Discover, error) {
	err := db.DB.Model(model.Discover{}).Create(&add).Error
	return add, err
}

func (r *discoverRepo) UpdateDiscover(id int, update model.Discover) (discover model.Discover, err error) {
	updates, _ := util.StructToMap(update, "json")
	err = db.DB.Model(model.Discover{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return
	}
	return r.GetDiscoverByID(id)
}

func (r *discoverRepo) GetDiscoverByID(id int) (discover model.Discover, err error) {
	err = db.DB.Model(model.Discover{}).Where("id = ?", id).First(&discover).Error
	return
}

func (r *discoverRepo) DeleteDiscover(id int) (err error) {
	err = db.DB.Model(model.Discover{}).Where("id = ?", id).Delete(&model.Discover{}).Error
	return
}

func (r *discoverRepo) FetchPrize(id int64) (data model.PrizeList, err error) {
	err = db.DB.Model(model.PrizeList{}).Where("id = ?", id).
		Where("status = ?", 1).
		Find(&data).Error
	return
}

func (r *discoverRepo) FetchRedeemPrizeLog(id int64) (data model.RedeemPrizeLog, err error) {
	err = db.DB.Model(model.RedeemPrizeLog{}).Where("id = ?", id).Find(&data).Error
	return
}

func (r *discoverRepo) AddPrize(add []model.PrizeList) error {
	return db.DB.Model(model.PrizeList{}).CreateInBatches(&add, len(add)).Error
}

func (r *discoverRepo) UpdatePrize(req model.UpdatePrizeReq) error {
	var prizeList model.PrizeList
	_ = util.CopyStructFields(&prizeList, &req)
	return db.DB.Model(model.PrizeList{}).Where("id = ?", req.ID).Updates(&prizeList).Error
}

func (r *discoverRepo) DeletePrize(id int64) (err error) {
	err = db.DB.Model(model.PrizeList{}).Where("id = ?", id).Updates(&model.PrizeList{Status: 2}).Error
	return
}

func (r *discoverRepo) AddRedeemPrize(userId string, data apiUserModel.RedeemPrizeReq) error {
	tx := db.DB.Begin()
	var (
		err   error
		prize model.PrizeList
	)
	if err = db.DB.Model(&model.PrizeList{}).Where("id = ?", data.PrizeID).Find(&prize).Error; err != nil {
		tx.Rollback()
		return err
	}
	if prize.ID == 0 {
		return code.ErrDataNotExists
	}
	status := 1
	if prize.IsFictitious == 2 {
		status = 21
	}
	logData := model.RedeemPrizeLog{
		UserId:       userId,
		UserName:     data.UserName,
		Address:      data.Address,
		Mobile:       data.Mobile,
		PrizeId:      data.PrizeID,
		Cost:         prize.Cost,
		Status:       int64(status),
		Name:         prize.Name,
		Icon:         prize.Icon,
		IsFictitious: prize.IsFictitious,
	}
	if err = db.DB.Model(&model.RedeemPrizeLog{}).Create(&logData).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err = db.DB.Model(&apiUserModel.User{}).Where("user_id = ?", userId).
		Where("balance > 0").
		UpdateColumn("balance", gorm.Expr(fmt.Sprintf("balance %s ?", "-"), prize.Cost)).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *discoverRepo) RedeemPrizeLog(req model.RedeemPrizeLogReq) (prizeLog []model.RedeemPrizeLog, count int64, err error) {
	req.Pagination.Check()
	tx := db.DB.Model(model.RedeemPrizeLog{})
	if len(req.UserId) != 0 {
		tx = tx.Where("user_id = ?", req.UserId)
	}
	if len(req.UserName) != 0 {
		tx = tx.Where(fmt.Sprintf("user_name like %q", ("%" + req.UserName + "%")))
	}
	if len(req.Mobile) != 0 {
		tx = tx.Where("mobile = ?", req.Mobile)
	}
	if req.PrizeId != 0 {
		tx = tx.Where("prize_id = ?", req.PrizeId)
	}
	if len(req.ExpressNumber) != 0 {
		tx = tx.Where("express_number = ?", req.ExpressNumber)
	}
	if req.Status != 0 {
		tx = tx.Where("status = ?", req.Status)
	}
	if req.StartTime != 0 {
		tx = tx.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime != 0 {
		tx = tx.Where("created_at <= ?", req.EndTime)
	}

	err = tx.Offset(req.Offset).Limit(req.Limit).Order("created_at desc").Find(&prizeLog).Limit(-1).Offset(-1).Count(&count).Error
	return
}

func (r *discoverRepo) RedeemPrizeLogForAPI(req apiUserModel.RedeemPrizeListReq) (prizeLog []model.RedeemPrizeLog, count int64, err error) {
	req.Pagination.Check()
	tx := db.DB.Model(model.RedeemPrizeLog{})
	if len(req.UserID) != 0 {
		tx = tx.Where("user_id = ?", req.UserID)
	}
	prizeData := model.PrizeList{}
	if err = db.DB.Model(&model.PrizeList{}).
		Where("name like ?", "%"+req.Key+"%").
		Where("status = ?", 1).
		Find(&prizeData).Error; err != nil {
		return nil, 0, err
	}

	if len(req.ExpressNumber) != 0 {
		tx = tx.Where("express_number = ?", req.ExpressNumber)
	}

	if req.StartTime != 0 {
		tx = tx.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime != 0 {
		tx = tx.Where("created_at <= ?", req.EndTime)
	}
	err = tx.Offset(req.Offset).Limit(req.Limit).Order("created_at desc").Find(&prizeLog).Limit(-1).Offset(-1).Count(&count).Error
	return
}

func (r *discoverRepo) PrizeList(req model.PrizeListReq) (users []model.PrizeList, count int64, err error) {
	req.Pagination.Check()
	tx := db.DB.Model(model.PrizeList{})
	if len(req.Name) != 0 {
		tx = tx.Where(fmt.Sprintf("name like %q", "%"+req.Name+"%"))
	}
	if len(req.Describe) != 0 {
		tx = tx.Where(fmt.Sprintf("describe like %q", "%"+req.Describe+"%"))
	}

	if req.IsFictitious != 0 {
		tx = tx.Where("is_fictitious = ?", req.IsFictitious)
	}

	if req.Cost != 0 {
		tx = tx.Where("cost = ?", req.Cost*100)
	}

	if req.Status == 0 {
		req.Status = 1
	}
	tx = tx.Where("status = ?", req.Status)
	if req.StartTime != 0 {
		tx = tx.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime != 0 {
		tx = tx.Where("created_at <= ?", req.EndTime)
	}

	err = tx.Offset(req.Offset).Limit(req.Limit).Order("created_at desc").Find(&users).Limit(-1).Offset(-1).Count(&count).Error
	return
}

func (r *discoverRepo) SetRedeemPrize(id int64, req model.SetRedeemPrizeReq) (prizeLog model.RedeemPrizeLog, err error) {
	_ = util.CopyStructFields(&prizeLog, &req)
	err = db.DB.Model(model.RedeemPrizeLog{}).Where("id = ?", id).Updates(&prizeLog).Error
	return
}
