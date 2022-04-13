package storage

import (
	"fmt"
	"github.com/duanchi/min"
	"github.com/duanchi/min-gateway/cache"
	"github.com/duanchi/min-gateway/mapper"
	"github.com/duanchi/min/abstract"
)

type ServiceStorage struct {
	abstract.Service

	CacheService *cache.CacheService `autowired:"true"`
	CACHE_PREFIX string              `value:"HASH:SERVICE"`
}

func (this *ServiceStorage) GetByCode(code string) (service mapper.Service, ok bool) {
	ok = this.CacheService.Get(this.CACHE_PREFIX, code, &service)
	return
}

func (this *ServiceStorage) GetFromDB(serviceId string) (service mapper.Service, ok bool) {
	service.Code = serviceId
	ok, _ = min.Db.Get(&service)
	return
}

func (this *ServiceStorage) GetAll() (services []mapper.Service, ok bool) {
	ok = this.CacheService.GetList(this.CACHE_PREFIX, &services)
	return
}

func (this *ServiceStorage) Update(service mapper.Service) (ok bool) {
	fmt.Println("asdfasdfasdfasfd")
	_, err := min.Db.ID(service.Id).Cols("code", "name", "load_balance_type").Update(service)
	if err == nil {
		this.DataToCache()
		return true
	} else {
		fmt.Println(err)
		return false
	}

}

func (this *ServiceStorage) Add(service mapper.Service) (ok bool) {
	_, err := min.Db.Insert(service)
	if err != nil {
		return false
	} else {
		this.DataToCache()
		return true
	}
}

func (this *ServiceStorage) Remove(code string) {
	service, _ := this.GetByCode(code)
	min.Db.Where("code = ?", code).Delete(service)
	this.DataToCache()
}

func (this *ServiceStorage) DataToCache() {
	var services []mapper.Service
	min.Db.Find(&services)

	this.CacheService.DelPrefix(this.CACHE_PREFIX)

	for _, service := range services {
		this.CacheService.Set(this.CACHE_PREFIX, service.Code, service)
	}
}
