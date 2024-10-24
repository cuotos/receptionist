package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Router       *mux.Router
	DockerClient *Client
	Config       *Config
	Version      string
}

func (a *App) handleIndex() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		containers, err := a.DockerClient.getRunningContainers(a.Config.Label)
		if err != nil {
			log.Printf("[ERROR] %s\n", err)
			http.Error(writer, http.StatusText(500), 500)
			return
		}

		sortContainers(containers)

		model := Model{
			Version:    a.Version,
			Containers: containers,
		}

		buf := &bytes.Buffer{}

		err = getIndexTpl().Execute(buf, model)

		if err != nil {
			log.Printf("[ERROR] failed to render template: %v", err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if _, err := writer.Write(buf.Bytes()); err != nil {
			log.Printf("[ERROR] failed to write response: %v\n", err)
		}
	}
}

func (a *App) handleAbout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "http://github.com/cuotos/receptionist")
	}
}
