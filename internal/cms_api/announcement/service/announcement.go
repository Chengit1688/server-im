package service

import (
	"im/internal/cms_api/announcement/model"
	configRepo "im/internal/cms_api/config/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var AnnouncementService = new(announcementService)

type announcementService struct{}

func (s *announcementService) GetAnnouncement(c *gin.Context) {
	req := new(model.GetAnnouncementInfoReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	config, err := configRepo.ConfigRepo.GetAnnouncementConfig()
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret := new(model.GetAnnouncementInfoResp)
	err = util.JsonUnmarshal([]byte(config.Content), &ret)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}
	http.Success(c, ret)
}

func (s *announcementService) UpdateAnnouncement(c *gin.Context) {
	req := new(model.UpdateAnnouncementInfoReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	regx, _ := regexp.Compile(`^\d{1,2}:\d{1,2}:\d{1,2}`)
	if *req.Start != "" {
		if !regx.MatchString(*req.Start) {
			logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "start 时间格式错误"), zap.String("operation_id", req.OperationID))
			http.Failed(c, code.ErrBadRequest)
			return
		}
	}
	if *req.End != "" {
		if !regx.MatchString(*req.End) {
			logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "end 时间格式错误"), zap.String("operation_id", req.OperationID))
			http.Failed(c, code.ErrBadRequest)
			return
		}
	}
	configModel := new(model.GetAnnouncementInfoResp)
	util.CopyStructFields(&configModel, &req)
	config, err := util.JsonMarshal(req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}
	err = configRepo.ConfigRepo.UpdateAnnouncementConfig(string(config))
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}
	http.Success(c, configModel)
}
