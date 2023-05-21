package layer

import (
	"net/http"
	"regexp"
)

type Restrictions map[string]string

type Layer struct {
	Handlers []http.HandlerFunc

	Name         string
	Path         string
	RegExp       *regexp.Regexp
	Priority     int
	Methods      []string
	Restrictions Restrictions

	Meta map[string]any
}

func (l Layer) WithHandlers(handlers ...http.HandlerFunc) Layer {
	l.Handlers = handlers
	return l
}
