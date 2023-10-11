package push

import (
	"im/pkg/logger"
	"im/pkg/util"
	"time"

	"github.com/Scorpio69t/jpush-api-golang-client"
)

var (
	JPushClient  *jpush.JPushClient
	errLimitTime time.Time
	ch           chan *JPushData
)

type JPushData struct {
	OperationID string
	UserIDs     []string
	Alert       string
	Title       string
}

func Init(appkey, masterSecret string) {
	JPushClient = jpush.NewJPushClient(appkey, masterSecret)
	ch = make(chan *JPushData, 10000)
	go func() {
		for {
			data := <-ch
			doPush(data)
		}
	}()
}

// 极光推送
func Jpush(UserIDs []string, alert, title, operationID string) {
	data := &JPushData{OperationID: operationID, UserIDs: UserIDs, Alert: alert, Title: title}
	ch <- data
}

// 极光推送
func doPush(data *JPushData) (err error) {
	// 推送平台设置
	var pf jpush.Platform
	pf.Add(jpush.ANDROID)
	//pf.Add(jpush.IOS)

	// 推送对象设置，这里使用用户ID
	var at jpush.Audience
	//at.SetID(UserIDs)
	//alias 别名 极光推送匹配用户的方式之一
	at.SetAlias(data.UserIDs)

	// Notification 通知
	var n jpush.Notification
	n.SetAlert(data.Alert)
	n.SetAndroid(&jpush.AndroidNotification{Alert: data.Alert, Title: data.Title, DisplayForeground: "1"})
	// Message 消息
	//var m jpush.Message
	//m.MsgContent = "This is a message"
	//m.Title = "Hello"

	payload := jpush.NewPayLoad()
	payload.SetPlatform(&pf)
	payload.SetAudience(&at)
	payload.SetNotification(&n)
	//payload.SetMessage(&m)

	payloadData, err := payload.Bytes()
	if err != nil {
		logger.Sugar.Errorw(data.OperationID, "func", util.GetSelfFuncName(), "error", err)
		return
	}
	res, err := JPushClient.Push(payloadData)
	if err != nil {
		logger.Sugar.Errorw(data.OperationID, "func", util.GetSelfFuncName(), "error", err, "message", "[JPUSH]")
		return
	} else {
		logger.Sugar.Debugw(data.OperationID, "func", util.GetSelfFuncName(), "msg", string(res))
		return
	}
}
