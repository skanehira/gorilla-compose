package compose

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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
	// get parent dir
	var path string
	var err error
	var sep string

	if runtime.GOOS == "windows" {
		sep = "\\"
	} else {
		sep = "/"
	}

	dir := filepath.Dir(m.Name)

	if dir == "." {
		path, err = os.Getwd()
		if err != nil {
			return err
		}
		dir = path
	}
	dirs := strings.Split(dir, sep)
	path = dirs[len(dirs)-1]

	// TODO pull image

	c.client.CreateNetwork(path)
	for n, s := range m.Services {
		fmt.Println(n)
		_, err := c.client.CreateContainer(path, n, s)
		if err != nil {
			return err
		}
	}

	return nil
}
