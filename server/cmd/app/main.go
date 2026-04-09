package main

import (
	"Library/internal/config"
	"Library/internal/handler/author"
	"Library/internal/handler/book"
	authorRepo "Library/internal/repo/author/postgres"
	bookRepo "Library/internal/repo/book/postgresql"
	serviceauthor "Library/internal/service/author"
	servicebook "Library/internal/service/book"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ================= METRICS =================

// Счётчик HTTP запросов
var httpRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"method", "endpoint", "status"},
)

// Гистограмма latency
var httpDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request latency in seconds",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method", "endpoint"},
)

// ================= MIDDLEWARE =================

// Обёртка для захвата status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Метрики middleware
func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := &responseWriter{
			ResponseWriter: w,
			statusCode:     200,
		}

		next.ServeHTTP(ww, r)

		duration := time.Since(start).Seconds()

		// (без query параметров)
		routePattern := chi.RouteContext(r.Context()).RoutePattern()
		if routePattern == "" {
			routePattern = r.URL.Path
		}

		httpRequests.WithLabelValues(
			r.Method,
			routePattern,
			strconv.Itoa(ww.statusCode),
		).Inc()

		httpDuration.WithLabelValues(
			r.Method,
			routePattern,
		).Observe(duration)
	})
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ================= MAIN =================

func main() {
	cnf := config.MustLoad()

	// Регистрируем метрики
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(httpDuration)

	// ===== REPOSITORIES =====
	authorRepository, err := authorRepo.NewAuthorRepositoryPostgres(cnf.StorageUrl)
	if err != nil {
		fmt.Println("Author repo error:", err)
		return
	}
	authorService := serviceauthor.NewAuthorServiceImpl(authorRepository)

	bookRepository, err := bookRepo.NewBookPostgresRepository(cnf.StorageUrl)
	if err != nil {
		fmt.Println("Book repo error:", err)
		return
	}
	bookService := servicebook.NewServiceBook(bookRepository, authorService)

	// ===== ROUTER =====
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(corsMiddleware)
	r.Use(metricsMiddleware)

	// endpoint для Prometheus
	r.Handle("/metrics", promhttp.Handler())

	// ===== ROUTES =====

	r.Route("/book", func(r chi.Router) {
		r.Put("/create", book.Create(bookService))
		r.Put("/update", book.Update(bookService))
		r.Get("/get", book.Get(bookService))
		r.Delete("/delete", book.Delete(bookService))
	})

	r.Route("/author", func(r chi.Router) {
		r.Put("/create", author.Create(authorService))
		r.Put("/update", author.Update(authorService))
		r.Get("/get", author.Get(authorService))
		r.Delete("/delete", author.Delete(authorService))
	})

	// ===== SERVER =====

	server := &http.Server{
		Addr:         cnf.HttpServer.Address,
		ReadTimeout:  cnf.HttpServer.Timeout,
		WriteTimeout: cnf.HttpServer.Timeout,
		IdleTimeout:  cnf.HttpServer.IdleTimeout,
		Handler:      r,
	}

	fmt.Printf("Server started at http://%s\n", cnf.HttpServer.Address)
	fmt.Printf("Metrics available at http://%s/metrics\n", cnf.HttpServer.Address)

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server error:", err)
		os.Exit(1)
	}
}