package middleware

import (
	"github.com/duanchi/min-gateway/native_api"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/server"
	"github.com/duanchi/min/server/handler"
	"github.com/gin-gonic/gin"
)

type NativeApiMiddleware struct {
	abstract.Middleware
}

func (this *NativeApiMiddleware) Execute(ctx *gin.Context) {

	resource, _ := ctx.Get("NATIVE_API_RESOURCE")

	if bean, ok := native_api.NativeApiBeans[resource.(string)]; ok {
		handler.RestfulHandle(resource.(string), bean, ctx, server.HttpServer)
	}
}
