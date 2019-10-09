package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func (d *Docker) Images(opts map[string]interface{}) ([]types.ImageSummary, error) {
	options := types.ImageListOptions{
		All: true,
	}

	args := filters.NewArgs()
	for k, v := range opts {
		args.Add(k, v)
	}

	options.Filters = args

	return d.Client.ImageList(context.Background(), options)
}
