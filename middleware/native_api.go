package middleware

import (
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/server"
	"github.com/duanchi/min/server/handler"
	"github.com/duanchi/min/types"
	"github.com/gin-gonic/gin"
	"github.com/duanchi/min-gateway/native_api"
	"net/http"
)

type NativeApiMiddleware struct {
	abstract.Middleware
	Token string `value:"${Gateway.NativeApi.AccessToken}"`
}

func (this *NativeApiMiddleware) Execute (ctx *gin.Context) {

	resource, _ := ctx.Get("NATIVE_API_RESOURCE")
	token := ctx.GetHeader("X-Heron-Authorization")

	if this.Token == token {
		if bean, ok := native_api.NativeApiBeans[resource.(string)]; ok {
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
