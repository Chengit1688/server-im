package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	chatModel "im/internal/api/chat/model"
	apiChatUseCase "im/internal/api/chat/usecase"
	groupModel "im/internal/api/group/model"
	groupRepo "im/internal/api/group/repo"
	groupUsecase "im/internal/api/group/usecase"
	userUsecase "im/internal/api/user/usecase"
	"im/internal/api/wallet/model"
	cmsWalletModel "im/internal/cms_api/wallet/model"
	cmsWalletRepo "im/internal/cms_api/wallet/repo"
	"im/pkg/common"
	"im/pkg/db"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/mqtt"
	"im/pkg/response"
	"im/pkg/util"
	"time"
)

func (s *walletService) RedpackGroupSend(c *gin.Context) {
	req := new(model.RedpackGroupSendReq)
	err := c.ShouldBind(&req)
	lang := c.GetHeader("Locale")
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	if req.Type == cmsWalletModel.TypeRedpackGroupNormal {
		if req.PreAmount == 0 {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(c, response.GetError(response.ErrRedPackPreAmount, lang))
			return
		}
	}
	userId := c.GetString("user_id")
	var info groupModel.GroupInfo
	if info, err = groupUsecase.GroupUseCase.GroupInfo(req.GroupID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrGroupNotExist, lang))
		return
	}
	if info.Status != 1 {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("群被解散,%+v", info))
		http.Failed(c, response.GetError(response.ErrGroupNotExist, lang))
		return
	}
	key := fmt.Sprintf(cmsWalletRepo.UserWalletKey, userId)
	lock := util.NewLock(db.RedisCli, key)
	if err = lock.Lock(); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("userId:%s ,err: %s", userId, err))
		http.Failed(c, response.GetError(response.ErrUnknown, lang))
		return
	}
	defer lock.Unlock()
	user, err := userUsecase.UserUseCase.GetInfo(userId)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("userId:%s ,err: %s", userId, err))
		http.Failed(c, response.GetError(response.ErrUserNotFound, lang))
		return
	}
	if m, _ := groupRepo.GroupMemberRepo.GetMember(req.GroupID, userId); m == nil || err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("群不存在此user_id:%s,err: %v", userId, err))
		http.Failed(c, response.GetError(response.ErrNotInGroup, lang))
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
	redpackID, e2 := cmsWalletRepo.WalletRepo.RedpackGroupSend(model.RedpackGroupSendReq{
		Amount:  req.Amount,
		GroupID: req.GroupID,
		Count:   req.Count,
		SendID:  userId,
		Type:    req.Type,
		Remark:  req.Remark,
		MsgType: req.MsgType,
	}, user.Balance)
	if e2 != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrUnknown, lang))
		return
	}
	ret := new(model.RedpackGroupSendResp)
	ret.RedpackGroupID = redpackID
	ret.Status = cmsWalletModel.StatusRedpackSingleRecv
	ret.Amount = req.Amount
	http.Success(c, ret)
	return
}

func (s *walletService) RedpackGroupRecv(c *gin.Context) {
	req := new(model.RedpackGroupRecvReq)
	err := c.ShouldBind(&req)
	lang := c.GetHeader("Locale")
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	userId := c.GetString("user_id")
	walletKey := fmt.Sprintf(cmsWalletRepo.UserWalletKey, userId)
	walletLock := util.NewLock(db.RedisCli, walletKey)
	if err = walletLock.Lock(); err != nil {
		return
	}
	defer walletLock.Unlock()
	redpackGroupKey := fmt.Sprintf(cmsWalletRepo.RedpackGroupKey, req.RedpackGroupID)
	redpackGroupLock := util.NewLock(db.RedisCli, redpackGroupKey)
	if err = redpackGroupLock.Lock(); err != nil {
		return
	}
	defer redpackGroupLock.Unlock()
	user, err := userUsecase.UserUseCase.GetInfo(userId)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrUserNotFound, lang))
		return
	}
	redpack, err := cmsWalletRepo.WalletRepo.RedpackGroupGet(req.RedpackGroupID)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	if redpack.RemainderAmount == 0 || redpack.Count == redpack.ReceiveCount ||
		redpack.Status == cmsWalletModel.StatusRedpackSingleReturn {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrRedPackReceiveFinish, lang))
		return
	}
	recv, err := cmsWalletRepo.WalletRepo.RedpackGroupRecvGet(userId)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrRedPackReceive, lang))
		return
	}

	if recv.UserID == userId {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrRedPackReceiveRepeat, lang))
		return
	}
	var recvAmount, timeInt64 int64
	if redpack.Type == cmsWalletModel.TypeRedpackGroupNormal {

		recvAmount = redpack.Amount / redpack.Count
	} else if redpack.Type == cmsWalletModel.TypeRedpackGroupRandom {

		rC := redpack.Count - redpack.ReceiveCount
		recvAmount = util.DoubleAverage(rC, redpack.RemainderAmount)
	}
	timeInt64 = time.Now().Unix()
	ret := new(model.RedpackGroupRecvResp)
	red := cmsWalletModel.RedpackGroupRecvs{
		RecvAt:         timeInt64,
		SendAt:         redpack.SendAt,
		SenderID:       redpack.SenderID,
		UserID:         userId,
		RedpackGroupID: req.RedpackGroupID,
		RedpackGroup:   cmsWalletModel.RedpackGroupRecords{GroupID: redpack.GroupID},
		Amount:         recvAmount,
		Status:         cmsWalletModel.StatusRedpackSingleRecvd,
	}
	ret.RecvAt = timeInt64
	ret.Status = cmsWalletModel.StatusRedpackSingleRecv
	if redpack.Status == cmsWalletModel.StatusRedpackSingleRecv {
		timeInt64, err = cmsWalletRepo.WalletRepo.RedpackGroupRecv(red, user.Balance)
		if err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(c, response.GetError(response.ErrUnknown, lang))
			return
		}
		ret.Status = cmsWalletModel.StatusRedpackSingleRecvd
		conversationID := apiChatUseCase.ConversationUseCase.GetConversationID(chatModel.ConversationTypeGroup, redpack.SenderID, redpack.GroupID)
		pushData := model.RedpackGroupMessagePush{
			ConversationID: conversationID,
			Timestamp:      timeInt64,

			ReceiverNickname: recv.User.NickName,
			GroupId:          redpack.GroupID,
			RedpackGroupID:   redpack.ID,
			Type:             redpack.MsgType,
		}
		err = mqtt.SendMessageToUsers(req.OperationID, common.RedpackGroupRecvPush, pushData, []string{redpack.SenderID, userId}...)
		if err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		}
	}
	ret.Amount = recvAmount
	ret.SendAt = redpack.SendAt
	ret.RedpackGroupID = req.RedpackGroupID
	http.Success(c, ret)
	return
}

func (s *walletService) RedpackGroupGetInfo(c *gin.Context) {
	req := new(model.RedpackGroupGetInfoReq)
	err := c.ShouldBindJSON(&req)
	lang := c.GetHeader("Locale")
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	redpack, err := cmsWalletRepo.WalletRepo.RedpackGroupGet(req.RedpackGroupID)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrFailRequest, lang))
		return
	}
	ret := new(model.RedpackGroupRecvResp)
	err = util.CopyStructFields(&ret, redpack)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrFailRequest, lang))
		return
	}
	ret.Amount = redpack.Amount
	ret.SendAt = redpack.SendAt
	ret.Remark = redpack.Remark
	ret.Type = redpack.Type
	ret.RedpackGroupID = redpack.ID
	ret.Status = redpack.Status
	http.Success(c, ret)
	return
}
