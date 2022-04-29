package service

import (
	"github.com/duanchi/min-gateway/mapper"
	"github.com/duanchi/min-gateway/storage"
	"github.com/duanchi/min-gateway/types/request"
	"github.com/duanchi/min-gateway/types/response"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/util"
)

type Service struct {
	abstract.Bean

	StorageKey             string                          `value:"GATEWAY:SERVICES"`
	KEY                    string                          `value:"services"`
	ServiceStorage         *storage.ServiceStorage         `bean:"autowired"`
	ServiceInstanceStorage *storage.ServiceInstanceStorage `bean:"autowired"`
}

func (this *Service) Init() {}

func (this *Service) GetAll() []response.ServiceResponse {
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
						Uri:         instance.Url,
						Id:          instance.InstanceId,
						IsEphemeral: mapper.CONSTANT.BOOLEAN_TYPE[instance.EphemeralFlag],
						IsOnline:    mapper.CONSTANT.BOOLEAN_TYPE[instance.OnlineFlag],
					})
				} else {
					service.Instances = append(service.Instances, response.Instance{
						Uri:         instance.Url,
						Id:          instance.InstanceId,
						IsEphemeral: mapper.CONSTANT.BOOLEAN_TYPE[instance.EphemeralFlag],
						IsOnline:    mapper.CONSTANT.BOOLEAN_TYPE[instance.OnlineFlag],
					})
				}

			}
		}

		serviceArray = append(serviceArray, service)
	}

	return serviceArray
}

func (this *Service) Add(service request.Service) {
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
					OnlineFlag:    mapper.CONSTANT.BOOLEAN_TYPE_REVERSE[instance.IsOnline],
					Weight:        0,
					Healthy:       0,
					Url:           instance.Uri,
					ServiceId:     code,
					EphemeralFlag: mapper.CONSTANT.BOOLEAN_TYPE_REVERSE[instance.IsEphemeral],
					CreateType:    0,
				})
			}
		}

		if len(service.Gray) > 0 {
			for _, instance := range service.Gray {
				insertInstances = append(insertInstances, mapper.ServiceInstance{
					InstanceId:    instance.Id,
					GrayFlag:      1,
					OnlineFlag:    mapper.CONSTANT.BOOLEAN_TYPE_REVERSE[instance.IsOnline],
					Weight:        0,
					Healthy:       0,
					Url:           instance.Uri,
					ServiceId:     code,
					EphemeralFlag: mapper.CONSTANT.BOOLEAN_TYPE_REVERSE[instance.IsEphemeral],
					CreateType:    0,
				})
			}
		}

		if len(insertInstances) > 0 {
			this.ServiceInstanceStorage.AddList(insertInstances)
		}
	}
}

func (this *Service) Modify(id string, modifiedService request.Service) {
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
						OnlineFlag:    mapper.CONSTANT.BOOLEAN_TYPE_REVERSE[instance.IsOnline],
						Weight:        0,
						Healthy:       0,
						Url:           instance.Uri,
						ServiceId:     service.Code,
						EphemeralFlag: mapper.CONSTANT.BOOLEAN_TYPE_REVERSE[instance.IsEphemeral],
						CreateType:    0,
					})
				}
			}

			if len(modifiedService.Gray) > 0 {
				for _, instance := range modifiedService.Gray {
					insertInstances = append(insertInstances, mapper.ServiceInstance{
						InstanceId:    instance.Id,
						GrayFlag:      1,
						OnlineFlag:    mapper.CONSTANT.BOOLEAN_TYPE_REVERSE[instance.IsOnline],
						Weight:        0,
						Healthy:       0,
						Url:           instance.Uri,
						ServiceId:     service.Code,
						EphemeralFlag: mapper.CONSTANT.BOOLEAN_TYPE_REVERSE[instance.IsEphemeral],
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

func (this *Service) Remove(id string) {
	this.ServiceStorage.Remove(id)
	this.ServiceInstanceStorage.RemoveByServiceId(id)
}
