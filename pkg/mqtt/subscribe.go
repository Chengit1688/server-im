package mqtt

import (
	"regexp"
	"strings"
)

// 客户端上下线信息结构体
type ClientConnPayload struct {
	Username string `json:"username"` //站点标识+用户ID
	ClientID string `json:"clientid"`
	Ts       int64  `json:"ts"` //13位 毫秒级时间戳
}

// 订阅主题
const ClientConnectedTopic = "$SYS/brokers/+/clients/+/connected"       //客户端上线
const ClientDisconnectedTopic = "$SYS/brokers/+/clients/+/disconnected" //客户端离线线

// 主题正则统配
var ClientConnectedTopicRx = regexp.MustCompile(topicFormat(ClientConnectedTopic))
var ClientDisconnectedTopicRx = regexp.MustCompile(topicFormat(ClientDisconnectedTopic))

// 主题名正则处理
func topicFormat(str string) string {
	str = strings.Replace(str, "+", ".*", -1) //处理加号
	str = strings.Replace(str, "$", ".", -1)  //处理系统主题首字符
	return str
}
