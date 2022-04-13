package storage

import (
	"github.com/duanchi/min"
	cache2 "github.com/duanchi/min-gateway/cache"
	"github.com/duanchi/min-gateway/mapper"
	"github.com/duanchi/min/abstract"
	"strings"
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
	this.CacheService.Set(this.CACHE_PREFIX, instance.InstanceId, instance)
}

func (this *ServiceInstanceStorage) Add(instance mapper.ServiceInstance) {
	min.Db.Insert(instance)
	this.CacheService.Set(this.CACHE_PREFIX, instance.InstanceId, instance)
}

func (this *ServiceInstanceStorage) AddList(instances []mapper.ServiceInstance) {
	_, err := min.Db.Insert(&instances)
	if err == nil {
		for _, instance := range instances {
			this.CacheService.Set(this.CACHE_PREFIX, instance.InstanceId, instance)
		}
	}
}

func (this *ServiceInstanceStorage) Remove(instanceId string) {
	var instance mapper.ServiceInstance
	min.Db.Where("instance_id = ?", instanceId).Delete(instance)
	this.CacheService.Del(this.CACHE_PREFIX, instance.InstanceId)
}

func (this *ServiceInstanceStorage) RemoveList(instanceIds []string) {
	var instance mapper.ServiceInstance
	min.Db.Where("instance_id in (?)", strings.Join(instanceIds, ",")).Delete(instance)

	for _, id := range instanceIds {
		this.CacheService.Del(this.CACHE_PREFIX, id)
	}
}

func (this *ServiceInstanceStorage) RemoveByServiceId(serviceId int64) {
	var instances []mapper.ServiceInstance
	min.Db.Where("service_id = ?", serviceId).Find(&instances)

	for _, instance := range instances {
		this.CacheService.Del(this.CACHE_PREFIX, instance.InstanceId)
	}
}

func (this *ServiceInstanceStorage) DataToCache() {
	var instances []mapper.ServiceInstance
	min.Db.Find(&instances)

	for _, instance := range instances {
		this.CacheService.Set(this.CACHE_PREFIX, instance.InstanceId, instance)
	}
}
