package repo

import (
	"im/internal/api/chat/model"
	"im/pkg/db"
)

var LinkRepo = new(Link)

type Link struct{}

func (r *Link) GetListByStatus(status int) (list []model.Link, err error) {
	err = db.DB.Model(model.Link{}).Where("status = ?", status).Find(&list).Error
	return list, err
}

func (r *Link) BatchUpdateStatus(msgId string, status int) error {
	return db.DB.Model(model.Link{}).Where("msg_id = ?", msgId).Update("status", status).Error
}

func (r *Link) Updates(col map[string]interface{}, msgId string) error {
	return db.DB.Model(model.Link{}).Where("msg_id = ?", msgId).Updates(col).Error
}
