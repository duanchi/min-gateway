package middleware

import (
	"fmt"
	"github.com/duanchi/min-gateway/mapper"
	"github.com/duanchi/min-gateway/storage"
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

	EnableGrayInstance bool `value:"${Gateway.EnableGrayInstance}"`

	Route           *storage.RouteStorage           `autowired:"true"`
	RouteRewrite    *storage.RouteRewriteStorage    `autowired:"true"`
	Service         *storage.ServiceStorage         `autowired:"true"`
	ServiceInstance *storage.ServiceInstanceStorage `autowired:"true"`
	NativeApi       *NativeApiMiddleware            `autowired:"true"`
	ConsoleApi      *ConsoleApiMiddleware           `autowired:"true"`
}

func (this *RouterMiddleware) AfterRoute(ctx *gin.Context) {

	routes := this.Route.GetAllRoutes()
	url := ctx.Request.URL
	method := ctx.Request.Method
	requestUrl := url.Path
	requestId := ctx.Request.Header.Get("Request-Id")

	if requestId == "" {
		requestId = util2.GenerateUUID().String()
	}

	nativeApiPrefix := config.Get("Gateway.NativeApi.Prefix").(string)
	consoleApiPrefix := config.Get("Gateway.ConsoleApi.Prefix").(string)

	ctx.Set("REQUEST_ID", requestId)

	if nativeApiPrefix != "" && strings.HasPrefix(requestUrl, nativeApiPrefix) {
		ctx.Set("NATIVE_API_RESOURCE", requestUrl[len(nativeApiPrefix)+1:])
		this.NativeApi.Execute(ctx)
		ctx.Abort()
	} else if consoleApiPrefix != "" && strings.HasPrefix(requestUrl, consoleApiPrefix) {
		ctx.Set("CONSOLE_API_RESOURCE", requestUrl[len(consoleApiPrefix)+1:])
		this.ConsoleApi.Execute(ctx)
		ctx.Abort()
	} else {
		if url.RawQuery != "" {
			requestUrl += "?" + url.RawQuery
		}

		if url.Fragment != "" {
			requestUrl += "#" + url.Fragment
		}

		if len(routes) > 0 {
			for _, stack := range routes {

				methodMatch := false
				urlMatch := false

				switch mapper.CONSTANT.URL_TYPE[stack.UrlType] {
				case "regex":
					regex, _ := regexp.Compile(stack.Pattern)
					urlMatch = regex.MatchString(requestUrl)

				case "fnmatch":
					urlMatch = util.Fnmatch(stack.Pattern, requestUrl, 0)

				case "path":
					urlMatch = strings.HasPrefix(requestUrl, stack.Pattern)
				}

				_, matchAll := arrays.ContainsString(strings.Split(stack.Methods, ","), "ALL")
				_, match := arrays.ContainsString(strings.Split(stack.Methods, ","), method)
				if matchAll || match {
					methodMatch = true
				} else {
					if _, has := arrays.ContainsString(strings.Split(stack.Methods, ","), "WEBSOCKET"); has && method == "GET" {
						upgradeRequest := ctx.Request.Header.Get("Connection")
						upgradeProtocol := ctx.Request.Header.Get("Upgrade")

						if strings.ToLower(upgradeRequest) == "upgrade" && strings.ToLower(upgradeProtocol) == "websocket" {
							methodMatch = true
						}
					}
				}

				if methodMatch && urlMatch {

					// service, _ := this.Service.Get(stack.ServiceId)
					instances := this.ServiceInstance.GetByServiceId(stack.ServiceId)
					rewrites := this.RouteRewrite.GetByRouteId(stack.Id)

					if len(rewrites) > 0 {
						for _, rewrite := range rewrites {
							replacePattern, _ := regexp.Compile(rewrite.Pattern)
							requestUrl = replacePattern.ReplaceAllString(requestUrl, rewrite.Rewrite)
						}
					}

					ctxUrl := ""
					ctxRoute := stack

					instanceId := ctx.GetHeader("X-Instance-Id")

					if instanceId == "" {
						instanceId = ctx.GetHeader("Client-Id")
					}

					if instanceId != "" {
						for _, instance := range instances {
							if instance.InstanceId == instanceId {
								ctxUrl = instance.Url + requestUrl
								ctxRoute = stack
								ctx.Set("GRAY_INSTANCE", instance.Id)
								fmt.Println("[" + requestId + "] Force switch to gray service " + instance.InstanceId + " at " + instance.Url + " !!!")
								break
							}
						}
					} else {
						liveInstances := []mapper.ServiceInstance{}
						for _, instance := range instances {
							if instance.GrayFlag != 1 && instance.OnlineFlag == 1 {
								liveInstances = append(liveInstances, instance)
							}
						}

						if len(liveInstances) > 0 {
							n := 0
							total := len(instances)

							if total > 1 {
								n = rand.Intn(total)
							}

							ctxUrl = liveInstances[n].Url + requestUrl
							ctxRoute = stack
						} else {
							ctx.AbortWithStatusJSON(404,
								types.Response{
									RequestId: util.CtxGet("REQUEST_ID", ctx).(string),
									Code:      100404,
									Message:   "No instance provided on service",
									Data:      nil,
								})
						}

					}

					ctx.Set("URL", ctxUrl)
					ctx.Set("ROUTE", ctxRoute)

					ctx.Next()

					return
				}
			}
		}

		ctx.AbortWithStatusJSON(404,
			types.Response{
				RequestId: requestId,
				Code:      100404,
				Message:   "No service provided at request \"" + requestUrl + "\"",
				Data:      nil,
			})
	}
}
