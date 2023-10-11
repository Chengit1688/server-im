package model

type UploadReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	FileType    int    `json:"file_type" form:"file_type" binding:"required" msg:"文件类型不能为空"`
}

type UploadResp struct {
	OldName     string `json:"old_name"`
	NewName     string `json:"new_name"`
	Url         string `json:"url"`
	Thumbnail   string `json:"thumbnail"`
	ContentType string `json:"content_type"`
}

type GetSTSReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required,gte=2"   msg:"日志id必须传递"`
}

type Credentials struct {
	Region        string `json:"region"`
	BucketName    string `json:"bucketName"`
	AccessId      string `json:"accessId"`
	AccessSecret  string `json:"accessSecret"`
	SecurityToken string `json:"securityToken"`
	Expiration    uint   `json:"expiration"`
	SavePath      string `json:"savePath"`
}

type GetUrlReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Filename    string `json:"file_name" form:"file_name" binding:"required" msg:"文件名称不能为空"`
	FileType    int    `json:"file_type" form:"file_type" binding:"required" msg:"文件类型不能为空"`
	Width       int    `json:"width" form:"width"  msg:"文件的宽"`
	Height      int    `json:"height" form:"height"  msg:"文件的高"`
}
