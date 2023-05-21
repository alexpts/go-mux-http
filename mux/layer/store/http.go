package store

import (
	"net/http"

	"github.com/alexpts/go-mux-http/mux/layer"
)

type Layer = layer.Layer

func (s *LayersStore) Method(method string, path string, layer Layer, handlers ...http.HandlerFunc) *LayersStore {
	layer.Path = path
	layer.Methods = []string{method}

	return s.Use(layer, handlers...)
}

func (s *LayersStore) Get(path string, options Layer, handlers ...http.HandlerFunc) *LayersStore {
	return s.Method(`GET`, path, options, handlers...)
}

func (s *LayersStore) Post(path string, options Layer, handlers ...http.HandlerFunc) *LayersStore {
	return s.Method(`POST`, path, options, handlers...)
}

func (s *LayersStore) Put(path string, options Layer, handlers ...http.HandlerFunc) *LayersStore {
	return s.Method(`PUT`, path, options, handlers...)
}

func (s *LayersStore) Patch(path string, options Layer, handlers ...http.HandlerFunc) *LayersStore {
	return s.Method(`PATCH`, path, options, handlers...)
}

func (s *LayersStore) Delete(path string, options Layer, handlers ...http.HandlerFunc) *LayersStore {
	return s.Method(`DELETE`, path, options, handlers...)
}
