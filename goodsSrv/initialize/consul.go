package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"srv/goodsSrv/global"
)

type RegArgs struct {
	Name    string
	ID      string
	Address string
	Port    int
	Tags    []string
}

var consulAddr string

func RegisterConsul(arg *RegArgs) *api.Client {
	consulAddr = fmt.Sprintf("%s:%d", global.SrvConfig.Consul.Host, global.SrvConfig.Consul.Port)
	config := api.DefaultConfig()
	config.Address = consulAddr
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	zap.S().Infoln(consulAddr)
	health := fmt.Sprintf("%s:%d", "host.docker.internal", arg.Port)
	zap.S().Infoln(health)
	zap.S().Infof("%#v", arg)
	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      arg.ID,
		Name:    arg.Name,
		Tags:    arg.Tags,
		Port:    arg.Port,
		Address: arg.Address,
		Check: &api.AgentServiceCheck{
			Interval:                       "5s",
			Timeout:                        "5s",
			GRPC:                           health,
			DeregisterCriticalServiceAfter: "10s",
		},
	})
	if err != nil {
		panic(err)
	}
	return client
}
