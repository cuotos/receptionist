package main

import (
	"log"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

var (
	appVersion = "unset"
)

type Port struct {
	PublicPort  uint16
	PrivatePort uint16
	Name        string
	Path        string
}

type Container struct {
	Ports []*Port
	Name  string
	Image string
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	c := &Config{}
	err := envconfig.Process("", c)
	if err != nil {
		log.Fatal(err)
	}

	cl, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

	app := App{
		Router:       mux.NewRouter(),
		Config:       c,
		DockerClient: &Client{cl},
		Version:      appVersion,
	}
	app.routes()

	log.Printf("listening on :8080")
	log.Printf(`using receptionist label "%v"`, c.Label)

	log.Fatal(http.ListenAndServe(":8080", app.Router))
}
