package storage

import (
	"github.com/duanchi/min"
	cache2 "github.com/duanchi/min-gateway/cache"
	"github.com/duanchi/min-gateway/mapper"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/cache"
)

type ServiceInstanceStorage struct {
	abstract.Service
	CacheService *cache2.CacheService `autowired:"true"`
	CACHE_PREFIX string               `value:"HASH:SERVICE_INSTANCE"`
}

func (this *ServiceInstanceStorage) Get(instanceId string) (instance mapper.ServiceInstance, ok bool) {
	ok = this.CacheService.Get(this.CACHE_PREFIX, instanceId, &instance)
	return
}

func (this *ServiceInstanceStorage) GetByServiceId(id int64) (instances []mapper.ServiceInstance) {
	allInstances := []mapper.ServiceInstance{}
	this.CacheService.GetList(this.CACHE_PREFIX, &allInstances)

	for _, instance := range allInstances {
		if instance.ServiceId == id {
			instances = append(instances, instance)
		}
	}

	return
}

func (this *ServiceInstanceStorage) Update(instance mapper.ServiceInstance) {
	min.Db.ID(instance.Id).Update(instance)
	this.CacheService.Set(this.CACHE_PREFIX, instance.InstanceId, instance)
}

func (this *ServiceInstanceStorage) Add(instance mapper.ServiceInstance) {
	min.Db.Insert(instance)
	this.CacheService.Set(this.CACHE_PREFIX, instance.InstanceId, instance)
}

func (this *ServiceInstanceStorage) Remove(instanceId string) {
	var instance mapper.ServiceInstance
	min.Db.Where("instance_id = ?", instanceId).Delete(instance)
	this.CacheService.Del(this.CACHE_PREFIX, instance.InstanceId)
}

func (this *ServiceInstanceStorage) DataToCache() {
	var instances []mapper.ServiceInstance
	min.Db.Find(&instances)

	for _, instance := range instances {
		this.CacheService.Set(this.CACHE_PREFIX, instance.InstanceId, instance)
	}
}

func (this *ServiceInstanceStorage) cacheGet(key string) interface{} {
	if cache.Has(this.CACHE_PREFIX + key) {
		return cache.Get(this.CACHE_PREFIX + key)
	}
	return interface{}(nil)
}

func (this *ServiceInstanceStorage) cacheSet(key string, value interface{}) {
	cache.Set(this.CACHE_PREFIX+key, value)
}

func (this *ServiceInstanceStorage) cacheDel(key string) {
	cache.Del(this.CACHE_PREFIX + key)
}
