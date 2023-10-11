package service

import (
	"fmt"
	chatModel "im/internal/api/chat/model"
	apiChatUseCase "im/internal/api/chat/usecase"
	friendUsecase "im/internal/api/friend/usecase"
	userModel "im/internal/api/user/model"
	userRepo "im/internal/api/user/repo"
	userUsecase "im/internal/api/user/usecase"
	"im/internal/api/wallet/model"
	configModel "im/internal/cms_api/config/model"
	configUsecase "im/internal/cms_api/config/usecase"
	cmsWalletModel "im/internal/cms_api/wallet/model"
	cmsWalletRepo "im/internal/cms_api/wallet/repo"
	"im/pkg/code"
	"im/pkg/common"
	"im/pkg/db"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/mqtt"
	"im/pkg/response"
	"im/pkg/util"

	"github.com/gin-gonic/gin"
)

var WalletService = new(walletService)

type walletService struct{}

func (s *walletService) GetWalletInfo(c *gin.Context) {
	req := new(model.GetWalletReq)
	err := c.ShouldBindQuery(&req)
	lang := c.GetHeader("Locale")
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	user_id := c.GetString("user_id")
	user, err := userUsecase.UserUseCase.GetInfo(user_id)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrUserNotFound, lang))
		return
	}
	deposit, err := configUsecase.ConfigUseCase.GetDepositConfig()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrUnknown, lang))
		return
	}
	withdraw, err := configUsecase.ConfigUseCase.GetWithdrawConfig()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrUnknown, lang))
		return
	}
	ret := new(model.GetWalletResp)
	ret.PayPasswdSet = 2
	ret.Balance = user.Balance
	ret.Deposit = *deposit
	min := withdraw.Min
	max := withdraw.Max
	ret.Withdraw = configModel.WithdrawApigResp{Min: min, Max: max, Columns: withdraw.Columns}
	if len(user.PayPassword) != 0 {
		ret.PayPasswdSet = 1
	}
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrUnknown, lang))
		return
	}
	http.Success(c, ret)
}

func (s *walletService) RedpackSingleSend(c *gin.Context) {
	req := new(model.RedpackSingleSendReq)
	err := c.ShouldBind(&req)
	lang := c.GetHeader("Locale")
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	user_id := c.GetString("user_id")
	if user_id == req.RecvID {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "个人红包不允许自己给自己发")
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	if !friendUsecase.FriendUseCase.CheckFriend(user_id, req.RecvID) {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}
	key := fmt.Sprintf(cmsWalletRepo.UserWalletKey, user_id)
	lock := util.NewLock(db.RedisCli, key)
	if err = lock.Lock(); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrUnknown, lang))
		return
	}
	defer lock.Unlock()
	user, err := userUsecase.UserUseCase.GetInfo(user_id)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrUserNotFound, lang))
		return
	}
	_, err = userUsecase.UserUseCase.GetInfo(req.RecvID)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrUserNotFound, lang))
		return
	}
	checkPwd := util.CheckPassword(user.PayPassword, req.PayPasswd, user.Salt)
	if !checkPwd {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "支付密码错误")
		http.Failed(c, response.GetError(response.ErrPayPasswdWrong, lang))
		return
	}
	if user.Balance < req.Amount {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "余额不足")
		http.Failed(c, response.GetError(response.ErrBalanceNotEnough, lang))
		return
	}
	redpackID, err := cmsWalletRepo.WalletRepo.RedpackSingleSend(model.RedpackSingleSendReq{
		Amount:  req.Amount,
		RecvID:  req.RecvID,
		SendID:  user_id,
		Remark:  req.Remark,
		MsgType: req.MsgType,
	}, user.Balance)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrUnknown, lang))
		return
	}
	ret := new(model.RedpackSingleSendResp)
	ret.RedpackSingleID = redpackID
	ret.Status = cmsWalletModel.StatusRedpackSingleRecv
	ret.Amount = req.Amount
	http.Success(c, ret)
	return
}

func (s *walletService) RedpackSingleRecv(c *gin.Context) {
	req := new(model.RedpackSingleRecvReq)
	err := c.ShouldBind(&req)
	lang := c.GetHeader("Locale")
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	user_id := c.GetString("user_id")
	walletKey := fmt.Sprintf(cmsWalletRepo.UserWalletKey, user_id)
	walletLock := util.NewLock(db.RedisCli, walletKey)
	if err = walletLock.Lock(); err != nil {
		return
	}
	defer walletLock.Unlock()

	redpackSingleKey := fmt.Sprintf(cmsWalletRepo.RedpackSingleKey, req.RedpackSingleID)
	redpackSingleLock := util.NewLock(db.RedisCli, redpackSingleKey)
	if err = redpackSingleLock.Lock(); err != nil {
		return
	}
	defer redpackSingleLock.Unlock()

	user, err := userUsecase.UserUseCase.GetInfo(user_id)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrUserNotFound, lang))
		return
	}
	redpack, err := cmsWalletRepo.WalletRepo.RedpackSingleGet(req.RedpackSingleID)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrRedPackReceive, lang))
		return
	}

	if redpack.ReceiverID == user_id {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrRedPackReceiveRepeat, lang))
		return
	}
	ret := new(model.RedpackSingleRecvResp)
	if redpack.SenderID == user_id {
		ret.Amount = redpack.Amount
		ret.RecvAt = redpack.RecvAt
		ret.RedpackSingleID = req.RedpackSingleID
		ret.Status = redpack.Status
		http.Success(c, ret)
		return
	}

	var timeInt64 int64
	ret.RecvAt = redpack.RecvAt
	ret.Status = redpack.Status
	if redpack.Status == cmsWalletModel.StatusRedpackSingleRecv {
		if timeInt64, err = cmsWalletRepo.WalletRepo.RedpackSingleRecv(redpack, user.Balance); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(c, response.GetError(response.ErrUnknown, lang))
			return
		}
		ret.RecvAt = &timeInt64
		ret.Status = cmsWalletModel.StatusRedpackSingleRecvd
		conversationID := apiChatUseCase.ConversationUseCase.GetConversationID(chatModel.ConversationTypeSingle, redpack.SenderID, redpack.ReceiverID)
		pushData := model.RedpackSingleMessagePush{
			ConversationID: conversationID, Timestamp: timeInt64,
			SenderNickname:   redpack.Sender.NickName,
			ReceiverNickname: redpack.Receiver.NickName,
			RedpackID:        redpack.ID,
			Type:             redpack.MsgType,
		}
		if err = mqtt.SendMessageToUsers(req.OperationID, common.RedpackSingleRecvPush, pushData, []string{redpack.SenderID, redpack.ReceiverID}...); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		}
	}
	ret.Amount = redpack.Amount
	ret.RedpackSingleID = req.RedpackSingleID
	http.Success(c, ret)
	return
}

func (s *walletService) RedpackSingleGetInfo(c *gin.Context) {
	req := new(model.RedpackSingleGetInfoReq)
	lang := c.GetHeader("Locale")
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	redpack, err := cmsWalletRepo.WalletRepo.RedpackSingleGet(req.RedpackSingleID)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	ret := new(model.RedpackSingleRecvResp)
	util.CopyStructFields(&ret, redpack)
	ret.Amount = redpack.Amount

	ret.RedpackSingleID = redpack.ID
	ret.Remark = redpack.Remark

	http.Success(c, ret)
	return
}

func (s *walletService) WalletSetPayPass(c *gin.Context) {
	req := new(model.WalletSetPayPassReqs)
	lang := c.GetHeader("Locale")
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(c, code.GetError(err, req))
		return
	}

	user_id := c.GetString("user_id")
	user, err := userUsecase.UserUseCase.GetInfo(user_id)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrUserNotFound, lang))
		return
	}

	data := &userModel.User{}
	data.PayPassword = util.GetPassword(req.PayPasswd, user.Salt)
	opt := userRepo.WhereOption{
		UserId: user.UserID,
	}
	if _, err = userRepo.UserRepo.UpdateBy(opt, data); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrUpdateAccount, lang))
		return
	}
	http.Success(c)
	return
}

func (s *walletService) WithdrawCommit(c *gin.Context) {
	req := new(model.WithdrawCommitReq)
	lang := c.GetHeader("Locale")
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	var amount int64
	for _, col := range req.Columns {
		if col.Name == "提现金额" || col.Name == "withdrawal Amount" || col.Name == "出金額" {
			amount = int64(util.String2Int(col.Value))
		}
	}

	if !s.WithdrawCheck(amount) {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "不满足最小提现金额或最大提现金额要求")
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	user_id := c.GetString("user_id")

	walletKey := fmt.Sprintf(cmsWalletRepo.UserWalletKey, user_id)
	walletLock := util.NewLock(db.RedisCli, walletKey)
	if err = walletLock.Lock(); err != nil {
		return
	}
	defer walletLock.Unlock()

	user, err := userUsecase.UserUseCase.GetInfo(user_id)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrUserNotFound, lang))
		return
	}
	if user.Balance < amount {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "余额不足")
		http.Failed(c, response.GetError(response.ErrBalanceNotEnough, lang))
		return
	}
	colsByte, err := util.JsonMarshal(req.Columns)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrUnknown, lang))
		return
	}
	add := new(cmsWalletModel.WithdrawRecords)
	add.Amount = amount
	add.UserID = user_id
	add.Columns = string(colsByte)
	ret := new(model.WithdrawCommitResp)
	record, err := cmsWalletRepo.WalletRepo.WithdrawRecordAdd(*add)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}
	ret.ID = record.ID
	ret.BillingID = record.BillingID
	http.Success(c, ret)
}

func (s *walletService) WithdrawCheck(amount int64) bool {
	config, err := configUsecase.ConfigUseCase.GetWithdrawConfig()
	if err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		return false
	}
	if amount >= config.Min && amount <= config.Max {
		return true
	}
	return false
}

func (s *walletService) BillingRecordsList(c *gin.Context) {
	req := new(model.BillingRecordsListReq)
	lang := c.GetHeader("Locale")
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	user_id := c.GetString("user_id")
	req.UserID = user_id
	records, count, err := cmsWalletRepo.WalletRepo.BillingRecordPagingByUser(*req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}

	ret := new(model.BillingRecordsListResp)
	ret.Count = count
	for _, record := range records {
		var name string
		switch lang {
		case "en_US":
			name = s.getBillingEnName(record)
		case "zh_CN":
			name = s.getBillingName(record)
		case "ja":
			name = s.getBillingJaName(record)
		default:
			name = s.getBillingName(record)
		}

		ret.List = append(ret.List, model.BillingRecordsListItem{Amount: record.Amount, Name: name, Type: record.Type, CreatedAt: record.CreatedAt, UpdatedAt: record.UpdatedAt})
	}
	ret.Page = req.Page
	ret.PageSize = req.PageSize
	http.Success(c, ret)
	return
}

func (s *walletService) getBillingName(record cmsWalletModel.BillingRecords) (name string) {
	switch record.Type {
	case cmsWalletModel.TypeBillingRecordCmsDeposit:
		name = "后台充值"
	case cmsWalletModel.TypeBillingRecordCmsChange:
		name = "后台调整"
	case cmsWalletModel.TypeBillingRecordSendRedPackSingle:
		name = fmt.Sprintf("红包 发给%s", record.Receiver.NickName)
	case cmsWalletModel.TypeBillingRecordRecvRedPackSingle:
		name = fmt.Sprintf("红包 收到%s", record.SenderApiUser.NickName)
	case cmsWalletModel.TypeBillingRecordSendRedPackGroup:
		name = fmt.Sprintf("红包 发给群%s", record.Group.Name)
	case cmsWalletModel.TypeBillingRecordRecvRedPackGroup:
		name = fmt.Sprintf("红包 收到群%s %s", record.Group.Name, record.SenderApiUser.NickName)
	case cmsWalletModel.TypeBillingRecordRedPackGroupReturn:
		name = "群红包退款"
	case cmsWalletModel.TypeBillingRecordRedPackSingleReturn:
		name = "个人红包退款"
	case cmsWalletModel.TypeBillingRecordWithdraw:
		name = "提现"
	case cmsWalletModel.TypeBillingRecordWithdrawSuccess:
		name = "提现成功"
	case cmsWalletModel.TypeBillingRecordWithdrawFailed:
		name = "提现失败"
	case cmsWalletModel.TypeBillingRecordWithdrawRollback:
		name = "提现回退"
	case cmsWalletModel.TypeBillingRecordWithdrawSignAward:
		name = "每日签到奖励"
	}
	return
}

func (s *walletService) getBillingEnName(record cmsWalletModel.BillingRecords) (name string) {
	switch record.Type {
	case cmsWalletModel.TypeBillingRecordCmsDeposit:
		name = "recharge"
	case cmsWalletModel.TypeBillingRecordCmsChange:
		name = "Adjustment"
	case cmsWalletModel.TypeBillingRecordSendRedPackSingle:
		name = fmt.Sprintf("red packet sent to %s", record.Receiver.NickName)
	case cmsWalletModel.TypeBillingRecordRecvRedPackSingle:
		name = fmt.Sprintf("red packet  receive by %s", record.SenderApiUser.NickName)
	case cmsWalletModel.TypeBillingRecordSendRedPackGroup:
		name = fmt.Sprintf("red packet sent to the group %s", record.Group.Name)
	case cmsWalletModel.TypeBillingRecordRecvRedPackGroup:
		name = fmt.Sprintf("red packet received group by %s %s", record.Group.Name, record.SenderApiUser.NickName)
	case cmsWalletModel.TypeBillingRecordRedPackGroupReturn:
		name = "Group Red Received Refund"
	case cmsWalletModel.TypeBillingRecordRedPackSingleReturn:
		name = "Personal red received refund"
	case cmsWalletModel.TypeBillingRecordWithdraw:
		name = "withdraw"
	case cmsWalletModel.TypeBillingRecordWithdrawSuccess:
		name = "Successful withdrawal"
	case cmsWalletModel.TypeBillingRecordWithdrawFailed:
		name = "Withdrawal failed"
	case cmsWalletModel.TypeBillingRecordWithdrawRollback:
		name = "Withdraw back"
	case cmsWalletModel.TypeBillingRecordWithdrawSignAward:
		name = "Sign In Reward"
	}
	return
}

func (s *walletService) getBillingJaName(record cmsWalletModel.BillingRecords) (name string) {
	switch record.Type {
	case cmsWalletModel.TypeBillingRecordCmsDeposit:
		name = "充電する"
	case cmsWalletModel.TypeBillingRecordCmsChange:
		name = "調整"
	case cmsWalletModel.TypeBillingRecordSendRedPackSingle:
		name = fmt.Sprintf("赤い封筒が送られてきました %s", record.Receiver.NickName)
	case cmsWalletModel.TypeBillingRecordRecvRedPackSingle:
		name = fmt.Sprintf("受け取る %s", record.SenderApiUser.NickName)
	case cmsWalletModel.TypeBillingRecordSendRedPackGroup:
		name = fmt.Sprintf("グループに送信する%s", record.Group.Name)
	case cmsWalletModel.TypeBillingRecordRecvRedPackGroup:
		name = fmt.Sprintf("受信したグループ%s %s", record.Group.Name, record.SenderApiUser.NickName)
	case cmsWalletModel.TypeBillingRecordRedPackGroupReturn:
		name = "グループ赤い封筒の払い戻し"
	case cmsWalletModel.TypeBillingRecordRedPackSingleReturn:
		name = "個人的な赤い封筒の払い戻し"
	case cmsWalletModel.TypeBillingRecordWithdraw:
		name = "撤退"
	case cmsWalletModel.TypeBillingRecordWithdrawSuccess:
		name = "出金成功"
	case cmsWalletModel.TypeBillingRecordWithdrawFailed:
		name = "出金に失敗しました"
	case cmsWalletModel.TypeBillingRecordWithdrawRollback:
		name = "撤退"
	case cmsWalletModel.TypeBillingRecordWithdrawSignAward:
		name = "サインイン報酬"
	}
	return
}
