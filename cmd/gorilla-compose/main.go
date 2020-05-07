package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/skanehira/go-compose/compose"
	"github.com/skanehira/go-compose/docker"
	"github.com/skanehira/go-compose/model"
	"gopkg.in/yaml.v2"
)

var (
	composeFile = flag.String("f", "docker-compose.yaml", "yaml file, default is docker-compose.yaml")
	host        = flag.String("host", "unix:///var/run/docker.sock", "docker engine host")
	api         = flag.String("api", "1.40", "docker engine api version")
)

var notFoundComposeFile = `Can't find a suitable configuration file in this directory or any
parent. Are you in the right directory?

Supported filenames: docker-compose.yaml
`

func init() {
	flag.Parse()
}

func parseComposeFile() model.DockerCompose {
	file, err := os.Open(*composeFile)
	if err != nil {
		fmt.Fprint(os.Stderr, notFoundComposeFile)
		os.Exit(1)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't read file: %s\n", err)
		os.Exit(1)
	}

	var compose model.DockerCompose
	if err := yaml.Unmarshal(data, &compose); err != nil {
		fmt.Fprintf(os.Stderr, "can't unmarshal yaml file: %s\n", err)
		os.Exit(1)
	}

	for name, co := range compose.Services {
		co.Name = name
	}

	compose.Name = *composeFile
	return compose
}

func main() {
	com := compose.NewCompose(docker.NewClient(docker.ClientConfig{Host: *host, ApiVersion: *api}))
	if err := com.Ping(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	compose := parseComposeFile()
	if err := com.CreateService(compose); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
