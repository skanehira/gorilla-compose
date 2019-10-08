package docker

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
)

func (c *Docker) CreateNetwork(name string) (string, error) {
	res, err := c.Client.NetworkCreate(context.Background(), name, types.NetworkCreate{Attachable: true})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return res.ID, nil
}
