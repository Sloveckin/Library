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

//go:generate mockery --name=CreateService
type CreateService interface {
	Create(name string) (*model.Author, error)
}

type createRequest struct {
	Name string `json:"name" validate:"required"`
}

type CreateResponse struct {
	v.Response
	Id string `json:"id"`
}

func Create(service CreateService) http.HandlerFunc {
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
		book, err := service.Create(name)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, v.Error(err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, CreateResponse{
			Response: v.Ok(),
			Id:       book.Id,
		})
	}
}
