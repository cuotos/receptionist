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

		ports, err := doWeWantThisContainer(cont, config.Prefix)
		if err != nil {
			return nil, err
		}

		if ports != nil {
			model = append(model, Container{ports, strings.TrimPrefix(cont.Name, "/"), cont})
		}
	}

	return model, nil
}

func doWeWantThisContainer(c types.ContainerJSON, label string) ([]Port, error) {
	labels := c.Config.Labels

	if ports, wanted := labels[label]; wanted {
		ps, _, err := parsePortsLabel(ports)
		if err != nil {
			return nil, err
		}
		return ps, nil
	}

	return nil, nil
}

func parsePortsLabel(portsString string) (ports []Port, listAllPorts bool, err error) {

	// split the argument by comma into port "element"
	portStrings := strings.Split(portsString, ",")

	for _, s := range portStrings {

		// clean the element up
		s = strings.TrimSpace(s)

		// if it is not empty
		if s != "" {

			var port, name string

			// the element contains only ALL, so turn on the flag that the user wants to list all exposed ports
			if s == "ALL" {
				listAllPorts = true

			} else {
				// If the element contains a colon, then get the port number and name from the element
				if strings.Contains(s, ":") {
					splitPort := strings.Split(s, ":")
					port = splitPort[1]
					name = splitPort[0]
				} else {
					port = s
					name = ""
				}

				ports = append(ports, Port{port, name})
			}
		}
	}

	return

}