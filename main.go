package main

import (
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
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
	Path        string
}

type Container struct {
	Ports []*Port
	Name  string
	Image string
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

	app := App{
		Router: mux.NewRouter(),
	}
	app.routes()

	log.Fatal(http.ListenAndServe(":8080", app.Router))
}
