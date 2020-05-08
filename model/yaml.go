package model

// NOTE go-compose is only support 3,7
// Reference https://docs.docker.com/compose/compose-file/
type DockerCompose struct {
	Name     string             `yaml:"-"`
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
	Project  string             `yaml:"-"`
}

type Service struct {
	Name            string            `yaml:"-"`
	Image           string            `yaml:"image"`
	ContainerName   string            `yaml:"container_name"`
	Ports           []string          `yaml:"ports"`
	Volumes         []string          `yaml:"volumes"`
	Command         []string          `yaml:"command"`
	Entrypoint      []string          `yaml:"entrypoint"`
	Environment     map[string]string `yaml:"environment"`
	Labels          map[string]string `yaml:"labels"`
	Restart         string            `yaml:"restart"`
	StopGracePeriod string            `yaml:"stop_grace_period"`
	StopSignal      string            `yaml:"stop_signal"`
	User            string            `yaml:"user"`
	WorkingDir      string            `yaml:"working_dir"`
	Domainname      string            `yaml:"domainname"`
	HostName        string            `yaml:"host_name"`
	MacAddress      string            `yaml:"mac_address"`
	StdinOpen       bool              `yaml:"stdin_open"`
	Tty             bool              `yaml:"tty"`

	// Privileged      bool // TODO suppor privileged
	// EnvFile         []string // TODO support EnvFile
	// DependsOn       []string // TODO support DependsOn
	// Build           Build // TODO support build
	// Expose          []int // TODO support Expose
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
