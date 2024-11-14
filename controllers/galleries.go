package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/alorents/lenslocked/context"
	"github.com/alorents/lenslocked/errors"
	"github.com/alorents/lenslocked/models"
)

type GalleriesController struct {
	Templates struct {
		New   Template
		Edit  Template
		Index Template
		Show  Template
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

	type Image struct {
		GalleryID       int
		Filename        string
		FilenameEscaped string
	}
	var data struct {
		ID     int
		Title  string
		Images []Image
	}
	data.ID = gallery.ID
	data.Title = gallery.Title
	images, err := c.GalleryService.Images(gallery.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryID:       image.GalleryID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
	}
	c.Templates.Edit.Execute(w, r, data)
}

func (c GalleriesController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusNotFound)
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
		http.Error(w, "You are not authorized to edit this gallery", http.StatusForbidden)
		return
	}

	title := r.FormValue("title")
	gallery.Title = title
	err = c.GalleryService.Update(gallery)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (c GalleriesController) Index(w http.ResponseWriter, r *http.Request) {
	type Gallery struct {
		ID    int
		Title string
	}
	var data struct {
		Galleries []Gallery
	}
	user := context.User(r.Context())
	galleries, err := c.GalleryService.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	for _, gallery := range galleries {
		data.Galleries = append(data.Galleries, Gallery{
			ID:    gallery.ID,
			Title: gallery.Title,
		})
	}

	c.Templates.Index.Execute(w, r, data)
}

func (c GalleriesController) Show(w http.ResponseWriter, r *http.Request) {
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
	type Image struct {
		GalleryID       int
		Filename        string
		FilenameEscaped string
	}
	var data struct {
		ID     int
		Title  string
		Images []Image
	}
	data.ID = gallery.ID
	data.Title = gallery.Title
	images, err := c.GalleryService.Images(gallery.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryID:       image.GalleryID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
	}
	c.Templates.Show.Execute(w, r, data)
}

func (c GalleriesController) Delete(w http.ResponseWriter, r *http.Request) {
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

	err = c.GalleryService.Delete(id)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (c GalleriesController) Image(w http.ResponseWriter, r *http.Request) {
	filename := c.filename(r)
	galleryID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return
	}
	image, err := c.GalleryService.Image(galleryID, filename)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Image not found", http.StatusNotFound)
			return
		}
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, image.Path)
}

func (c GalleriesController) UploadImage(w http.ResponseWriter, r *http.Request) {
	gallery, err := c.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	err = r.ParseMultipartForm(5 << 20) // 5MB
	if err != nil {
		return
	}
	fileHeaders := r.MultipartForm.File["images"]
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		err = c.GalleryService.CreateImage(gallery.ID, fileHeader.Filename, file)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
		editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
		http.Redirect(w, r, editPath, http.StatusFound)
	}
}

func (c GalleriesController) DeleteImage(w http.ResponseWriter, r *http.Request) {
	filename := c.filename(r)
	gallery, err := c.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	err = c.GalleryService.DeleteImage(gallery.ID, filename)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

type galleryOpt func(http.ResponseWriter, *http.Request, *models.Gallery) error

func (c GalleriesController) galleryByID(w http.ResponseWriter, r *http.Request, opts ...galleryOpt) (*models.Gallery, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return nil, err
	}
	gallery, err := c.GalleryService.ByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Gallery not found", http.StatusNotFound)
			return nil, err
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return nil, err
	}
	for _, opt := range opts {
		err = opt(w, r, gallery)
		if err != nil {
			return nil, err
		}
	}
	return gallery, nil
}

func userMustOwnGallery(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You are not authorized to edit this gallery", http.StatusForbidden)
		return fmt.Errorf("user does not hav access to this gallery")
	}
	return nil
}

func (c GalleriesController) filename(r *http.Request) string {
	filename := chi.URLParam(r, "filename")
	filename = filepath.Base(filename)
	return filename
}
