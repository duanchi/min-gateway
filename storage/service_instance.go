package storage

import (
	"github.com/duanchi/min"
	cache2 "github.com/duanchi/min-gateway/cache"
	"github.com/duanchi/min-gateway/mapper"
	"github.com/duanchi/min/abstract"
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

func (this *ServiceInstanceStorage) GetByServiceId(id string) (instances []mapper.ServiceInstance) {
	allInstances := []mapper.ServiceInstance{}
	this.CacheService.GetList(this.CACHE_PREFIX, &allInstances)

	for _, instance := range allInstances {
		if instance.ServiceId == id {
			instances = append(instances, instance)
		}
	}

	return
}

func (this *ServiceInstanceStorage) GetAllGroupByServiceId() (instances map[string][]mapper.ServiceInstance) {
	instances = map[string][]mapper.ServiceInstance{}
	allInstances := []mapper.ServiceInstance{}
	this.CacheService.GetList(this.CACHE_PREFIX, &allInstances)

	for _, instance := range allInstances {
		if _, has := instances[instance.ServiceId]; !has {
			instances[instance.ServiceId] = []mapper.ServiceInstance{}
		}
		instances[instance.ServiceId] = append(instances[instance.ServiceId], instance)
	}

	return
}

func (this *ServiceInstanceStorage) Update(instance mapper.ServiceInstance) {
	min.Db.ID(instance.Id).Update(instance)
	this.DataToCache()
}

func (this *ServiceInstanceStorage) Add(instance mapper.ServiceInstance) {
	min.Db.Insert(instance)
	this.DataToCache()
}

func (this *ServiceInstanceStorage) AddList(instances []mapper.ServiceInstance) {
	min.Db.Insert(&instances)
	this.DataToCache()
}

func (this *ServiceInstanceStorage) Remove(instanceId string) {
	var instance mapper.ServiceInstance
	min.Db.Where("instance_id = ?", instanceId).Delete(instance)
	this.DataToCache()
}

func (this *ServiceInstanceStorage) RemoveList(instanceIds []string) {
	var instance mapper.ServiceInstance
	min.Db.In("instance_id", instanceIds).Delete(instance)
	this.DataToCache()
}

func (this *ServiceInstanceStorage) RemoveByServiceId(serviceId string) {
	var instance mapper.ServiceInstance
	min.Db.Where("service_id = ?", serviceId).Delete(&instance)
	this.DataToCache()
}

func (this *ServiceInstanceStorage) DataToCache() {
	var instances []mapper.ServiceInstance
	min.Db.Find(&instances)

	this.CacheService.DelPrefix(this.CACHE_PREFIX)

	for _, instance := range instances {
		this.CacheService.Set(this.CACHE_PREFIX, instance.InstanceId, instance)
	}
}
