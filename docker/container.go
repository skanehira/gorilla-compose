package docker

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
	"github.com/skanehira/go-compose/model"
)

func parseEnvironment(e map[string]string) []string {
	var envs []string
	for k, v := range e {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}
	return envs
}

func parseConfig(s model.Service) *container.Config {
	return &container.Config{
		Hostname:     s.HostName,
		Domainname:   s.Domainname,
		User:         s.User,
		AttachStdin:  s.StdinOpen, // when stdin_open is true, AttachStdin is true
		AttachStdout: true,        // default is true
		AttachStderr: true,        // default is true
		ExposedPorts: map[nat.Port]struct{}{
			"": {},
		},
		Tty:         s.Tty,
		OpenStdin:   s.StdinOpen,
		StdinOnce:   s.StdinOpen,                     // AttachStdin and OpenStdin is true, StdinOnce is true
		Env:         parseEnvironment(s.Environment), // TODO read from env file
		Cmd:         strslice.StrSlice(s.Command),
		ArgsEscaped: false, // TODO how to parse this option
		Image:       s.Image,
		//Volumes: map[string]struct{}{ // TODO how to parse this option
		//	"": {},
		//},
		WorkingDir: s.WorkingDir,
		Entrypoint: s.Entrypoint,
		MacAddress: s.MacAddress,
		// OnBuild:    nil, // TODO how to parse this option
		Labels:     s.Labels,
		StopSignal: s.StopSignal,
	}
}

func parseMounts(volumes []string) []mount.Mount {
	var mounts []mount.Mount
	for _, v := range volumes {
		kv := strings.SplitN(v, ":", 2)
		if len(kv) != 2 {
			fmt.Fprintf(os.Stderr, "go-compose is not support this syntax yet: %s\n", kv)
			os.Exit(1)
		}
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: kv[0],
			Target: kv[1],
		})
	}
	return mounts
}

func parsePorts(ports []string) map[nat.Port][]nat.PortBinding {
	_, portBinding, err := nat.ParsePortSpecs(ports)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return portBinding
}

func parseHostConfig(s model.Service) *container.HostConfig {
	return &container.HostConfig{
		Mounts:       parseMounts(s.Volumes),
		PortBindings: parsePorts(s.Ports),
	}
}

func parseNetworkingConfig(networkName, serviceName string, s model.Service) *network.NetworkingConfig {
	return &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			networkName: {Aliases: []string{serviceName}},
		},
	}
}

func (c *Docker) CreateContainer(networkName, serviceName string, s model.Service) (string, error) {
	config := parseConfig(s)
	hostConfig := parseHostConfig(s)
	networkingConfig := parseNetworkingConfig(networkName, serviceName, s)

	res, err := c.ContainerCreate(context.Background(), config, hostConfig, networkingConfig, s.ContainerName)
	if err != nil {
		return "", err
	}

	for _, w := range res.Warnings {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("WARNING: %s", w))
	}

	return res.ID, nil
}
