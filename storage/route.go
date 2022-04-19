package storage

import (
	"github.com/duanchi/min"
	cache2 "github.com/duanchi/min-gateway/cache"
	"github.com/duanchi/min-gateway/mapper"
	"github.com/duanchi/min-gateway/types"
	"github.com/duanchi/min/abstract"
	"sort"
	"strconv"
)

type RouteStorage struct {
	abstract.Service

	CacheService *cache2.CacheService `bean:"autowired"`
	CACHE_PREFIX string               `value:"HASH:ROUTE"`
}

func (this *RouteStorage) GetRoutesByServiceId(id string) (routes []mapper.Route) {
	allRoutes := []mapper.Route{}
	this.CacheService.GetList(this.CACHE_PREFIX, &allRoutes)

	for _, route := range allRoutes {
		if route.ServiceId == id {
			routes = append(routes, route)
		}
	}

	return
}

func (this *RouteStorage) GetFromDB(id string) (route mapper.Route, ok bool) {
	route.RouteId = id
	ok, _ = min.Db.Get(&route)
	return
}

func (this *RouteStorage) GetAll() (routes []mapper.Route) {
	this.CacheService.GetList(this.CACHE_PREFIX, &routes)
	return
}

func (this *RouteStorage) GetAllSorted() (routes types.RouteSlice) {
	this.CacheService.GetList(this.CACHE_PREFIX, &routes)
	sort.Sort(routes)
	return
}

func (this *RouteStorage) GetByRouteId(id string) (route mapper.Route) {
	this.CacheService.Get(this.CACHE_PREFIX, id, &route)
	return
}

func (this *RouteStorage) GetLastSort() (sort int64) {
	sort = 1
	route := mapper.Route{}

	if route.Id > 0 {
		sort += route.Id + 1
	}

	return
}

func (this *RouteStorage) Update(route mapper.Route) (ok bool) {
	_, err := min.Db.Where("route_id = ?", route.RouteId).Cols("pattern", "url_type", "methods", "service_id", "need_authorize", "authorize_prefix", "authorize_type_key", "custom_token", "description", "sort").Update(route)

	if err == nil {
		this.DataToCache()
		return true
	} else {
		return false
	}

}

func (this *RouteStorage) Add(route mapper.Route) (id string, ok bool) {
	_, err := min.Db.Insert(&route)
	if err != nil {
		return "", false
	} else {
		this.DataToCache()
		return route.RouteId, true
	}
}

func (this *RouteStorage) Remove(id string) {
	route := mapper.Route{}
	min.Db.Where("route_id = ?", route.RouteId).Delete(route)
	this.DataToCache()
}

func (this *RouteStorage) Sort(sorts []string) {
	for index, routeId := range sorts {
		min.Db.Where("route_id = ?", routeId).Table(new(mapper.Route)).Update(map[string]interface{}{"sort": index})
	}

	this.DataToCache()
}

func (this *RouteStorage) DataToCache() {
	var routes []mapper.Route
	min.Db.Find(&routes)

	this.CacheService.DelPrefix(this.CACHE_PREFIX)

	for _, route := range routes {
		this.CacheService.Set(this.CACHE_PREFIX, strconv.FormatInt(route.Id, 10), route)
	}
}
