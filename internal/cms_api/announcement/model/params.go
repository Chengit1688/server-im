package model

type GetAnnouncementInfoReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
}

type GetAnnouncementInfoResp struct {
	IsOpen  int    `json:"is_open"`
	Start   string `json:"start"`
	End     string `json:"end"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateAnnouncementInfoReq struct {
	OperationID string  `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	IsOpen      int     `json:"is_open" binding:"required,min=1,max=2" msg:"公告开关"`
	Start       *string `json:"start" binding:"required" msg:"公告开始时间"`
	End         *string `json:"end" binding:"required" msg:"公告结束时间"`
	Title       string  `json:"title" binding:"required,min=1,max=20" msg:"标题必输，长度1～20"`
	Content     string  `json:"content" binding:"required,min=1,max=200" msg:"内容必输，长度1～200"`
}
