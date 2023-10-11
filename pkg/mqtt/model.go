package mqtt

type ConnStateType string

const (
	ConnStateTypeConnected    = "connected"
	ConnStateTypeIdle         = "idle"
	ConnStateTypeDisconnected = "disconnected"
)

type Meta struct {
	Page    int  `json:"page"`
	Limit   int  `json:"limit"`
	Count   int  `json:"count"`
	HasNext bool `json:"hasnext"`
}

type Resp struct {
	Code int  `json:"code"`
	Meta Meta `json:"meta"`
}

type Client struct {
	Username      string `json:"username"`     // 认证账号
	ClientID      string `json:"clientid"`     // 客户端唯一ID
	IsOnline      bool   `json:"connected"`    // 是否在线
	LastLoginTime string `json:"connected_at"` // 最后上线时间
	LoginIP       string `json:"ip_address"`   // 登录IP
}

type ClientResp struct {
	Resp
	Data []Client `json:"data"`
}

type Auth struct {
	Username string `json:"username"` // 认证用户名
}

type AuthResp struct {
	Resp
	Data []Auth
}
