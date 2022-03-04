package console_api

import (
	"github.com/duanchi/min-gateway/routes"
	types2 "github.com/duanchi/min-gateway/types"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/types"
	"github.com/gin-gonic/gin"
)

type RoutesController struct {
	abstract.RestController

	Routes   *routes.Routes   `autowired:"true"`
	Services *routes.Services `autowired:"true"`
}

func (this *RoutesController) Fetch(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {
	this.Routes.Refresh()
	return this.Routes.GetAll(), nil
}

func (this *RoutesController) Create(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {
	var route types2.Route
	ctx.ShouldBindJSON(&route)
	this.Routes.Add(route)
	return true, nil
}

func (this *RoutesController) Remove(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {
	id = ctx.Query("id")
	this.Routes.Remove(id)
	return true, nil
}

func (this *RoutesController) Update(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {
	scope := ctx.Query("scope")
	action := ctx.Query("action")

	if action == "refresh" {
		this.Init()
		this.Services.Init()
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
		this.Routes.Sort(order)
	} else {
		var route types2.Route
		id = ctx.Query("id")
		bindErr := ctx.ShouldBindJSON(&route)
		if bindErr != nil {
			err = types.RuntimeError{
				Message: bindErr.Error(),
			}
		}
		this.Routes.Modify(id, route)
	}
	return true, nil
}
