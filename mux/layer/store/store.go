package store

import (
	"net/http"
	"sort"

	"github.com/alexpts/go-mux-http/mux/layer"
)

type LayersStore struct {
	normalizer layer.INormalizer
	layers     []layer.Layer
	sorted     bool
}

func New(normalizer layer.INormalizer) *LayersStore {
	if normalizer == nil {
		normalizer = &layer.StdNormalizer{
			RegExpMaker: &layer.StdRegExpMaker{},
		}
	}

	return &LayersStore{
		normalizer: normalizer,
	}
}

func (s *LayersStore) AddLayer(layer layer.Layer) *LayersStore {
	nLayer := s.normalizer.Normalize(layer)
	s.layers = append(s.layers, nLayer)
	s.sorted = false

	return s
}

func (s *LayersStore) Use(layer layer.Layer, handlers ...http.HandlerFunc) *LayersStore {
	return s.AddLayer(
		layer.WithHandlers(handlers...),
	)
}

func (s *LayersStore) HandleFunc(path string, layer layer.Layer, handlers ...http.HandlerFunc) *LayersStore {
	layer.Path = path
	return s.Use(layer, handlers...)
}

func (s *LayersStore) GetLayers() []layer.Layer {
	if !s.sorted {
		s.sortByPriority()
	}

	return s.layers
}

func (s *LayersStore) sortByPriority() {
	sort.SliceStable(s.layers, func(i, j int) bool {
		return s.layers[i].Priority > s.layers[j].Priority
	})

	s.sorted = true
}
