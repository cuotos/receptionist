package main

import (
	"embed"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//go:embed static
var static embed.FS

func (a *App) routes() {
	a.Router.HandleFunc("/", a.handleIndex())
	a.Router.HandleFunc("/about", a.handleAbout())
	a.Router.PathPrefix("/static/").Handler(http.FileServer(http.FS(static)))
	a.Router.Handle("/metrics", promhttp.Handler())
}
