package book

import (
	v "Library/internal/handler"
	"Library/internal/model"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type createService interface {
	Create(name string, authors ...model.Author) (*model.Book, error)
}

type createRequest struct {
	Name     string   `json:"name" validate:"required"`
	AuthorId []string `json:"author_id" validate:"required"`
}

type createResponse struct {
	v.Response
	Id string `json:"id"`
}

func Create(service createService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createRequest
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

		name := strings.TrimSpace(req.Name)
		authors := make([]model.Author, 0, len(req.AuthorId))
		for _, id := range req.AuthorId {
			authors = append(authors, model.Author{Id: id})
		}

		book, err := service.Create(name, authors...)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, v.Error(err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, createResponse{
			Response: v.Ok(),
			Id:       book.Id,
		})
	}
}
