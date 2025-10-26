package grpc

import (
	"fmt"
	"time"

	"QuickStone/src/config"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func Connect(serviceName string) (conn *grpc.ClientConn) {
	kacp := keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             time.Second,
		PermitWithoutStream: false,
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("consul://%s:%d/%s?wait=15s", config.EnvCfg.ConsulAddr, config.EnvCfg.ConsulPort, config.EnvCfg.ConsulNamePrefix+serviceName),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithKeepaliveParams(kacp),
	)

	logrus.Debugf("connect %s", serviceName)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"service": config.EnvCfg.ConsulNamePrefix + serviceName,
			"err":     err,
		}).Errorf("Cannot connect to %v service", config.EnvCfg.ConsulNamePrefix+serviceName)
	}
	return
}
