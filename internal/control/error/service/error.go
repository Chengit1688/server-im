package service

import (
	"fmt"
	"im/internal/control/error/model"
	"im/internal/control/error/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var ErrorService = new(errorService)

type errorService struct{}

func (s *errorService) HandlerUpload(c *gin.Context) {
	params := model.UploadInfo{}
	if err := c.BindJSON(&params); err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "params error"), zap.String("operation_id", params.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}

	info := model.ErrLog{
		AppName:   params.AppName,
		UserId:    params.UserId,
		MacType:   params.MacType,
		PhoneType: params.PhoneType,
		CreatTime: time.Now().Unix(),
		Info:      params.Info,
		Extra:     params.Extra,
	}

	if err := repo.ErrorRepo.InsertIntoErrLog(info); err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "insert error"), zap.String("operation_id", params.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}

	go writelog(info)

	http.Success(c, nil)
	return
}

func writelog(info model.ErrLog) {
	fmt.Printf("\n------insert info  ----%+v\n", info)
}

func (s *errorService) HandlerSearch(c *gin.Context) {
	params := model.SearchInfo{}

	if err := c.BindJSON(&params); err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "params error"), zap.String("operation_id", params.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}

	params.Pagination.Check()

	infos, count, err := repo.ErrorRepo.QueryErrlogInfo(params.Offset, params.Limit, params.Stime, params.Etime, params.UserId, params.MacType)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "db query error"), zap.String("operation_id", params.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}
	resp := new(model.SearchInfoResp)

	for _, info := range infos {
		var errInfo model.ErrlogInfo
		errInfo.AppName = info.AppName
		errInfo.UserId = info.UserId
		errInfo.MacType = info.MacType
		errInfo.Info = info.Info
		errInfo.CreateTime = info.CreatTime
		errInfo.PhoneType = info.PhoneType
		errInfo.Extra = info.Extra
		resp.List = append(resp.List, errInfo)
	}
	resp.Page = params.Page
	resp.PageSize = params.PageSize
	resp.Count = count
	http.Success(c, resp)
	return
}
