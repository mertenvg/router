package router

import (
	"net/http"
	"strings"
	"sync"
)

const (
	Default = "/*"
	Error   = "/!"
)

// New creates a new Map that can be used as a http.Handler
func New() *Map {
	m := &Map{
		mutex:    sync.RWMutex{},
		handlers: map[string]http.HandlerFunc{},
	}
	m.handler = m.serveHTTP
	return m
}

type Middleware func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

// Map of paths to http.HandlerFunc
type Map struct {
	pathPrefix string
	mutex      sync.RWMutex
	handlers   map[string]http.HandlerFunc
	handler    http.HandlerFunc
}

// Middleware adds a middleware handler to all requests for the Map
// using LIFO (Last In First Out) order.
func (m *Map) Middleware(middleware ...Middleware) *Map {
	m.mutex.Lock()
	for _, mw := range middleware {
		m.handler = wrap(mw, m.handler)
	}
	m.mutex.Unlock()
	return m
}

// Handle maps the given path to the given http.HandlerFunc
func (m *Map) Handle(path string, h http.Handler) *Map {
	m.mutex.Lock()
	m.handlers[path] = h.ServeHTTP
	m.mutex.Unlock()
	return m
}

// HandleFunc maps the given path to the given http.HandlerFunc
func (m *Map) HandleFunc(path string, hf http.HandlerFunc) *Map {
	m.mutex.Lock()
	m.handlers[path] = hf
	m.mutex.Unlock()
	return m
}

// Default sets the catch all handler to the given http.Handler
func (m *Map) Default(h http.Handler) *Map {
	return m.Handle(Default, h)
}

// DefaultFunc sets the catch all handler to the given http.HandlerFunc
func (m *Map) DefaultFunc(hf http.HandlerFunc) *Map {
	return m.HandleFunc(Default, hf)
}

// Error sets the error handler to the given http.Handler
func (m *Map) Error(h http.Handler) *Map {
	return m.Handle(Error, h)
}

// ErrorFunc sets the error handler to the given http.HandlerFunc
func (m *Map) ErrorFunc(hf http.HandlerFunc) *Map {
	return m.HandleFunc(Error, hf)
}

// ServeHTTP implements http.Handler
func (m *Map) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.handler(w, r)
}

// serveHTTP handles the request
func (m *Map) serveHTTP(w http.ResponseWriter, r *http.Request) {
	// iterate the possible matches
	for _, k := range matches(strings.TrimPrefix(r.URL.Path, m.pathPrefix)) {
		m.mutex.RLock()
		hf, ok := m.handlers[k]
		m.mutex.RUnlock()
		if ok {
			hf(w, r)
			return
		}
	}
	// if no handlers were triggered then default to a 404
	w.WriteHeader(http.StatusNotFound)
}

// Sub returns a new router that will be triggered for the given path
func (m *Map) Sub(path string) *Map {
	sub := New()
	sub.pathPrefix = path
	m.Handle(path+"/*", sub)
	return sub
}

// matches splits a given path into possible route matches
func matches(path string) []string {
	p := strings.Split(path, "/")
	c := len(p)
	// precreating the array and setting the values is about 50% faster than using append.
	result := make([]string, c+1)
	// Set the first match to the
	result[0] = path
	// The ones in between are wildcard matches
	for l := c - 1; l > 0; l-- {
		result[c-l] = strings.Join(p[:l], "/") + "/*"
	}
	// set the last match to the error
	result[c] = Error
	return result
}

// wrap a middlewareFunc into a handleFunc
func wrap(mw Middleware, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mw(w, r, next)
	}
}
