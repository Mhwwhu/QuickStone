package main

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	. "github.com/mhwwhu/QuickStone/src/constant/config"
	"github.com/mhwwhu/QuickStone/src/rpc/auth"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", EnvCfg.PodIpAddr, Config.AuthServerPort))

	if err != nil {
		logrus.Panicf("Rpc %s listen error: %v", Config.AuthServerName, err)
	}

	s := grpc.NewServer()
	var srv AuthService
	auth.RegisterAuthServiceServer(s, srv)
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	srv.Init()
	s.Serve(lis)
}
