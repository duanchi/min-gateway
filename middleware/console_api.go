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

	if this.Token == token {
		if bean, ok := console_api.ConsoleApiBeans[resource.(string)]; ok {
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
