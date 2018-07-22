package gotemplate

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"flamingo.me/flamingo/framework/router"
	flamingotemplate "flamingo.me/flamingo/framework/template"
	"flamingo.me/flamingo/framework/web"
	"go.opencensus.io/trace"
)

const pathSeparatorString = string(os.PathSeparator)

type (
	engine struct {
		templatesBasePath  string
		layoutTemplatesDir string
		debug              bool
		templateFunctions  *flamingotemplate.FunctionRegistry
		templates          map[string]*template.Template
	}

	// urlFunc allows templates to access the routers `URL` helper method
	urlFunc struct {
		Router *router.Router `inject:""`
	}

	// getFunc allows templates to access the router's `get` method
	dataFunc struct {
		Router *router.Router `inject:""`
	}

	getFunc struct {
		Router *router.Router `inject:""`
	}
)

var (
	_    flamingotemplate.Function        = new(urlFunc)
	_    flamingotemplate.ContextFunction = new(getFunc)
	_    flamingotemplate.ContextFunction = new(dataFunc)
	lock                                  = &sync.Mutex{}
)

// Inject engine dependencies
func (e *engine) Inject(
	templateFunctions *flamingotemplate.FunctionRegistry,
	config *struct {
		TemplatesBasePath  string `inject:"config:gotemplates.engine.templates.basepath"`
		LayoutTemplatesDir string `inject:"config:gotemplates.engine.layout.dir"`
		Debug              bool   `inject:"config:debug.mode"`
	},
) {
	e.templateFunctions = templateFunctions
	e.templatesBasePath = config.TemplatesBasePath
	e.layoutTemplatesDir = config.LayoutTemplatesDir
	e.debug = config.Debug
}

func (e *engine) Render(ctx context.Context, name string, data interface{}) (io.Reader, error) {
	ctx, span := trace.StartSpan(ctx, "gotemplate/Render")
	defer span.End()

	lock.Lock()
	if e.debug || e.templates == nil {
		e.loadTemplates(ctx)
	}
	lock.Unlock()

	_, span = trace.StartSpan(ctx, "gotemplate/Execute")
	buf := &bytes.Buffer{}
	err := e.templates[name+".html"].Execute(buf, data)
	defer span.End()

	return buf, err
}

func (e *engine) loadTemplates(ctx context.Context) {
	//done := context.Profile("template engine", "load templates")
	ctx, span := trace.StartSpan(ctx, "gotemplate/loadTemplates")
	defer span.End()

	e.templates = make(map[string]*template.Template, 0)

	functionsMap := template.FuncMap{
		"Upper": strings.ToUpper,
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		"map": func(p ...interface{}) map[string]interface{} {
			res := make(map[string]interface{})
			for i := 0; i < len(p); i += 2 {
				res[fmt.Sprint(p[i])] = p[i+1]
			}
			return res
		},
	}

	funcs := e.templateFunctions.Populate()
	for k, f := range e.templateFunctions.ContextAware {
		funcs[k] = f(web.ToContext(ctx))
	}

	layoutTemplate := template.Must(e.parseLayoutTemplates(functionsMap, funcs))

	err := e.parseSiteTemplateDirectory(layoutTemplate, e.templatesBasePath)
	if err != nil {
		panic(err)
	}

	//done()
}

// parses all layout templates in a template instance which is the base instance for all other templates
func (e *engine) parseLayoutTemplates(functionsMap template.FuncMap, funcs template.FuncMap) (*template.Template, error) {
	tpl := template.New("").Funcs(functionsMap).Funcs(funcs)

	if e.layoutTemplatesDir == "" {
		return tpl, nil
	}

	dir := e.templatesBasePath + pathSeparatorString + e.layoutTemplatesDir
	layoutFilesInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	layoutFilesNames := make([]string, 0)
	for _, file := range layoutFilesInfo {
		if file.IsDir() {
			continue
		}
		layoutFilesNames = append(layoutFilesNames, dir+pathSeparatorString+file.Name())
	}

	return tpl.ParseFiles(layoutFilesNames...)
}

// parses all templates from a given directory into a clone of the given layout template, so that all layouts are available
func (e *engine) parseSiteTemplateDirectory(layoutTemplate *template.Template, dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		t := template.Must(layoutTemplate.Clone())
		fullName := dir + pathSeparatorString + f.Name()
		if f.IsDir() {
			err = e.parseSiteTemplateDirectory(layoutTemplate, fullName)
			if err != nil {
				return err
			}
			continue
		}
		tContent, err := ioutil.ReadFile(fullName)
		if err != nil {
			return err
		}

		templateName := strings.TrimPrefix(fullName, e.templatesBasePath+pathSeparatorString)
		e.templates[templateName] = template.Must(t.Parse(string(tContent)))
	}

	return nil
}

// Name alias for use in template
func (g getFunc) Name() string {
	return "get"
}

// Func as implementation of get method
func (g *getFunc) Func(ctx web.Context) interface{} {
	return func(what string, params ...map[string]interface{}) interface{} {
		var p = make(map[interface{}]interface{})
		if len(params) == 1 {
			for k, v := range params[0] {
				p[k] = fmt.Sprint(v)
			}
		}
		return g.Router.Get(what, ctx, p)
	}
}

// Name alias for use in template
func (d dataFunc) Name() string {
	return "data"
}

// Func as implementation of get method
func (d *dataFunc) Func(ctx web.Context) interface{} {
	return func(what string, params ...map[string]interface{}) interface{} {
		var p = make(map[interface{}]interface{})
		if len(params) == 1 {
			for k, v := range params[0] {
				p[k] = fmt.Sprint(v)
			}
		}
		return d.Router.Get(what, ctx, p)
	}
}

// Name alias for use in template
func (u urlFunc) Name() string {
	return "url"
}

// Func as implementation of url method
func (u *urlFunc) Func() interface{} {
	return func(where string, params ...map[string]interface{}) template.URL {
		var p = make(map[string]string)
		if len(params) == 1 {
			for k, v := range params[0] {
				p[k] = fmt.Sprint(v)
			}
		}
		return template.URL(u.Router.URL(where, p).String())
	}
}
