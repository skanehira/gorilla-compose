package docker

import (
	"context"
	"os"

	"github.com/docker/cli/cli/streams"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/pkg/jsonmessage"
)

type ImageOpts map[string]string

func (d *Docker) Images(opts ImageOpts) ([]types.ImageSummary, error) {
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

func (d *Docker) PullImage(name string) error {
	res, err := d.Client.ImagePull(context.Background(), name, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	out := streams.NewOut(os.Stdout)
	return jsonmessage.DisplayJSONMessagesToStream(res, out, nil)
}
