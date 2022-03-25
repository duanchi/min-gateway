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

func (this *RouteRewriteStorage) GetByRouteId(id int64) (rewrites []mapper.RouteRewrite) {
	allRewrites := []mapper.RouteRewrite{}
	this.CacheService.GetList(this.CACHE_PREFIX, &allRewrites)

	for _, rewrite := range allRewrites {
		if rewrite.RouteId == id {
			rewrites = append(rewrites, rewrite)
		}
	}

	return
}

func (this *RouteRewriteStorage) DataToCache() {
	var rewrites []mapper.RouteRewrite
	min.Db.Find(&rewrites)

	for _, rewrite := range rewrites {
		this.CacheService.Set(this.CACHE_PREFIX, strconv.FormatInt(rewrite.Id, 10), rewrite)
	}
}
