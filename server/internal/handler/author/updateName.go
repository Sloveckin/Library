package author

import (
	"fmt"
	"net/http"
	"strings"

	v "Library/internal/handler"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type UpdateService interface {
	UpdateName(id, name string) error
}

type updateNameRequest struct {
	Name string
}

func UpdateName(service UpdateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, v.Error("id is required"))
			return
		}

		var req updateNameRequest
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

		err = service.UpdateName(id, name)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, v.Error(err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, v.Ok())
	}
}
