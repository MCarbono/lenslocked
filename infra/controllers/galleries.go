package controllers

import (
	"fmt"
	"lenslocked/application/usecases"
	"lenslocked/context"
	"net/http"

	"github.com/go-chi/chi"
)

type Galleries struct {
	Templates struct {
		Show  Template
		New   Template
		Edit  Template
		Index Template
	}
	*usecases.CreateGalleryUseCase
	*usecases.UpdateGalleryUseCase
	*usecases.FindGalleryUseCase
	*usecases.FindGalleriesUseCase
	*usecases.DeleteGalleryUseCase
	*usecases.FindImageUseCase
	*usecases.DeleteImageUseCase
	*usecases.CreateImageUsecase
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	g.Templates.New.Execute(w, r, data)
}

func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var input = &usecases.CreateGalleryInput{
		UserID: context.User(r.Context()).ID,
		Title:  r.FormValue("title"),
	}
	gallery, err := g.CreateGalleryUseCase.Execute(input)
	if err != nil {
		g.Templates.New.Execute(w, r, input, err)
		return
	}
	editPath := fmt.Sprintf("/galleries/%s/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	gallery, err := g.FindGalleryUseCase.Execute(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	g.Templates.Edit.Execute(w, r, gallery)
}

func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	title := r.FormValue("title")
	err := g.UpdateGalleryUseCase.Execute(&usecases.UpdateGalleryInput{ID: id, Title: title})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	editPath := fmt.Sprintf("/galleries/%s/edit", id)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) Index(w http.ResponseWriter, r *http.Request) {
	type Gallery struct {
		ID    string
		Title string
	}
	var data struct {
		Galleries []Gallery
	}

	user := context.User(r.Context())
	galleries, err := g.FindGalleriesUseCase.Execute(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dto := make([]Gallery, len(galleries))
	for i := range dto {
		dto[i] = Gallery{
			ID:    galleries[i].ID,
			Title: galleries[i].Title,
		}
	}
	data.Galleries = dto
	g.Templates.Index.Execute(w, r, data)
}

func (g Galleries) Show(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	gallery, err := g.FindGalleryUseCase.Execute(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	g.Templates.Show.Execute(w, r, gallery)
}

func (g Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := g.DeleteGalleryUseCase.Execute(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (g Galleries) Image(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	galleryID := chi.URLParam(r, "id")
	image, err := g.FindImageUseCase.Execute(galleryID, filename)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, image.Path)
}

func (g Galleries) DeleteImage(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	id := chi.URLParam(r, "id")
	err := g.DeleteImageUseCase.Execute(id, filename)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	editPath := fmt.Sprintf("/galleries/%s/edit", id)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) UploadImage(w http.ResponseWriter, r *http.Request) {
	//validate the userID and the gallery ID - if the user own it
	galleryID := chi.URLParam(r, "id")
	err := r.ParseMultipartForm(5 << 20)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
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
		err = g.CreateImageUsecase.Execute(&usecases.CreateImageInput{GalleryID: galleryID, Filename: fileHeader.Filename, Contents: file})
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		editPath := fmt.Sprintf("/galleries/%s/edit", galleryID)
		http.Redirect(w, r, editPath, http.StatusFound)
	}
}
