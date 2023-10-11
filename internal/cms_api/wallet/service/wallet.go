package service

import (
	"fmt"
	userModel "im/internal/api/user/model"
	userRepo "im/internal/api/user/repo"
	userUseCase "im/internal/api/user/usecase"
	"im/internal/cms_api/wallet/model"
	"im/internal/cms_api/wallet/repo"
	"im/pkg/code"
	"im/pkg/db"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"
	http2 "net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

var WalletService = new(walletService)

type walletService struct{}

func (s *walletService) BillingRecordsList(c *gin.Context) {
	req := new(model.BillingRecordsListReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	records, count, err := repo.WalletRepo.BillingRecordPaging(*req)

	ret := new(model.BillingRecordsListResp)
	util.CopyStructFields(&ret.List, &records)
	for index := range ret.List {
		if ret.List[index].Type == model.TypeBillingRecordCmsChange {
			ret.List[index].SenderNickName = records[index].SenderCmsUser.Nickname
		} else {
			ret.List[index].SenderNickName = records[index].SenderApiUser.NickName
		}
		ret.List[index].ReceiverNickName = records[index].Receiver.NickName
		ret.List[index].Amount = records[index].Amount
		ret.List[index].ChangeBefore = records[index].ChangeBefore
		ret.List[index].ChangeAfter = records[index].ChangeAfter
	}
	ret.Count = count
	ret.Page = req.Page
	ret.PageSize = req.PageSize
	http.Success(c, ret)
}

func (s *walletService) BillingRecordsExport(c *gin.Context) {
	req := new(model.BillingRecordsListReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	records, err := repo.WalletRepo.BillingRecordExport(*req)

	ret := new(model.BillingRecordsListResp)
	util.CopyStructFields(&ret.List, &records)

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		}
	}()

	index, err := f.NewSheet("Sheet1")
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}

	err = f.SetColWidth("Sheet1", "A", "K", 20)

	sheetHeader := []interface{}{"账单ID", "发送者ID", "发送者昵称", "接收者ID", "接收者昵称", "类型", "金额", "账变前", "账变后", "备注", "完成时间"}
	cell, err := excelize.CoordinatesToCellName(1, 1)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}
	f.SetSheetRow("Sheet1", cell, &sheetHeader)
	var row []interface{}
	var recordType string
	var finishTime string
	for index = range ret.List {
		if ret.List[index].Type == model.TypeBillingRecordCmsChange {
			ret.List[index].SenderNickName = records[index].SenderCmsUser.Nickname
		} else {
			ret.List[index].SenderNickName = records[index].SenderApiUser.NickName
		}
		ret.List[index].ReceiverNickName = records[index].Receiver.NickName
		ret.List[index].Amount = records[index].Amount
		ret.List[index].ChangeBefore = records[index].ChangeBefore
		ret.List[index].ChangeAfter = records[index].ChangeAfter
		cell, err := excelize.CoordinatesToCellName(1, index+2)
		if err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(c, code.ErrUnknown)
			return
		}

		switch ret.List[index].Type {
		case model.TypeBillingRecordCmsDeposit:
			recordType = "后台充值"
		case model.TypeBillingRecordCmsChange:
			recordType = "后台调整"
		case model.TypeBillingRecordSendRedPackSingle:
			recordType = "发个人红包"
		case model.TypeBillingRecordRecvRedPackSingle:
			recordType = "收个人红包"
		case model.TypeBillingRecordSendRedPackGroup:
			recordType = "发群红包"
		case model.TypeBillingRecordRecvRedPackGroup:
			recordType = "收群红包"
		case model.TypeBillingRecordRedPackGroupReturn:
			recordType = "群红包退款"
		case model.TypeBillingRecordRedPackSingleReturn:
			recordType = "个人红包退款"
		case model.TypeBillingRecordWithdraw:
			recordType = "提现"
		case model.TypeBillingRecordWithdrawSuccess:
			recordType = "提现成功"
		case model.TypeBillingRecordWithdrawFailed:
			recordType = "提现失败"
		case model.TypeBillingRecordWithdrawRollback:
			recordType = "提现回退"
		}
		timeLayout := "2006-01-02 15:04:05"
		finishTime = time.UnixMilli(ret.List[index].CreatedAt).Format(timeLayout)

		row = []interface{}{ret.List[index].ID, ret.List[index].SenderID, ret.List[index].SenderNickName, ret.List[index].ReceiverID, ret.List[index].ReceiverNickName, recordType, ret.List[index].Amount, ret.List[index].ChangeBefore, ret.List[index].ChangeAfter, ret.List[index].Note, finishTime}
		f.SetSheetRow("Sheet1", cell, &row)
	}

	f.SetActiveSheet(index)

	buf, err := f.WriteToBuffer()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}
	c.Writer.WriteHeader(http2.StatusOK)
	filename := url.QueryEscape("账单记录.xlsx")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Writer.Write(buf.Bytes())
	return
}

func (s *walletService) WalletChangeAmount(c *gin.Context) {
	req := new(model.WalletChangeAmountReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	senderID := c.GetString("o_user_id")
	key := fmt.Sprintf(repo.UserWalletKey, req.UserID)
	lock := util.NewLock(db.RedisCli, key)
	if err = lock.Lock(); err != nil {
		return
	}
	defer lock.Unlock()
	user, err := userUseCase.UserUseCase.GetInfo(req.UserID)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName())
		http.Failed(c, code.ErrUserIdNotExist)
		return
	}
	AmountNew := user.Balance + req.Amount
	if AmountNew < 0 {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName())
		http.Failed(c, code.ErrBalanceNotEnough)
		return
	}
	record := &model.BillingRecords{
		SenderID:     senderID,
		ReceiverID:   req.UserID,
		Type:         model.TypeBillingRecordCmsChange,
		Amount:       req.Amount,
		ChangeBefore: user.Balance,
		ChangeAfter:  AmountNew,
		Note:         req.Note,
	}
	record.CreatedAt = time.Now().UnixMilli()
	err = repo.WalletRepo.AmountChange("+", *record)
	http.Success(c)
	return
}

func (s *walletService) RedpackSingleRecordsList(c *gin.Context) {
	req := new(model.RedpackSingleRecordsListReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	records, count, err := repo.WalletRepo.RedpackSingleRecordPaging(*req)

	ret := new(model.RedpackSingleRecordsListResp)
	util.CopyStructFields(&ret.List, &records)
	for index := range ret.List {
		ret.List[index].SenderNickName = records[index].Sender.NickName
		ret.List[index].ReceiverNickName = records[index].Receiver.NickName
		ret.List[index].Amount = records[index].Amount
	}
	ret.Count = count
	ret.Page = req.Page
	ret.PageSize = req.PageSize
	http.Success(c, ret)
}

func (s *walletService) RedpackSingleRecordsExport(c *gin.Context) {
	req := new(model.RedpackSingleRecordsListReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	records, err := repo.WalletRepo.RedpackSingleRecordExport(*req)

	ret := new(model.RedpackSingleRecordsListResp)
	util.CopyStructFields(&ret.List, &records)

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		}
	}()

	index, err := f.NewSheet("Sheet1")
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}

	err = f.SetColWidth("Sheet1", "A", "H", 20)

	sheetHeader := []interface{}{"发送者账号", "发送者昵称", "接收者账号", "接收者昵称", "金额", "状态", "发送时间", "领取时间"}
	cell, err := excelize.CoordinatesToCellName(1, 1)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}
	f.SetSheetRow("Sheet1", cell, &sheetHeader)
	var row []interface{}
	var recordStatus, sendTime, recvTime string
	for index = range ret.List {
		ret.List[index].SenderNickName = records[index].Sender.NickName
		ret.List[index].ReceiverNickName = records[index].Receiver.NickName
		ret.List[index].Amount = records[index].Amount
		cell, err = excelize.CoordinatesToCellName(1, index+2)
		if err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(c, code.ErrUnknown)
			return
		}

		switch ret.List[index].Status {
		case 1:
			recordStatus = "待领取"
		case 2:
			recordStatus = "已领取"
		case 3:
			recordStatus = "已退回"
		}
		timeLayout := "2006-01-02 15:04:05"
		sendTime = time.Unix(ret.List[index].SendAt, 0).Format(timeLayout)
		if ret.List[index].RecvAt != nil {
			recvTime = time.Unix(*ret.List[index].RecvAt, 0).Format(timeLayout)
		}
		row = []interface{}{ret.List[index].SenderID, ret.List[index].SenderNickName, ret.List[index].ReceiverID, ret.List[index].ReceiverNickName, ret.List[index].Amount, recordStatus, sendTime, recvTime}
		f.SetSheetRow("Sheet1", cell, &row)
	}

	f.SetActiveSheet(index)

	buf, err := f.WriteToBuffer()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}
	c.Writer.WriteHeader(http2.StatusOK)
	filename := url.QueryEscape("个人红包记录.xlsx")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Writer.Write(buf.Bytes())
	return
}

func (s *walletService) WithdrawRecordsList(c *gin.Context) {
	req := new(model.WithdrawRecordsListReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	records, count, err := repo.WalletRepo.WithdrawRecordPaging(*req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret := new(model.WithdrawRecordsListResp)
	util.CopyStructFields(&ret.List, &records)
	for index := range ret.List {
		ret.List[index].NickName = records[index].User.NickName
		ret.List[index].Amount = records[index].Amount
	}
	ret.Count = count
	ret.Page = req.Page
	ret.PageSize = req.PageSize
	http.Success(c, ret)
}

func (s *walletService) WithdrawRecordsCountPending(c *gin.Context) {
	req := new(model.GetWithdrawRecordsNotDoneReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	count, err := repo.WalletRepo.WithdrawRecordCountPending()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}
	http.Success(c, count)
}

func (s *walletService) WithdrawRecordsDescribe(c *gin.Context) {
	req := new(model.GetWithdrawRecordsNotDoneReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	id := c.Param("id")
	record, err := repo.WalletRepo.WithdrawRecordDescribeByID(id)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret := new(model.GetWithdrawRecordsDescribeResp)
	ret.ID = record.ID
	err = util.JsonUnmarshal([]byte(record.Columns), &ret.Columns)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret.Note = record.Note
	http.Success(c, ret)
}

func (s *walletService) WithdrawRecordsStatusSet(c *gin.Context) {
	req := new(model.SetWithdrawRecordsStatusReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	if req.Status == model.StatusWithdrawRefused && len(req.Note) == 0 {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "提现拒绝时，备注为必填")
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = repo.WalletRepo.WithdrawRecordStatusSet(*req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret := new(model.SetWithdrawRecordsStatusResp)
	ret.ID = req.ID
	ret.Note = req.Note
	ret.Status = req.Status
	http.Success(c, ret)
}

func (s *walletService) WalletSetPayPass(c *gin.Context) {
	req := new(model.WalletSetPayPassReq)
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(c, code.GetError(err, req))
		return
	}

	user, err := userUseCase.UserUseCase.GetInfo(req.UserID)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUserNotFound)
		return
	}
	data := &userModel.User{}
	data.PayPassword = util.GetPassword(req.PayPasswd, user.Salt)
	opt := userRepo.WhereOption{
		UserId: user.UserID,
	}
	if _, err = userRepo.UserRepo.UpdateBy(opt, data); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUpdateAccount)
		return
	}
	http.Success(c)
	return
}
