package console_api

import (
	"fmt"
	"github.com/duanchi/min-gateway/service"
	"github.com/duanchi/min-gateway/types/request"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/types"
	"github.com/gin-gonic/gin"
)

type ServicesController struct {
	abstract.RestController

	Service *service.Service `bean:"autowired"`
}

func (this *ServicesController) Fetch(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {

	return this.Service.GetAll(), nil
}

func (this *ServicesController) Create(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {

	var service request.Service
	bindErr := ctx.BindJSON(&service)
	fmt.Println(bindErr)
	this.Service.Add(service)
	return true, nil
}

func (this *ServicesController) Remove(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {
	id = ctx.Query("id")
	this.Service.Remove(id)
	return true, nil
}

func (this *ServicesController) Update(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {
	var service request.Service
	ctx.BindJSON(&service)
	id = ctx.Query("id")
	this.Service.Modify(id, service)
	return true, nil
}
