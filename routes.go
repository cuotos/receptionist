package main

import (
	"github.com/markbates/pkger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func (a *App) routes() {
	a.Router.HandleFunc("/", a.handleIndex())
	a.Router.HandleFunc("/about", a.handleAbout())
	a.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(pkger.Dir("/static"))))
	a.Router.Handle("/metrics", promhttp.Handler())
}
