package model

import (
	"fmt"
	"time"
)

type GetDashboardReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
}

type GetDashboardResp struct {
	RegisterNum int64       `json:"register_num"`
	LoginNum    int64       `json:"login_num"`
	OnlineMax   int64       `json:"online_max"`
	SigleMsgNum int64       `json:"sigle_msg_num"`
	GroupNum    int64       `json:"group_num"`
	GroupMsgNum int64       `json:"group_msg_num"`
	RegisterBar DataBarList `json:"register_bar"`
	LoginBar    DataBarList `json:"login_bar"`
	SigleMsgBar DataBarList `json:"sigle_msg_bar"`
	GroupMsgBar DataBarList `json:"group_msg_bar"`
}

type DataBar struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type DataBarList []DataBar

func (m DataBarList) Len() int {
	return len(m)
}

func (m DataBarList) Less(i, j int) bool {
	itemI, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s 00:00:00", m[i].Date))
	itemJ, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s 00:00:00", m[j].Date))
	return itemI.Unix() < itemJ.Unix()
}

func (m DataBarList) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

type GetDailyDataReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	BeginDate   int64  `json:"begin_date" form:"begin_date" binding:"required"`
	EndDate     int64  `json:"end_date" form:"end_date" binding:"required"`
}
