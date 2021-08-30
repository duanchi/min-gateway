package routes

import (
	"github.com/duanchi/min/abstract"
	util2 "github.com/duanchi/min/util"
	"github.com/duanchi/min-gateway/service/storage"
	"github.com/duanchi/min-gateway/types"
)



type ServicesArray []types.Service

type Services struct {
	abstract.Bean

	Maps types.ServicesMap

	ConfigFile string `value:"/services.json"`

	StorageKey string `value:"GATEWAY:SERVICES"`
	KEY string `value:"services"`
	StorageService *storage.StorageService `autowired:"true"`
}

func (this *Services) Init () {
	this.Maps = types.ServicesMap{}
	// this.StorageService.HGetAll(this.StorageKey, &this.Maps)
	data := this.StorageService.Get(this.KEY)
	if data == nil {
		this.Maps = types.ServicesMap{}
	} else {
		this.Maps = data.(types.ServicesMap)
	}
}

func (this *Services) GetAll () types.ServicesMap {
	this.Init()
	return this.Maps
}

func (this *Services) Add (service types.Service) {
	service.Id = util2.GenerateUUID().String()
	if this.Maps == nil {
		this.Maps = types.ServicesMap{}
	}
	this.Maps[service.Id] = service
	// this.StorageService.HSet(this.StorageKey, service.Id, service, -1)
	this.StorageService.Save(service.Id, service, this.KEY)
}

func (this *Services) Modify (id string, modifiedService types.Service) {
	this.Maps[id] = modifiedService
	// this.StorageService.HSet(this.StorageKey, id, modifiedService, -1)
	this.StorageService.Save(id, modifiedService, this.KEY)
}

func (this *Services) Remove (id string) {
	delete(this.Maps, id)
	// this.StorageService.HRemove(this.StorageKey, id)
	this.StorageService.Remove(id, this.KEY)
}