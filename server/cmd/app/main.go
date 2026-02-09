package main

import (
	"Library/internal/config"
	"Library/internal/handler/author"
	"Library/internal/handler/book"
	author2 "Library/internal/repo/author/postgres"
	book2 "Library/internal/repo/book/postgresql"
	serviceauthor "Library/internal/service/author"
	servicebook "Library/internal/service/book"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cnf := config.MustLoad()

	authorRepo, err := author2.NewAuthorRepositoryPostgres(cnf.StorageUrl)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	authorService := serviceauthor.NewAuthorServiceImpl(authorRepo)

	bookRepo, err := book2.NewBookPostgresRepository(cnf.StorageUrl)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
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
		Addr:        cnf.HttpServer.Address,
		ReadTimeout: cnf.HttpServer.Timeout,
		IdleTimeout: cnf.HttpServer.IdleTimeout,
		Handler:     r,
	}

	if err := server.ListenAndServe(); err != nil {
		os.Exit(1)
	}

}
