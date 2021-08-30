package routes

import (
	"github.com/duanchi/min/abstract"
	util2 "github.com/duanchi/min/util"
	"github.com/duanchi/min-gateway/service/storage"
	"github.com/duanchi/min-gateway/types"
	"sort"
)

type Routes struct {
	abstract.Bean

	Maps types.RoutesArray
	Raw types.RoutesMap

	StorageKey string `value:"GATEWAY:ROUTES"`
	KEY string `value:"routes"`
	StorageService *storage.StorageService `autowired:"true"`
}

func (this *Routes) Init () {
	this.Raw = types.RoutesMap{}
	// this.StorageService.HGetAll(this.StorageKey, &this.Raw)
	data := this.StorageService.Get(this.KEY)

	if data == nil {
		this.Raw = types.RoutesMap{}
	} else {
		this.Raw = data.(types.RoutesMap)
	}
	this.Refresh()
}

func (this *Routes) GetAll () types.RoutesArray {
	this.Init()
	return this.Maps
}

func (this *Routes) Add (route types.Route) {

	route.Id = util2.GenerateUUID().String()
	if this.Maps == nil {
		this.Maps = types.RoutesArray{}
	}
	if this.Raw == nil {
		this.Raw = types.RoutesMap{}
	}
	route.Order = len(this.Maps)
	this.Raw[route.Id] = route
	// this.StorageService.HSet(this.StorageKey, route.Id, route, -1)
	this.StorageService.Save(route.Id, route, this.KEY)
	this.Refresh()
}

func (this *Routes) Modify (id string, modifiedRoute types.Route) {
	this.Raw[id] = modifiedRoute
	// this.StorageService.HSet(this.StorageKey, id, modifiedRoute, -1)
	this.StorageService.Save(id, modifiedRoute, this.KEY)
	this.Refresh()
}

func (this *Routes) Sort (orders []string) {
	for index, id := range orders {
		for key, route := range this.Raw {
			if id == key {
				route.Order = index
				this.Raw[key] = route
				break
			}
		}
	}
	this.StorageService.Update()
	this.Refresh()
}

func (this *Routes) Remove (id string) {
	delete(this.Raw, id)
	// this.StorageService.HRemove(this.StorageKey, id)
	this.StorageService.Remove(id, this.KEY)
	this.Refresh()
}

func (this *Routes) Refresh () {
	newMap := types.RoutesArray{}

	for _, route := range this.Raw {
		newMap = append(newMap, route)
	}

	sort.Sort(newMap)

	this.Maps = newMap
}