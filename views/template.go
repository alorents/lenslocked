package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Template struct {
	HtmlTpl *template.Template
}

func Parse(filepath string) (Template, error) {
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		return Template{}, fmt.Errorf("error parsing template: %w", err)
	}
	return Template{
		HtmlTpl: tpl,
	}, nil
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.HtmlTpl.Execute(w, nil)
	if err != nil {
		log.Printf("error executing template: %v", err)
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}
}
