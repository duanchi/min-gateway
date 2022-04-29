package event

import (
	"fmt"
	"github.com/duanchi/min-gateway/service"
	"github.com/duanchi/min-gateway/types/request"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/microservice/discovery"
	"github.com/duanchi/min/types"
	"strconv"
)

type DiscoveryServiceEvent struct {
	abstract.Event
	ServiceService *service.Service `bean:"autowired"`
}

func (this *DiscoveryServiceEvent) Conditions() (conditions []string) {

	conditions = []string{
		"DISCOVERED",
		"CACHED",
	}

	return
}

func (this *DiscoveryServiceEvent) Run(event types.Event, arguments ...interface{}) {
	fmt.Println("Updating services from discovery...")

	discoveryServices := discovery.GetServiceList()
	services := this.ServiceService.GetAll()

	if len(services) > 0 && len(discoveryServices) > 0 {
		for _, service := range services {
			serviceInstances := []request.Instance{}
			grayInstances := []request.Instance{}
			existInstances := map[string]string{}

			for _, tempInstance := range service.Instances {
				if !tempInstance.IsEphemeral {
					serviceInstances = append(serviceInstances, request.Instance{
						Uri:         tempInstance.Uri,
						Id:          tempInstance.Id,
						IsEphemeral: false,
					})
				}
				existInstances[tempInstance.Uri] = tempInstance.Id
			}

			for _, tempInstance := range service.Gray {
				if !tempInstance.IsEphemeral {
					grayInstances = append(grayInstances, request.Instance{
						Uri:         tempInstance.Uri,
						Id:          tempInstance.Id,
						IsEphemeral: false,
					})
				}
				existInstances[tempInstance.Uri] = tempInstance.Id
			}

			if discoveryService, ok := discoveryServices[service.Name]; ok {
				discoveryTotal := len(discoveryService.Instances)
				matchCount := 0

				for _, instance := range discoveryService.Instances {
					instanceId := ""
					if id, has := existInstances["http://"+instance.Ip+":"+strconv.FormatUint(instance.Port, 10)]; has {
						matchCount += 1
						instanceId = id
					}
					if id, ok := instance.Metadata["instance_id"]; ok {
						grayInstances = append(grayInstances, request.Instance{
							Uri:         "http://" + instance.Ip + ":" + strconv.FormatUint(instance.Port, 10),
							Id:          id,
							IsEphemeral: instance.Ephemeral,
						})
					} else {
						serviceInstances = append(serviceInstances, request.Instance{
							Uri:         "http://" + instance.Ip + ":" + strconv.FormatUint(instance.Port, 10),
							Id:          instanceId,
							IsEphemeral: instance.Ephemeral,
						})
					}
				}

				if discoveryTotal > matchCount {

					this.ServiceService.Modify(service.Id, request.Service{
						Id:        service.Id,
						Name:      service.Name,
						Instances: serviceInstances,
						Gray:      grayInstances,
					})
				} else {
					fmt.Println("No Service Instance Update!")
				}
			}
		}
	}
}
