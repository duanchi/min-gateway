package native_api

import (
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/types"
	"github.com/gin-gonic/gin"
)

type AuthorizeController struct {
	abstract.RestController
}

func (this *AuthorizeController) Fetch (id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {
	return "Ok", nil
}