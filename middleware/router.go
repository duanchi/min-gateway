package middleware

import (
	"fmt"
	"github.com/duanchi/min-gateway/routes"
	"github.com/duanchi/min-gateway/util"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/config"
	"github.com/duanchi/min/types"
	util2 "github.com/duanchi/min/util"
	"github.com/duanchi/min/util/arrays"
	"github.com/gin-gonic/gin"
	"math/rand"
	"regexp"
	"strings"
)

type RouterMiddleware struct {
	abstract.Middleware

	Routes *routes.Routes `autowired:"true"`
	Services *routes.Services `autowired:"true"`
	NativeApi *NativeApiMiddleware `autowired:"true"`
}

func (this *RouterMiddleware) AfterRoute (ctx *gin.Context) {

	url := ctx.Request.URL
	method := ctx.Request.Method
	requestUrl := url.Path
	requestId := ctx.Request.Header.Get("Request-Id")

	if requestId == "" {
		requestId = util2.GenerateUUID().String()
	}

	nativeApiPrefix := config.Get("Gateway.NativeApi.Prefix").(string)

	ctx.Set("REQUEST_ID", requestId)

	if nativeApiPrefix != "" && strings.HasPrefix(requestUrl, nativeApiPrefix) {
		ctx.Set("NATIVE_API_RESOURCE", requestUrl[len(nativeApiPrefix) + 1:])
		this.NativeApi.Execute(ctx)
		ctx.Abort()
	} else {
		if url.RawQuery != "" {
			requestUrl += "?" + url.RawQuery
		}

		if url.Fragment != "" {
			requestUrl += "#" + url.Fragment
		}

		fmt.Println("=====================================", this.Routes.Maps)

		if len(this.Routes.Maps) > 0 {
			for _, stack := range this.Routes.Maps {

				methodMatch := false
				urlMatch := false

				switch stack.Url.Type {
				case "regex":
					regex, _ := regexp.Compile(stack.Url.Match)
					urlMatch = regex.MatchString(requestUrl)

				case "fnmatch":
					urlMatch = util.Fnmatch(stack.Url.Match, requestUrl, 0)

				case "path":
					urlMatch = strings.HasPrefix(requestUrl, stack.Url.Match)
				}

				_, matchAll := arrays.ContainsString(stack.Method, "ALL")
				_, match := arrays.ContainsString(stack.Method, method)
				if matchAll || match {
					methodMatch = true
				} else {
					if _, has := arrays.ContainsString(stack.Method, "WEBSOCKET"); has && method == "GET" {
						upgradeRequest := ctx.Request.Header.Get("Connection")
						upgradeProtocol := ctx.Request.Header.Get("Upgrade")

						if strings.ToLower(upgradeRequest) == "upgrade"  && strings.ToLower(upgradeProtocol) == "websocket" {
							methodMatch = true
						}
					}
				}

				if methodMatch && urlMatch {

					for _, service := range this.Services.Maps {
						if service.Name == stack.Service {
							if len(stack.Rewrite) > 0 {
								for match, replace := range stack.Rewrite {
									replacePattern, _ := regexp.Compile(match)
									requestUrl = replacePattern.ReplaceAllString(requestUrl, replace)
								}
							}

							n := 0

							if len(service.Instances) > 1 {
								n = rand.Intn(len(service.Instances) - 1)
							}

							clientId := ctx.GetHeader("Client-Id")

							ctxUrl := service.Instances[n] + requestUrl
							ctxRoute := stack

							if clientId != "" && len(service.Gray) > 0 {
								for _, client := range service.Gray {
									if client.Id == clientId {
										ctxUrl = client.Uri + requestUrl
										ctxRoute = stack
										break
									}
								}
							}

							ctx.Set("url", ctxUrl)
							ctx.Set("route", ctxRoute)

							ctx.Next()

							return
						}
					}
				}
			}
		}

		ctx.AbortWithStatusJSON(404,
			types.Response{
				RequestId: requestId,
				Code: 100404,
				Message:   "No service provided at request \"" + requestUrl + "\"",
				Data:      nil,
			})
	}
}
