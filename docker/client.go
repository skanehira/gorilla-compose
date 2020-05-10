package docker

import (
	"context"
	"fmt"
	"os"

	ctypes "github.com/docker/cli/cli/compose/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Client interface {
	CreateContainer(ctypes.ServiceConfig) (string, error)
	CreateNetwork(name string) (string, error)
	Images(opts ImageOpts) ([]types.ImageSummary, error)
	PullImage(name string) error
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
