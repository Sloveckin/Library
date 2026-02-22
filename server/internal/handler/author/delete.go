package author

import (
	v "Library/internal/handler"
	"net/http"

	"github.com/go-chi/render"
)

type deleteService interface {
	Delete(id string) error
}

type deleteResponse struct {
	v.Response
}

func Delete(service deleteService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, v.Error("id is required"))
			return
		}

		err := service.Delete(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, v.Error(err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, deleteResponse{
			Response: v.Ok(),
		})
	}
}
