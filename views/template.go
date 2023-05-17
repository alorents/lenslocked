package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/csrf"

	"github.com/alorents/lenslocked/context"
	"github.com/alorents/lenslocked/models"
)

type Template struct {
	htmlTpl *template.Template
}

// publicError cab be used to determine if an error provides the Public method.
type publicError interface {
	Public() string
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	htmlTpl := template.New(path.Base(patterns[0]))

	htmlTpl.Funcs(template.FuncMap{
		// returns a placeholder CSRF field - this will be replaced when the template is executed
		// using the placeholder allows us to parse the template when the application starts
		"csrfField": func() (template.HTML, error) {
			return "", fmt.Errorf("csrfField called but not implemented")
		},
		"currentUser": func() (template.HTML, error) {
			return "", fmt.Errorf("currentUser called but not implemented")
		},
		"errors": func() []string {
			return nil
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

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	htmlTpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("error cloning template: %v", err)
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
		return
	}

	// Call the errMessages func before the closures.
	errMsgs := errMessages(errs...)
	htmlTpl = htmlTpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
			"errors": func() []string {
				// return the pre-processed err messages inside the closure.
				return errMsgs
			},
		},
	)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer
	err = htmlTpl.Execute(&buf, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(w, &buf)
	if err != nil {
		panic(err)
	}
}

func errMessages(errs ...error) []string {
	var msgs []string
	for _, err := range errs {
		var pubErr publicError
		if errors.As(err, &pubErr) {
			msgs = append(msgs, pubErr.Public())
		} else {
			fmt.Println(err)
			msgs = append(msgs, "Something went wrong.")
		}
	}
	return msgs
}
