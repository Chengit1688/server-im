package model

type SetCmsUIDataNormalReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Type        int    `json:"type" binding:"required,min=1,max=5" msg:"类型不能为空，1站点名，2登录页图标，3登录页背景图，4页签图标，5菜单栏图标"`
	Value       string `json:"value"`
}

type SetCmsUIDataSiteNameReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Value       string `json:"value" binding:"omitempty,max=20" msg:"站点名最长20"`
}
