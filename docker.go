package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Client struct {
	Client *client.Client
}

func (c *Client) getRunningContainers(receptionistLabel string) ([]Container, error) {

	model := []Container{}

	containers, err := c.Client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return model, err
	}

	for _, c := range containers {

		var ports []*Port

		if l, found := c.Labels[receptionistLabel]; found {
			ports, err = getAllWantedPortsFromContainer(c, l)
			if err != nil {
				return nil, err
			}
		}

		if len(ports) > 0 {
			sortPorts(ports)
			model = append(model, Container{ports, strings.TrimPrefix(c.Names[0], "/"), c.Image})
		}
	}
	return model, nil
}

// getAllWantedPortsFromContainer extracts the publicly shared ports of a container.
// It returns a Port, which contains the publicly mounted port, the private container port, the name if provided
//   via the RECEPTIONIST label, and the default path that requests should be routed.
func getAllWantedPortsFromContainer(c types.Container, l string) ([]*Port, error) {
	//TODO: get the label string out of here, Parse the label value at startup and use that to hold the parsed data
	// that can be injected into the Port when required.
	var allPorts []*Port

	for _, p := range c.Ports {

		if p.PublicPort != 0 {

			port := &Port{
				PublicPort:  p.PublicPort,
				PrivatePort: p.PrivatePort,
				Path:        "/",
			}

			// Parse the RECEPTIONIST label add name and path to port if provided
			err := populatePortMetaData(port, l)
			if err != nil {
				return nil, fmt.Errorf("unable to populate port name: %w", err)
			}

			allPorts = append(allPorts, port)
		}
	}

	return allPorts, nil
}

type LabelElement struct {
	Name string
	Port string
	Path string
}

func extractElementsFromLabel(s string) ([]LabelElement, error) {
	regex := regexp.MustCompile(`(?P<Name>[^:]*):(?P<Port>[^:]+):?(?P<Path>[^:]*)`)

	var elements []LabelElement

	ports := strings.Split(s, ",")

	for _, p := range ports {

		els := regex.FindStringSubmatch(p)

		if len(els) > 0 {

			el := LabelElement{}
			el.Name = els[1]
			el.Port = els[2]
			el.Path = els[3]

			if el.Path == "" {
				el.Path = "/"
			}

			if !strings.HasPrefix(el.Path, "/") {
				el.Path = fmt.Sprintf("/%v", el.Path)
			}

			elements = append(elements, el)
		}
	}

	return elements, nil
}

func populatePortMetaData(p *Port, label string) error {
	labelElements, err := extractElementsFromLabel(label)
	if err != nil {
		return err
	}

	for _, e := range labelElements {
		portUint, err := strconv.ParseUint(e.Port, 10, 16)
		if err != nil {
			return fmt.Errorf("unable to parse port number from string: %w", err)
		}

		if portUint == uint64(p.PrivatePort) {
			p.Name = e.Name
			p.Path = e.Path
		}
	}

	return nil
}

func sortPorts(ports []*Port) {
	sort.Slice(ports, func(i, j int) bool {
		return ports[i].PublicPort < ports[j].PublicPort
	})
}

func sortContainers(cs []Container) {
	sort.Slice(cs, func(i, j int) bool {
		return cs[i].Name < cs[j].Name
	})
}
