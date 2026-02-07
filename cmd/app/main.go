package main

import (
	"Library/internal/handler/book"
	book2 "Library/internal/repo/book/memory"
	servicebook "Library/internal/service/book"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	repo := book2.NewRepositoryInMemory()
	bookService := servicebook.NewServiceBook(repo)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Put("/create", book.Create(bookService))
	r.Get("/get", book.Get(bookService))
	r.Delete("/delete", book.Delete(bookService))

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		os.Exit(1)
	}

}
