package docker

import (
	"context"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
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
	cmd := exec.Command("docker", "pull", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
