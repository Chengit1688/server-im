package middleware

import (
	"bytes"
	"encoding/json"
	"im/pkg/logger"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
)

type Op struct {
	OperationId string `json:"operation_id"`
}

func Cost() gin.HandlerFunc {
	return func(context *gin.Context) {
		nowTime := time.Now()
		// ioutil.ReadAll读取到的是[]byte,读完body就没有了
		body, _ := ioutil.ReadAll(context.Request.Body)
		// 使用ioutil.NopCloser重新赋值给body
		context.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
		context.Next()
		costTime := time.Since(nowTime)
		if costTime > time.Duration(1*time.Second) {
			url := context.Request.URL.String()
			opInfo := Op{}
			json.Unmarshal(body, &opInfo)
			logger.Sugar.Debugf("the request URL %s operation_id %s cost %v\n", url, opInfo.OperationId, costTime)
		}
	}
}
