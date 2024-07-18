package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	TLS struct {
		Enabled  bool   `yaml:"enabled"`
		CertFile string `yaml:"cert_file"`
		KeyFile  string `yaml:"key_file"`
	} `yaml:"tls"`
	Port       int    `yaml:"port"`
	RootDirSRV string `yaml:"root_directory"`
}
type GlobalConfig struct {
	Portrange struct {
		Min int `yaml:"min"`
		Max int `yaml:"max"`
	} `yaml:"portrange"`
}

type GameImage struct {
	Nom   string `yaml:"nom"`
	Tag   string `yaml:"tag"`
	Ports struct {
		TCP []string `yaml:"tcp"`
		UDP []string `yaml:"udp"`
	} `yaml:"ports"`
	Spec struct {
		Core int `yaml:"core"`
		RAM  int `yaml:"ram"`
	} `yaml:"spec"`
}

type ConfigCephal struct {
	Server     ServerConfig `yaml:"server"`
	Global     GlobalConfig `yaml:"global"`
	GameImages []GameImage  `yaml:"gameimages"`
}

func LoadConfig(filename string) (*ConfigCephal, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config ConfigCephal
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
