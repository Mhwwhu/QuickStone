package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"syscall"

	"QuickStone/src/config"
	trans "QuickStone/src/rpc/transmission"
	"QuickStone/src/utils/consul"

	"github.com/oklog/run"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var ServerId string

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.EnvCfg.PodIpAddr, config.TransmissionServerPort))

	if err != nil {
		logrus.Panicf("Rpc %s listen error: %v", config.TransmissionServerName, err)
	}

	s := grpc.NewServer()

	if ServerId, err = consul.RegisterConsul(config.TransmissionServerName, config.TransmissionServerPort); err != nil {
		logrus.Panicf("Rpc %s registering consul error: %v", config.TransmissionServerName, err)
	}

	var srv TransmissionService
	trans.RegisterTransmissionServiceServer(s, srv)
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	srv.Init()
	g := &run.Group{}
	g.Add(func() error {
		return s.Serve(lis)
	}, func(err error) {
		s.GracefulStop()
		s.Stop()
		logrus.Errorf("Rpc %s listening error: %v", config.TransmissionServerName, err)
	})

	g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	if err := g.Run(); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Error when runing http server")
		os.Exit(1)
	}
}
