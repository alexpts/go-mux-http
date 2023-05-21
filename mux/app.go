package mux

import (
	"context"
	"net/http"

	"github.com/alexpts/go-mux-http/mux/layer/resolver"
	"github.com/alexpts/go-mux-http/mux/layer/runner"
	"github.com/alexpts/go-mux-http/mux/layer/store"
)

type RunnerCtxKey struct{}

type App struct {
	*store.LayersStore
	resolver resolver.IResolver
}

func NewApp(layersStore *store.LayersStore, resolverDep resolver.IResolver) App {
	if layersStore == nil {
		layersStore = store.New(nil)
	}

	if resolverDep == nil {
		resolverDep = &resolver.RequestResolver{}
	}

	return App{
		LayersStore: layersStore,
		resolver:    resolverDep,
	}
}

// ServeHTTP - implement mux/http handler.
func (app *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// работаем со своей копией данных, чтобы не было гонок
	runnerCtx, req := app.getRunner(req)
	runnerCtx.ServeHTTP(w, req)
}

// @todo использовать pool для уменьшения аллокаций (нужно попрофилировать потом).
func (app *App) getRunner(req *http.Request) (*runner.Runner, *http.Request) {
	ctx := req.Context()
	runnerCtx, _ := GetMuxRunnerCtx(ctx)

	if runnerCtx == nil {
		runnerCtx = &runner.Runner{
			Layers:     app.GetLayers(),
			Resolver:   app.resolver,
			UriParams:  make(map[string]string),
			UserParams: make(map[string]any),
		}

		ctx = context.WithValue(ctx, RunnerCtxKey{}, runnerCtx)
		req = req.WithContext(ctx)
	}

	return runnerCtx, req
}

func GetMuxRunnerCtx(ctx context.Context) (*runner.Runner, bool) {
	runnerCtx, isSet := ctx.Value(RunnerCtxKey{}).(*runner.Runner)
	return runnerCtx, isSet
}

// Mount - mount layers from another application to current application.
func (app *App) Mount(app2 App, prefix string) *App {
	for _, l := range app2.GetLayers() {
		newLayer := l

		newLayer.Path = prefix + l.Path
		if l.Path == `` {
			newLayer.Path += `/.*`
		}

		app.AddLayer(newLayer)
	}

	return app
}
