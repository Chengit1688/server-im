package repo

import (
	"fmt"
	apiGroupModel "im/internal/api/group/model"
	apiUserUsecase "im/internal/api/user/usecase"
	apiWalletModel "im/internal/api/wallet/model"
	"im/internal/cms_api/wallet/model"
	"im/pkg/db"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *walletRepo) RedpackGroupRecordPaging(req model.RedpackGroupRecordsListReq) (records []model.RedpackGroupRecords, count int64, err error) {
	req.Pagination.Check()
	tx := db.DB.Model(model.RedpackGroupRecords{})
	if len(req.SenderID) != 0 {
		tx = tx.Where("sender_id = ?", req.SenderID)
	}
	if req.GroupName != "" {
		group := apiGroupModel.Group{}
		if err = db.DB.Model(apiGroupModel.Group{}).Select("group_id").
			Where("name = ?", req.GroupName).First(&group).Error; err != nil {
			return nil, 0, err
		}
		tx = tx.Where("group_id = ?", group.GroupId)
	}
	if req.SendTimeStart != 0 {
		tx = tx.Where("send_at >= ?", req.SendTimeStart)
	}
	if req.SendTimeEnd != 0 {
		tx = tx.Where("send_at <= ?", req.SendTimeEnd)
	}
	if req.Status != 0 {
		tx = tx.Where("status = ?", req.Status)
	}
	tx.Preload(clause.Associations)

	err = tx.Offset(req.Offset).Limit(req.Limit).
		Order("id desc").Find(&records).Limit(-1).Offset(-1).Count(&count).Error
	return
}

func (r *walletRepo) RedpackGroupRecordExport(req model.RedpackGroupRecordsListReq) (records []model.RedpackGroupRecords, err error) {
	tx := db.DB.Model(model.RedpackGroupRecords{})
	if len(req.SenderID) != 0 {
		tx = tx.Where("sender_id = ?", req.SenderID)
	}
	if req.GroupName != "" {
		group := apiGroupModel.Group{}
		if err = db.DB.Model(apiGroupModel.Group{}).Select("group_id").
			Where("name = ?", req.GroupName).First(&group).Error; err != nil {
			return
		}
		tx = tx.Where("group_id = ?", group.GroupId)
	}
	if req.SendTimeStart != 0 {
		tx = tx.Where("send_at >= ?", req.SendTimeStart)
	}
	if req.SendTimeEnd != 0 {
		tx = tx.Where("send_at <= ?", req.SendTimeEnd)
	}
	if req.Status != 0 {
		tx = tx.Where("status = ?", req.Status)
	}
	tx.Preload(clause.Associations)
	err = tx.Order("id desc").Find(&records).Error
	return
}

func (r *walletRepo) RedpackGroupSend(req apiWalletModel.RedpackGroupSendReq, balance int64) (id int64, err error) {
	timeInt64 := time.Now().Unix()
	err = db.DB.Transaction(func(tx *gorm.DB) error {

		err = tx.Table("users").Where("user_id = ? ", req.SendID).UpdateColumn("balance", gorm.Expr(fmt.Sprintf("balance %s ?", "-"), req.Amount)).Error
		if err != nil {
			return err
		}

		redpackRecordes := model.RedpackGroupRecords{
			SendAt:          timeInt64,
			SenderID:        req.SendID,
			GroupID:         req.GroupID,
			Count:           req.Count,
			Type:            req.Type,
			Amount:          req.Amount,
			RemainderAmount: req.Amount,
			ReceiveCount:    0,
			Remark:          req.Remark,
			MsgType:         req.MsgType,
		}
		err = tx.Create(&redpackRecordes).Error
		if err != nil {
			return err
		}
		id = redpackRecordes.ID

		billing := model.BillingRecords{
			Type:           model.TypeBillingRecordSendRedPackSingle,
			SenderID:       req.SendID,
			Amount:         req.Amount,
			GroupID:        req.GroupID,
			RedpackGroupID: redpackRecordes.ID,
			ChangeBefore:   balance,
			ChangeAfter:    balance - req.Amount,
			Note:           "发送群红包",
		}
		billing.CreatedAt = timeInt64
		billing.UpdatedAt = timeInt64
		err = tx.Create(&billing).Error
		if err != nil {
			return err
		}
		return nil
	})
	return id, err
}

func (r *walletRepo) RedpackGroupGet(recordID int64) (model.RedpackGroupRecords, error) {
	redpack := model.RedpackGroupRecords{}
	err := db.DB.Model(model.RedpackGroupRecords{}).Where("id = ?", recordID).First(&redpack).Error
	return redpack, err
}

func (r *walletRepo) RedpackGroupRecvGet(userId string) (model.RedpackGroupRecvs, error) {
	redpack := model.RedpackGroupRecvs{}
	err := db.DB.Model(model.RedpackGroupRecvs{}).Where("user_id = ?", userId).First(&redpack).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.RedpackGroupRecvs{}, nil
		}
		return model.RedpackGroupRecvs{}, err
	}
	return redpack, nil
}
func (r *walletRepo) UpdateRedpackGroupStatus(redpackGroupId int64, status int) (err error) {
	err = db.DB.Model(model.RedpackGroupRecords{}).Where("id = ?", redpackGroupId).UpdateColumn("status", status).Error
	return
}

func (r *walletRepo) RedpackGroupRecv(redpack model.RedpackGroupRecvs, balance int64) (timeInt64 int64, err error) {
	timeInt64 = time.Now().Unix()
	redpack.RecvAt = timeInt64
	err = db.DB.Transaction(func(tx *gorm.DB) error {

		err = tx.Table("users").Where("user_id = ? ", redpack.UserID).
			UpdateColumn("balance", gorm.Expr(fmt.Sprintf("balance %s ?", "+"), redpack.Amount)).Error
		if err != nil {
			return err
		}

		if err = tx.Table("redpack_group_records").Where("id = ?", redpack.RedpackGroupID).
			Updates(map[string]interface{}{"remainder_amount": gorm.Expr("remainder_amount - ?", redpack.Amount),
				"receive_count": gorm.Expr("receive_count + ?", 1)}).Error; err != nil {
			return err
		}

		tx.Save(&redpack)
		if err != nil {
			return err
		}

		billing := model.BillingRecords{
			Type:           model.TypeBillingRecordRecvRedPackSingle,
			SenderID:       redpack.SenderID,
			ReceiverID:     redpack.UserID,
			Amount:         redpack.Amount,
			ChangeBefore:   balance,
			GroupID:        redpack.RedpackGroup.GroupID,
			RedpackGroupID: redpack.RedpackGroupID,
			ChangeAfter:    balance + redpack.Amount,
			Note:           "领取群红包",
		}
		billing.CreatedAt = timeInt64
		billing.UpdatedAt = timeInt64
		err = tx.Create(&billing).Error
		if err != nil {
			return err
		}

		return nil
	})
	return timeInt64, err
}

func (r *walletRepo) RedpackGroupNeedReturns() (records []model.RedpackGroupReturnsRecords, err error) {
	var (
		data    []model.RedpackGroupReturnsRecords
		timeNow = time.Now()
	)
	columns := "id as group_records_id, send_at,sender_id,group_id,amount,remainder_amount"
	if err = db.DB.Model(model.RedpackGroupRecords{}).Select(columns).Where("status != ?", model.StatusRedpackSingleReturn).Find(&data).Error; err != nil {
		return
	}
	for _, item := range data {

		dbTime := time.Unix(item.SendAt, 0)

		if timeNow.Sub(dbTime).Hours() > 24 {
			records = append(records, item)
		}
	}
	return
}

func (r *walletRepo) RedpackGroupReturn(redpack model.RedpackGroupRecvs) (err error) {
	user, _ := apiUserUsecase.UserUseCase.GetInfo(redpack.SenderID)
	err = db.DB.Transaction(func(tx *gorm.DB) error {

		err = tx.Table("users").Where("user_id = ? ", redpack.SenderID).UpdateColumn("balance", gorm.Expr(fmt.Sprintf("balance %s ?", "+"), redpack.Amount)).Error
		if err != nil {
			return err
		}

		if err = tx.Model(model.RedpackGroupRecords{}).Where("id = ?", redpack.RedpackGroupID).
			UpdateColumn("status", model.StatusRedpackSingleReturn).Error; err != nil {
			return err
		}

		billing := model.BillingRecords{
			Type:           model.TypeBillingRecordRedPackSingleReturn,
			SenderID:       redpack.SenderID,
			ReceiverID:     redpack.UserID,
			Amount:         redpack.Amount,
			ChangeBefore:   user.Balance,
			ChangeAfter:    user.Balance + redpack.Amount,
			GroupID:        redpack.RedpackGroup.GroupID,
			RedpackGroupID: redpack.RedpackGroupID,
			Note:           "退回未领取的群红包",
		}
		timeInt64 := time.Now().Unix()
		billing.CreatedAt = timeInt64
		billing.UpdatedAt = timeInt64
		err = tx.Create(&billing).Error
		if err != nil {
			return err
		}

		return nil
	})
	return
}
