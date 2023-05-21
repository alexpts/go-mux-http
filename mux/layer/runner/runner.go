package runner

import (
	"fmt"
	"net/http"

	"github.com/alexpts/go-mux-http/mux/layer"
	"github.com/alexpts/go-mux-http/mux/layer/resolver"
)

type Runner struct {
	Resolver resolver.IResolver
	Layers   []layer.Layer

	curLayer   *layer.Layer
	layerPos   int
	handlerPos int

	UriParams  map[string]string
	UserParams map[string]any
}

func (r *Runner) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handlerFunc := r.getNextHandler(req)
	handlerFunc(w, req)
}

func (r *Runner) getNextHandler(request *http.Request) http.HandlerFunc {
	if r.curLayer == nil {
		r.curLayer = r.getNextLayer(request)
	}

	// Если обработчики в слое закончились, то идем в следующий слой
	if r.handlerPos == len(r.curLayer.Handlers) {
		r.curLayer = nil
		return r.getNextHandler(request)
	}

	handler := r.curLayer.Handlers[r.handlerPos]
	r.handlerPos++

	return handler
}

func (r *Runner) getNextLayer(request *http.Request) *layer.Layer {
	r.handlerPos = 0

	for r.layerPos < len(r.Layers) {
		refLayer := &r.Layers[r.layerPos]
		r.layerPos++

		refLayer, uriParams := r.Resolver.ForRequest(refLayer, request, true)
		if refLayer != nil {
			r.addUriParams(uriParams)
			return refLayer
		}
	}

	// (нужно сделать чек лист, по которому нужно ставить fallback, чтобы в ядро не тащить)
	panic(fmt.Errorf("can`t delegate to layer by index %d", r.layerPos))
}

func (r *Runner) addUriParams(params map[string]string) {
	for k, v := range params {
		r.UriParams[k] = v
	}
}
