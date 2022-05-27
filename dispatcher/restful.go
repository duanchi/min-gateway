package dispatcher

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/duanchi/min-gateway/mapper"
	"github.com/duanchi/min-gateway/service"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/types"
	"github.com/duanchi/min/types/gateway"
	"github.com/duanchi/min/util/arrays"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type RestfulDispatcher struct {
	abstract.Router

	AuthorizationService *service.AuthorizationService `autowired:"true"`

	DefaultSingleton     bool  `value:"${Authorization.DefaultSingleton}"`
	GlobalRequestTimeout int64 `value:"${Gateway.GlobalRequestTimeout}"`
}

func (this *RestfulDispatcher) Handle(path string, method string, params gin.Params, ctx *gin.Context) {
	if urlValue, has := ctx.Get("URL"); has {
		url := urlValue.(string)
		rawRequestId, _ := ctx.Get("REQUEST_ID")
		requestId := rawRequestId.(string)
		rawRoute, _ := ctx.Get("ROUTE")
		route := rawRoute.(mapper.Route)
		rawInstanceId, _ := ctx.Get("GRAY_INSTANCE")
		instanceId := rawInstanceId.(string)

		var gatewayData = gateway.Data{}

		/*if url, has := ctx.Get("url"); has {
			gatewayData.Url = url.(string)
		}*/

		if token, has := ctx.Get("TOKEN"); has {
			gatewayData.Data.Token = token.(string)
		}

		if user, has := ctx.Get("USER"); has {
			gatewayData.Data.User = user.(string)
		}

		if more, has := ctx.Get("MORE"); has {
			gatewayData.Data.More = more.(map[string]interface{})
		}

		if method == "GET" {
			upgradeRequest := ctx.Request.Header.Get("Connection")
			upgradeProtocol := ctx.Request.Header.Get("Upgrade")

			if upgradeRequest == "Upgrade" && strings.ToLower(upgradeProtocol) == "websocket" {
				//websocket proxy
				err := WebsocketProxy(url, ctx, gatewayData)
				if err != nil {
					panic(types.RuntimeError{
						Message:   err.Error(),
						ErrorCode: 10500,
					})
				}
				return
			}
		}

		requestHeader := ctx.Request.Header.Clone()
		additionalHeader := map[string]string{}
		responseStatus := http.StatusInternalServerError

		if route.NeedAuthorize == 1 {
			if data, err := json.Marshal(gatewayData); err == nil {
				additionalHeader["X-Gateway-Data"] = base64.URLEncoding.EncodeToString(data)
			} else {
				additionalHeader["X-Gateway-Data"] = base64.URLEncoding.EncodeToString([]byte("{}"))
			}
		}

		if ctx.Request.RemoteAddr != "" {
			additionalHeader["X-Forward-For"] = ctx.Request.RemoteAddr
		}

		additionalHeader["X-Gateway-Request-Id"] = requestId

		parseRequestHeader(&requestHeader, additionalHeader)
		body, _ := ioutil.ReadAll(ctx.Request.Body)

		responseStatus, responseHeaders, responseBody, err := restfulRequest(
			requestId,
			url,
			ctx.Request.Method,
			requestHeader,
			body,
			0,
		)

		contentType := ""

		if err != nil {

			responseBody, _ = json.Marshal(types.Response{
				RequestId: requestId,
				Status:    false,
				Message:   err.Error(),
				Data:      nil,
			})
			responseStatus = http.StatusInternalServerError

		} else {
			contentType = responseHeaders.Get("Content-Type")

			if _, has := arrays.ContainsString([]string{"CREATE", "REFRESH", "REMOVE"}, strings.ToUpper(responseHeaders.Get("X-Gateway-Authorization-Action"))); route.NeedAuthorize == 1 && has {
				//进入授权流程
				singleton := this.DefaultSingleton
				if !this.DefaultSingleton && responseHeaders.Get("X-Gateway-Authorization-Singleton") == "true" {
					singleton = true
				}

				if this.DefaultSingleton && responseHeaders.Get("X-Gateway-Authorization-Singleton") == "false" {
					singleton = false
				}

				authorizeType := ""
				if route.AuthorizeTypeKey != "" {
					stack := strings.SplitN(route.AuthorizeTypeKey, ":", 2)

					if len(stack) == 2 && stack[0] == "HEADER" {
						authorizeType = ctx.GetHeader(stack[1])
					} else if len(stack) == 2 {
						authorizeType, _ = ctx.GetQuery(stack[1])
					} else {
						authorizeType, _ = ctx.GetQuery("platform")
					}
				}
				if authorizeType == "" {
					authorizeType = "default"
				}

				_, _, response, err := this.AuthorizationService.Handle(&responseHeaders, gatewayData, body, singleton, authorizeType, ctx)

				if err != nil {
					panic(err)
				}
				responseBody, _ = json.Marshal(types.Response{
					RequestId: requestId,
					Status:    true,
					Code:      0,
					Message:   "Ok",
					Data:      response,
				})
				responseStatus = http.StatusOK
			} else {
				for headerKey, header := range responseHeaders {
					for _, headerValue := range header {
						ctx.Writer.Header().Add(headerKey, headerValue)
					}
				}
			}
		}

		if contentType == "" {
			contentType = "text/plain"
		}

		ctx.Writer.Header().Set("X-Request-Id", requestId)
		if instanceId != "" {
			ctx.Writer.Header().Set("X-Instance-Id", instanceId)
		}

		ctx.Data(responseStatus, contentType, responseBody)
	}
}

func restfulRequest(
	requestId string,
	requestUrl string,
	method string,
	requestHeaders http.Header,
	requestBody []byte,
	timeout int64,
) (
	status int,
	responseHeaders http.Header,
	responseBody []byte,
	err error,
) {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	request, err := http.NewRequest(method, requestUrl, bytes.NewReader(requestBody))
	if err != nil {
		// handle error
		responseBody, _ = json.Marshal(types.Response{
			RequestId: requestId,
			Status:    false,
			Code:      100450,
			Message:   "Gateway request error",
			Data:      nil,
		})
		status = http.StatusInternalServerError
		responseHeaders = http.Header{
			"Content-Type": []string{"application/json"},
		}
		return
	}

	request.Header.Del("Authorization")
	request.Header.Del("Host")
	request.Header.Del("X-Gateway-Request-Id")
	request.Header.Del("X-Gateway-Data")

	request.Header = requestHeaders

	request.Header.Set("X-Gateway-Request-Id", requestId)

	response, requestErr := client.Do(request)

	if requestErr != nil {
		responseBody, _ = json.Marshal(types.Response{
			RequestId: requestId,
			Status:    false,
			Code:      100550,
			Message:   "Gateway response error, " + requestErr.(*url.Error).Error(),
			Data:      nil,
		})
		status = http.StatusInternalServerError
		responseHeaders = http.Header{
			"Content-Type": []string{"application/json"},
		}
		return
	} else {
		if reflect.ValueOf(response).IsValid() {
			responseBody, err = ioutil.ReadAll(response.Body)
			responseHeaders = response.Header
			status = response.StatusCode

			return
		}

		responseBody, _ = json.Marshal(types.Response{
			RequestId: requestId,
			Status:    false,
			Code:      100551,
			Message:   "Gateway response error, nil response",
			Data:      nil,
		})
		status = http.StatusInternalServerError
		responseHeaders = http.Header{
			"Content-Type": []string{"application/json"},
		}
	}

	defer response.Body.Close()

	return
}

func parseRequestHeader(header *http.Header, headerMap map[string]string) {
	header.Del("Authorization")
	header.Del("Host")

	for key, value := range headerMap {
		header.Del(key)
		header.Set(key, value)
	}
}
