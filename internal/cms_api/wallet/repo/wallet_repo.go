package repo

import (
	"fmt"
	apiUserUsecase "im/internal/api/user/usecase"
	apiWalletModel "im/internal/api/wallet/model"
	"im/internal/cms_api/wallet/model"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/util"
	"time"

	"encoding/json"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var WalletRepo = new(walletRepo)

type walletRepo struct{}

func (r *walletRepo) BillingRecordPaging(req model.BillingRecordsListReq) (records []model.BillingRecords, count int64, err error) {
	req.Pagination.Check()
	tx := db.DB.Model(model.BillingRecords{})
	if req.ID != 0 {
		tx = tx.Where("id = ?", req.ID)
	}
	if len(req.SenderID) != 0 {
		tx = tx.Where("sender_id = ?", req.SenderID)
	}
	if len(req.ReceiverID) != 0 {
		tx = tx.Where("receiver_id = ?", req.ReceiverID)
	}
	if req.CreatedTimeStart != 0 {
		tx = tx.Where("created_at >= ?", req.CreatedTimeStart)
	}
	if req.CreatedTimeEnd != 0 {
		tx = tx.Where("created_at <= ?", req.CreatedTimeEnd)
	}
	if req.Type != 0 {
		tx = tx.Where("type = ?", req.Type)
	}
	if len(req.Direction) != 0 {
		switch req.Direction {
		case "in":
			tx = tx.Where("amount > ?", 0)
		case "out":
			tx = tx.Where("amount < ?", 0)
		}
	}
	tx.Preload(clause.Associations)
	err = tx.Offset(req.Offset).Limit(req.Limit).Order("created_at desc").Find(&records).Limit(-1).Offset(-1).Count(&count).Error
	return
}

func (r *walletRepo) BillingRecordExport(req model.BillingRecordsListReq) (records []model.BillingRecords, err error) {
	tx := db.DB.Model(model.BillingRecords{})
	if req.ID != 0 {
		tx = tx.Where("id = ?", req.ID)
	}
	if len(req.SenderID) != 0 {
		tx = tx.Where("sender_id = ?", req.SenderID)
	}
	if len(req.ReceiverID) != 0 {
		tx = tx.Where("receiver_id = ?", req.ReceiverID)
	}
	if req.CreatedTimeStart != 0 {
		tx = tx.Where("created_at >= ?", req.CreatedTimeStart)
	}
	if req.CreatedTimeEnd != 0 {
		tx = tx.Where("created_at <= ?", req.CreatedTimeEnd)
	}
	if req.Type != 0 {
		tx = tx.Where("type = ?", req.Type)
	}
	if len(req.Direction) != 0 {
		switch req.Direction {
		case "in":
			tx = tx.Where("amount > ?", 0)
		case "out":
			tx = tx.Where("amount < ?", 0)
		}
	}
	tx.Preload(clause.Associations)
	err = tx.Order("created_at desc").Find(&records).Error
	return
}

func (r *walletRepo) BillingRecordPagingByUser(req apiWalletModel.BillingRecordsListReq) (records []model.BillingRecords, count int64, err error) {
	req.Pagination.Check()
	defaultRecords := `[
			{
				"amount": 100,
				"created_at": 1683373930,
				"name": 1,
				"type": 13,
                "updated_at": 1683373930
			},
			{
       			"amount": 100,
				"created_at": 1683373930,
				"name": 1,
				"type": 6,
                "updated_at": 1683373930
			},
			{
       		    "amount": 100,
				"created_at": 1683373930,
				"name": 1,
				"type": 4
                "updated_at": 1683373930,
			}
		]`
	tx := db.DB.Model(model.BillingRecords{})
	if len(req.Direction) != 0 {
		switch req.Direction {
		case "in":
			tx = tx.Where("amount > ?", 0)
		case "out":
			tx = tx.Where("amount < ?", 0)
		}

		if req.Direction == "all" {
			tx = tx.Where("type in ?", []model.TypeBillingRecord{
				model.TypeBillingRecordRecvRedPackGroup,
				model.TypeBillingRecordRecvRedPackSingle,
				model.TypeBillingRecordSendRedPackSingle,
				model.TypeBillingRecordSendRedPackGroup,
				model.TypeBillingRecordRedPackGroupReturn,
				model.TypeBillingRecordRedPackSingleReturn,
				model.TypeBillingRecordWithdrawSignAward,
				model.TypeBillingRecordWithdraw,
			})
		} else {
			tx = tx.Where("type in ?", []model.TypeBillingRecord{
				model.TypeBillingRecordRecvRedPackGroup,
				model.TypeBillingRecordRecvRedPackSingle,
				model.TypeBillingRecordSendRedPackSingle,
				model.TypeBillingRecordSendRedPackGroup,
				model.TypeBillingRecordRedPackGroupReturn,
				model.TypeBillingRecordRedPackSingleReturn,
			})
		}

	}
	tx = tx.Where("sender_id = ?", req.UserID).Or("receiver_id = ?", req.UserID)
	tx.Preload(clause.Associations)
	err = tx.Offset(req.Offset).Limit(req.Limit).Order("created_at desc").Find(&records).Limit(-1).Offset(-1).Count(&count).Error
	if records == nil || len(records) <= 0 {
		if err = json.Unmarshal([]byte(defaultRecords), &records); err != nil {
			logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("make json, error: %v", err))
			return
		}
	}
	return
}

func (r *walletRepo) AmountChange(flag string, record model.BillingRecords) (err error) {
	decimal.DivisionPrecision = 2

	err = db.DB.Transaction(func(tx *gorm.DB) error {

		err = apiUserUsecase.UserUseCase.UpdateUserWallet(record.ReceiverID, flag, record.Amount)
		if err != nil {
			return err
		}
		err = tx.Create(&record).Error
		if err != nil {
			return err
		}
		return err
	})
	return err
}

func (r *walletRepo) RedpackSingleRecordPaging(req model.RedpackSingleRecordsListReq) (records []model.RedpackSingleRecords, count int64, err error) {
	req.Pagination.Check()
	tx := db.DB.Model(model.RedpackSingleRecords{})
	if len(req.SenderID) != 0 {
		tx = tx.Where("sender_id = ?", req.SenderID)
	}
	if len(req.ReceiverID) != 0 {
		tx = tx.Where("receiver_id = ?", req.ReceiverID)
	}
	if req.SendTimeStart != 0 {
		tx = tx.Where("send_at >= ?", req.SendTimeStart)
	}
	if req.SendTimeEnd != 0 {
		tx = tx.Where("send_at <= ?", req.SendTimeEnd)
	}
	if req.RecvTimeStart != 0 {
		tx = tx.Where("recv_at >= ?", req.RecvTimeStart)
	}
	if req.RecvTimeEnd != 0 {
		tx = tx.Where("recv_at <= ?", req.RecvTimeEnd)
	}
	if req.Status != 0 {
		tx = tx.Where("status = ?", req.Status)
	}
	tx.Preload(clause.Associations)
	err = tx.Offset(req.Offset).Limit(req.Limit).Order("id desc").Find(&records).Limit(-1).Offset(-1).Count(&count).Error
	return
}

func (r *walletRepo) RedpackSingleRecordExport(req model.RedpackSingleRecordsListReq) (records []model.RedpackSingleRecords, err error) {
	tx := db.DB.Model(model.RedpackSingleRecords{})
	if len(req.SenderID) != 0 {
		tx = tx.Where("sender_id = ?", req.SenderID)
	}
	if len(req.ReceiverID) != 0 {
		tx = tx.Where("receiver_id = ?", req.ReceiverID)
	}
	if req.SendTimeStart != 0 {
		tx = tx.Where("send_at >= ?", req.SendTimeStart)
	}
	if req.SendTimeEnd != 0 {
		tx = tx.Where("send_at <= ?", req.SendTimeEnd)
	}
	if req.RecvTimeStart != 0 {
		tx = tx.Where("recv_at >= ?", req.RecvTimeStart)
	}
	if req.RecvTimeEnd != 0 {
		tx = tx.Where("recv_at <= ?", req.RecvTimeEnd)
	}
	if req.Status != 0 {
		tx = tx.Where("status = ?", req.Status)
	}
	tx.Preload(clause.Associations)
	err = tx.Order("id desc").Find(&records).Error
	return
}

func (r *walletRepo) RedpackSingleSend(req apiWalletModel.RedpackSingleSendReq, balance int64) (id int64, err error) {
	timeInt64 := time.Now().Unix()
	err = db.DB.Transaction(func(tx *gorm.DB) error {

		err = tx.Table("users").Where("user_id = ? ", req.SendID).UpdateColumn("balance", gorm.Expr(fmt.Sprintf("balance %s ?", "-"), req.Amount)).Error
		if err != nil {
			return err
		}

		redpack := model.RedpackSingleRecords{
			SendAt:     timeInt64,
			SenderID:   req.SendID,
			ReceiverID: req.RecvID,
			Status:     model.StatusRedpackSingleRecv,
			Amount:     req.Amount,
			Remark:     req.Remark,
			MsgType:    req.MsgType,
		}
		err1 := tx.Create(&redpack).Error
		if err1 != nil {
			return err1
		}
		id = redpack.ID

		billing := model.BillingRecords{
			Type:         model.TypeBillingRecordSendRedPackSingle,
			SenderID:     req.SendID,
			ReceiverID:   req.RecvID,
			Amount:       req.Amount,
			ChangeBefore: balance,
			ChangeAfter:  balance - req.Amount,
			Note:         "发送个人红包",
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

func (r *walletRepo) RedpackSingleGet(recordID int64) (redpack model.RedpackSingleRecords, err error) {
	err = db.DB.Model(model.RedpackSingleRecords{}).First(&redpack, recordID).Error
	return
}

func (r *walletRepo) RedpackSingleRecv(redpack model.RedpackSingleRecords, balance int64) (timeInt64 int64, err error) {
	timeInt64 = time.Now().Unix()
	redpack.RecvAt = &timeInt64
	redpack.Status = model.StatusRedpackSingleRecvd
	err = db.DB.Transaction(func(tx *gorm.DB) error {

		err = tx.Table("users").Where("user_id = ? ", redpack.ReceiverID).UpdateColumn("balance", gorm.Expr(fmt.Sprintf("balance %s ?", "+"), redpack.Amount)).Error
		if err != nil {
			return err
		}

		tx.Save(&redpack)
		if err != nil {
			return err
		}

		billing := model.BillingRecords{
			Type:         model.TypeBillingRecordRecvRedPackSingle,
			SenderID:     redpack.SenderID,
			ReceiverID:   redpack.ReceiverID,
			Amount:       redpack.Amount,
			ChangeBefore: balance,
			ChangeAfter:  balance + redpack.Amount,
			Note:         "领取个人红包",
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

func (r *walletRepo) RedpackSingleNeedReturns() (records []model.RedpackSingleRecords, err error) {
	timeNow := time.Now()
	var data []model.RedpackSingleRecords
	err = db.DB.Model(model.RedpackSingleRecords{}).Where("status = ?", model.StatusRedpackSingleRecv).Find(&data).Error
	if err != nil {
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

func (r *walletRepo) RedpackSingleReturn(redpack model.RedpackSingleRecords) (err error) {
	user, _ := apiUserUsecase.UserUseCase.GetInfo(redpack.SenderID)
	err = db.DB.Transaction(func(tx *gorm.DB) error {

		err = tx.Table("users").Where("user_id = ? ", redpack.SenderID).UpdateColumn("balance", gorm.Expr(fmt.Sprintf("balance %s ?", "+"), redpack.Amount)).Error
		if err != nil {
			return err
		}
		redpack.Status = model.StatusRedpackSingleReturn

		tx.Save(&redpack)
		if err != nil {
			return err
		}

		billing := model.BillingRecords{
			Type:         model.TypeBillingRecordRedPackSingleReturn,
			SenderID:     redpack.SenderID,
			ReceiverID:   redpack.ReceiverID,
			Amount:       redpack.Amount,
			ChangeBefore: user.Balance,
			ChangeAfter:  user.Balance + redpack.Amount,
			Note:         "退回未领取的个人红包",
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

func (r *walletRepo) WithdrawRecordAdd(add model.WithdrawRecords) (record model.WithdrawRecords, err error) {
	user, _ := apiUserUsecase.UserUseCase.GetInfo(add.UserID)
	err = db.DB.Transaction(func(tx *gorm.DB) error {

		err = tx.Table("users").Where("user_id = ? ", add.UserID).UpdateColumn("balance", gorm.Expr(fmt.Sprintf("balance %s ?", "-"), add.Amount)).Error
		if err != nil {
			return err
		}

		timeInt64 := time.Now().Unix()
		add.BillingID = tx.Table("billing_records").Create(map[string]interface{}{
			"type":          model.TypeBillingRecordWithdraw,
			"sender_id":     add.UserID,
			"receiver_id":   add.UserID,
			"amount":        add.Amount,
			"change_before": user.Balance,
			"change_after":  user.Balance - add.Amount,
			"note":          "签到奖励",
			"created_at":    timeInt64,
			"updated_at":    timeInt64,
		}).RowsAffected
		err = tx.Create(&add).Error
		record = add
		return nil
	})
	return
}

func (r *walletRepo) WithdrawRecordPaging(req model.WithdrawRecordsListReq) (records []model.WithdrawRecords, count int64, err error) {
	req.Pagination.Check()
	tx := db.DB.Model(model.WithdrawRecords{})
	if len(req.UserID) != 0 {
		tx = tx.Where("user_id = ?", req.UserID)
	}
	if len(req.NickName) != 0 {
		tx.Preload("User", fmt.Sprintf("nick_name like %q", ("%"+req.NickName+"%")))
	} else {
		tx.Preload(clause.Associations)
	}
	if req.BillingID != 0 {
		tx = tx.Where("billing_id = ?", req.BillingID)
	}
	if req.CreatedTimeStart != 0 {
		tx = tx.Where("created_at >= ?", req.CreatedTimeStart)
	}
	if req.CreatedTimeEnd != 0 {
		tx = tx.Where("created_at <= ?", req.CreatedTimeEnd)
	}
	if req.Status != 0 {
		tx = tx.Where("status = ?", req.Status)
	}
	if req.IsDone == "yes" {
		tx = tx.Where("status not null")
	} else if req.IsDone == "no" {
		tx = tx.Where("status is null")
	}
	err = tx.Offset(req.Offset).Limit(req.Limit).Order("id desc").Find(&records).Limit(-1).Offset(-1).Count(&count).Error
	return
}

func (r *walletRepo) WithdrawRecordCountPending() (count int64, err error) {
	err = db.DB.Model(model.WithdrawRecords{}).Where("status is null").Count(&count).Error
	return
}

func (r *walletRepo) WithdrawRecordDescribeByID(id string) (record model.WithdrawRecords, err error) {
	err = db.DB.Model(model.WithdrawRecords{}).Where("id = ?", id).First(&record).Error
	return
}

func (r *walletRepo) WithdrawRecordStatusSet(params model.SetWithdrawRecordsStatusReq) (err error) {
	record := model.WithdrawRecords{Status: &params.Status, Note: params.Note}
	record.UpdatedAt = time.Now().UnixMilli()
	err = db.DB.Model(model.WithdrawRecords{}).Where("id = ?", params.ID).UpdateColumns(record).Error
	return
}
