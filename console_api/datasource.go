package console_api

import (
	"github.com/duanchi/min-gateway/service/storage"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DatasourceController struct {
	abstract.RestController

	StorageService *storage.StorageService `bean:"autowired"`
}

func (this *DatasourceController) Fetch(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {

	data, exportError := this.StorageService.Export()

	if exportError == nil {
		fileContentDisposition := "attachment;filename=\"min-gateway.json\""
		ctx.Header("Content-Type", "application/json") // 这里是压缩文件类型 .zip
		ctx.Header("Content-Disposition", fileContentDisposition)
		ctx.Data(http.StatusOK, "application/json", data)
	}

	return "Ok", nil
}

func (this *DatasourceController) Create(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {
	return "error", nil
}
