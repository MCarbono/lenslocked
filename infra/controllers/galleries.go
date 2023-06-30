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
		New  Template
		Edit Template
	}
	*usecases.CreateGalleryUseCase
	*usecases.UpdateGalleryUseCase
	*usecases.FindGalleryUseCase
	*usecases.FindGalleriesUseCase
	*usecases.DeleteGalleryUseCase
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
	// input := struct {
	// 	ID    string
	// 	Title string
	// }{
	// 	ID:    gallery.ID,
	// 	Title: gallery.Title,
	// }
	g.Templates.Edit.Execute(w, r, gallery)
}
