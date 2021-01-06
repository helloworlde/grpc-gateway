package server

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"

	"github.com/hashicorp/consul/api"
)

func RegisterToConsul(consulAddr string, port int) {
	config := api.DefaultConfig()
	config.Address = consulAddr
	client, err := api.NewClient(config)

	if err != nil {
		fmt.Println("注册Consul 失败:", err)
	}
	agent := client.Agent()

	address, err := getClientIp()

	checkServiceName := "grpc.health.v1.HealthCheck"
	grpcHealthCheck := fmt.Sprintf("%s:%d/%s", address, port, checkServiceName)

	weight := rand.Intn(10)

	err = agent.ServiceRegister(&api.AgentServiceRegistration{
		ID:   "server" + strconv.Itoa(weight),
		Name: "server",
		Tags: []string{"server"},
		Meta: nil,
		Port: port,
		Weights: &api.AgentWeights{
			Passing: weight,
			Warning: 0,
		},
		Address: address,
		Check: &api.AgentServiceCheck{
			CheckID:                        "Check" + strconv.Itoa(weight),
			Name:                           "Check",
			Interval:                       "5s",
			Timeout:                        "5s",
			GRPC:                           grpcHealthCheck,
			GRPCUseTLS:                     false,
			DeregisterCriticalServiceAfter: "1m",
		},
	})
	if err != nil {
		fmt.Println("注册 Consul 失败:", err)
	}
}

func getClientIp() (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), err
			}

		}
	}

	return "", errors.New("can not find the client ip address")

}
