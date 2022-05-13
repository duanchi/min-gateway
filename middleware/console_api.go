package middleware

import (
	"github.com/duanchi/min-gateway/console_api"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/server"
	"github.com/duanchi/min/server/handler"
	"github.com/duanchi/min/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ConsoleApiMiddleware struct {
	abstract.Middleware
	Token string `value:"${Gateway.ConsoleApi.AccessToken}"`
}

func (this *ConsoleApiMiddleware) Execute(ctx *gin.Context) {

	resource, _ := ctx.Get("CONSOLE_API_RESOURCE")
	token := ctx.GetHeader("X-Min-Gateway-Authorization")

	if ctx.Request.Method == "OPTIONS" {
		ctx.Header("Access-Control-Allow-Origin", ctx.GetHeader("Origin"))
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization,X-Min-Gateway-Authorization")
		return
	}

	if this.Token == token {
		if bean, ok := console_api.ConsoleApiBeans[resource.(string)]; ok {
			ctx.Header("Access-Control-Allow-Origin", ctx.GetHeader("Origin"))
			ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			ctx.Header("Access-Control-Allow-Headers", "DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization,X-Min-Gateway-Authorization")
			handler.RestfulHandle(resource.(string), bean, ctx, server.HttpServer)
		}
	} else {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, types.RuntimeError{
			Message:   "控制台授权认证失败",
			ErrorCode: http.StatusUnauthorized,
			ErrorData: nil,
		})
	}
}

func (this *ConsoleApiMiddleware) BeforeResponse(ctx *gin.Context) {
	/*if ctx.GetHeader("X-Min-Gateway-Authorization") != "" {
		ctx.Header("Access-Control-Allow-Origin", ctx.GetHeader("Origin"))
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization,X-Min-Gateway-Authorization")
	}*/
}
