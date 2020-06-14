package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"receptionist/templates"
)

type App struct {
	Router       *mux.Router
	DockerClient *Client
	Config       *Config
}

func (a *App) handleIndex() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		containers, err := a.DockerClient.getRunningContainers(a.Config.Label)
		if err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(500), 500)
			return
		}

		sortContainers(containers)

		model := struct {
			Containers []Container
		}{
			containers,
		}

		err = templates.Tpl.Execute(writer, model)

		if err != nil {
			log.Printf("unable to render template: %v", err)
			return
		}
	}
}

func (a *App) handleAbout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "http://github.com/cuotos/receptionist")
	}
}
