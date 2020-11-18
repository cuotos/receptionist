package main

import (
	"html/template"
	"io/ioutil"
	"log"

	"github.com/markbates/pkger"
)

type Model struct {
	Containers []Container
	Version    string
}

func getIndexTpl() *template.Template {

	tplFile, err := pkger.Open("/assets/index.html.tpl")
	if err != nil {
		log.Fatal("failed to open index template file")
	}
	defer tplFile.Close()

	tplBody, err := ioutil.ReadAll(tplFile)
	if err != nil {
		log.Fatal("unable to read index template file")
	}
	return template.Must(template.New("IndexTpl").Parse(string(tplBody)))
}

var Tpl = template.Must(template.New("tpl").Parse(`
<html lang="en">
  <head>
	<title>Receptionist</title>
    <link rel="stylesheet" href="/static/css/bootstrap.min.css" integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous"/>
    <meta name="google" value="notranslate"/>
  <link rel="shortcut icon" href="/static/img/favicon.ico"/>
  <script src="/static/js/receptionist.js"></script>
  </head>
  <body>
    <div class="container mt-4">
      <h1>Receptionist</h1>
      <h3>"one moment caller, putting you through now..."</h3>
      <table class="table table-striped">
        <thead>
          <th>Container Name</th>
          <th>Ports</th>
          <th>Image</th>
        </thead>
        <tbody>
        {{ range .Containers }}
        <tr>
          <td>{{ .Name }}</td>
		  <td>
			{{ range .Ports }}
				{{ if .Name }}
					<a class="rec-link" href="http://localhost:{{ .PublicPort }}{{ .Path }}" target="_blank">{{.PublicPort}} {{ if .Name }} - {{ .Name }} {{ end }}</a><br/>
				{{ else }}
					<a class="rec-link" href="http://localhost:{{ .PublicPort }}{{ .Path }}" target="_blank">{{ .PublicPort }}</a><br/>
				{{ end }}
			{{ end }}
		  </td>
          <td>{{ .Image }}</td>
        </tr>
        {{ end }}
        </tbody>
      </table>
    </div>
    <script type="text/javascript">
		var links = Array.from(document.getElementsByClassName("rec-link"));
		for ( let l of links ) {
			l.hostname = location.hostname;
		}
    </script>
  </body>
</html>
`))
