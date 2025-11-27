package handlers

import (
	"bannerService/internals/middleware"
	"net/http"
)

type RouteGroup struct {
	mux         *http.ServeMux
	middlewares []middleware.Middleware
}

func (rg *RouteGroup) With(middlewares ...middleware.Middleware) *RouteGroup {
	return &RouteGroup{
		mux:         rg.mux,
		middlewares: append(rg.middlewares, middlewares...),
	}
}

func (rg *RouteGroup) Handle(pattern string, handler http.Handler) {
	wrapped := handler
	for i := len(rg.middlewares) - 1; i >= 0; i-- {
		wrapped = rg.middlewares[i](wrapped)
	}
	rg.mux.Handle(pattern, wrapped)
}

func (rg *RouteGroup) HandleFunc(pattern string, handler http.HandlerFunc) {
	rg.Handle(pattern, handler)
}
