package config

import (
	"encoding/json"
	"io/ioutil"

	"gopx.io/gopx-vcs/pkg/log"
)

// ServiceVCSConfigPath holds VCS service related configuration file path.
const ServiceVCSConfigPath = "./config/service.vcs.json"

// ServiceVCSConfig represents VCS service related configurations.
type ServiceVCSConfig struct {
	Host      string `json:"host"`
	UseHTTP   bool   `json:"useHTTP"`
	HTTPPort  int    `json:"HTTPPort"`
	UseHTTPS  bool   `json:"useHTTPS"`
	HTTPSPort int    `json:"HTTPSPort"`
	CertFile  string `json:"certFile"`
	KeyFile   string `json:"keyFile"`
}

// VCSService holds loaded VCS service related configurations.
var VCSService = new(ServiceVCSConfig)

func init() {
	bytes, err := ioutil.ReadFile(ServiceVCSConfigPath)
	if err != nil {
		log.Fatal("Error: %s", err)
	}
	err = json.Unmarshal(bytes, VCSService)
	if err != nil {
		log.Fatal("Error: %s", err)
	}
}
