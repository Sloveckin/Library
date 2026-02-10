package book

import (
	v "Library/internal/handler"
	"Library/internal/model"
	"net/http"

	"github.com/go-chi/render"
)

type getRequest struct {
	Id string `json:"id" validate:"required"`
}

type getResponse struct {
	v.Response
	Id      string   `json:"id" validate:"required"`
	Name    string   `json:"name" validate:"required"`
	Authors []string `json:"authors" validate:"required"`
}

type getService interface {
	Get(id string) (*model.Book, error)
}

func Get(service getService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, v.Error("id is required"))
			return
		}

		book, err := service.Get(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, v.Error(err.Error()))
			return
		}

		authors := make([]string, 0, len(book.Authors))
		for _, a := range book.Authors {
			authors = append(authors, a.Id)
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, getResponse{
			Response: v.Ok(),
			Id:       book.Id,
			Name:     book.Name,
			Authors:  authors,
		})
	}
}
