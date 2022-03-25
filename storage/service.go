package storage

import (
	"github.com/duanchi/min"
	cache2 "github.com/duanchi/min-gateway/cache"
	"github.com/duanchi/min-gateway/mapper"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/cache"
	"strconv"
)

type ServiceStorage struct {
	abstract.Service

	CacheService    *cache2.CacheService `autowired:"true"`
	CACHE_PREFIX    string               `value:"HASH:SERVICE"`
	CACHE_ID_PREFIX string               `value:"HASH:SERVICE_ID"`
}

func (this *ServiceStorage) GetByCode(code string) (provider mapper.Service, ok bool) {
	ok = this.CacheService.Get(this.CACHE_PREFIX, code, &provider)
	return
}

func (this *ServiceStorage) Get(id int64) (provider mapper.Service, ok bool) {
	ok = this.CacheService.Get(this.CACHE_ID_PREFIX, strconv.FormatInt(id, 10), &provider)
	return
}

func (this *ServiceStorage) DataToCache() {
	var services []mapper.Service
	min.Db.Find(&services)

	for _, service := range services {
		this.CacheService.Set(this.CACHE_PREFIX, service.Code, service)
		this.CacheService.Set(this.CACHE_ID_PREFIX, strconv.FormatInt(service.Id, 10), service)
	}
}

func (this *ServiceStorage) cacheGet(key string) interface{} {
	if cache.Has(this.CACHE_PREFIX + key) {
		return cache.Get(this.CACHE_PREFIX + key)
	}
	return interface{}(nil)
}

func (this *ServiceStorage) cacheSet(key string, value interface{}) {
	cache.Set(this.CACHE_PREFIX+key, value)
}

func (this *ServiceStorage) cacheDel(key string) {
	cache.Del(this.CACHE_PREFIX + key)
}
