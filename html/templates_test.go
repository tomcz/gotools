package html

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"
	"testing"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
	"gotest.tools/v3/assert"
)

func TestRender(t *testing.T) {
	files := http.Dir("testdata")
	funcMap := template.FuncMap{"uppercase": strings.ToUpper}
	tmpl := New(files, true, funcMap)

	data := map[string]any{
		"Title":   "wibble",
		"Content": "wobble",
	}

	var out bytes.Buffer
	err := tmpl.Render(&out, "content.gohtml", data)
	assert.NilError(t, err)

	doc, err := html.Parse(strings.NewReader(out.String()))
	assert.NilError(t, err)

	sel, err := cascadia.Parse("title")
	assert.NilError(t, err)
	node := cascadia.Query(doc, sel)
	assert.Equal(t, node.FirstChild.Data, "wibble")

	sel, err = cascadia.Parse("p")
	assert.NilError(t, err)
	node = cascadia.Query(doc, sel)
	assert.Equal(t, node.FirstChild.Data, "WOBBLE")
}
