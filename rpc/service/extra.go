package service

import (
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/rpc"
)



type ExtraService struct {
	abstract.Rpc

	BeanName string `value:"asdfgh"`
	PackageName string `value:"device-platform-service"`
	ApplicationName string `value:"127.0.0.1:9082/rpc"`
}

func (this *ExtraService) Init () {
	this.Rpc.Init()
	this.SetName(this.BeanName)
	this.SetPackageName(this.PackageName)
	this.SetApplicationName(this.ApplicationName)
}

func (this *ExtraService) Test (name string, kind string) (value int, err error) {
	err = rpc.Call(rpc.IN{name, kind}, &rpc.OUT{&value, &err}, this)

	if err != nil {
		panic(err)
	}
	return
}
