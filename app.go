package main

import (
	"bytes"
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

		buf := &bytes.Buffer{}

		err = templates.Tpl.Execute(buf, model)

		if err != nil {
			log.Printf("failed to render template: %v", err)
			http.Error(writer, fmt.Sprintf("failed to render template: %v", err), http.StatusInternalServerError)
			return
		}

		if _, err := writer.Write(buf.Bytes()); err != nil {
			log.Printf("failed to write response: %v\n", err)
		}

		return
	}
}

func (a *App) handleAbout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "http://github.com/cuotos/receptionist")
	}
}
