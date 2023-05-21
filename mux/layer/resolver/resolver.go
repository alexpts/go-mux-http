package resolver

import (
	"net/http"

	"github.com/alexpts/go-mux-http/internal"
	"github.com/alexpts/go-mux-http/mux/layer"
)

type RequestResolver struct{}

func (r *RequestResolver) ForRequest(
	l *layer.Layer,
	request *http.Request,
	checkMethod bool,
) (*layer.Layer, map[string]string) {
	if checkMethod && !isAllowMethod(l, request) {
		return nil, nil
	}

	if l.Path == `` {
		return l, nil
	}

	l, matched := matchRegexpLayer(l, request)
	if l == nil {
		return nil, nil
	}

	params := make(map[string]string)
	fillRequestUriParams(params, matched, l)

	return l, params
}

func isAllowMethod(l *layer.Layer, req *http.Request) bool {
	if len(l.Methods) == 0 {
		return true
	}

	return internal.InSlice(l.Methods, req.Method)
}

// matchRegexpLayer compare request and Layer config.
func matchRegexpLayer(l *layer.Layer, req *http.Request) (*layer.Layer, []string) {
	matched := l.RegExp.FindStringSubmatch(req.URL.Path)
	if len(matched) == 0 {
		return nil, matched
	}

	return l, matched
}

// [mutable].
func fillRequestUriParams(params map[string]string, matched []string, l *layer.Layer) {
	groups := l.RegExp.SubexpNames()
	for i, name := range groups {
		if name != `` {
			params[name] = matched[i]
		}
	}
}
