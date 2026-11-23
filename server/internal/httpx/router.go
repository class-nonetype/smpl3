package httpx

import (
	"net/http"
	"slices"
)

type Middleware func(http.Handler) http.Handler

type Router struct {
	mux        *http.ServeMux
	prefix     string
	middleware []Middleware
	routes     *[]string
}

func NewRouter() *Router {
	return &Router{mux: http.NewServeMux(), routes: &[]string{}}
}

func (rt *Router) Group(prefix string, mw ...Middleware) *Router {
	return &Router{
		mux:        rt.mux,
		prefix:     rt.prefix + prefix,
		middleware: slices.Concat(rt.middleware, mw),
		routes:     rt.routes,
	}
}

func (rt *Router) Handle(method, path string, h HandlerFunc) {
	var handler http.Handler = h
	for i := len(rt.middleware) - 1; i >= 0; i-- {
		handler = rt.middleware[i](handler)
	}
	rt.mux.Handle(method+" "+rt.prefix+path, handler)
	*rt.routes = append(*rt.routes, method+" "+rt.prefix+path)
}

func (rt *Router) Routes() []string {
	return *rt.routes
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.mux.ServeHTTP(w, r)
}
