package controllers

import (
	"fmt"
	"net/http"

	"github.com/alorents/lenslocked/context"
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

func (c GalleriesController) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID int
		Title  string
	}
	data.UserID = context.User(r.Context()).ID
	data.Title = r.FormValue("title")

	gallery, err := c.GalleryService.Create(data.Title, data.UserID)
	if err != nil {
		c.Templates.New.Execute(w, r, data, err)
		return
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}
