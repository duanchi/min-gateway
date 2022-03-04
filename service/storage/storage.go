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

/*func (this *StorageService) Set (key string, value interface{}, ttl int64) (err error) {
	valueString,_ := json.Marshal(value)
	this.Instance().Write(valueString)

	return
}*/

/*func (this *StorageService) Set (key string, value interface{}, ttl int64) (err error) {
	valueString,_ := json.Marshal(value)
	_, err = this.Instance().Set(
		key,
		valueString,
		time.Duration(ttl) * time.Second,
		).Result()
	this.instance().Write(valueString)

	return
}*/

/*func (this *StorageService) HSet (key string, field string, value interface{}, ttl int64) (err error) {
	valueString,_ := json.Marshal(value)
	if ttl == -1 {
		_, err = this.Instance().HSet(
			key,
			field,
			valueString,
		).Result()
	} else {
		_, err = this.Instance().HSet(
			key,
			field,
			valueString,
			time.Duration(ttl) * time.Second,
		).Result()
	}


	return
}

func (this *StorageService) HGet (key string, field string, value interface{}) (has bool, err error) {
	has = false
	valueString, err := this.Instance().HGet(key, field).Result()

	if err != nil {
		return
	}
	has = true
	err = json.Unmarshal([]byte(valueString), value)
	return
}

func (this *StorageService) HGetAll (key string, value interface{}) (has bool, err error) {
	has = false
	// *value, err = this.Instance().HGetAll(key).Result()

	valueStringList, err := this.Instance().HGetAll(key).Result()

	if err != nil {
		return
	}
	has = true

	for key, valueString := range valueStringList {
		val := reflect.New(reflect.TypeOf(value).Elem().Elem())

		err = json.Unmarshal([]byte(valueString), val.Interface())
		if err != nil {
			log.Fatal(err)
		}
		reflect.ValueOf(value).Elem().SetMapIndex(reflect.ValueOf(key), val.Elem())
	}
	// err = json.Unmarshal([]byte(valueString), value)
	return
}

func (this *StorageService) Get (key string, value interface{}) (has bool, err error) {
	has = false
	valueString, err := this.Instance().Get(key).Result()

	if err != nil {
		return
	}
	has = true
	err = json.Unmarshal([]byte(valueString), value)

	return
}

func (this *StorageService) Remove (key string) {
	this.Instance().Del(key)
}

func (this *StorageService) HRemove (key string, field string) {
	this.Instance().HDel(key, field)
}

func (this *StorageService) Has (key string) (has bool) {
	num, _ := this.Instance().Exists(key).Result()

	return num > 0
}

func (this *StorageService) HHas (key string, field string) (has bool) {
	has, _ = this.Instance().HExists(key, field).Result()
	return
}*/
