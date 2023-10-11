package sms

import (
	"im/config"
	"im/pkg/util"

	unisms "github.com/apistd/uni-go-sdk/sms"
)

func SendSmsCode(phone_number, code string) error {
	smsInfo := config.Config.Sms
	// 初始化
	client := unisms.NewClient(smsInfo.AppID, smsInfo.AppSecret)

	// 构建信息
	message := unisms.BuildMessage()
	message.SetTo(phone_number)
	message.SetSignature(smsInfo.Signature)
	message.SetTemplateId(smsInfo.TemplateId)
	message.SetTemplateData(map[string]string{"code": code, "ttl": util.IntToString(smsInfo.ExpireTTL)}) // 设置自定义参数 (变量短信)

	// 发送短信
	_, err := client.Send(message)
	if err != nil {
		return err
	}
	return nil
}
