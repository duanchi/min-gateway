package console_api

import (
	"github.com/duanchi/min-gateway/service"
	types2 "github.com/duanchi/min-gateway/types/request"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/types"
	"github.com/gin-gonic/gin"
)

type IntegrationController struct {
	abstract.RestController

	RouteService   *service.Route   `bean:"autowired"`
	ServiceService *service.Service `autowired:"true"`
}

func (this *IntegrationController) Fetch(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {
	return this.RouteService.GetAll(), nil
}

func (this *IntegrationController) Create(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {
	var route types2.RouteRequest
	ctx.ShouldBindJSON(&route)
	this.RouteService.Add(route)
	return true, nil
}

func (this *IntegrationController) Remove(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {
	id = ctx.Query("id")
	this.RouteService.Remove(id)
	return true, nil
}

func (this *IntegrationController) Update(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {
	scope := ctx.Query("scope")
	action := ctx.Query("action")

	if action == "refresh" {
		this.Init()
		this.ServiceService.Init()
		return true, nil
	}

	if scope == "order" {

		var order []string
		bindErr := ctx.ShouldBindJSON(&order)
		if bindErr != nil {
			err = types.RuntimeError{
				Message: bindErr.Error(),
			}
		}
		this.RouteService.Sort(order)
	} else {
		var route types2.RouteRequest
		id = ctx.Query("id")
		bindErr := ctx.ShouldBindJSON(&route)
		if bindErr != nil {
			err = types.RuntimeError{
				Message: bindErr.Error(),
			}
		}
		this.RouteService.Modify(id, route)
	}
	return true, nil
}
