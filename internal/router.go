package metrics

// import (
// 	// "net/http"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-chi/chi/v5/middleware"
// )

// func NewRouter() chi.Router {
// 	r := chi.NewRouter()

// 	r.Use(middleware.RequestID)
// 	r.Use(middleware.RealIP)
// 	r.Use(middleware.Logger)
// 	r.Use(middleware.Recoverer)

// 	r.Route("/update", func(r chi.Router) {
// 		r.Get("/{type}{name}", GetMetricHandler)
// 	})
// 	return r
// }
