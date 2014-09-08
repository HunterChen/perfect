package perfect

import (
	"errors"
	"github.com/vpetrov/perfect/orm"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	TEMPLATE_DIR = "templates"
	TEMPLATE_EXT = ".html"
)

//A Perfect Module is a standalone component that can be mounted on
//URL paths and can decide how requests are routed.
type Module struct {
	Mux
	Name           string
	MountPoint     string
	Path           string
	SessionTimeout time.Duration

	Db  orm.Database
	Log *log.Logger

	Templates *template.Template
}

func (m *Module) abs(p string) string {
	return m.MountPoint + p
}

func (m *Module) asset(p string) string {
	return m.MountPoint + m.StaticPrefix() + p
}

func (m *Module) _string(i int) string {
	return strconv.Itoa(i)
}

//parses all template files from the 'templates' folder of the module
func (m *Module) ParseTemplates() error {

	log.Println("Parsing templates from", m.Path)

	m.Templates = template.New(m.Name)
	//set start/end tags (delimiters)
	_ = m.Templates.Delims("<%", "%>")

	pathlen := len(m.Path)
	tpldirlen := len(TEMPLATE_DIR)
	tplextlen := len(TEMPLATE_EXT)

	moduleFuncs := map[string]interface{}{
		"abs":    m.abs,
		"asset":  m.asset,
		"string": m._string,
	}

	tplParser := func(currentPath string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(currentPath) == TEMPLATE_EXT {
			//the template name is anything after 'path/templates/'
			// '2' is for the 2 slashes in the path above
			from := pathlen + tpldirlen + 2
			to := len(currentPath) - tplextlen
			requestPath := currentPath[from:to]

			//read the template file
			data, err := ioutil.ReadFile(currentPath)
			if err != nil {
				return err
			}

			t, err := m.Templates.New(requestPath).Parse(string(data))
			if err != nil {
				return err
			}

			t.Funcs(moduleFuncs)
		}

		return nil
	}

	err := filepath.Walk(m.Path, tplParser)

	return err
}

// renders a template file
func (m *Module) RenderTemplate(w http.ResponseWriter, r *Request, path string, data interface{}) {
	tpl := m.Templates.Lookup(path)
	if tpl == nil {
		Error(w, r, errors.New("Template not found: "+path))
		return
	}

	err := tpl.Execute(w, data)
	if err != nil {
		LogError(r, err)
		return
	}
}
