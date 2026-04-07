package book

import (
	v "Library/internal/handler"
	"Library/internal/model"
	"log"
	"net/http"

	"github.com/go-chi/render"
)

type allResponse struct {
	v.Response
	Books []model.Book `json:"books"`
}

type allService interface {
	GetAll() ([]model.Book, error)
}

func All(service allService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Handling GET /book/all, URL: %s", r.URL.String())

		books, err := service.GetAll()
		if err != nil {
			log.Println("Error while fetching all books:", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, v.Error("Failed to fetch books"))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, allResponse{
			Response: v.Ok(),
			Books:    books,
		})
	}
}
