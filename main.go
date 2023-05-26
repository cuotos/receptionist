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
	port       = ":8080"
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

func main() {

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"INFO", "ERROR", "FATAL"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	c := &Config{}
	err := envconfig.Process("", c)
	if err != nil {
		log.Fatalf("[FATAL] failed to parse env vars: %s", err)
	}

	cl, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
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

	log.Printf("[INFO] using receptionist label: %v", c.Label)
	log.Printf("[INFO] listening on %s", port)

	if c.TLSCertFile == "" || c.TLSKeyFile == "" {
		log.Printf("[INFO] running insecure")
		if err := http.ListenAndServe(port, app.Router); err != nil {
			log.Fatalf("[FATAL] %s", err)
		}
	} else {
		log.Printf("[INFO] running with tls")
		if err = http.ListenAndServeTLS(port, c.TLSCertFile, c.TLSKeyFile, app.Router); err != nil {
			log.Fatalf("[FATAL] %s", err)
		}
	}
}
