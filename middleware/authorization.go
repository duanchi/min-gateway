package middleware

import (
	"fmt"
	"github.com/duanchi/min-gateway/mapper"
	"github.com/duanchi/min-gateway/service"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/types"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strings"
)

type AuthorizationMiddleware struct {
	abstract.Middleware
	TokenService *service.TokenService `autowired:"true"`
}

func (this *AuthorizationMiddleware) AfterRoute(ctx *gin.Context) {

	/**
	需要进行token验证的
	*/
	routeValue, has := ctx.Get("route")

	if has {
		route := routeValue.(mapper.Route)

		if mapper.CONSTANT.IS_AUTHORIZE[route.NeedAuthorize] {
			accessTokenRaw := ctx.GetHeader("Authorization")
			accessToken := ""
			if accessTokenRaw == "" {
				accessToken = ctx.Query("token")
			} else if mapper.CONSTANT.IS_CUSTOM_TOKEN[route.CustomToken] {
				accessToken = accessTokenRaw
			} else {
				accessTokenStack := strings.Split(accessTokenRaw, " ")
				if len(accessTokenStack) > 1 && accessTokenStack[0] == "Bearer" {
					accessToken = accessTokenStack[1]
				}
			}

			prefix := "0000"
			if len(route.AuthorizePrefix) > 0 && len(route.AuthorizePrefix) < 4 {
				prefix = fmt.Sprintf("%0*s", 4, route.AuthorizePrefix)
			} else if len(route.AuthorizePrefix) >= 4 {
				prefix = route.AuthorizePrefix[0:4]
			}

			if mapper.CONSTANT.IS_CUSTOM_TOKEN[route.CustomToken] {

				if ok, more, err := this.TokenService.CustomAuth(accessToken, prefix); ok {
					if err != nil {
						panic(err)
					}
					ctx.Set("user", accessToken)
					ctx.Set("token", accessToken)
					ctx.Set("more", more)
					ctx.Next()
					return
				} else {
					requestId, _ := ctx.Get("REQUEST_ID")
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, types.Response{
						RequestId: requestId.(string),
						Code:      100401,
						Message:   "授权过期或错误",
						Data:      nil,
					})

					return
				}
			} else {
				if accessToken != "" && this.TokenService.Auth(accessToken, prefix) {

					claims, err := this.TokenService.Parse(accessToken)

					if err != nil {
						panic(err)
					}

					userId, more, err := this.TokenService.Fetch(claims.Id, prefix)

					if err != nil {
						panic(err)
					}

					ctx.Set("user", userId)
					ctx.Set("token", claims.Id)
					ctx.Set("more", more)
					ctx.Next()
					return

				} else {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, types.Response{
						RequestId: uuid.NewV4().String(),
						Code:      100401,
						Message:   "授权过期或错误",
						Data:      nil,
					})

					return
				}
			}
		}
	}

	ctx.Next()
	return
}
