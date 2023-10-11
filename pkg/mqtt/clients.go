package mqtt

import (
	"errors"
	"fmt"
	"im/config"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"
	"net/url"
	"strconv"
	"time"
)

const defaultTimeout = 30

// 获取客户端信息
func GetClients(username string, likeUsername string, likeClientID string, connState string, page int, limit int) (clients []Client, err error) {
	start := time.Now()
	defer func() {
		duration := time.Now().Sub(start).Milliseconds()
		if duration > 1000 {
			logger.Sugar.Debugw("", "func", "GetClients", "info", fmt.Sprintf("emqx cost: %d ms", duration))
		}
	}()

	values := url.Values{}

	if username != "" {
		values.Add("username", username)
	}

	if likeUsername != "" {
		values.Add("_like_username", likeUsername)
	}

	if likeClientID != "" {
		values.Add("_like_clientid", likeClientID)
	}

	if connState != "" {
		values.Add("conn_state", connState)
	}

	// 查询所有数据
	if page == 0 && limit == 0 {
		page = 1
		limit = 99999999
	}

	values.Add("_page", strconv.Itoa(page))
	values.Add("_limit", strconv.Itoa(limit))

	cfg := config.Config.EMQXServer
	u := fmt.Sprintf("%s%s?%s", cfg.APIAddress, "/api/clients", values.Encode())
	data, err := http.GetWithBasicAuth(u, cfg.APIUsername, cfg.APIPassword, defaultTimeout)
	if err != nil {
		err = errors.New(fmt.Sprintf("get clients error, http get error, error: %v", err))
		return
	}

	var resp ClientResp
	if err = util.JsonUnmarshal(data, &resp); err != nil {
		err = errors.New(fmt.Sprintf("get clients error, json unmarshal error, error: %v", err))
		return
	}

	if resp.Code != 0 {
		err = errors.New(fmt.Sprintf("get clients error, code: %d", resp.Code))
		return
	}

	clients = resp.Data
	return
}

// 删除指定客户端，连接和会话都会终结
func DeleteClient(clientID string) (err error) {
	cfg := config.Config.EMQXServer

	u := fmt.Sprintf("%s%s/%s", cfg.APIAddress, "/api/clients", clientID)
	data, err := http.DeleteWithBasicAuth(u, cfg.APIUsername, cfg.APIPassword, defaultTimeout)
	if err != nil {
		err = errors.New(fmt.Sprintf("delete client error, http delete with basic auth error, error: %v", err))
		return
	}

	resp := struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}{}
	if err = util.JsonUnmarshal(data, &resp); err != nil {
		err = errors.New(fmt.Sprintf("delete client error, json unmarshal error, error: %v", err))
		return
	}

	if resp.Code != 0 {
		err = errors.New(fmt.Sprintf("delete client error, code: %d", resp.Code))
		return
	}
	return
}
