package main

import (
	"log"
	"net/http"
	"os"

	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
	"github.com/hashicorp/logutils"
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
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"INFO", "ERROR", "FATAL"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	c := &Config{}
	err := envconfig.Process("", c)
	if err != nil {
		log.Fatalf("[FATAL] failed to parse env vars: %s", err)
	}

	cl, err := client.NewEnvClient()
	if err != nil {
		log.Fatalf("[FATAL] %s", err)
	}

	app := App{
		Router:       mux.NewRouter(),
		Config:       c,
		DockerClient: &Client{cl},
		Version:      appVersion,
	}
	app.routes()

	log.Printf("[INFO] listening on :8080")
	log.Printf("[INFO] using receptionist label: %v", c.Label)

	if err := http.ListenAndServe(":8080", app.Router); err != nil {
		log.Fatalf("[FATAL] %s", err)
	}
}
