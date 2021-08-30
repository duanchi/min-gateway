package dispatcher

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/duanchi/min/types"
	"github.com/duanchi/min/types/gateway"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebsocketProxy(url string, ctx *gin.Context, gatewayData gateway.Data) (err types.Error) {

	if strings.Contains(url, "?") {
		url += "&"
	} else {
		url += "?"
	}
	url += "__X-GATEWAY-DATA__="
	if data, err := json.Marshal(gatewayData); err == nil {
		url += base64.URLEncoding.EncodeToString(data)
	} else {
		url += base64.URLEncoding.EncodeToString([]byte("{}"))
	}

	if strings.Index(url, "https") == 0 {
		url = strings.Replace(url, "https", "wss", 1)
	} else if strings.Index(url, "http") == 0 {
		url = strings.Replace(url, "http", "ws", 1)
	}

	proxyConnection, wsErr := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if wsErr != nil {
		return types.RuntimeError{
			Message:   wsErr.Error(),
			ErrorCode: 500,
			ErrorData: nil,
		}
	}
	defer proxyConnection.Close()

	websocketClient, _, clientErr := websocket.DefaultDialer.Dial(url, nil)

	if clientErr != nil {
		proxyConnection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseGoingAway, "Websocket Proxy Error, " + clientErr.Error()))
		return nil
	}
	defer websocketClient.Close()

	go func() {
		for {
			//读取ws中的数据
			messageType, message, err := proxyConnection.ReadMessage()
			if err != nil {
				websocketClient.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
				break
			}
			websocketClient.WriteMessage(messageType, message)
			// proxyConnection.WriteMessage(messageType, message)
			//写入ws数据
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	done := make(chan struct{})

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			err := websocketClient.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return types.RuntimeError{
					Message:   err.Error(),
					ErrorCode: 500,
				}
			}
			return nil

		default:
			messageType, message, err := websocketClient.ReadMessage()
			if err != nil {
				proxyConnection.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
				websocketClient.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
				return nil
			}
			proxyConnection.WriteMessage(messageType, message)
		}
	}
	// websocketClient, _, clientErr := websocket.DefaultDialer.Dial(url, nil)

	/*if clientErr != nil {
		return types.RuntimeError{
			Message:   "Websocket Proxy Error, " + clientErr.Error(),
			ErrorCode: http.StatusInternalServerError,
		}
	}
	defer websocketClient.Close()*/

	/*proxyConnection, websocketError := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if websocketError != nil {
		return types.RuntimeError{
			Message:   websocketError.Error(),
			ErrorCode: http.StatusInternalServerError,
		}
	}
	defer proxyConnection.Close()

	for {
		messageType, message, err := proxyConnection.ReadMessage()
		if err != nil {
			// websocketClient.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
			fmt.Println(err)
			return types.RuntimeError{
				Message:   err.Error(),
				ErrorCode: 500,
				ErrorData: nil,
			}
		}
		// websocketClient.WriteMessage(messageType, message)
		proxyConnection.WriteMessage(messageType, message)
	}*/

	/*for {
		err := handleFunction(connection, id, resource, parameters, ctx)
		if err != nil {
			break
		}
	}*/

	/*interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	done := make(chan struct{})*/

	/*go func() {
		defer close(done)
		for {
			messageType, message, err := proxyConnection.ReadMessage()
			if err != nil {
				websocketClient.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
				return
			}
			websocketClient.WriteMessage(messageType, message)
		}
	}()*/

	/*for {
		select {
		case <-done:
			return
		case <-interrupt:
			err := websocketClient.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return types.RuntimeError{
					Message:   err.Error(),
					ErrorCode: 500,
				}
			}
			return nil

		default:
			messageType, message, err := websocketClient.ReadMessage()
			if err != nil {
				proxyConnection.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
				return types.RuntimeError{
					Message:   err.Error(),
					ErrorCode: 500,
				}
			}
			proxyConnection.WriteMessage(messageType, message)
		}
	}*/
}

func ping(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	for {
		//读取ws中的数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if string(message) == "ping" {
			message = []byte("pong")
		}
		//写入ws数据
		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}