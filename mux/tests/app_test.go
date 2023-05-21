package tests_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alexpts/go-mux-http/mux"
	"github.com/alexpts/go-mux-http/mux/layer"
)

type Layer = layer.Layer

func createRequest(method string, path string) *http.Request {
	return httptest.NewRequest(method, path, nil)
}

func createApp() mux.App {
	return mux.NewApp(nil, nil)
}

func TestMinimalApp(t *testing.T) {
	app := createApp()

	app.Use(Layer{}, func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("Hello"))
	})

	w := httptest.NewRecorder()
	app.ServeHTTP(w, createRequest(`GET`, `/api/`))

	assert.Equal(t, `Hello`, w.Body.String())
}

func TestMultiHandler(t *testing.T) {
	app := createApp()

	app.Use(layer.Layer{}, func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("Hello"))
		app.ServeHTTP(writer, request)
	})

	app.Get(
		"/api/",
		layer.Layer{},
		func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(" World"))
			app.ServeHTTP(writer, request)
		},
		func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(" Alex"))
		},
	)

	w := httptest.NewRecorder()
	app.ServeHTTP(w, createRequest(`GET`, `/api/`))

	assert.Equal(t, `Hello World Alex`, w.Body.String())
}

func TestMultiLayers(t *testing.T) {
	app := createApp()

	app.AddLayer(Layer{}.
		WithHandlers(
			func(writer http.ResponseWriter, request *http.Request) {
				_, _ = writer.Write([]byte(`Hello`))
				app.ServeHTTP(writer, request)
			},
		),
	)

	app.AddLayer(Layer{}.
		WithHandlers(
			func(writer http.ResponseWriter, request *http.Request) {
				_, _ = writer.Write([]byte(` World`))
			},
		),
	)

	w := httptest.NewRecorder()
	app.ServeHTTP(w, createRequest(`GET`, `/`))
	assert.Equal(t, `Hello World`, w.Body.String())
}

func TestLayerPriority(t *testing.T) {
	app := createApp()

	app.
		Use(Layer{Priority: 100}, func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(`1-`)) // run second
		}).
		Use(Layer{Priority: 200}, func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(`2-`)) // run first
			app.ServeHTTP(writer, request)
		})

	w := httptest.NewRecorder()
	app.ServeHTTP(w, createRequest(`GET`, `/`))
	assert.Equal(t, `2-1-`, w.Body.String())
}

func TestDelegateToNotDefinedLayer(t *testing.T) {
	app := createApp()

	app.Use(Layer{}, func(writer http.ResponseWriter, request *http.Request) {
		app.ServeHTTP(writer, request)
	})

	w := httptest.NewRecorder()

	assert.Panics(t, func() {
		app.ServeHTTP(w, createRequest(`GET`, `/`))
	}, "can`t delegate to layer by index 1")
}

func TestFilterByHttpMethod(t *testing.T) {
	app := createApp()

	app.
		Use(Layer{}, func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(`1-`))
			app.ServeHTTP(writer, request)
		}).
		// Disallow POST
		Use(Layer{Methods: []string{`POST`}}, func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(`2-`))
			app.ServeHTTP(writer, request)
			_, _ = writer.Write([]byte(`2_2`))
		}).
		// Allow one of GET
		Use(Layer{Methods: []string{`POST`, `GET`}}, func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(`3-`))
			app.ServeHTTP(writer, request)
		}).
		Use(Layer{Methods: []string{`GET`}}, func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(`4-`))
		})

	w := httptest.NewRecorder()
	app.ServeHTTP(w, createRequest(`GET`, `/`))
	assert.Equal(t, `1-3-4-`, w.Body.String())
}

func TestFilterByPath(t *testing.T) {
	app := createApp()

	app.
		Use(Layer{Path: `/users/`}, func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(`layer1-`))
			app.ServeHTTP(writer, request)
		}).
		Use(Layer{Path: `/admin/`}, func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(`layer2-`))
			app.ServeHTTP(writer, request)
		}).
		Use(Layer{}, func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(`layer3-`))
		})

	w := httptest.NewRecorder()
	app.ServeHTTP(w, createRequest(`GET`, `/admin/`))

	assert.Equal(t, `layer2-layer3-`, w.Body.String())
}

func TestMatchUrlParam(t *testing.T) {
	app := createApp()

	app.Use(Layer{Path: `/city/{slug}/`}, func(writer http.ResponseWriter, request *http.Request) {
		runnerCtx, _ := mux.GetMuxRunnerCtx(request.Context())
		slug, ok := runnerCtx.UriParams["slug"]
		if ok {
			_, _ = writer.Write([]byte(slug))
		}
	})

	w := httptest.NewRecorder()
	app.ServeHTTP(w, createRequest(`GET`, `/city/london/`))
	assert.Equal(t, "london", w.Body.String())
}

func TestFastHttpMethod(t *testing.T) {
	type testProvider struct {
		method   string
		expected string
	}

	tests := map[string]testProvider{
		"GET": {
			method:   `GET`,
			expected: `GET`,
		},
		"POST": {
			method:   `POST`,
			expected: `POST`,
		},
		"PUT": {
			method:   `PUT`,
			expected: `PUT`,
		},
		"PATCH": {
			method:   `PATCH`,
			expected: `PATCH`,
		},
		"DELETE": {
			method:   `DELETE`,
			expected: `DELETE`,
		},
	}

	for name, provider := range tests {
		t.Run(name, func(t *testing.T) {
			app := createApp()

			handler := func(writer http.ResponseWriter, request *http.Request) {
				_, _ = writer.Write([]byte(request.Method))
			}

			switch provider.method {
			case `GET`:
				app.Get(`/`, Layer{}, handler)
			case `POST`:
				app.Post(`/`, Layer{}, handler)
			case `PUT`:
				app.Put(`/`, Layer{}, handler)
			case `PATCH`:
				app.Patch(`/`, Layer{}, handler)
			case `DELETE`:
				app.Delete(`/`, Layer{}, handler)
			}

			w := httptest.NewRecorder()
			app.ServeHTTP(w, createRequest(provider.method, `/`))
			assert.Equal(t, provider.expected, w.Body.String())
		})
	}
}

func TestMount(t *testing.T) {
	apiV1 := createApp()

	apiV1.
		Use(Layer{}, func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(`api v1 md; `))
			apiV1.ServeHTTP(writer, request)
		}).
		Get(`/users/`, Layer{}, func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(`api v1 - users`))
		})

	apiV2 := createApp()
	apiV2.Get(`/users/`, Layer{}, func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`api v2 - users`))
	})

	reuseApp := createApp()
	reuseApp.Get(`/users/`, Layer{}, func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`reuse - users`))
	})

	app := createApp()
	app.
		Mount(apiV1, `/api/v1`).
		Mount(apiV2, `/api/v2`).
		Mount(reuseApp, ``)

	// asserts
	w := httptest.NewRecorder()
	app.ServeHTTP(w, createRequest(`GET`, `/api/v1/users/`))
	assert.Equal(t, `api v1 md; api v1 - users`, w.Body.String())

	w = httptest.NewRecorder()
	app.ServeHTTP(w, createRequest(`GET`, `/api/v2/users/`))
	assert.Equal(t, `api v2 - users`, w.Body.String())

	w = httptest.NewRecorder()
	app.ServeHTTP(w, createRequest(`GET`, `/users/`))
	assert.Equal(t, `reuse - users`, w.Body.String())
}

// mux -> app.
func TestDelegateFromOldMux(t *testing.T) {
	app := createApp()

	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(`old mux / `))
			next.ServeHTTP(writer, request)
		})
	}

	oldMux := http.ServeMux{}
	oldMux.Handle("/", middleware(&app))

	app.Get("/", Layer{}, func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`app mux`))
	})

	w := httptest.NewRecorder()
	oldMux.ServeHTTP(w, createRequest(`GET`, `/`))
	assert.Equal(t, "old mux / app mux", w.Body.String())
}

// app -> mux.
func TestDelegateToOldMux(t *testing.T) {
	app := createApp()

	oldMux := http.ServeMux{}
	oldMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`old mux`))
	})

	app.Get("/", Layer{}, func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`app mux / `))
		oldMux.ServeHTTP(writer, request)
	})

	w := httptest.NewRecorder()
	app.ServeHTTP(w, createRequest(`GET`, `/`))
	assert.Equal(t, "app mux / old mux", w.Body.String())
}

func TestRestrictionParam(t *testing.T) {
	app := createApp()

	app.
		Use(Layer{
			Path:     `/users/{name}/`,
			Name:     `User Page`,
			Methods:  []string{`GET`, `POST`},
			Priority: 100,
			Restrictions: layer.Restrictions{
				`name`: `[a-z]+`,
			},
		}, func(writer http.ResponseWriter, request *http.Request) {
			runnerCtx, _ := mux.GetMuxRunnerCtx(request.Context())
			slug := runnerCtx.UriParams["name"]

			_, _ = writer.Write([]byte(`user: ` + slug))
		}).
		Use(Layer{}, func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(`other`))
		})

	w := httptest.NewRecorder()
	app.ServeHTTP(w, createRequest(`GET`, `/users/alex/`))
	assert.Equal(t, `user: alex`, w.Body.String())

	w = httptest.NewRecorder()
	app.ServeHTTP(w, createRequest(`GET`, `/users/alex-2/`))
	assert.Equal(t, `other`, w.Body.String())
}

func TestInlineRestriction(t *testing.T) {
	app := createApp()

	app.
		Use(Layer{
			Path: `/users/{name:[a-z]+}/`,
		}, func(writer http.ResponseWriter, request *http.Request) {
			runnerCtx, _ := mux.GetMuxRunnerCtx(request.Context())
			slug := runnerCtx.UriParams["name"]

			_, _ = writer.Write([]byte(`user: ` + slug))
		}).
		Use(Layer{}, func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(`other`))
		})

	w := httptest.NewRecorder()
	app.ServeHTTP(w, createRequest(`GET`, `/users/alex/`))
	assert.Equal(t, `user: alex`, w.Body.String())

	w = httptest.NewRecorder()
	app.ServeHTTP(w, createRequest(`GET`, `/users/alex-2/`))
	assert.Equal(t, `other`, w.Body.String())
}
