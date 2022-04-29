package service

import (
	"github.com/duanchi/min-gateway/mapper"
	"github.com/duanchi/min-gateway/storage"
	"github.com/duanchi/min-gateway/types/request"
	"github.com/duanchi/min-gateway/types/response"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/util"
)

type Integration struct {
	abstract.Bean

	StorageKey             string                          `value:"GATEWAY:INTEGRATION"`
	KEY                    string                          `value:"services"`
	ServiceStorage         *storage.ServiceStorage         `bean:"autowired"`
	ServiceInstanceStorage *storage.ServiceInstanceStorage `bean:"autowired"`
}

func (this *Integration) Init() {}

func (this *Integration) GetAll() []response.ServiceResponse {
	rawArray, _ := this.ServiceStorage.GetAll()
	serviceArray := []response.ServiceResponse{}
	instanceMap := this.ServiceInstanceStorage.GetAllGroupByServiceId()

	for _, value := range rawArray {
		service := response.ServiceResponse{
			Id:        value.Code,
			Name:      value.Name,
			Instances: []response.Instance{},
			Gray:      []response.Instance{},
		}

		instances, ok := instanceMap[value.Code]

		if ok {
			for _, instance := range instances {
				if instance.GrayFlag == 1 {
					service.Gray = append(service.Gray, response.Instance{
						Uri: instance.Url,
						Id:  instance.InstanceId,
					})
				} else {
					service.Instances = append(service.Instances, response.Instance{
						Uri: instance.Url,
						Id:  instance.InstanceId,
					})
				}

			}
		}

		serviceArray = append(serviceArray, service)
	}

	return serviceArray
}

func (this *Integration) Add(service request.Service) {
	code := util.GenerateUUID().String()
	ok := this.ServiceStorage.Add(mapper.Service{
		Code: code,
		Name: service.Name,
	})

	if ok {
		insertInstances := []mapper.ServiceInstance{}
		if len(service.Instances) > 0 {

			for _, instance := range service.Instances {
				insertInstances = append(insertInstances, mapper.ServiceInstance{
					InstanceId:    util.GenerateUUID().String(),
					GrayFlag:      0,
					OnlineFlag:    0,
					Weight:        0,
					Healthy:       0,
					Url:           instance.Uri,
					ServiceId:     code,
					EphemeralFlag: 0,
					CreateType:    0,
				})
			}
		}

		if len(service.Gray) > 0 {
			for _, instance := range service.Gray {
				insertInstances = append(insertInstances, mapper.ServiceInstance{
					InstanceId:    instance.Id,
					GrayFlag:      1,
					OnlineFlag:    0,
					Weight:        0,
					Healthy:       0,
					Url:           instance.Uri,
					ServiceId:     code,
					EphemeralFlag: 0,
					CreateType:    0,
				})
			}
		}

		if len(insertInstances) > 0 {
			this.ServiceInstanceStorage.AddList(insertInstances)
		}
	}
}

func (this *Integration) Modify(id string, modifiedService request.Service) {
	service, ok := this.ServiceStorage.GetFromDB(id)

	if ok {
		service.Name = modifiedService.Name
		this.ServiceInstanceStorage.RemoveByServiceId(service.Code)
		updateOk := this.ServiceStorage.Update(service)

		if updateOk {
			insertInstances := []mapper.ServiceInstance{}
			if len(modifiedService.Instances) > 0 {
				for _, instance := range modifiedService.Instances {

					instanceId := instance.Id

					if instanceId == "" {
						instanceId = util.GenerateUUID().String()
					}

					insertInstances = append(insertInstances, mapper.ServiceInstance{
						InstanceId:    instanceId,
						GrayFlag:      0,
						OnlineFlag:    1,
						Weight:        0,
						Healthy:       0,
						Url:           instance.Uri,
						ServiceId:     service.Code,
						EphemeralFlag: 0,
						CreateType:    0,
					})
				}
			}

			if len(modifiedService.Gray) > 0 {
				for _, instance := range modifiedService.Gray {
					insertInstances = append(insertInstances, mapper.ServiceInstance{
						InstanceId:    instance.Id,
						GrayFlag:      1,
						OnlineFlag:    1,
						Weight:        0,
						Healthy:       0,
						Url:           instance.Uri,
						ServiceId:     service.Code,
						EphemeralFlag: 0,
						CreateType:    0,
					})
				}
			}

			if len(insertInstances) > 0 {
				this.ServiceInstanceStorage.AddList(insertInstances)
			}
		}
	}

}

func (this *Integration) Remove(id string) {
	this.ServiceStorage.Remove(id)
	this.ServiceInstanceStorage.RemoveByServiceId(id)
}
