package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/skanehira/go-compose/docker"
	"github.com/skanehira/go-compose/model"
	"gopkg.in/yaml.v2"
)

var (
	composeFile = flag.String("f", "docker-compose.yaml", "yaml file, default is docker-compose.yaml")
	host        = flag.String("host", "unix:///var/run/docker.sock", "docker engine host")
	api         = flag.String("api", "1.40", "docker engine api version")
)

var notFoundComposeFile = `
Can't find a suitable configuration file in this directory or any
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
		fmt.Fprintf(os.Stderr, "Can't read file: %s\n", err)
		os.Exit(1)
	}

	var compose model.DockerCompose
	if err := yaml.Unmarshal(data, &compose); err != nil {
		fmt.Fprintf(os.Stderr, "Can't unmarshal yaml file: %s\n", err)
		os.Exit(1)
	}

	return compose
}

func main() {
	docker.NewDocker(*host, *api)
	//if _, err := docker.Client.Ping(context.Background()); err != nil {
	//	fmt.Fprintln(os.Stderr, err)
	//	os.Exit(1)
	//}
	compose := parseComposeFile()

	networkName := "compose_test"
	docker.Client.CreateNetwork(networkName)
	for n, s := range compose.Services {
		fmt.Println(n)
		_, err := docker.Client.CreateContainer(networkName, n, s)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
