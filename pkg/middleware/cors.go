package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin","*")
		context.Header("Access-Control-Allow-Methods", "*")
		context.Header("Access-Control-Allow-Headers", "im_site,Referer,Access-Control-Allow-Origin,Accept,Token,User-Agent,Accept-Encoding,Accept-Language,Content-Length,Content-Type,Connection,Locale")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
		context.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
		context.Header("Access-Control-Allow-Credentials", "true")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
		//context.Header("content-type", "application/json,text/html,multipart/form-data")                                                                                                                                                           // 设置返回格式是json
		//Release all option pre-requests
		if context.Request.Method == http.MethodOptions {
			context.JSON(http.StatusOK, "Options Request!")
		}
		context.Next()
	}
}
