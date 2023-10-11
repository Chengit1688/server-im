package model

type JpushReq struct {
	OperationID string   `json:"operation_id" form:"operation_id"  binding:"required,gte=1" msg:"operation_id required"`
	UserIDs     []string `json:"user_ids"`
	Title       string   `json:"title"`
	Alert       string   `json:"alert"`
}

type JpushResp struct {
	SendNo string `json:"sendno"`
	MsgID  string `json:"msg_id"`
}

type MsgText struct {
	Text string `json:"text"`
}
