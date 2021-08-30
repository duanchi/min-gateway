package dispatcher

import (
	"github.com/duanchi/min-gateway/routes"
	"github.com/duanchi/min-gateway/rpc/service"
	"github.com/gin-gonic/gin"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/util/arrays"
	"regexp"
	"strings"
)

type RpcDispatcher struct {
	abstract.Router

	Routes *routes.Routes `autowired:"true"`
	Services *routes.Services `autowired:"true"`
	ExtraService *service.ExtraService `autowired:"true"`
}

func (this *RpcDispatcher) Handle (path string, method string, params gin.Params, ctx *gin.Context) {
	url, _ := params.Get("url")
	method = strings.ToUpper(method)

	// res, rpcErr := this.ExtraService.Test("DEF", "STRING")

	if len(this.Routes.Maps) > 0 {
		for _, stack := range this.Routes.Maps {
			switch stack.Url.Type {
			case "regex":
				regex, _ := regexp.Compile(stack.Url.Match)
				methodMatch := false

				if _, has := arrays.ContainsString(stack.Method, "ALL"); has {
					methodMatch = true
				} else if _, has := arrays.ContainsString(stack.Method, method); has {
					methodMatch = true
				}

				if methodMatch && regex.MatchString(url) {

					for _, service := range this.Services.Maps {
						if service.Name == stack.Service {
							/*serviceUrl := service.Instances[rand.Intn(len(service.Instances) - 1)]
							serviceName := "rpc"
							serviceClass := "RpcService"
							serviceMethod := "Execute"*/

						}
					}
				}
			}
		}
	}
}