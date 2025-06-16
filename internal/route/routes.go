package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	config "github.com/gvillela7/temperature/configs"
	"github.com/gvillela7/temperature/internal/handler"
	"github.com/gvillela7/temperature/util"
	"net/http"
)

func Run() {
	err := config.Load(".")
	if err != nil {
		util.Log(true, false, "error", "failed to initialize environment variables:", "error", err)
		panic(err)
	}
	cfg := config.GetAPIConfig()
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/temperature", handler.GetCep)
	})

	util.Log(true, false, "info", "Server running on", "Port", cfg.Port)
	http.ListenAndServe(":8080", r)
}
