package scheduled

import (
	"github.com/duanchi/min-gateway/service"
	"github.com/duanchi/min-gateway/types/request"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/util"
	"os"
	"strings"
)

type AutoInitServiceSchedule struct {
	abstract.Scheduled
	ServiceService *service.Service `bean:"autowired"`
}

func (this *AutoInitServiceSchedule) Run() {
	initServices := map[string][]string{}
	for _, env := range os.Environ() {
		envPair := strings.SplitN(env, "=", 2)
		if strings.HasPrefix(envPair[0], "GATEWAY_AUTO_SERVICE_") && envPair[0] != "GATEWAY_AUTO_SERVICE_" {
			initServices[envPair[0][21:]] = strings.Split(envPair[1], ",")
		}
	}

	if len(initServices) > 0 {
		services := this.ServiceService.GetAll()
		for _, oneService := range services {
			if initInstances, has := initServices[oneService.Name]; has {
				grayInstance := []request.Instance{}
				staticInstance := []request.Instance{}

				for _, instance := range oneService.Instances {
					if !instance.IsStatic {
						staticInstance = append(staticInstance, request.Instance{
							Uri:         instance.Uri,
							Id:          instance.Id,
							IsEphemeral: instance.IsEphemeral,
							IsOnline:    instance.IsOnline,
							IsStatic:    false,
						})
					}
				}

				for _, initInstance := range initInstances {
					staticInstance = append(staticInstance, request.Instance{
						Uri:         initInstance,
						Id:          util.GenerateUUID().String(),
						IsEphemeral: false,
						IsOnline:    true,
						IsStatic:    true,
					})
				}

				for _, gray := range oneService.Gray {
					grayInstance = append(grayInstance, request.Instance{
						Uri:         gray.Uri,
						Id:          gray.Id,
						IsEphemeral: gray.IsEphemeral,
						IsOnline:    gray.IsOnline,
						IsStatic:    false,
					})
				}
				this.ServiceService.Modify(oneService.Id, request.Service{
					Id:        oneService.Id,
					Name:      oneService.Name,
					Instances: staticInstance,
					Gray:      grayInstance,
				})

				delete(initServices, oneService.Name)
			}
		}

		if len(initServices) > 0 {
			for serviceName, instances := range initServices {

				initInstances := []request.Instance{}

				for _, uri := range instances {
					initInstances = append(initInstances, request.Instance{
						Uri:         uri,
						Id:          util.GenerateUUID().String(),
						IsEphemeral: false,
						IsOnline:    true,
						IsStatic:    true,
					})
				}

				this.ServiceService.Add(request.Service{
					Id:        util.GenerateUUID().String(),
					Name:      serviceName,
					Instances: initInstances,
					Gray:      nil,
				})
			}
		}
	}
}
