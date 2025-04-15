package proxy

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Route struct {
	PathPrefix string `yaml:"path_prefix"`
	Backend    string `yaml:"backend"`
	Server     string `yaml:"server"`
}

type Config struct {
	Routes []Route `yaml:"routes"`
}

var ProxyConfig Config

func LoadConfig(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &ProxyConfig)
}
