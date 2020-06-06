package templates

import "html/template"

var Tpl  = template.Must(template.New("tpl").Parse(`
<html lang="en">
  <head>
	<title>Receptionist</title>
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css" integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous">
    <meta name="google" value="notranslate">
  </head>
  </body>
    <div class="container mt-4">
      <h1>Receptionist - {{ .Hostname }}</h1>
      <h3>"one moment caller, putting you through now..."</h3>
      <table class="table table-striped">
        <thead>
          <th>Container Name</th>
          <th>Ports ( Host - [ Name ] (Container) )</th>
          <th>Internal Port</th>
          <th>Image</th>
        </thead>
        <tbody>
        {{ range .Containers }}
        <tr>
          <td>{{ .Name }}</td>
		  <td>
			{{ range .Ports }}
				<a class="rec-link" href="http://localhost:{{ .PublicPort }}" target="_blank">{{ .PublicPort }} {{ if .Name }} - {{ .Name }} {{ end }}</a></br>
			{{ end }}
		  </td>
          <td>
            {{ range .Ports }}
              {{.PrivatePort}}</br>
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
