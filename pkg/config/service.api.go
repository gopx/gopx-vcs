package config

import (
	"encoding/json"
	"io/ioutil"

	"gopx.io/gopx-vcs/pkg/log"
)

// ServiceAPIConfigPath holds API service related configuration file path.
const ServiceAPIConfigPath = "./config/service.api.json"

// ServiceAPIConfig represents API service related configurations.
type ServiceAPIConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// APIService holds loaded API service related configurations.
var APIService = new(ServiceAPIConfig)

func init() {
	bytes, err := ioutil.ReadFile(ServiceAPIConfigPath)
	if err != nil {
		log.Fatal("Error: %s", err)
	}
	err = json.Unmarshal(bytes, APIService)
	if err != nil {
		log.Fatal("Error: %s", err)
	}
}
