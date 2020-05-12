package docker

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/cli/cli/compose/types"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
)

func convertEnvironment(e types.MappingWithEquals) []string {
	var envs []string
	for k, v := range e {
		if v != nil && *v != "" {
			envs = append(envs, k+"="+*v)
		}
	}
	return envs
}

func convertPosePort(ports []string) (natPorts map[nat.Port]struct{}) {
	for _, port := range ports {
		natPorts[nat.Port(port)] = struct{}{}
	}
	return
}

func convertContainerConfig(s types.ServiceConfig) *container.Config {
	return &container.Config{
		Hostname:     s.Hostname,
		Domainname:   s.DomainName,
		User:         s.User,
		AttachStdin:  s.StdinOpen,
		ExposedPorts: convertPosePort(s.Expose),
		Tty:          s.Tty,
		OpenStdin:    s.StdinOpen,
		StdinOnce:    s.StdinOpen,
		Env:          convertEnvironment(s.Environment), // TODO read from env file
		Cmd:          strslice.StrSlice(s.Command),
		ArgsEscaped:  false, // TODO parse this option
		Image:        s.Image,
		//Volumes: map[string]struct{}{ // TODO parse this option
		//	"": {},
		//},
		WorkingDir: s.WorkingDir,
		Entrypoint: strslice.StrSlice(s.Entrypoint),
		MacAddress: s.MacAddress,
		// OnBuild:    nil, // TODO how to parse this option
		Labels:     s.Labels,
		StopSignal: s.StopSignal,
	}
}

func parseMounts(configs []types.ServiceVolumeConfig) []mount.Mount {
	var mounts []mount.Mount
	for _, v := range configs {
		mounts = append(mounts, mount.Mount{
			Type:        mount.Type(v.Type),
			Source:      v.Source,
			Target:      v.Target,
			ReadOnly:    v.ReadOnly,
			Consistency: mount.Consistency(v.Consistency),
		})
		// TODO implement options
	}
	return mounts
}

func parsePorts(configs []types.ServicePortConfig) map[nat.Port][]nat.PortBinding {
	ports := []string{}
	for _, c := range configs {
		port := fmt.Sprintf("%d:%d", c.Published, c.Target)
		ports = append(ports, port)
	}

	_, portBinding, err := nat.ParsePortSpecs(ports)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return portBinding
}

func convertHostConfig(s types.ServiceConfig) *container.HostConfig {
	return &container.HostConfig{
		Mounts:       parseMounts(s.Volumes),
		PortBindings: parsePorts(s.Ports),
	}
}

func convertNetworkConfig(configs map[string]*types.ServiceNetworkConfig) *network.NetworkingConfig {
	endpoints := map[string]*network.EndpointSettings{}
	for name, c := range configs {
		setting := &network.EndpointSettings{
			Aliases: c.Aliases,
		}
		endpoints[name] = setting
	}
	return &network.NetworkingConfig{EndpointsConfig: endpoints}
}

func (c *Docker) CreateContainer(service types.ServiceConfig) (string, error) {
	containerConfig := convertContainerConfig(service)
	hostConfig := convertHostConfig(service)
	networkingConfig := convertNetworkConfig(service.Networks)

	res, err := c.ContainerCreate(context.Background(), containerConfig, hostConfig, networkingConfig, service.ContainerName)
	if err != nil {
		return "", err
	}

	for _, w := range res.Warnings {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("WARNING: %s", w))
	}

	return res.ID, nil
}
