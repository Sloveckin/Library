package book

import (
	v "Library/internal/handler"
	"Library/internal/model"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type getRequest struct {
	Name string `json:"name" validate:"required"`
}

type getResponse struct {
	v.Response
	Id      string   `json:"id" validate:"required"`
	Name    string   `json:"name" validate:"required"`
	Authors []string `json:"authors" validate:"required"`
}

type getService interface {
	Get(name string) (*model.Book, error)
}

func Get(service getService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req getRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			fmt.Println("Error while decoding request:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		validate := validator.New()
		err = validate.Struct(req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, v.Error(err.Error()))
			return
		}

		book, err := service.Get(req.Name)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, v.Error(err.Error()))
			return
		}

		authors := make([]string, 0, len(book.Authors))
		for _, a := range book.Authors {
			authors = append(authors, a.Name)
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
