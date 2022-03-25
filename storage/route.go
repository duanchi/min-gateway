package storage

import (
	"github.com/duanchi/min"
	cache2 "github.com/duanchi/min-gateway/cache"
	"github.com/duanchi/min-gateway/mapper"
	"github.com/duanchi/min/abstract"
	"strconv"
)

type RouteStorage struct {
	abstract.Service

	CacheService *cache2.CacheService `autowired:"true"`
	CACHE_PREFIX string               `value:"HASH:ROUTE"`
}

func (this *RouteStorage) GetRoutesByProviderId(id int64) (routes []mapper.Route) {
	allRoutes := []mapper.Route{}
	this.CacheService.GetList(this.CACHE_PREFIX, &allRoutes)

	for _, route := range allRoutes {
		if route.ServiceId == id {
			routes = append(routes, route)
		}
	}

	return
}

func (this *RouteStorage) GetAllRoutes() (routes []mapper.Route) {
	this.CacheService.GetList(this.CACHE_PREFIX, &routes)
	return
}

func (this *RouteStorage) DataToCache() {
	var routes []mapper.Route
	min.Db.Find(&routes)

	for _, route := range routes {
		this.CacheService.Set(this.CACHE_PREFIX, strconv.FormatInt(route.Id, 10), route)
	}
}
