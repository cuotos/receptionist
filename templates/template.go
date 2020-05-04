package templates

import "html/template"

var Tpl  = template.Must(template.New("tpl").Parse(`
<html>
  <head>
	<title>Receptionist</title>
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css" integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous">
  </head>
  </body>
    <div class="container mt-4">
      <h1>Receptionist</h1>
      <h3>"one moment caller, putting you through now..."</h3>
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
		  <td>
			{{ range .Ports }}
			<a class="rec-link" href="http://localhost:{{ .Port }}" target="_blank">{{ .Port }}</a>
			{{ end }}
		  </td>
          <td>{{ .Config.Image }}</td>
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
