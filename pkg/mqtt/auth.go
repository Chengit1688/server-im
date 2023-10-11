package mqtt

import (
	"errors"
	"fmt"
	"im/config"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"
	"time"
)

// 获取客户端认证信息
func GetAuthUsernames() (auths []Auth, err error) {
	start := time.Now()
	defer func() {
		duration := time.Now().Sub(start).Milliseconds()
		if duration > 1000 {
			logger.Sugar.Debugw("", "func", "GetAuthUsernames", "info", fmt.Sprintf("emqx cost: %d ms", duration))
		}
	}()

	cfg := config.Config.EMQXServer
	u := fmt.Sprintf("%s%s", cfg.APIAddress, "/api/auth_username")
	data, err := http.GetWithBasicAuth(u, cfg.APIUsername, cfg.APIPassword, defaultTimeout)
	if err != nil {
		err = errors.New(fmt.Sprintf("get auth usernames error, http get error, error: %v", err))
		return
	}

	var resp AuthResp
	if err = util.JsonUnmarshal(data, &resp); err != nil {
		err = errors.New(fmt.Sprintf("get auth usernames error, json unmarshal error, error: %v", err))
		return
	}

	if resp.Code != 0 {
		err = errors.New(fmt.Sprintf("get auth usernames error, code: %d", resp.Code))
		return
	}

	auths = resp.Data
	return
}

// 创建客户端认证账号
func CreateAuthUsername(username string, password string) (err error) {
	start := time.Now()
	defer func() {
		duration := time.Now().Sub(start).Milliseconds()
		if duration > 1000 {
			logger.Sugar.Debugw("", "func", "CreateAuthUsername", "info", fmt.Sprintf("emqx cost: %d ms", duration))
		}
	}()

	data := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	data.Username = username
	data.Password = password

	cfg := config.Config.EMQXServer
	u := fmt.Sprintf("%s%s", cfg.APIAddress, "/api/auth_username")
	content, err := http.PostWithBasicAuth(u, data, cfg.APIUsername, cfg.APIPassword, defaultTimeout)
	if err != nil {
		err = errors.New(fmt.Sprintf("emqx create auth username error, post error, error: %v", err))
		return
	}

	resp := struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}{}

	if err = util.JsonUnmarshal(content, &resp); err != nil {
		err = errors.New(fmt.Sprintf("emqx create auth username error, json unmarshal error, error: %v", err))
		return
	}

	if resp.Code != 0 {
		err = errors.New(fmt.Sprintf("emqx create auth username error, code: %d, message: %s", resp.Code, resp.Message))
	}
	return
}

// 删除客户端认证账号
func DeleteAuthUsername(username string) (err error) {
	cfg := config.Config.EMQXServer

	u := fmt.Sprintf("%s%s/%s", cfg.APIAddress, "/api/auth_username", username)
	data, err := http.DeleteWithBasicAuth(u, cfg.APIUsername, cfg.APIPassword, defaultTimeout)
	if err != nil {
		err = errors.New(fmt.Sprintf("delete auth username error, http delete with basic auth error, error: %v", err))
		return
	}

	resp := struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}{}
	if err = util.JsonUnmarshal(data, &resp); err != nil {
		err = errors.New(fmt.Sprintf("delete auth username error, json unmarshal error, error: %v", err))
		return
	}

	if resp.Code != 0 {
		err = errors.New(fmt.Sprintf("delete auth username error, code: %d", resp.Code))
		return
	}
	return
}
