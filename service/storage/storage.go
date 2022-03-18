package storage

import (
	"encoding/json"
	"fmt"
	"github.com/duanchi/min-gateway/types"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/config"
	"github.com/duanchi/min/util"
	"io/ioutil"
	"os"
)

type Configuration struct {
	Routes   types.RoutesMap
	Services types.ServicesMap
}

type StorageService struct {
	abstract.Service

	DataPath        string `value:"${Gateway.DataPath}"`
	JwtSignatureKey string `value:"${Authorization.SignatureKey}"`
	JwtExpiresIn    int64  `value:"${Authorization.Ttl}"`
	Dsn             string `value:"${Authorization.Dsn}"`
	random          string
	instance        *os.File
	Configuration   Configuration
	Inited          bool
}

func (this *StorageService) Init() {
	this.random = util.GenerateUUID().String()
	this.Inited = false
	this.Instance()
}

func (this *StorageService) Instance() {

	dataPath := config.Get("Gateway.DataPath").(string)

	if this.instance == nil {
		defer this.instance.Close()
		var err error
		WaitGroup.Add(1)
		fmt.Println("Configuration file locate: " + dataPath + "/configuration.json")
		this.instance, err = os.OpenFile(dataPath+"/configuration.json", os.O_RDWR, 0755)
		this.Inited = true
		WaitGroup.Done()
		if err != nil {
			fmt.Printf("Cannot Read Configuration File, Create it!\r\n")
			this.instance, _ = os.Create(dataPath + "/configuration.json")
			data, _ := json.Marshal(this.Configuration)
			this.instance.Write(data)
		} else {
			valueString, _ := ioutil.ReadAll(this.instance)
			json.Unmarshal(valueString, &this.Configuration)
		}
	}
	// return this.instance
}

func (this *StorageService) Save(key string, value interface{}, field string) {
	if field == "services" {
		if this.Configuration.Services == nil {
			this.Configuration.Services = types.ServicesMap{
				key: value.(types.Service),
			}
		} else {
			this.Configuration.Services[key] = value.(types.Service)
		}

	} else if field == "routes" {
		if this.Configuration.Routes == nil {
			this.Configuration.Routes = types.RoutesMap{
				key: value.(types.Route),
			}
		} else {
			this.Configuration.Routes[key] = value.(types.Route)
		}
	}

	this.Update()
}

func (this *StorageService) Update() {
	valueString, _ := json.Marshal(this.Configuration)
	// this.Instance().Write(valueString)
	f, err := os.OpenFile(this.DataPath+"/configuration.json", os.O_WRONLY|os.O_TRUNC, 0755)
	if err == nil {
		f.Write(valueString)
		f.Close()
	}
}

func (this *StorageService) Remove(key string, field string) {
	if field == "services" {
		delete(this.Configuration.Services, key)
	} else if field == "routes" {
		delete(this.Configuration.Routes, key)
	}
	this.Update()
}

func (this *StorageService) Get(field string) interface{} {
	if field == "services" {
		return this.Configuration.Services
	} else if field == "routes" {
		return this.Configuration.Routes
	}

	return nil
}

func (this *StorageService) Export() (data []byte, err error) {
	jsonFile, err := os.Open(this.DataPath + "/configuration.json")
	defer jsonFile.Close()

	data, err = ioutil.ReadAll(jsonFile)

	return
}

func (this *StorageService) Import(data []byte) {
	ioutil.WriteFile(this.DataPath+"/configuration.json", data, 0755)
}
