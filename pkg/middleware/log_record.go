package middleware

import (
	"bytes"
	"encoding/json"
	"im/config"
	adminModel "im/internal/cms_api/admin/model"
	"im/internal/cms_api/log/model"
	"im/internal/cms_api/log/usecase"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UrlLogMapItem struct {
	Url            string
	Method         string
	Remark         string
	Record         bool
	SkipRecordBody bool
}

var CmsUrlLogMap = []UrlLogMapItem{
	{Url: "/cms_api/admin/login", Method: "POST", Remark: "管理员登录", Record: true, SkipRecordBody: true},
	{Url: "/cms_api/admin/refresh_token", Method: "POST", Remark: "管理员TOKEN刷新", Record: true},
	{Url: "/cms_api/admin/get_info", Method: "GET", Remark: "管理员获取菜单"},
	{Url: "/cms_api/admin", Method: "GET", Remark: "管理员列表分页"},
	{Url: "/cms_api/admin", Method: "POST", Remark: "添加管理员", Record: true},
	{Url: "/cms_api/admin/update_info", Method: "PUT", Remark: "修改管理员账号基本信息", Record: true},
	{Url: "/cms_api/admin/update_password", Method: "PUT", Remark: "修改管理员账号密码", Record: true},
	{Url: "/cms_api/admin", Method: "DELETE", Remark: "删除管理员", Record: true},
	{Url: "/cms_api/admin/get_google_secret", Method: "GET", Remark: "获取谷歌验证码密钥", Record: true},
	{Url: "/cms_api/admin/ip_white_list", Method: "GET", Remark: "ip白名单分页查询"},
	{Url: "/cms_api/admin/ip_white_list", Method: "POST", Remark: "ip白名单新增", Record: true},
	{Url: "/cms_api/admin/ip_white_list/:id", Method: "PUT", Remark: "ip白名单修改", Record: true},
	{Url: "/cms_api/admin/ip_white_list/:id", Method: "DELETE", Remark: "ip白名单删除", Record: true},
	{Url: "/cms_api/role/:id", Method: "GET", Remark: "获取角色信息"},
	{Url: "/cms_api/role/:id", Method: "PUT", Remark: "修改角色信息", Record: true},
	{Url: "/cms_api/role/:id", Method: "DELETE", Remark: "删除角色", Record: true},
	{Url: "/cms_api/role", Method: "POST", Remark: "添加角色", Record: true},
	{Url: "/cms_api/role", Method: "GET", Remark: "角色列表分页"},
	{Url: "/cms_api/config/invite_add", Method: "POST", Remark: "新增邀请码", Record: true},
	{Url: "/cms_api/config/invite_delete", Method: "POST", Remark: "删除邀请码", Record: true},
	{Url: "/cms_api/config/invite_update", Method: "POST", Remark: "编辑邀请码", Record: true},
	{Url: "/cms_api/config/invite_update_friends", Method: "POST", Remark: "编辑邀请码", Record: true},
	{Url: "/cms_api/config/invite_update_groups", Method: "POST", Remark: "编辑邀请码", Record: true},
	{Url: "/cms_api/config/invite_list", Method: "POST", Remark: "邀请码列表"},
	{Url: "/cms_api/config/invite_update_status", Method: "POST", Remark: "修改邀请码状态", Record: true},
	{Url: "/cms_api/config/version_list", Method: "POST", Remark: "获取版本列表"},
	{Url: "/cms_api/config/version_add", Method: "POST", Remark: "新增版本", Record: true},
	{Url: "/cms_api/config/version_delete", Method: "POST", Remark: "删除版本", Record: true},
	{Url: "/cms_api/config/version_update", Method: "POST", Remark: "编辑版本", Record: true},
	{Url: "/cms_api/config/version_update_status", Method: "POST", Remark: "编辑版本状态", Record: true},
	{Url: "/cms_api/config/default_account_add", Method: "POST", Remark: "新增默认好友", Record: true},
	{Url: "/cms_api/config/default_account_update", Method: "POST", Remark: "编辑默认好友", Record: true},
	{Url: "/cms_api/config/default_account_delete", Method: "POST", Remark: "删除默认好友", Record: true},
	{Url: "/cms_api/config/default_account_list", Method: "POST", Remark: "获取默认好友列表"},
	{Url: "/cms_api/config/shield_add", Method: "POST", Remark: "添加敏感词", Record: true},
	{Url: "/cms_api/config/shield_update", Method: "POST", Remark: "编辑敏感词", Record: true},
	{Url: "/cms_api/config/shield_list", Method: "POST", Remark: "敏感词列表"},
	{Url: "/cms_api/config/shield_delete", Method: "POST", Remark: "删除敏感词", Record: true},
	{Url: "/cms_api/config/info", Method: "GET", Remark: "获取后管通用配置"},
	{Url: "/cms_api/config/login_config", Method: "GET", Remark: "获取登录配置"},
	{Url: "/cms_api/config/register_config", Method: "GET", Remark: "获取注册配置"},
	{Url: "/cms_api/config/login_config", Method: "POST", Remark: "更新登录配置", Record: true},
	{Url: "/cms_api/config/register_config", Method: "POST", Remark: "更新注册配置", Record: true},
	{Url: "/cms_api/config/sign_config_handle", Method: "POST", Remark: "更新签到配置", Record: true},
	{Url: "/cms_api/config/get_sign_config", Method: "GET", Remark: "获取签到配置"},
	{Url: "/cms_api/config/google_code_is_open", Method: "POST", Remark: "修改谷歌验证开关", Record: true},
	{Url: "/cms_api/config/site_ui", Method: "POST", Remark: "设置界面UI", Record: true},
	{Url: "/cms_api/config/jpush", Method: "POST", Remark: "设置极光推送配置", Record: true},
	{Url: "/cms_api/config/jpush", Method: "GET", Remark: "获取极光推送配置"},
	{Url: "/cms_api/config/feihu", Method: "POST", Remark: "设置飞虎配置", Record: true},
	{Url: "/cms_api/config/feihu", Method: "GET", Remark: "获取飞虎配置"},
	{Url: "/cms_api/config/parameter_config", Method: "GET", Remark: "获取参数配置"},
	{Url: "/cms_api/config/parameter_config_update", Method: "POST", Remark: "修改参数配置", Record: true},
	{Url: "/cms_api/config/deposite", Method: "GET", Remark: "获取充值配置"},
	{Url: "/cms_api/config/deposite", Method: "POST", Remark: "修改充值配置", Record: true},
	{Url: "/cms_api/config/withdraw", Method: "GET", Remark: "获取提现配置"},
	{Url: "/cms_api/config/withdraw", Method: "POST", Remark: "修改提现配置", Record: true},
	{Url: "/cms_api/config/about_us", Method: "GET", Remark: "获取关于我们"},
	{Url: "/cms_api/config/about_us", Method: "POST", Remark: "修改关于我们", Record: true},
	{Url: "/cms_api/config/privacy_policy", Method: "GET", Remark: "获取隐私政策"},
	{Url: "/cms_api/config/privacy_policy", Method: "POST", Remark: "修改隐私政策", Record: true},
	{Url: "/cms_api/config/user_agreement", Method: "GET", Remark: "获取用户协议"},
	{Url: "/cms_api/config/user_agreement", Method: "POST", Remark: "修改用户协议", Record: true},
	{Url: "/cms_api/config/ip_white_list_is_open", Method: "GET", Remark: "获取IP白名单开关"},
	{Url: "/cms_api/config/ip_white_list_is_open", Method: "POST", Remark: "修改IP白名单开关", Record: true},
	{Url: "/cms_api/menu", Method: "GET", Remark: "获取全量菜单"},
	{Url: "/cms_api/group/create_group", Method: "POST", Remark: "创建群", Record: true},
	{Url: "/cms_api/group/remove_group_member", Method: "POST", Remark: "删除群成员", Record: true},
	{Url: "/cms_api/group/group_info", Method: "POST", Remark: "获取群信息"},
	{Url: "/cms_api/group/group_member_list", Method: "POST", Remark: "获取群成员列表"},
	{Url: "/cms_api/group/group_list", Method: "POST", Remark: "群列表"},
	{Url: "/cms_api/group/set_admin", Method: "POST", Remark: "群设置管理员", Record: true},
	{Url: "/cms_api/group/set_owner", Method: "POST", Remark: "群设置创建者", Record: true},
	{Url: "/cms_api/group/set_robot", Method: "POST", Remark: "群设置机器人", Record: true},
	{Url: "/cms_api/group/add_group_members", Method: "POST", Remark: "添加群成员", Record: true},
	{Url: "/cms_api/group/group_merge", Method: "POST", Remark: "合并群", Record: true},
	{Url: "/cms_api/group/batch_join_group", Method: "POST", Remark: "批量入群", Record: true},
	{Url: "/cms_api/group/search", Method: "GET", Remark: "群搜索"},
	{Url: "/cms_api/user", Method: "GET", Remark: "用户列表分页"},
	{Url: "/cms_api/user/export", Method: "GET", Remark: "用户列表导出"},
	{Url: "/cms_api/user/search", Method: "GET", Remark: "用户搜索"},
	{Url: "/cms_api/user", Method: "POST", Remark: "用户批量新增", Record: true},
	{Url: "/cms_api/user/details", Method: "GET", Remark: "用户详情"},
	{Url: "/cms_api/user/updates", Method: "POST", Remark: "用户编辑资料", Record: true},
	{Url: "/cms_api/user/freeze", Method: "POST", Remark: "用户冻结", Record: true},
	{Url: "/cms_api/user/unfreeze", Method: "POST", Remark: "用户解冻", Record: true},
	{Url: "/cms_api/user/set_password", Method: "POST", Remark: "用户密码修改", Record: true},
	{Url: "/cms_api/user/sign_log_list", Method: "POST", Remark: "签到记录列表"},
	{Url: "/cms_api/user/privilege_user_list", Method: "POST", Remark: "特权用户列表"},
	{Url: "/cms_api/user/privilege_user_add", Method: "POST", Remark: "新增特权用户", Record: true},
	{Url: "/cms_api/user/privilege_user_remove", Method: "POST", Remark: "删除特权用户", Record: true},
	{Url: "/cms_api/user/disabled/user", Method: "GET", Remark: "封禁管理-用户分页查询"},
	{Url: "/cms_api/user/disabled/device", Method: "GET", Remark: "封禁管理-设备分页查询"},
	{Url: "/cms_api/user/disabled/device/disable", Method: "POST", Remark: "封禁管理-设备封禁", Record: true},
	{Url: "/cms_api/user/disabled/device/enable", Method: "POST", Remark: "封禁管理-设备解封", Record: true},
	{Url: "/cms_api/user/disabled/ip", Method: "GET", Remark: "封禁管理-IP分页查询"},
	{Url: "/cms_api/third/upload", Method: "POST", Remark: "文件上传v1", SkipRecordBody: true},
	{Url: "/cms_api/third/upload/v2", Method: "POST", Remark: "文件上传v2", SkipRecordBody: true},
	{Url: "/cms_api/announcement", Method: "GET", Remark: "公告获取"},
	{Url: "/cms_api/announcement", Method: "POST", Remark: "公告修改", Record: true},
	{Url: "/cms_api/discover", Method: "GET", Remark: "获取发现页"},
	{Url: "/cms_api/discover", Method: "POST", Remark: "新增发现页", Record: true},
	{Url: "/cms_api/discover/:id", Method: "PUT", Remark: "修改发现页", Record: true},
	{Url: "/cms_api/discover/:id", Method: "DELETE", Remark: "删除发现页", Record: true},
	{Url: "/cms_api/discover/status", Method: "GET", Remark: "获取发现页开关"},
	{Url: "/cms_api/discover/status", Method: "POST", Remark: "修改发现页开关", Record: true},
	{Url: "/cms_api/friend/user_friend_list", Method: "POST", Remark: "用户好友列表"},
	{Url: "/cms_api/friend/user_add_friend", Method: "POST", Remark: "用户添加好友", Record: true},
	{Url: "/cms_api/friend/user_remove_friend", Method: "POST", Remark: "用户删除好友", Record: true},
	{Url: "/cms_api/operation/reg_statistics", Method: "POST", Remark: "用户注册统计"},
	{Url: "/cms_api/operation/online_statistics", Method: "GET", Remark: "在线用户统计"},
	{Url: "/cms_api/operation/single_msg_statistics", Method: "GET", Remark: "单聊消息统计"},
	{Url: "/cms_api/operation/group_msg_statistics", Method: "GET", Remark: "群聊消息统计"},
	{Url: "/cms_api/operation/invite_code_statistics_list", Method: "POST", Remark: "渠道码统计列表"},
	{Url: "/cms_api/operation/invite_code_statistics_details", Method: "POST", Remark: "渠道码统计详情"},
	{Url: "/cms_api/operation/suggestion_list", Method: "GET", Remark: "在线举报列表"},
	{Url: "/cms_api/operation/ip_black_list", Method: "GET", Remark: "IP黑名单列表查询"},
	{Url: "/cms_api/operation/ip_black_list", Method: "POST", Remark: "添加IP黑名单", Record: true},
	{Url: "/cms_api/operation/ip_black_list/batch", Method: "POST", Remark: "批量添加IP黑名单", Record: true},
	{Url: "/cms_api/operation/ip_black_list/:id", Method: "PUT", Remark: "修改IP黑名单", Record: true},
	{Url: "/cms_api/operation/ip_black_list/batch", Method: "DELETE", Remark: "批量删除IP黑名单", Record: true},
	{Url: "/cms_api/operation/ip_black_list/:id", Method: "DELETE", Remark: "删除IP黑名单", Record: true},
	{Url: "/cms_api/dashboard", Method: "GET", Remark: "首页"},
	{Url: "/cms_api/chat/message_history", Method: "GET", Remark: "消息历史记录"},
	{Url: "/cms_api/chat/message_change", Method: "POST", Remark: "消息状态修改", Record: true},
	{Url: "/cms_api/chat/message_clear", Method: "POST", Remark: "消息客户端删除", Record: true},
	{Url: "/cms_api/chat/multi_send", Method: "POST", Remark: "消息群发", Record: true},
	{Url: "/cms_api/chat/multi_send_records", Method: "GET", Remark: "消息群发-群发记录"},
	{Url: "/cms_api/wallet/billing_records", Method: "GET", Remark: "账单记录"},
	{Url: "/cms_api/wallet/billing_records/export", Method: "GET", Remark: "账单记录导出"},
	{Url: "/cms_api/wallet/change_amount", Method: "POST", Remark: "余额调整", Record: true},
	{Url: "/cms_api/wallet/redpack_single_records", Method: "GET", Remark: "个人红包记录"},
	{Url: "/cms_api/wallet/redpack_single_records/export", Method: "GET", Remark: "个人红包记录导出"},
	{Url: "/cms_api/wallet/withdraw_records", Method: "GET", Remark: "提现记录"},
	{Url: "/cms_api/wallet/withdraw_records/count_pending", Method: "GET", Remark: "提现记录待处理数"},
	{Url: "/cms_api/wallet/withdraw_records/set_status", Method: "POST", Remark: "提现记录审核", Record: true},
	{Url: "/cms_api/wallet/withdraw_records/:id", Method: "GET", Remark: "通过ID获取提现详情"},
	{Url: "/cms_api/wallet/set_paypasswd", Method: "POST", Remark: "设置支付密码", Record: true},
}

func CmsOperateLog() gin.HandlerFunc {
	var serviceIP, servicePort, serviceID, env, logLevel string
	var statusCode int
	serviceIP, _ = util.LocalIP()
	cfg := config.Config
	servicePort = cfg.Server.CmsApiListenAddr
	servicePort = strings.Replace(servicePort, "0.0.0.0", serviceIP, 1)
	serviceID = "im_cms_api"
	env = cfg.Station
	logLevel = cfg.Log.Level
	return func(context *gin.Context) {
		nowTime := time.Now()
		var url, operationID, params, userAgent, remark, method, rBody, userID string
		var record, skipRecordBody bool
		ctx := context.Copy()
		body, _ := ioutil.ReadAll(context.Request.Body)
		context.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
		blw := &CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: context.Writer}
		context.Writer = blw
		context.Next()
		url = ctx.Request.URL.Path
		userAgent = ctx.Request.UserAgent()
		params = ctx.Request.URL.RawQuery
		method = ctx.Request.Method
		userID = ctx.GetString("o_user_id")
		statusCode = context.Writer.Status()
		// 登录接口从请求上下文上拿不到管理员ID 改从响应里拿
		if userID == "" {
			hResp := http.Resp{}
			var lResp adminModel.LoginResp
			err := json.Unmarshal(blw.body.Bytes(), &hResp)
			if err == nil && hResp.Code == 0 {
				v, _ := json.Marshal(hResp.Data)
				json.Unmarshal(v, &lResp)
				userID, _, _, _ = util.CmsParseToken(lResp.Token)
			}
		}
		for _, i := range CmsUrlLogMap {
			if util.KeyMatch(url, i.Url) && method == i.Method {
				remark = i.Remark
				record = i.Record
				skipRecordBody = i.SkipRecordBody
				break
			}
		}
		if !record {
			return
		}
		switch method {
		case "GET":
			operationID = ctx.Request.URL.Query().Get("operation_id")
		case "POST", "PUT", "DELETE":
			opInfo := Op{}
			json.Unmarshal(body, &opInfo)
			operationID = opInfo.OperationId
			if !skipRecordBody {
				rBody = string(body)
			}
		default:
			return
		}
		log := model.OperateLogs{
			CreatedAt:          nowTime.UnixMilli(),
			UserID:             userID,
			ServiceID:          serviceID,
			Env:                env,
			LogLevel:           logLevel,
			RequestUrl:         url,
			RequestParams:      params,
			RequestBody:        rBody,
			RequestMethod:      method,
			RequestUserAgent:   userAgent,
			ResponseStatusCode: statusCode,
			ServiceIp:          serviceIP,
			ServiceHost:        servicePort,
			OperationID:        operationID,
			LogRemark:          remark,
		}

		err := usecase.LogUseCase.AddOperateLog(log)
		if err != nil {
			logger.Sugar.Errorw(util.GetSelfFuncName(), "AddOperateLog error:", err)
		}

	}
}

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
