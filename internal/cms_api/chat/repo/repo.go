package repo

import (
	"fmt"
	"im/internal/cms_api/chat/model"
	"im/pkg/db"

	"gorm.io/gorm"
)

var MessageRepo = new(messageRepo)

type messageRepo struct{}

func (r *messageRepo) MultiSendPaging(req model.GetMultiSendPagingReq) (items []model.MultiSendRecord, count int64, err error) {

	req.Pagination.Check()
	tx := db.DB.Table("cms_multi_send_record AS t1").
		Joins("JOIN cms_admins AS t2 ON t1.operate_id = t2.user_id").
		Joins("JOIN cms_multi_send_user AS t3 ON t1.id = t3.record_id").
		Joins("JOIN users AS t4 ON t3.sender_id = t4.user_id").
		Joins("inner join (SELECT id FROM cms_multi_send_record ORDER BY created_at DESC  LIMIT ? OFFSET ?) as o on o.id = t1.id ", req.Pagination.Limit, req.Pagination.Offset)
	if len(req.Content) != 0 {
		tx = tx.Where(fmt.Sprintf("t1.content like %q", ("%" + req.Content + "%")))
	}
	if len(req.SenderID) != 0 {
		tx = tx.Where("t3.sender_id = ?", req.SenderID)
	}
	if len(req.SenderNickname) != 0 {
		tx = tx.Where(fmt.Sprintf("t4.nick_name like %q", ("%" + req.SenderNickname + "%")))
	}
	if len(req.Operate) != 0 {
		tx = tx.Where(fmt.Sprintf("t2.username like %q", ("%" + req.Operate + "%")))
	}
	if req.OperateTimeStart != 0 {
		tx = tx.Where("t1.created_at >= ?", req.OperateTimeStart)
	}
	if req.OperateTimeEnd != 0 {
		tx = tx.Where("t1.created_at <= ?", req.OperateTimeEnd)
	}
	tx.Select("t1.id,t1.operate_id,t1.content,t1.created_at,t2.username as username,t4.user_id as sender_id,t4.nick_name as sender_nickname").Order("t1.created_at DESC")
	err = tx.Find(&items).Error
	if err != nil {
		return
	}
	var countItems []model.MultiSendRecord
	var countSumMap map[int]bool
	countSumMap = make(map[int]bool)
	countTx := db.DB.Table("cms_multi_send_record AS t1").
		Joins("JOIN cms_admins AS t2 ON t1.operate_id = t2.user_id").
		Joins("JOIN cms_multi_send_user AS t3 ON t1.id = t3.record_id").
		Joins("JOIN users AS t4 ON t3.sender_id = t4.user_id")
	if len(req.Content) != 0 {
		tx = countTx.Where(fmt.Sprintf("t1.content like %q", ("%" + req.Content + "%")))
	}
	if len(req.SenderID) != 0 {
		tx = countTx.Where("t3.sender_id = ?", req.SenderID)
	}
	if len(req.SenderNickname) != 0 {
		tx = countTx.Where(fmt.Sprintf("t4.nick_name like %q", ("%" + req.SenderNickname + "%")))
	}
	if len(req.Operate) != 0 {
		tx = countTx.Where(fmt.Sprintf("t2.username like %q", ("%" + req.Operate + "%")))
	}
	if req.OperateTimeStart != 0 {
		tx = countTx.Where("t1.created_at >= ?", req.OperateTimeStart)
	}
	if req.OperateTimeEnd != 0 {
		tx = countTx.Where("t1.created_at <= ?", req.OperateTimeEnd)
	}
	countTx.Select("t1.id,t1.operate_id,t1.content,t1.created_at")
	err = countTx.Find(&countItems).Error
	for _, row := range countItems {
		if _, ok := countSumMap[row.ID]; !ok {
			countSumMap[row.ID] = true
		}
	}
	count = int64(len(countSumMap))
	return
}

func (r *messageRepo) MultiSendAdd(item model.MultiSendRecord, userIDs []string) (err error) {
	db.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(&item).Find(&item).Error; err != nil {
			return err
		}
		id := int(item.ID)
		for _, userID := range userIDs {

			add := new(model.MultiSendUser)
			add.SenderID = userID
			add.RecordID = id
			if err = tx.Create(&add).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return
}
