package author

import (
	v "Library/internal/handler"
	"Library/internal/model"
	"net/http"

	"github.com/go-chi/render"
)

type getResponse struct {
	v.Response
	Id   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

//go:generate mockery --name=GetService
type GetService interface {
	Get(id string) (*model.Author, error)
}

func Get(service GetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, v.Error("id is required"))
			return
		}

		author, err := service.Get(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, v.Error(err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, getResponse{
			Response: v.Ok(),
			Id:       author.Id,
			Name:     author.Name,
		})
	}
}
