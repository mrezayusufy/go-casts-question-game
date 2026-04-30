package handler

import (
	"net/http"
	"strings"
)

type Router struct {
	mux         *http.ServeMux
	authHandler *Auth
	middleware  []func(http.Handler) http.Handler
}

func NewRouter(authHandler *Auth) *Router {
	return &Router{
		mux:         http.NewServeMux(),
		authHandler: authHandler,
	}
}

func (r *Router) Use(middleware func(http.Handler) http.Handler) {
	r.middleware = append(r.middleware, middleware)
}

func (r *Router) Handle(pattern string, handler http.Handler) {
	r.mux.Handle(pattern, handler)
}

func (r *Router) HandleFunc(pattern string, handler http.HandlerFunc) {
	r.mux.HandleFunc(pattern, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var handler http.Handler = r.mux
	for i := len(r.middleware) - 1; i >= 0; i-- {
		handler = r.middleware[i](handler)
	}
	handler.ServeHTTP(w, req)
}

type SimpleMux struct {
	routes map[string]map[string]http.HandlerFunc
}

func NewSimpleMux() *SimpleMux {
	return &SimpleMux{
		routes: make(map[string]map[string]http.HandlerFunc),
	}
}
func (m *SimpleMux) HandleFunc(method, path string, handler http.HandlerFunc) {
	if m.routes[path] == nil {
		m.routes[path] = make(map[string]http.HandlerFunc)
	}
	m.routes[path][strings.ToUpper(method)] = handler
}

func (m *SimpleMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := strings.ToUpper(r.Method)
	path := r.URL.Path
	if handlers, ok := m.routes[path]; ok {
		if handler, ok := handlers[method]; ok {
			handler(w, r)

			return
		}
		// method not allowed
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "method not allowed"}`))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "not found"}`))
}
