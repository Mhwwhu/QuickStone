package rabbitmq

import (
	"fmt"

	"QuickStone/src/common"
	"QuickStone/src/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectMQ() *amqp.Connection {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/%s", config.EnvCfg.RabbitMQUserName, config.EnvCfg.RabbitMQPassword,
		config.EnvCfg.RabbitMQAddr, config.EnvCfg.RabbitMQPort, config.EnvCfg.RabbitMQVHostPrefix))
	common.ExitOnErr(err)
	return conn
}
