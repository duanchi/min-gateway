package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func ReadJson (fileDir string, fileName string, bind interface{}, defaults string) (err error) {
	filePath := fileDir + fileName

	_, err = os.Stat(fileDir)

	if err != nil {
		err = os.MkdirAll(fileDir, 0755)

		if err != nil {
			log.Fatal(err)
		} else {
			ioutil.WriteFile(filePath, []byte(defaults), 0755)
		}
	}

	data, err := ioutil.ReadFile(filePath)
	jsonErr := json.Unmarshal(data, bind)

	if err != nil || jsonErr != nil {
		err := ioutil.WriteFile(filePath, []byte(defaults), 0755)
		if err != nil {
			log.Fatal(err)
		}
		data = []byte("[]")
	}

	err = json.Unmarshal(data, bind)

	return
}


func WriteJson (fileDir string, fileName string, bind interface{}) (err error) {

	filePath := fileDir + fileName
	data, err := json.Marshal(bind)

	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(fileDir)

	if err != nil {
		err = os.MkdirAll(fileDir, 0755)

		if err != nil {
			log.Fatal(err)
		}
	}

	ioutil.WriteFile(filePath, data, os.ModeAppend)

	return
}