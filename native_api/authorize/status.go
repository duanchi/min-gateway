package authorize

import (
	"github.com/duanchi/min-gateway/service"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/types"
	"github.com/duanchi/min/types/error"
	"github.com/gin-gonic/gin"
	"strings"
)

type StatusController struct {
	abstract.RestController

	TokenService *service.TokenService `autowired:"true"`
}

func (this *StatusController) Fetch(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {
	isCustomToken := ctx.Query("custom_token") == "true"
	prefix := ctx.Query("prefix")
	accessTokenRaw := ctx.GetHeader("Authorization")
	accessToken := ""
	if accessTokenRaw == "" {
		accessToken = ctx.Query("token")
	} else if isCustomToken {
		accessToken = accessTokenRaw
	} else {
		accessTokenStack := strings.Split(accessTokenRaw, " ")
		if len(accessTokenStack) > 1 && accessTokenStack[0] == "Bearer" {
			accessToken = accessTokenStack[1]
		}
	}

	if isCustomToken {

		if ok, more, err := this.TokenService.CustomAuth(accessToken, prefix); ok {
			if err != nil {
				panic(err)
			}

			return map[string]interface{}{
				"user":  accessToken,
				"token": accessToken,
				"more":  more,
			}, nil
		} else {
			return nil, error.AuthorizeError{
				Message:    "授权过期或错误",
				ErrorCode:  100401,
				StatusCode: 401,
				ErrorData:  nil,
			}
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

			return map[string]interface{}{
				"user":  userId,
				"token": claims.Id,
				"more":  more,
			}, nil

		} else {
			return nil, error.AuthorizeError{
				Message:    "授权过期或错误",
				ErrorCode:  100401,
				StatusCode: 401,
				ErrorData:  nil,
			}
		}
	}
	return "Ok", nil
}
