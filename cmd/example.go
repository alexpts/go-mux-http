package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/alexpts/go-mux-http/mux"
	"github.com/alexpts/go-mux-http/mux/layer"
)

type Layer = layer.Layer

func main() {
	app := mux.NewApp(nil, nil)

	action := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello"))
	})
	actionAlex := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello Alex"))
	})

	middlewareForAll := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// some middleware
		app.ServeHTTP(w, r)
	})

	middlewareForAdmin := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// some middleware not for all paths
		app.ServeHTTP(w, r)
	})

	app.
		Use(Layer{Path: "/admin/{slug}/"}, middlewareForAll). // middleware activates only `/admin/*/`
		Use(Layer{}, middlewareForAdmin).
		Get("/", Layer{}, action).
		Get("/alex/", Layer{}, actionAlex)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/alex/", nil)
	req = req.WithContext(context.Background())
	app.ServeHTTP(w, req)

	fmt.Println(w.Body.String())
}
