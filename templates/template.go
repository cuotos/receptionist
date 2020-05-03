package templates

import "html/template"

var Tpl  = template.Must(template.New("tpl").Parse(`
<html>
  <head>
	<title>Receptionist</title>
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css" integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous">
  </head>
  </body>
    <div class="container">
      <h1>"please hold..."
      <table class="table table-striped">
        <thead>
          <th>Name</th>
          <th>Port</th>
          <th>Image</th>
        </thead>
        <tbody>
        {{ range . }}
        <tr>
          <td>{{ .ModelName }}</td>
		  <td><a href="http://localhost:{{ .Port }}" target="_blank">{{ .Port }}</a></td>
          <td>{{ .Config.Image }}</td>
        </tr>
        {{ end }}
        </tbody>
      </table>
    </div>
  </body>
</html>
`))
