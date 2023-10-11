package model

type ListReq struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type DomainInfo struct {
	ID     uint   `json:"id"`
	Site   string `json:"site"`
	Domain string `json:"domain"`
}

type AddDomainReq struct {
	Site   string `json:"site" binding:"required"`
	Domain string `json:"domain" binding:"required"`
}

type AddDomainResp struct{}

type RemoveDomainReq struct {
	Site   string `json:"site" binding:"required"`
	Domain string `json:"domain" binding:"required"`
}

type RemoveDomainResp struct{}

type DomainListReq struct {
	ListReq
	Site   string `json:"site"`
	Domain string `json:"domain"`
}
type HttpDomainListReturn struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Resp    DomainListResp `json:"data"`
}

type DomainListResp struct {
	Count    int          `json:"count"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
	List     []DomainInfo `json:"list"`
}

type AppDomainListReq struct {
	Site string `json:"site" binding:"required"`
}

type AppDomainListResp []string

type WarningInfo struct {
	ID     uint   `json:"id"`
	Domain string `json:"domain"`

	Ip        string `json:"ip"`
	CreatedAt int64  `json:"created_at"`
}

type AddWarningReq struct {
	Domain string `json:"domain" binding:"required"`
}

type AddWarningResp struct{}

type WarningListReq struct {
	ListReq
	Domain string `json:"domain"`
}

type WarningListResp struct {
	Count    int           `json:"count"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
	List     []WarningInfo `json:"list"`
}
