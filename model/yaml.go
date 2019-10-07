package model

// NOTE go-compose is only support 3,7
// Reference https://docs.docker.com/compose/compose-file/
type DockerCompose struct {
	Version  string
	Services map[string]Service
}

type Service struct {
	Image           string
	ContainerName   string
	Build           Build
	Ports           []string
	Volumes         []string
	Command         []string
	DependsOn       []string
	Entrypoint      []string
	EnvFile         []string
	Environment     map[string]string
	Expose          []int
	Labels          []string
	Restart         string
	StopGracePeriod string
	StopSignal      string
	User            string
	WorkingDir      string
	Domainname      string
	HostName        string
	MacAddress      string
	Privileged      bool
	ReadOnly        bool
	StdinOpen       bool
	Tty             bool

	// TODO CapAdd          []string
	// TODO CadDrop         []string
	// TODO CgroupParent    string
	// TODO support links
	// TODO support userns_mode
	// TODO support ulimits
	// TODO support tmpfs
	// TODO support sysctls
	// TODO support secrets
	// TODO support configs
	// TODO support credential_spec
	// TODO support deploy
	// TODO support dns
	// TODO support dns_search
	// TODO support external_links
	// TODO support extra_hosts
	// TODO support healthcheck
	// TODO support logging
	// TODO support network_mode
	// TODO support network aliases
	// TODO support networks
	// TODO support volumes long syntax
	// TODO support pid
}

type Build struct {
	Context    string
	Dockerfile string
	Args       []string
	Labels     []string
}
