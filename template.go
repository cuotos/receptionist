package main

import (
	_ "embed"
	"html/template"
)

//go:embed assets/index.html.tpl
var indexTpl []byte

type Model struct {
	Containers []Container
	Version    string
}

func getIndexTpl() *template.Template {

	return template.Must(template.New("IndexTpl").Parse(string(indexTpl)))
}
