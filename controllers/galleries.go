package controllers

import (
	"net/http"

	"github.com/alorents/lenslocked/models"
)

type GalleriesController struct {
	Templates struct {
		New Template
	}
	GalleryService *models.GalleryService
}

func (c GalleriesController) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	c.Templates.New.Execute(w, r, data)
}
