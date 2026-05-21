package handler

import (
	"net/http"

	"github.com/go-chi/render"
)

const idRequiredMessage = "id is required"

func DeleteByID(deleteFunc func(id string) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Error(idRequiredMessage))
			return
		}

		if err := deleteFunc(id); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Error(err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, Ok())
	}
}
