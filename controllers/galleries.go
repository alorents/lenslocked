package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/alorents/lenslocked/context"
	"github.com/alorents/lenslocked/errors"
	"github.com/alorents/lenslocked/models"
)

type GalleriesController struct {
	Templates struct {
		New  Template
		Edit Template
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

func (c GalleriesController) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusNotFound)
		return
	}
	gallery, err := c.GalleryService.ByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Gallery not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You do not have permissions to edit this gallery", http.StatusForbidden)
		return
	}

	var data struct {
		ID    int
		Title string
	}
	data.ID = gallery.ID
	data.Title = gallery.Title
	c.Templates.Edit.Execute(w, r, data)
}
