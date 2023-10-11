package sms

import (
	"bytes"
	"fmt"
	"html/template"
	"im/config"
	"im/pkg/common/constant"
	"im/pkg/http"
	"net/url"
)

// 短信宝
type Bao struct {
	UserName string
	ApiKey   string
}

func NewBao() *Bao {
	return &Bao{
		UserName: config.Config.Sms.DxbUsername,
		ApiKey:   config.Config.Sms.DxbApiKey,
	}
}

var contentTemplate = `【止正Talk】您的验证码是{{ .Code }}，如非本人操作，请忽略本短信。`

/*
// Send 发短信
// @param area 国家编码
// @param code 验证码
*/
func (s *Bao) Send(phoneNumber, area, code string) error {
	if len(phoneNumber) < 3 {
		return fmt.Errorf("duan xin bao send message error, phone number error, phone number: %s", phoneNumber)
	}

	var buffer bytes.Buffer
	tmpl, err := template.New("template").Parse(contentTemplate)
	if err != nil {
		return fmt.Errorf("短信宝 send message error, generate template error, err: %v", err)
	}

	if err = tmpl.Execute(&buffer, struct {
		Code string
	}{Code: code}); err != nil {
		return fmt.Errorf("短信宝 send message error, generate template error, execute error, err: %v", err)
	}

	//area := string([]byte(phoneNumber)[:3])
	//number := string([]byte(phoneNumber)[3:])
	content := buffer.String()
	//fmt.Println(content)
	switch area {
	case constant.ChinaCountryCode:
		return s.sendChinese(phoneNumber, content)
	default:
		return s.sendInternational(fmt.Sprintf("%s%s", area, phoneNumber), content)
	}
}

func (s *Bao) sendChinese(phoneNumber string, content string) error {
	values := url.Values{}
	values.Add("u", s.UserName)
	values.Add("p", s.ApiKey)
	//values.Add("g", s.ApiKey)
	values.Add("m", phoneNumber) //专用通道产品时，需要指定产品ID
	values.Add("c", content)

	u := fmt.Sprintf("%s?%s", "https://api.smsbao.com/sms", values.Encode())
	resp, err := http.Get(u)
	if err != nil {
		return fmt.Errorf("duan xin bao send message error, http get error, err: %v", err)
	}
	codeMap := map[string]string{
		"30": "错误密码",
		"40": "账号不存在",
		"41": "余额不足",
		"43": "IP地址限制",
		"50": "内容含有敏感词",
		"51": "手机号码不正确",
	}
	if string(resp) != "0" {
		return fmt.Errorf("duan xin bao send message error, resp error code: %s", codeMap[string(resp)])
	}
	return nil
}

// 国际短信
func (s *Bao) sendInternational(phoneNumber string, content string) error {
	values := url.Values{}
	values.Add("u", s.UserName)
	values.Add("p", s.ApiKey)
	values.Add("m", phoneNumber)
	values.Add("c", content)

	u := fmt.Sprintf("%s?%s", "https://api.smsbao.com/wsms", values.Encode())
	fmt.Printf("url:%v", u)
	resp, err := http.Get(u)
	if err != nil {
		return fmt.Errorf("短信宝http get error, err: %v phone: %s", err, phoneNumber)
	}
	codeMap := map[string]string{
		"30": "错误密码",
		"40": "账号不存在",
		"41": "余额不足",
		"43": "IP地址限制",
		"50": "内容含有敏感词",
		"51": "手机号码不正确",
	}
	if string(resp) != "0" {
		return fmt.Errorf("短信宝 resp error code: %s phone:%s", codeMap[string(resp)], phoneNumber)
	}
	return nil
}
