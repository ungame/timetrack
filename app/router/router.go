package router

import (
	muxHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ungame/timetrack/app/handlers"
	"github.com/ungame/timetrack/app/middlewares"
	"net/http"
)

type Router struct {
	router *mux.Router
}

func New(handlers ...handlers.Handler) *Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Use(middlewares.Logger)
	router.Handle("/metrics", promhttp.Handler())

	for _, handler := range handlers {
		handler.Register(router)
	}

	return &Router{router: router}
}

func (r *Router) Handler() http.Handler {
	return r.router
}

func (r *Router) WithCors() http.Handler {
	methods := muxHandlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete})
	headers := muxHandlers.AllowedHeaders([]string{"Content-Type", "X-Requested-with", "Accept"})
	origins := muxHandlers.AllowedOrigins([]string{"*"})
	return muxHandlers.CORS(methods, headers, origins)(r.router)
}
