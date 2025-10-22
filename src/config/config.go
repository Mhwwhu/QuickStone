package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var EnvCfg envConfig

type envConfig struct {
	PodIpAddr        string `env:"POD_IP" envDefault:"localhost"`
	ConsulAddr       string `env:"CONSUL_ADDR" envDefault:"localhost"`
	ConsulPort       uint32 `env:"CONSUL_PORT" envDefault:"8500"`
	ConsulNamePrefix string `env:"CONSUL_NAME_PREFIX" envDefault:""`
}

const WebServiceName = "Web-service"
const WebServicePort = 10001
const AuthServerName = "Auth-service"
const AuthServerPort = 10002
const TransmissionServerName = "Transmission-service"
const TransmissionServerPort = 10003

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Infof("Cannot read env from file system.")
	}

	EnvCfg = envConfig{}

	if err := env.Parse(&EnvCfg); err != nil {
		panic("Cannot parse env from file system, please check the env.")
	}
}
