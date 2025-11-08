package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"syscall"

	"QuickStone/src/config"

	"QuickStone/src/rpc/bucket"
	"QuickStone/src/utils/consul"

	"github.com/oklog/run"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var ServerId string

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.EnvCfg.PodIpAddr, config.BucketServerPort))

	if err != nil {
		logrus.Panicf("Rpc %s listen error: %v", config.BucketServerName, err)
	}

	s := grpc.NewServer()

	if ServerId, err = consul.RegisterConsul(config.BucketServerName, config.BucketServerPort); err != nil {
		logrus.Panicf("Rpc %s registering consul error: %v", config.BucketServerName, err)
	}

	var srv BucketService
	bucket.RegisterBucketServiceServer(s, srv)
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	srv.Init()
	g := &run.Group{}
	g.Add(func() error {
		return s.Serve(lis)
	}, func(err error) {
		s.GracefulStop()
		s.Stop()
		logrus.Errorf("Rpc %s listening error: %v", config.BucketServerName, err)
	})

	g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	if err := g.Run(); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Error when runing http server")
		os.Exit(1)
	}
}
