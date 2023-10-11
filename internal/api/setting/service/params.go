package service

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

type DomainInfo struct {
	ID     uint   `json:"id"`
	Site   string `json:"site"`
	Domain string `json:"domain"`
}
