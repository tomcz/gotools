package html

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"sync"
)

// ErrWrite when unable to write the template to the output writer.
var ErrWrite = errors.New("write error")

// Templates renders HTML templates from a filesystem with optional template caching.
// Caching of templates is very useful when using embedded filesystems or when the files
// are not expected to change at all. However, in development environments it is useful
// to not cache templates so that they can be changed by designers at runtime without
// needing to restart the service.
type Templates struct {
	files http.FileSystem
	funcs []template.FuncMap
	cache *sync.Map
}

// New constructor with optional template function maps that are available to every rendered template.
func New(files http.FileSystem, cacheTemplates bool, funcs ...template.FuncMap) *Templates {
	tmpl := &Templates{files: files, funcs: funcs}
	if cacheTemplates {
		tmpl.cache = new(sync.Map)
	}
	return tmpl
}

type renderCfg struct {
	layoutFile    string
	templateName  string
	templateFiles []string
	unbuffered    bool
}

type RenderOpt func(cfg *renderCfg)

// WithoutLayoutFile will not use a layout file at all for this render.
func WithoutLayoutFile() RenderOpt {
	return WithLayoutFile("")
}

// WithLayoutFile overrides the default "layout.gohtml" layout file for this render.
func WithLayoutFile(layoutFile string) RenderOpt {
	return func(cfg *renderCfg) {
		cfg.layoutFile = layoutFile
	}
}

// WithTemplateName overrides the default "main" template name for this render.
func WithTemplateName(templateName string) RenderOpt {
	return func(cfg *renderCfg) {
		cfg.templateName = templateName
	}
}

// WithTemplateFiles adds additional files to be used for this render.
func WithTemplateFiles(files ...string) RenderOpt {
	return func(cfg *renderCfg) {
		cfg.templateFiles = files
	}
}

// WithoutBuffer disables buffering of template execution.
// We buffer template execution output by default to avoid writing incomplete
// or malformed content to the response, but sometimes huge data sets are
// best rendered without any intermediate buffering.
func WithoutBuffer() RenderOpt {
	return func(cfg *renderCfg) {
		cfg.unbuffered = true
	}
}

// Render a template to a [io.Writer]. Template files generally have a common layout
// file ("layout.gohtml") with a common "main" template that renders the given template
// file within the context of the layout to produce well-structrured HTML output.
//
// Callers can use [ErrWrite] to distinguish between template preparation, execution, and
// output errors as this error is only returned when writing to the output writer has failed.
func (t *Templates) Render(w io.Writer, templateFile string, data map[string]any, opts ...RenderOpt) error {
	cfg := &renderCfg{
		layoutFile:   "layout.gohtml",
		templateName: "main",
	}
	for _, opt := range opts {
		opt(cfg)
	}
	if data == nil {
		data = map[string]any{}
	}
	tmpl, err := t.newTemplate(cfg, templateFile)
	if err != nil {
		return err
	}
	if cfg.unbuffered {
		return writeUnbuffered(w, cfg, tmpl, data)
	}
	return writeBuffered(w, cfg, tmpl, data)
}

func (t *Templates) newTemplate(cfg *renderCfg, templateFile string) (*template.Template, error) {
	templatePaths := []string{cfg.layoutFile, templateFile}
	if len(cfg.templateFiles) > 0 {
		templatePaths = append(templatePaths, cfg.templateFiles...)
	}
	var cacheKey string
	if t.cache != nil {
		cacheKey = strings.Join(templatePaths, ",")
		if cached, ok := t.cache.Load(cacheKey); ok {
			return cached.(*template.Template), nil
		}
	}
	tmpl := template.New("")
	for _, fm := range t.funcs {
		tmpl = tmpl.Funcs(fm)
	}
	for _, path := range templatePaths {
		if path == "" {
			continue
		}
		buf, err := t.readTemplate(path)
		if err != nil {
			return nil, err
		}
		tmpl, err = tmpl.Parse(string(buf))
		if err != nil {
			return nil, fmt.Errorf("%s parse failed: %w", path, err)
		}
	}
	if t.cache != nil {
		t.cache.Store(cacheKey, tmpl)
	}
	return tmpl, nil
}

func (t *Templates) readTemplate(path string) ([]byte, error) {
	in, err := t.files.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%s open failed: %w", path, err)
	}
	defer in.Close()

	buf, err := io.ReadAll(in)
	if err != nil {
		return nil, fmt.Errorf("%s read failed: %w", path, err)
	}
	return buf, nil
}

func writeUnbuffered(w io.Writer, cfg *renderCfg, tmpl *template.Template, data map[string]any) error {
	err := tmpl.ExecuteTemplate(w, cfg.templateName, data)
	if err != nil {
		return fmt.Errorf("unbuffered %w: %w", ErrWrite, err)
	}
	return nil
}

func writeBuffered(w io.Writer, cfg *renderCfg, tmpl *template.Template, data map[string]any) error {
	buf := bufBorrow()
	defer bufReturn(buf)

	err := tmpl.ExecuteTemplate(buf, cfg.templateName, data)
	if err != nil {
		return fmt.Errorf("template execute: %w", err)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		return fmt.Errorf("buffered %w: %w", ErrWrite, err)
	}
	return nil
}
