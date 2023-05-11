package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
)

type Template struct {
	htmlTpl *template.Template
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	htmlTpl := template.New(patterns[0])

	htmlTpl.Funcs(template.FuncMap{
		// returns a placeholder CSRF field - this will be replaced when the template is executed
		// using the placeholder allows us to parse the template when the application starts
		"csrfField": func() template.HTML {
			return `<input class="hidden">"placeholder CSRF Field"</input>`
		},
	})

	htmlTpl, err := htmlTpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	return Template{
		htmlTpl: htmlTpl,
	}, nil
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {
	htmlTpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("error cloning template: %v", err)
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
		return
	}

	htmlTpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.htmlTpl.Execute(w, data)
	if err != nil {
		log.Printf("error executing template: %v", err)
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}
}
