package service

import (
	chatModel "im/internal/api/chat/model"
	"im/internal/api/push/model"
	pushUseCase "im/internal/api/push/usecase"
	userUseCase "im/internal/api/user/usecase"
	"im/pkg/code"
	"im/pkg/common"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"

	"github.com/gin-gonic/gin"
)

var JpushService = new(jpushService)

type jpushService struct{}

func (s *jpushService) Jpush(c *gin.Context) {
	req := new(model.JpushReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	user_id := c.GetString("user_id")
	var alert model.MsgText
	alert.Text = req.Alert
	data, _ := util.JsonMarshal(alert)
	str, _ := util.Encrypt(data, common.ContentKey)
	for _, recv_id := range req.UserIDs {
		_, err = userUseCase.UserUseCase.GetBaseInfo(recv_id)
		if err != nil {
			http.Failed(c, code.ErrUserIdNotExist)
			return
		}
		pushUseCase.JpushUseCase.Push(req.OperationID, user_id, recv_id, str, chatModel.MessageText, chatModel.ConversationTypeSingle)
	}
	http.Success(c)
	return
}
