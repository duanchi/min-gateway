package native_api

import (
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/types"
	"github.com/gin-gonic/gin"
	"github.com/duanchi/min-gateway/routes"
	types2 "github.com/duanchi/min-gateway/types"
)

type ServicesController struct {
	abstract.RestController

	Services *routes.Services `autowired:"true"`
}

func (this *ServicesController) Fetch (id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {

	return this.Services.GetAll(), nil
}

func (this *ServicesController) Create (id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {

	var service types2.Service
	ctx.BindJSON(&service)
	this.Services.Add(service)
	return true, nil
}

func (this *ServicesController) Remove (id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {
	id = ctx.Query("id")
	this.Services.Remove(id)
	return true, nil
}

func (this *ServicesController) Update (id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {
	var service types2.Service
	ctx.BindJSON(&service)
	id = ctx.Query("id")
	this.Services.Modify(id, service)
	return true, nil
}