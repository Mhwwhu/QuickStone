package config

import (
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var EnvCfg envConfig
var Config config

type envConfig struct {
	PodIpAddr string `env:"POD_IP" envDefault:"localhost"`
}

type config struct {
	WebServiceName string `yaml:"WebServiceName"`
	WebServicePort uint32 `yaml:"WebServicePort"`

	AuthServerName string `yaml:"AuthServerName"`
	AuthServerPort uint32 `yaml:"AuthServerPort"`
}

func loadConfigFromFile(path string) (*config, error) {
	var cfg config
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		logrus.Errorf("Cannot parse config file")
		return nil, err
	}
	return &cfg, nil
}

func init() {
	loadConfigFromFile("./config.yaml")

	if err := godotenv.Load(); err != nil {
		logrus.Infof("Cannot read env from file system.")
	}

	EnvCfg = envConfig{}

	if err := env.Parse(&EnvCfg); err != nil {
		panic("Cannot parse env from file system, please check the env.")
	}
}
