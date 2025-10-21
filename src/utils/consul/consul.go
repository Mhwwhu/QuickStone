package consul

import (
	"fmt"

	"github.com/mhwwhu/QuickStone/src/config"

	"github.com/google/uuid"
	capi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

var consulClient *capi.Client

func init() {
	cfg := capi.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", config.EnvCfg.ConsulAddr, config.EnvCfg.ConsulPort)
	if c, err := capi.NewClient(cfg); err == nil {
		consulClient = c
		return
	} else {
		logrus.Panicf("Connecting consul error: %v", err)
	}
}

func RegisterConsul(name string, port uint32) (string, error) {
	logrus.WithFields(logrus.Fields{
		"name": name,
		"port": port,
	}).Infof("Services Register Consul")
	name = config.EnvCfg.ConsulNamePrefix + name

	id := fmt.Sprintf("%s-%s", name, uuid.New().String())
	reg := &capi.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Address: config.EnvCfg.PodIpAddr,
		Port:    int(port),
		Check: &capi.AgentServiceCheck{
			Interval:                       "5s",
			Timeout:                        "5s",
			GRPC:                           fmt.Sprintf("%s:%d", config.EnvCfg.PodIpAddr, port),
			GRPCUseTLS:                     false,
			DeregisterCriticalServiceAfter: "30s",
		},
	}
	if err := consulClient.Agent().ServiceRegister(reg); err != nil {
		return id, err
	}
	return id, nil
}
