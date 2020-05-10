package compose

import (
	"fmt"

	"github.com/docker/cli/cli/compose/types"

	"github.com/skanehira/go-compose/docker"
)

type Compose interface {
	CreateService(*types.Config) error
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

func (c *compose) CreateService(config *types.Config) error {
	for _, net := range config.Networks {
		c.client.CreateNetwork(net.Name)
	}

	for _, service := range config.Services {
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

		_, err = c.client.CreateContainer(service)
		if err != nil {
			return err
		}
		fmt.Println(service.Name)
	}

	return nil
}
