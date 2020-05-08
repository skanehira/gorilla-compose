package compose

import (
	"fmt"

	"github.com/skanehira/go-compose/docker"
	"github.com/skanehira/go-compose/model"
)

type Compose interface {
	CreateService(m model.DockerCompose) error
	Ping() error
}

type compose struct {
	client docker.Client
}

func NewCompose(c docker.Client) Compose {
	return &compose{client: c}
}

func (c *compose) Ping() error {
	return c.client.Ping()
}

func (c *compose) CreateService(m model.DockerCompose) error {
	networkName := m.Project + "_default"

	c.client.CreateNetwork(networkName)

	for serviceName, service := range m.Services {
		opts := docker.ImageOpts{"reference": service.Image}
		res, err := c.client.Images(opts)
		if err != nil {
			return err
		}
		if len(res) == 0 {
			if err := c.client.PullImage(service.Image); err != nil {
				return err
			}
		}

		_, err = c.client.CreateContainer(networkName, serviceName, service)
		if err != nil {
			return err
		}
		fmt.Println(serviceName)
	}

	return nil
}
