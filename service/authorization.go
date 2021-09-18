package service

import (
	"encoding/json"
	"fmt"
	"github.com/duanchi/min-gateway/service/storage"
	types2 "github.com/duanchi/min-gateway/types"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/types"
	"github.com/duanchi/min/types/gateway"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AuthorizationService struct {
	abstract.Service

	JwtExpiresIn int64 `value:"${Authorization.Ttl}"`
	TokenService *TokenService `autowired:"true"`
	StorageService *storage.StorageService `autowired:"true"`
	ValueService *storage.ValuesService `autowired:"true"`
}

type authorizationResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int64 `json:"expires_in"`
	ExpireAt int64	`json:"expire_at"`
	RefreshToken string `json:"refresh_token"`
	RefreshTokenExpireAt int64 `json:"refresh_token_expire_at"`
}

type refreshAuthorizationRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (this *AuthorizationService) Handle(
	header *http.Header,
	gatewayData gateway.Data,
	requestBody []byte,
	singleton bool,
	authorizeType string,
	ctx *gin.Context,
) (
	status int,
	responseHeaders http.Header,
	response interface{},
	err error,
) {
	data := header.Get("X-Gateway-Authorization-Data")
	expiresIn, _ := strconv.ParseInt(header.Get("X-Gateway-Authorization-ExpiresIn"), 10, 64)
	moreDataRaw := header.Get("X-Gateway-Authorization-More")
	action := strings.ToUpper(header.Get("X-Gateway-Authorization-Action"))
	multi := header.Get("X-Gateway-Authorization-Remove-Multi")
	removePrefix := header.Get("X-Gateway-Authorization-Remove-Prefix")

	responseHeaders = header.Clone()
	responseHeaders.Del("X-Gateway-Authorization-Data")
	responseHeaders.Del("X-Gateway-Authorization-More")
	responseHeaders.Del("X-Gateway-Authorization-ExpiresIn")
	responseHeaders.Del("X-Gateway-Authorization-Action")
	responseHeaders.Del("X-Gateway-Authorization-Action-Singleton")
	responseHeaders.Del("X-Gateway-Authorization-Remove-Multi")
	responseHeaders.Del("X-Gateway-Authorization-Remove-Prefix")
	status = http.StatusOK

	moreData := map[string]interface{}{}
	json.Unmarshal([]byte(moreDataRaw), &moreData)

	if routeValue, has := ctx.Get("route"); has {

		route := routeValue.(types2.Route)

		prefix := "AUTH"
		if len(route.AuthorizePrefix) > 0 && len(route.AuthorizePrefix) < 4  {
			prefix = fmt.Sprintf("%0*s",4, route.AuthorizePrefix)
		} else if len(route.AuthorizePrefix) >= 4 {
			prefix = route.AuthorizePrefix[0:4]
		}

		switch action {
		case "CREATE":

			if route.CustomToken {

				if expiresIn == 0 {
					expiresIn = this.JwtExpiresIn
				}

				expireAt, tokenErr := this.TokenService.CustomGenerate(data, prefix + ":", moreData, expiresIn)

				if tokenErr != nil {
					panic(types.RuntimeError{
						Message:    tokenErr.Error(),
						ErrorCode:  http.StatusInternalServerError,
						StatusCode: http.StatusInternalServerError,
					})
				}

				response = map[string]interface{}{
					"token": data,
					"expiresIn": expiresIn,
					"expiresAt": expireAt,
				}
			} else {

				accessToken, expireAt, refreshToken, refreshExpireAt, tokenErr := this.TokenService.Generate(data, prefix + ":", singleton, authorizeType, moreData)

				if tokenErr != nil {
					panic(types.RuntimeError{
						Message:    tokenErr.Error(),
						ErrorCode:  http.StatusInternalServerError,
						StatusCode: http.StatusInternalServerError,
					})
				}

				jwtExpiresIn := int64(math.MaxInt32)

				if this.JwtExpiresIn != -1 {
					jwtExpiresIn = this.JwtExpiresIn
				}

				response = authorizationResponse{
					AccessToken: accessToken,
					ExpiresIn: jwtExpiresIn,
					ExpireAt: expireAt,
					RefreshToken: refreshToken,
					RefreshTokenExpireAt: refreshExpireAt,
				}
			}



		case "REFRESH":
			refreshRequest := refreshAuthorizationRequest{}
			err = json.Unmarshal(requestBody, &refreshRequest)
			now := time.Now().Unix()
			token, tokenErr := this.TokenService.Parse(
				refreshRequest.RefreshToken)

			if tokenErr != nil {
				panic(types.RuntimeError{
					Message:   tokenErr.Error(),
					ErrorCode: http.StatusNotAcceptable,
				})
			}

			if now > token.ExpiresAt {
				panic(types.RuntimeError{
					Message:   "refresh token已过期",
					ErrorCode: http.StatusUnauthorized,
				})
			}

			/**
			新建refresh-token
			*/
			accessToken, expireAt, refreshToken, refreshExpireAt, refreshErr := this.TokenService.Refresh(token.Issuer, token.Id, prefix + ":")

			if refreshErr != nil {
				panic(types.RuntimeError{
					Message:   tokenErr.Error(),
					ErrorCode: http.StatusInternalServerError,
				})
			}

			jwtExpiresIn := int64(math.MaxInt32)

			if this.JwtExpiresIn != -1 {
				jwtExpiresIn = this.JwtExpiresIn
			}

			response = authorizationResponse{
				AccessToken: accessToken,
				ExpiresIn: jwtExpiresIn,
				ExpireAt: expireAt,
				RefreshToken: refreshToken,
				RefreshTokenExpireAt: refreshExpireAt,
			}
		case "REMOVE":
			if multi != "" {
				if removePrefix != "" {
					prefix = removePrefix
				}
				for _, token := range strings.Split(multi, ",") {
					fmt.Println("[REMOVE TOKEN] " + prefix + ":" + token)
					this.TokenService.Delete(token, prefix + ":")
				}
			} else {
				this.TokenService.Delete(gatewayData.Data.Token, prefix + ":")
			}
		}
	}


	return
}
