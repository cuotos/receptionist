package main

import (
	"context"
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
}

type Port struct {
	Port string
	Name string
}

type Container struct {
	Ports      []Port
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

	log.Printf("listening on :8080")
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

	log.Fatal(http.ListenAndServe(":8080", nil))
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

		ports, err := doWeWantThisContainer(cont)
		if err != nil {
			return nil, err
		}

		if ports != nil {
			model = append(model, Container{ports, strings.TrimPrefix(cont.Name, "/"), cont})
		}
	}

	return model, nil
}

func doWeWantThisContainer(c types.ContainerJSON) ([]Port, error) {
	labels := c.Config.Labels

	if ports, wanted := labels[config.Prefix]; wanted {
		ps, err := extractPorts(ports)
		if err != nil {
			return nil, err
		}
		return ps, nil
	}

	return nil, nil
}

func extractPorts(portsString string) ([]Port, error) {

	ports := []Port{}

	portStrings := strings.Split(portsString, ",")
	for _, s := range portStrings {
		if strings.TrimSpace(s) != "" {
			ports = append(ports, Port{s, ""})
		}
	}

	return ports, nil

}