package main

import (
	"Library/internal/handler/author"
	"Library/internal/handler/book"
	author2 "Library/internal/repo/author/memory"
	book2 "Library/internal/repo/book/memory"
	serviceauthor "Library/internal/service/author"
	servicebook "Library/internal/service/book"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	authorRepo := author2.NewAuthorRepositoryInMemory()
	authorService := serviceauthor.NewAuthorServiceImpl(authorRepo)

	bookRepo := book2.NewRepositoryInMemory()
	bookService := servicebook.NewServiceBook(bookRepo, authorService)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Route("/book", func(r chi.Router) {
		r.Put("/create", book.Create(bookService))
		r.Get("/get", book.Get(bookService))
		r.Delete("/delete", book.Delete(bookService))
	})

	r.Route("/author", func(r chi.Router) {
		r.Put("/create", author.Create(authorService))
		r.Get("/get", author.Get(authorService))
		r.Delete("/delete", author.Delete(authorService))
	})

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		os.Exit(1)
	}

}
