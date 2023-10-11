package service

import (
	"fmt"
	"im/internal/cms_api/wallet/model"
	"im/internal/cms_api/wallet/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"
	http2 "net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func (s *walletService) RedpackGroupRecordsList(c *gin.Context) {
	req := new(model.RedpackGroupRecordsListReq)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "err", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	records, count, err := repo.WalletRepo.RedpackGroupRecordPaging(*req)
	ret := new(model.RedpackGroupRecordsListResp)
	util.CopyStructFields(&ret.List, &records)
	for index := range ret.List {
		ret.List[index].SenderNickName = records[index].Sender.NickName
		ret.List[index].SenderUserId = records[index].Sender.UserID
		ret.List[index].Amount = records[index].Amount
		ret.List[index].GroupName = records[index].Group.Name
	}
	ret.Count = count
	ret.Page = req.Page
	ret.PageSize = req.PageSize
	http.Success(c, ret)
}

func (s *walletService) RedpackGroupRecordsExport(c *gin.Context) {
	req := new(model.RedpackGroupRecordsListReq)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	records, err := repo.WalletRepo.RedpackGroupRecordExport(*req)

	ret := new(model.RedpackGroupRecordsListResp)
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
	sheetHeader := []interface{}{"群名称", "发送者ID", "发送者昵称", "红包类型", "红包金额", "红包个数", "状态", "发送时间"}
	cell, err := excelize.CoordinatesToCellName(1, 1)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}
	f.SetSheetRow("Sheet1", cell, &sheetHeader)
	var row []interface{}
	var recordStatus, sendTime, packType string
	for index = range ret.List {
		ret.List[index].SenderNickName = records[index].Sender.NickName
		ret.List[index].SenderUserId = records[index].Sender.UserID
		ret.List[index].Amount = records[index].Amount
		ret.List[index].GroupName = records[index].Group.Name
		ret.List[index].Total = records[index].Count

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

		switch ret.List[index].Type {
		case 1:
			packType = "拼手气红包"
		case 2:
			packType = "普通红包"
		}
		timeLayout := "2006-01-02 15:04:05"
		sendTime = time.Unix(ret.List[index].SendAt, 0).Format(timeLayout)
		if ret.List[index].SendAt != 0 {
			sendTime = time.Unix(ret.List[index].SendAt, 0).Format(timeLayout)
		}
		row = []interface{}{ret.List[index].GroupName, ret.List[index].SenderUserId, ret.List[index].SenderNickName, packType, ret.List[index].Amount, ret.List[index].Total, recordStatus, sendTime}
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
	filename := url.QueryEscape("群红包记录.xlsx")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Writer.Write(buf.Bytes())
	return
}
