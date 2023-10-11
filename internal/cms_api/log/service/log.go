package service

import (
	"im/internal/cms_api/log/model"
	"im/internal/cms_api/log/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"

	"github.com/gin-gonic/gin"
)

var LogService = new(logService)

type logService struct{}

func (s *logService) OperateLogList(c *gin.Context) {
	req := new(model.OperateLogPagingReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	logs, count, err := repo.LogRepo.OperateLogPaging(*req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret := new(model.OperateLogPagingResp)
	util.CopyStructFields(&ret.List, &logs)
	for index, _ := range ret.List {
		ret.List[index].Username = logs[index].User.Username
		ret.List[index].NickName = logs[index].User.Nickname
	}
	ret.Count = count
	ret.Page = req.Page
	ret.PageSize = req.PageSize
	http.Success(c, ret)
}
