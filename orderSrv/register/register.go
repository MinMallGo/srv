package register

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

// SrvRegisterArgs 注册服务时必要的信息
type SrvRegisterArgs struct {
	Name string   // 服务名称
	ID   string   // 服务的id
	Host string   // 服务的ip
	Port int      // 服务的端口
	Tags []string // 服务的tag
}

// RegistrationCenter 定义注册中心需要实现的方法
// 实现即可替换
type RegistrationCenter interface {
	Register(args *SrvRegisterArgs) error
	Deregister(id string) error
}

// ConsulRegistry 注册中心的地址
type ConsulRegistry struct {
	Host string // 注册中心的ip
	Port int    // 注册中心的端口
}

// NewConsulRegistry 实例化
func NewConsulRegistry(host string, port int) *ConsulRegistry {
	return &ConsulRegistry{
		Host: host,
		Port: port,
	}
}

// Register 注册到服务
func (c *ConsulRegistry) Register(arg *SrvRegisterArgs) error {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", c.Host, c.Port)
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}

	health := fmt.Sprintf("%s:%d", "host.docker.internal", arg.Port)
	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      arg.ID,
		Name:    arg.Name,
		Tags:    arg.Tags,
		Port:    arg.Port,
		Address: arg.Host,
		Check: &api.AgentServiceCheck{
			Interval:                       "5s",
			Timeout:                        "5s",
			GRPC:                           health,
			DeregisterCriticalServiceAfter: "10s",
		},
	})

	if err != nil {
		return err
	}
	return nil
}

// Deregister 从服务中退出
func (c *ConsulRegistry) Deregister(id string) error {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", c.Host, c.Port)
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}
	return client.Agent().ServiceDeregister(id)
}
