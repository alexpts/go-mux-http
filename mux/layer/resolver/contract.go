package resolver

import (
	"net/http"

	"github.com/alexpts/go-mux-http/mux/layer"
)

type IResolver interface {
	ForRequest(
		*layer.Layer,
		*http.Request,
		bool,
	) (*layer.Layer, map[string]string)
}
