package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
	"receptionist/templates"
	"strings"
)

var (
	config *Config
)

type Config struct {
	Prefix string `envconfig:"WATCHVAR" default:"RECEPTIONIST"`
	Port   string `envconfig:"PORT" default:"8080"`
}

type Container struct {
	Port      string
	ModelName string
	types.ContainerJSON
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	config = &Config{}
	err := envconfig.Process("", config)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening on %v", config.Port)
	log.Printf(`using receptionist env var "%v"`, config.Prefix)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		model, err := getRunningContainers()
		if err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(500), 500)
			return
		}

		err = templates.Tpl.Execute(writer, model)

		if err != nil {
			return
		}
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.Port), nil))
}

func getRunningContainers() ([]Container, error) {

	model := []Container{}

	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return model, err
	}

	for _, c := range containers {
		cont, err := cli.ContainerInspect(context.Background(), c.ID)
		if err != nil {
			return nil, err
		}

		port, err := doWeWantThisContainer(cont)
		if err != nil {
			return nil, err
		}

		if port != "" {
			model = append(model, Container{port, strings.TrimPrefix(cont.Name, "/"), cont})
		}
	}

	return model, nil
}

func doWeWantThisContainer(c types.ContainerJSON) (string, error) {
	labels := c.Config.Labels

	if port, wanted := labels[config.Prefix]; wanted {
		return port, nil
	}

	return "", nil
}
