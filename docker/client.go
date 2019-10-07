package docker

import (
	"fmt"
	"os"

	"github.com/docker/docker/client"
)

// Client docker client
var Client *Docker

// Docker docker client
type Docker struct {
	*client.Client
}

// ClientConfig docker client config
type ClientConfig struct {
	host       string
	apiVersion string
}

// NewDocker create new docker client
func NewDocker(host, version string) *Docker {
	// TODO use tls
	// refer client.NewEnvClient
	c, err := client.NewClient(host, version, nil, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	Client = &Docker{c}

	return Client
}
