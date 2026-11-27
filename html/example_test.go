package html_test

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"testing/fstest"

	"github.com/tomcz/gotools/html"
)

// layout.gohtml: default layout template
const layoutTmpl = `
{{ define "main" -}}
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  {{- template "head" . }}
</head>
<body>
{{- template "body" . }}
</body>
</html>
{{- end }}
{{- define "head" }}{{ end }}
{{- define "body" }}{{ end }}
`

// content.gohtml: content of page to be merged into the layout
const contentTmpl = `
{{- define "head" }}
  <title>{{ .Title }}</title>
{{- end }}
{{- define "body" }}
  <p>{{ .Content | uppercase }}</p>
{{- end }}
`

func Example() {
	files := http.FS(fstest.MapFS{
		"layout.gohtml":  &fstest.MapFile{Data: []byte(layoutTmpl)},
		"content.gohtml": &fstest.MapFile{Data: []byte(contentTmpl)},
	})

	templates := html.New(files, true, template.FuncMap{"uppercase": strings.ToUpper})

	data := map[string]any{
		"Title":   "wibble",
		"Content": "wobble",
	}

	var out bytes.Buffer
	err := templates.Render(&out, "content.gohtml", data)
	if err != nil {
		panic(err)
	}
	fmt.Println(out.String())

	// Output:
	// <!DOCTYPE html>
	// <html lang="en">
	// <head>
	//   <meta charset="utf-8">
	//   <meta name="viewport" content="width=device-width, initial-scale=1">
	//   <title>wibble</title>
	// </head>
	// <body>
	//   <p>WOBBLE</p>
	// </body>
	// </html>
}
