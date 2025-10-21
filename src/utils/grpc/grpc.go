package grpc

import (
	"fmt"
	"time"

	"github.com/mhwwhu/QuickStone/src/constant/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func Connect(serviceName string) (conn *grpc.ClientConn) {
	kacp := keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: false,            // send pings even without active streams
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("consul://%s/%s?wait=15s", config.EnvCfg.ConsulAddr, config.EnvCfg.ConsulNamePrefix+serviceName),
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
