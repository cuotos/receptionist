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
	"sort"
	"strconv"
	"strings"
)

var (
	config *Config
)

type Config struct {
	Prefix string `envconfig:"WATCHLABEL" default:"RECEPTIONIST"`
}

type Port struct {
	PublicPort  uint16
	PrivatePort uint16
	Name        string
}

type Container struct {
	Ports []*Port
	Name  string
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config = &Config{}
	err := envconfig.Process("", config)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	log.Printf("listening on :8080")
	log.Printf(`using receptionist label "%v"`, config.Prefix)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		containers, err := getRunningContainers()
		if err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(500), 500)
			return
		}

		err = templates.Tpl.Execute(writer, containers)

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
		ports, err := getAllPortsFromContainer(c)
		if err != nil {
			return nil, err
		}

		if ports != nil {
			sortPorts(ports)
			model = append(model, Container{ports, strings.TrimPrefix(c.Names[0], "/")})
		}
	}
	return model, nil
}

func getAllPortsFromContainer(c types.Container) ([]*Port, error) {

	allPorts := []*Port{}

	if l, found := c.Labels["RECEPTIONIST"]; found {

		for _, p := range c.Ports {

			if p.PublicPort != 0 {

				newPort := &Port{
					PublicPort:  p.PublicPort,
					PrivatePort: p.PrivatePort,
				}

				err := populatePortName(newPort, l)
				if err != nil {
					return nil, fmt.Errorf("unable to populate port name: %w", err)
				}

				allPorts = append(allPorts, newPort)
			}
		}
	}

	return allPorts, nil
}

func populatePortName(p *Port, label string) error {
	labelElements := strings.Split(label, ",")

	for _, e := range labelElements {
		if strings.Contains(e, ":") {
			name := strings.Split(e, ":")[0]
			port := strings.Split(e, ":")[1]

			portUint, err := strconv.ParseUint(port, 10, 16)
			if err != nil {
				return fmt.Errorf("unable to parse port number from string: %w", err)
			}

			if p.PrivatePort == uint16(portUint) {
				p.Name = name
			}
		}
	}

	return nil
}

func sortPorts(ports []*Port){
	sort.Slice(ports, func(i, j int) bool {
		return ports[i].PublicPort < ports[j].PublicPort
	})
}