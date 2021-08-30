package dispatcher

import (
	"github.com/gin-gonic/gin"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/rpc"
	"net/http"
	"net/url"
)

type Request struct {
	abstract.Rpc
}

func (this *Request) execute (url url.URL, query gin.Param, requestHeader http.Header, requestBody string) (status int, responseBody string, responseHeader http.Header, err error) {

	err = rpc.Call(rpc.IN{url, query, requestHeader, requestBody}, &rpc.OUT{&status, &responseBody, &responseHeader}, this)

	if err != nil {
		panic(err)
	}
	return
}