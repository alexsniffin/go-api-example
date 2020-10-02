package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	httpMetrics "github.com/slok/go-http-metrics/metrics/prometheus"
	httpMiddleware "github.com/slok/go-http-metrics/middleware"
	nm "github.com/slok/go-http-metrics/middleware/negroni"
	"github.com/urfave/negroni"

	lHandler "github.com/alexsniffin/go-api-starter/internal/todo-api/handlers/logging"
	"github.com/alexsniffin/go-api-starter/internal/todo-api/handlers/todo"
	"github.com/alexsniffin/go-api-starter/internal/todo-api/models"
)

// Creates Chi based multiplexer router with middleware
func NewRouter(cfg models.HTTPRouterConfig, logger zerolog.Logger, todoHandler todo.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(lHandler.NewHandlerFunc(logger))
	r.Use(middleware.Timeout(time.Duration(cfg.TimeoutSec) * time.Second))

	httpMw := httpMiddleware.New(httpMiddleware.Config{
		DisableMeasureInflight: true,
		Recorder:               httpMetrics.NewRecorder(httpMetrics.Config{}),
	})

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   cfg.AllowedMethods,
		AllowedHeaders:   cfg.AllowedHeaders,
		AllowCredentials: false,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Route("/todo", func(r chi.Router) {
			r.Route("/{id}", func(r chi.Router) {
				idMetricHandler := nm.Handler("/api/todo/{id}", httpMw)
				r.Get("/", negroni.New(idMetricHandler, negroni.WrapFunc(todoHandler.Get)).ServeHTTP)
				r.Delete("/", negroni.New(idMetricHandler, negroni.WrapFunc(todoHandler.Delete)).ServeHTTP)
			})
			r.Post("/", negroni.New(nm.Handler("/api/todo", httpMw), negroni.WrapFunc(todoHandler.Post)).ServeHTTP)
		})
		r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	r.Route("/metrics", func(r chi.Router) {
		r.Get("/", promhttp.Handler().ServeHTTP)
	})
	return r
}
