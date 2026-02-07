package author

import (
	v "Library/internal/handler"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type deleteService interface {
	Delete(id string) error
}

type deleteRequest struct {
	Id string `json:"id" validate:"required"`
}

type deleteResponse struct {
	v.Response
}

func Delete(service deleteService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req deleteRequest
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

		err = service.Delete(req.Id)
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
