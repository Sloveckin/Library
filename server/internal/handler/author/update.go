package author

import (
	v "Library/internal/handler"
	"Library/internal/model"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

//go:generate mockery --name=UpdateService
type UpdateService interface {
	Update(id, name string) (*model.Author, error)
}

type updateRequest struct {
	Id   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type updateResponse struct {
	v.Response
	Id string `json:"id"`
}

func Update(service UpdateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req updateRequest
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
		author, err := service.Update(req.Id, name)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, v.Error(err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, updateResponse{
			Response: v.Ok(),
			Id:       author.Id,
		})
	}
}
