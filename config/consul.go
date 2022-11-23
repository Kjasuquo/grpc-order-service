package config

import (
	"encoding/json"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"

	//"github.com/hashicorp/vault/api"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Consul struct {
	client *consulapi.Client
}

type IP struct {
	Query string
}

func getip2() string {
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return err.Error()
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err.Error()
	}

	var ip IP
	json.Unmarshal(body, &ip)

	return ip.Query
}

func NewConsulClient(configAddress string) (*Consul, error) {
	config := consulapi.DefaultConfig()
	config.Address = configAddress
	// Get a new client
	c, err := consulapi.NewClient(config)
	if err != nil {
		log.Println("Error creating client: ", err)
		panic(err)
	}
	return &Consul{client: c}, nil
}

func (c *Consul) ServiceRegistryWithConsul(serviceId string, serviceName string, thePort string, tags []string) {

	var port int
	address := getip2()
	if thePort[0] == ':' {
		port, _ = strconv.Atoi(thePort[1:len(thePort)])
	} else {
		port, _ = strconv.Atoi(thePort)
	}

	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = serviceId
	registration.Name = serviceName
	registration.Port = port
	registration.Tags = tags
	registration.Address = address
	if tags[0] == "GRPC" {
		registration.Check = &consulapi.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", address, port),
			Interval:                       "5s",
			Timeout:                        "5s",
			Notes:                          "Check if the service is alive",
			GRPCUseTLS:                     false,
			DeregisterCriticalServiceAfter: "5s",
		}
	} else {
		registration.Check = &consulapi.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d", address, port),
			Interval:                       "5s",
			Timeout:                        "5s",
			Notes:                          "Check if the service is alive",
			DeregisterCriticalServiceAfter: "5s",
		}
	}

	err := c.client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Printf("Failed to register service: %s:%v ", address, port)
	} else {
		log.Printf("successfully register service: %s:%v", address, port)
	}

}

func (c *Consul) GetConsulService(service, tag string) (string, error) {

	var serviceAddress string
	addrs, _, err := c.client.Health().Service(service, tag, true, nil)
	if len(addrs) == 0 && err == nil {
		return "", fmt.Errorf("service ( %s ) was not found", service)
	}
	if err != nil {
		return "", err
	}

	for _, v := range addrs {
		serviceAddress = fmt.Sprintf("%s:%d", v.Service.Address, v.Service.Port)
	}

	return serviceAddress, nil
}
