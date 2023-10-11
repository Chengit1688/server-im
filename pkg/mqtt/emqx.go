package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eclipse/paho.golang/paho"
	"im/config"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"
	"net"
	"time"
)

// 连接EMQX服务器，创建永久session，成功后断开session连接
func Connect(userID string, clientID string) (err error) {
	start := time.Now()
	defer func() {
		duration := time.Now().Sub(start).Milliseconds()
		if duration > 1000 {
			logger.Sugar.Debugw("", "func", "Connect", "info", fmt.Sprintf("emqx cost: %d ms", duration))
		}
	}()

	cfg := config.Config.EMQXServer
	username := fmt.Sprintf("%s_%s", config.Config.Station, userID)
	password := "root1234"

	var conn net.Conn
	if conn, err = net.Dial("tcp", cfg.MQTTAddress); err != nil {
		err = errors.New(fmt.Sprintf("emqx connect error, dail error, error: %v", err))
		return
	}
	defer conn.Close()

	c := paho.NewClient(paho.ClientConfig{
		Conn: conn,
	})

	defer func() {
		d := paho.Disconnect{
			ReasonCode: 0,
		}
		if err = c.Disconnect(&d); err != nil {
			logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("emqx connect error, disconnect error, error: %v", err))
		}
	}()

	var interval uint32 = 0xFFFFFFFF
	cp := paho.Connect{
		KeepAlive:    30,
		UsernameFlag: true,
		Username:     username,
		PasswordFlag: true,
		Password:     []byte(password),
		CleanStart:   false,
		ClientID:     clientID,
		Properties: &paho.ConnectProperties{
			SessionExpiryInterval: &interval,
		},
	}

	resp, err := c.Connect(context.Background(), &cp)
	if err != nil {
		err = errors.New(fmt.Sprintf("emqx connect error, connect error, error: %v", err))
		return
	}

	if resp.ReasonCode != 0 {
		err = errors.New(fmt.Sprintf("emqx connect error, code error, code: %d", resp.ReasonCode))
		return
	}
	return
}

// 订阅
func Subscribe(topic string, qos int, clientID string) (err error) {
	start := time.Now()
	defer func() {
		duration := time.Now().Sub(start).Milliseconds()
		if duration > 1000 {
			logger.Sugar.Debugw("", "func", "Subscribe", "info", fmt.Sprintf("emqx cost: %d ms", duration))
		}
	}()

	data := struct {
		Topic    string `json:"topic"`
		Qos      int    `json:"qos"`
		ClientID string `json:"clientid"`
	}{}
	data.Topic = topic
	data.Qos = qos
	data.ClientID = clientID

	cfg := config.Config.EMQXServer
	url := fmt.Sprintf("%s%s", cfg.APIAddress, "/api/mqtt/subscribe")
	content, err := http.PostWithBasicAuth(url, data, cfg.APIUsername, cfg.APIPassword, defaultTimeout)
	if err != nil {
		err = errors.New(fmt.Sprintf("emqx subscribe error, error: %v", err))
		return
	}

	resp := struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}{}

	if err = json.Unmarshal(content, &resp); err != nil {
		err = errors.New(fmt.Sprintf("emqx subscribe error, json unmarshal error, error: %v", err))
		return
	}

	if resp.Code != 0 {
		err = errors.New(fmt.Sprintf("emqx subscribe error, code error, code: %d, message: %s", resp.Code, resp.Message))
		return
	}
	return
}

// 取消订阅
func Unsubscribe(topic string, clientID string) (err error) {
	start := time.Now()
	defer func() {
		duration := time.Now().Sub(start).Milliseconds()
		if duration > 1000 {
			logger.Sugar.Debugw("", "func", "Unsubscribe", "info", fmt.Sprintf("emqx cost: %d ms", duration))
		}
	}()

	data := struct {
		Topic    string `json:"topic"`
		ClientID string `json:"clientid"`
	}{}
	data.Topic = topic
	data.ClientID = clientID

	cfg := config.Config.EMQXServer
	url := fmt.Sprintf("%s%s", cfg.APIAddress, "/api/mqtt/unsubscribe")
	content, err := http.PostWithBasicAuth(url, data, cfg.APIUsername, cfg.APIPassword, defaultTimeout)
	if err != nil {
		err = errors.New(fmt.Sprintf("emqx unsubscribe error, error: %v", err))
		return
	}

	resp := struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}{}

	if err = json.Unmarshal(content, &resp); err != nil {
		err = errors.New(fmt.Sprintf("emqx unsubscribe error, json unmarshal error, error: %v", err))
		return
	}

	if resp.Code != 0 {
		err = errors.New(fmt.Sprintf("emqx unsubscribe error, code error, code: %d, message: %s", resp.Code, resp.Message))
		return
	}
	return
}

// 发布
func Publish(topic string, qos int, clientID string, payload string) (err error) {
	start := time.Now()
	defer func() {
		duration := time.Now().Sub(start).Milliseconds()
		if duration > 1000 {
			logger.Sugar.Debugw("", "func", "Publish", "info", fmt.Sprintf("emqx cost: %d ms", duration))
		}
	}()

	data := struct {
		Topic    string `json:"topic"`
		Qos      int    `json:"qos"`
		ClientID string `json:"clientid"`
		Payload  string `json:"payload"`
	}{}
	data.Topic = topic
	data.Qos = qos
	data.ClientID = clientID
	data.Payload = payload

	cfg := config.Config.EMQXServer

	url := fmt.Sprintf("%s%s", cfg.APIAddress, "/api/mqtt/publish")
	content, err := http.PostWithBasicAuth(url, data, cfg.APIUsername, cfg.APIPassword, defaultTimeout)
	if err != nil {
		err = errors.New(fmt.Sprintf("emqx publish error, error: %v", err))
		return
	}

	resp := struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}{}

	if err = json.Unmarshal(content, &resp); err != nil {
		err = errors.New(fmt.Sprintf("emqx publish error, json unmarshal error, error: %v", err))
		return
	}

	if resp.Code != 0 {
		err = errors.New(fmt.Sprintf("emqx publish error, code error, code: %d, message: %s", resp.Code, resp.Message))
		return
	}
	return
}

// 获取客户端唯一ID
func GetClientID(userID string, deviceID string) string {
	return fmt.Sprintf("%s_%s_%s", config.Config.Station, userID, deviceID)
}

// 获取用户唯一主题，消息推送
func GetUserTopic(userID string) string {
	return fmt.Sprintf("%s/users/%s", config.Config.Station, userID)
}

// 获取群主题，消息推送
func GetGroupTopic(groupID string) string {
	return fmt.Sprintf("%s/groups/%s", config.Config.Station, groupID)
}

// 获取系统主题，消息广播
func GetSystemTopic() string {
	return fmt.Sprintf("%s/system", config.Config.Station)
}
