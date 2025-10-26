package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var EnvCfg envConfig

type envConfig struct {
	PodIpAddr           string `env:"POD_IP" envDefault:"localhost"`
	ConsulAddr          string `env:"CONSUL_ADDR" envDefault:"localhost"`
	ConsulPort          uint32 `env:"CONSUL_PORT" envDefault:"8500"`
	ConsulNamePrefix    string `env:"CONSUL_NAME_PREFIX" envDefault:""`
	RabbitMQUserName    string `env:"RABBITMQ_USER_NAME" envDefault:"USER"`
	RabbitMQAddr        string `env:"RABBITMQ_ADDR" envDefault:"localhost"`
	RabbitMQPassword    string `env:"RABBITMQ_PASSWORD" envDefault:"123456"`
	RabbitMQPort        uint32 `env:"RABBITMQ_PORT" envDefault:"5672"`
	RabbitMQVHostPrefix string `env:"RABBITMQ_VHOST_PREFIX" envDefault:""`

	FsRootPath  string `env:"FS_ROOT_PATH" envDefault:"/tmp"`
	StorageType string `env:"STORAGE_TYPE" envDefault:"fs"`

	PostgreSQLAddr        string `env:"PG_ADDR" envDefault:"localhost"`
	PostgreSQLPort        uint32 `env:"PG_PORT" envDefault:"5432"`
	PostgreSQLTablePrefix string `env:"PG_TABLE_PREFIX" envDefault:""`
	PostgreSQLUser        string `env:"PG_USER" envDefault:"QuickStone"`
	PostgreSQLPassword    string `env:"PG_PWD" envDefault:"123456"`
	PostgreSQLDatabase    string `env:"PG_DATABASE" envDefault:"QuickStone"`

	JwtSecretKey string `env:"JWT_SECRET_KEY" envDefault:"123456"`
}

const WebServiceName = "Web-service"
const WebServicePort = 10001
const UserServerName = "User-service"
const UserServerPort = 10002
const TransmissionServerName = "Transmission-service"
const TransmissionServerPort = 10003
const MetadataServerName = "Metadata-service"
const MetadataServerPort = 10004

const GrpcStreamUploadSliceSize = 1024 * 256

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Infof("Cannot read env from file system.")
	}

	EnvCfg = envConfig{}

	if err := env.Parse(&EnvCfg); err != nil {
		panic("Cannot parse env from file system, please check the env.")
	}
}
