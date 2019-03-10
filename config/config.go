package config

import (
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var (
	// DeploySet - environmental parameters
	DeploySet DeploySetConfig
	// DefaultConfigPath - default config file path
	DefaultConfigPath = "/etc/golang-proxy-server/config.yaml"
)

// EnvConfig - yaml config formate struct
type EnvConfig struct {
	EnvHost    string `yaml:"env_host"`
	EnvPort    string `yaml:"env_port"`
	EnvTimeout int    `yaml:"env_timeout"`
}

// ExternalConfig - yaml config formate struct
type ExternalConfig struct {
	ExternalURL            string `yaml:"external_url"`
	ExternalMethod         string `yaml:"external_method"`
	ExternalLimitPer       int    `yaml:"external_limit_per"`
	ExternalRequestTimeout int    `yaml:"external_request_timeout"`
	ExternalRequestQueue   int    `yaml:"external_request_queue"`
}

// DeploySetConfig - config collection
type DeploySetConfig struct {
	Env      EnvConfig      `yaml:"env"`
	External ExternalConfig `yaml:"external"`
}

func init() {
	// configFile, err := filepath.Abs("./config/config.yaml")
	// if err != nil {
	// 	configFile = DefaultConfigPath
	// }

	var content []byte
	var ioErr error
	if content, ioErr = ioutil.ReadFile(DefaultConfigPath); ioErr != nil && strings.Index(ioErr.Error(), "no such file or directory") != -1 {
		panic("please copy config/config.yaml config file to " + DefaultConfigPath)
	} else if ioErr != nil {
		panic("read service config file error: " + ioErr.Error())
	}

	if ymlErr := yaml.Unmarshal(content, &DeploySet); ymlErr != nil {
		panic("error while unmarshal from db config: " + ymlErr.Error())
	}

}
