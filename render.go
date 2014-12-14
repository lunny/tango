package tango

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"path"
	"reflect"
)

type T map[string]interface{}

type Renderer struct {
	resp         ResponseWriter
	templateMgr  *TemplateMgr
	logger       Logger
	T            T
	RootTemplate *template.Template
	funcs        template.FuncMap
	beforeRender func()
	afterRender  func()
	action       interface{}
}

func NewRenderer(resp ResponseWriter,
	logger Logger,
	templateMgr *TemplateMgr,
	funcs template.FuncMap,
	action interface{}) *Renderer {
	return &Renderer{
		resp:        resp,
		templateMgr: templateMgr,
		logger:      logger,
		funcs:       funcs,
		action:      action,
	}
}

// add a name value for template
func (r *Renderer) AddTmplVar(name string, varOrFunc interface{}) {
	if varOrFunc == nil {
		r.T[name] = varOrFunc
		return
	}

	if reflect.ValueOf(varOrFunc).Type().Kind() == reflect.Func {
		r.funcs[name] = varOrFunc
	} else {
		r.T[name] = varOrFunc
	}
}

// add names and values for template
func (r *Renderer) AddTmplVars(t *T) {
	for name, value := range *t {
		r.AddTmplVar(name, value)
	}
}

func (r *Renderer) getTemplate(tmpl string) ([]byte, error) {
	if r.templateMgr != nil {
		return r.templateMgr.GetTemplate(tmpl)
	}
	path := getTemplatePath(r.templateMgr.RootDir, tmpl)
	if path == "" {
		return nil, errors.New(fmt.Sprintf("No template file %v found", path))
	}

	return ioutil.ReadFile(path)
}

func (r *Renderer) GetFuncs() template.FuncMap {
	return r.funcs
}

// render the template with vars map, you can have zero or one map
func (r *Renderer) Render(tmpl string, params ...*T) error {
	content, err := r.getTemplate(tmpl)
	if err == nil {
		err = r.NamedRender(tmpl, string(content), params...)
	}
	return err
}

// Include method provide to template for {{include "xx.tmpl"}}
func (r *Renderer) Include(tmplName string) interface{} {
	t := r.RootTemplate.New(tmplName)
	t.Funcs(r.GetFuncs())

	content, err := r.getTemplate(tmplName)
	if err != nil {
		r.logger.Errorf("RenderTemplate %v read err: %s", tmplName, err)
		return ""
	}

	constr := string(content)

	tmpl, err := t.Parse(constr)
	if err != nil {
		r.logger.Errorf("Parse %v err: %v", tmplName, err)
		return ""
	}
	newbytes := bytes.NewBufferString("")
	err = tmpl.Execute(newbytes, r.action)
	if err != nil {
		r.logger.Errorf("Parse %v err: %v", tmplName, err)
		return ""
	}

	tplcontent, err := ioutil.ReadAll(newbytes)
	if err != nil {
		r.logger.Errorf("Parse %v err: %v", tmplName, err)
		return ""
	}
	return template.HTML(string(tplcontent))
}

// render the template with vars map, you can have zero or one map
func (r *Renderer) NamedRender(name, content string, params ...*T) error {
	r.funcs["include"] = r.Include

	if len(params) > 0 {
		r.AddTmplVars(params[0])
	}

	r.RootTemplate = template.New(name)
	r.RootTemplate.Funcs(r.funcs)

	tmpl, err := r.RootTemplate.Parse(content)
	if err == nil {
		newbytes := bytes.NewBufferString("")
		err = tmpl.Execute(newbytes, r.action)
		if err == nil {
			tplcontent, err := ioutil.ReadAll(newbytes)
			if err == nil {
				_, err = r.resp.Write(tplcontent)
			}
		}
	}
	return err
}

func (r *Renderer) RenderString(content string, params ...*T) error {
	h := md5.New()
	h.Write([]byte(content))
	name := h.Sum(nil)
	return r.NamedRender(string(name), content, params...)
}

type RendererInterface interface {
	SetRenderer(render *Renderer)
}

func getTemplatePath(templateDir, name string) string {
	templateFile := path.Join(templateDir, name)
	if fileExists(templateFile) {
		return templateFile
	}
	return ""
}

type Render struct {
	templateDir string
	templateMgr *TemplateMgr
	FuncMaps    template.FuncMap
	VarMaps     T
	logger      Logger
}

func (render *Render) AddTmplVar(name string, varOrFun interface{}) {
	if reflect.TypeOf(varOrFun).Kind() == reflect.Func {
		render.FuncMaps[name] = varOrFun
	} else {
		render.VarMaps[name] = varOrFun
	}
}

func (render *Render) AddTmplVars(t *T) {
	for name, value := range *t {
		render.AddTmplVar(name, value)
	}
}

func (itor *Render) SetLogger(logger Logger) {
	itor.logger = logger
	itor.templateMgr.logger = logger
}

func NewRender(templateDir string,
	reloadTemplates, cacheTemplates bool) *Render {
	itor := &Render{
		templateDir: templateDir,
		templateMgr: new(TemplateMgr),
		FuncMaps:    defaultFuncs,
		VarMaps:     T{},
	}

	itor.templateMgr.Init(
		templateDir,
		reloadTemplates,
	)

	itor.VarMaps["TangoVer"] = Version

	return itor
}

func (render *Render) NewRenderer(rw ResponseWriter, logger Logger, action interface{}) *Renderer {
	renderer := NewRenderer(
		rw,
		render.logger,
		render.templateMgr,
		render.FuncMaps,
		action,
	)
	// copy func from app to renderer
	renderer.T = render.VarMaps
	return renderer
}

func (render *Render) Handle(ctx *Context) {
	if action := ctx.Action(); action != nil {
		if rd, ok := action.(RendererInterface); ok {
			rd.SetRenderer(render.NewRenderer(
				ctx, 
				render.logger, 
				action,
			))
		}
	}

	ctx.Next()
}
