package utils

import (
	"net/http"
	"strings"
)

type Handler func(http.ResponseWriter, *http.Request)

type Route struct {
	Method  string
	Path    string
	Handler Handler
}

type Router struct {
	prefix string
	routes []Route
}

func (r *Router) Get(path string, handler Handler) {
	r.routes = append(r.routes, Route{"GET", path, handler})
}

// Post adds a POST route to the router
func (r *Router) Post(path string, handler Handler) {
	r.routes = append(r.routes, Route{"POST", path, handler})
}

// Put adds a PUT route to the router
func (r *Router) Put(path string, handler Handler) {
	r.routes = append(r.routes, Route{"PUT", path, handler})
}

// Delete adds a DELETE route to the router
func (r *Router) Delete(path string, handler Handler) {
	r.routes = append(r.routes, Route{"DELETE", path, handler})
}

// ServeHTTP implements the http.Handler interface
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, router := range a.routers {
		if strings.HasPrefix(r.URL.Path, router.prefix) {
			for _, route := range router.routes {
				if r.Method == route.Method && r.URL.Path == router.prefix+route.Path {
					route.Handler(w, r)
					return
				}
			}
		}
	}
	http.NotFound(w, r)
}

func (a *App) NewRouter(prefix string) *Router {
	r := &Router{prefix: prefix}
	a.routers = append(a.routers, r)
	return r
}
