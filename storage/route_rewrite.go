package storage

import (
	"github.com/duanchi/min"
	cache2 "github.com/duanchi/min-gateway/cache"
	"github.com/duanchi/min-gateway/mapper"
	"github.com/duanchi/min/abstract"
	"strconv"
)

type RouteRewriteStorage struct {
	abstract.Service

	CacheService *cache2.CacheService `autowired:"true"`
	CACHE_PREFIX string               `value:"HASH:ROUTE_REWRITE"`
}

func (this *RouteRewriteStorage) GetByRouteId(id string) (rewrites []mapper.RouteRewrite) {
	allRewrites := []mapper.RouteRewrite{}
	this.CacheService.GetList(this.CACHE_PREFIX, &allRewrites)

	for _, rewrite := range allRewrites {
		if rewrite.RouteId == id {
			rewrites = append(rewrites, rewrite)
		}
	}

	return
}

func (this *RouteRewriteStorage) Get(id int64) (rewrite mapper.RouteRewrite) {
	this.CacheService.Get(this.CACHE_PREFIX, strconv.FormatInt(id, 10), &rewrite)
	return
}

func (this *RouteRewriteStorage) GetAllGroupByRouteId() (rewrites map[string][]mapper.RouteRewrite) {
	allRewrites := []mapper.RouteRewrite{}
	this.CacheService.GetList(this.CACHE_PREFIX, &allRewrites)
	rewrites = map[string][]mapper.RouteRewrite{}

	for _, rewrite := range allRewrites {
		if _, has := rewrites[rewrite.RouteId]; !has {
			rewrites[rewrite.RouteId] = []mapper.RouteRewrite{}
		}
		rewrites[rewrite.RouteId] = append(rewrites[rewrite.RouteId], rewrite)
	}

	return
}

func (this *RouteRewriteStorage) Add(rewrite mapper.RouteRewrite) (id int64, ok bool) {
	_, err := min.Db.Insert(&rewrite)
	if err != nil {
		return 0, false
	} else {
		this.DataToCache()
		return rewrite.Id, true
	}
}

func (this *RouteRewriteStorage) AddList(rewrites []mapper.RouteRewrite) {
	_, err := min.Db.Insert(&rewrites)
	if err == nil {
		this.DataToCache()
	}
}

func (this *RouteRewriteStorage) RemoveByRouteId(id string) {
	rewrite := mapper.RouteRewrite{}
	min.Db.Where("route_id = ?", id).Delete(rewrite)
	this.DataToCache()
}

func (this *RouteRewriteStorage) Remove(id int64) {
	rewrite := this.Get(id)
	min.Db.Where("code = ?", id).Delete(rewrite)
	this.DataToCache()
}

func (this *RouteRewriteStorage) DataToCache() {
	var rewrites []mapper.RouteRewrite
	min.Db.Find(&rewrites)

	this.CacheService.DelPrefix(this.CACHE_PREFIX)

	for _, rewrite := range rewrites {
		this.CacheService.Set(this.CACHE_PREFIX, strconv.FormatInt(rewrite.Id, 10), rewrite)
	}
}
