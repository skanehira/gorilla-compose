package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/cli/cli/compose/loader"
	"github.com/docker/cli/cli/compose/types"
	"github.com/skanehira/go-compose/compose"
	"github.com/skanehira/go-compose/docker"
	"github.com/urfave/cli/v2"
)

var notFoundComposeFile = `Can't find a suitable configuration file in this directory or any
parent. Are you in the right directory?

Supported filenames: docker-compose.yaml
`

const banner = `
      .""".
  .-./ _=_ \.-.
 {  (,(oYo),) }}
 {{ |   "   |} }
 { { \(---)/  }}
 {{  }'-=-'{ } }
 { { }._:_.{  }}
 {{  } -:- { } }
 {_{ }"==="{  _}
((((\)     (/))))

`

func currentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	return filepath.Base(dir)
}

func loadConfig(workDir, file string) (*types.Config, error) {
	yaml, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	source, err := loader.ParseYAML(yaml)
	if err != nil {
		return nil, err
	}

	detail := types.ConfigDetails{
		WorkingDir: workDir,
		ConfigFiles: []types.ConfigFile{
			{Filename: file, Config: source},
		},
	}

	config, err := loader.Load(detail)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func doUp(host, project, file string) error {
	config, err := loadConfig(project, file)
	if err != nil {
		return err
	}

	com := compose.NewCompose(docker.NewClient(docker.ClientConfig{Host: host, ApiVersion: "1.40"}))
	if err := com.Ping(); err != nil {
		return err
	}

	if err := com.CreateService(config); err != nil {
		return err
	}

	return nil
}

func main() {
	fmt.Print(banner)

	var (
		file    string
		project string
		host    string
	)

	app := &cli.App{
		Name:  "gorilla-compose",
		Usage: "Define and run multi-container applications with Docker.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				Value:       "docker-compose.yml",
				Usage:       "Specify an alternate compose file",
				Destination: &file,
			},
			&cli.StringFlag{
				Name:        "project-name",
				Aliases:     []string{"p"},
				Value:       currentDir(),
				Usage:       "Specify an alternate project name",
				Destination: &project,
			},
			&cli.StringFlag{
				Name:        "host",
				Aliases:     []string{"H"},
				Value:       "unix:///var/run/docker.sock",
				Usage:       "Daemon socket to connect to",
				Destination: &host,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "up",
				Usage: "Create and start application services",
				Action: func(c *cli.Context) error {
					return doUp(host, project, file)
				},
			},
			// TODO implement down
			//{
			//	Name:  "down",
			//	Usage: "Stop services created by `up`",
			//	Action: func(c *cli.Context) error {
			//		return nil
			//	},
			//},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
