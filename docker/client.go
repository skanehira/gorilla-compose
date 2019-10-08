package docker

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/client"
	"github.com/skanehira/go-compose/model"
)

type Client interface {
	CreateContainer(networkName, serviceName string, s model.Service) (string, error)
	CreateNetwork(name string) (string, error)
	Ping() error
}

type Docker struct {
	*client.Client
}

// ClientConfig docker client config
type ClientConfig struct {
	Host       string
	ApiVersion string
}

// NewClient create new docker client
func NewClient(config ClientConfig) *Docker {
	// TODO use tls
	// refer client.NewEnvClient
	c, err := client.NewClient(config.Host, config.ApiVersion, nil, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return &Docker{c}
}

func (d *Docker) Ping() error {
	_, err := d.Client.Ping(context.Background())
	return err
}
