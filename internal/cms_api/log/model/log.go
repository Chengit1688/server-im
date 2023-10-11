package model

import (
	cmsUserModel "im/internal/cms_api/admin/model"
)

type OperateLogs struct {
	ID                 int64              `gorm:"column:id;primarykey" json:"-"`
	CreatedAt          int64              `gorm:"column:created_at;index:idx_created_at" json:"created_at"`
	UserID             string             `gorm:"column:user_id;index:idx_user_id;size:32" json:"user_id"`
	User               cmsUserModel.Admin `gorm:"foreignKey:UserID;references:UserID" json:"-"`
	ServiceID          string             `gorm:"column:service_id;index:idx_service_id;size:20" json:"service_id"`
	Env                string             `gorm:"column:env;size:20" json:"env"`
	LogLevel           string             `gorm:"column:log_level;size:20" json:"log_level"`
	RequestUrl         string             `gorm:"column:request_url;" json:"request_url"`
	RequestParams      string             `gorm:"column:request_params;" json:"request_params"`
	RequestBody        string             `gorm:"column:request_body;" json:"request_body"`
	RequestMethod      string             `gorm:"column:request_method;" json:"request_method"`
	RequestUserAgent   string             `gorm:"column:request_user_agent;" json:"request_user_agent"`
	ResponseStatusCode int                `gorm:"column:response_status_code;" json:"response_status_code"`
	ServiceIp          string             `gorm:"column:service_ip;" json:"service_ip"`
	ServiceHost        string             `gorm:"column:service_host;" json:"service_host"`
	OperationID        string             `gorm:"column:operation_id;size:50;index:idx_operation_id;" json:"operation_id"`
	LogRemark          string             `gorm:"column:log_remark;" json:"log_remark"`
}

func (d *OperateLogs) TableName() string {
	return "operate_logs"
}
