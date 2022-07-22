package main

import (
	"context"
	"net/http"

	"github.com/HDRUK/dar-integration-sandbox/api"
	"github.com/gorilla/mux"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(mux.NewRouter),
		fx.Provide(api.Get),
		fx.Provide(api.ProvideDatabase),
		fx.Invoke(api.NewServer),
		fx.Invoke(registerHooks),
		api.LoggerFXModule,
	).Run()
}

func registerHooks(lifecycle fx.Lifecycle, mux *mux.Router, logger *zap.SugaredLogger, config *api.Config) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go http.ListenAndServe(":"+config.Port, mux)

			return nil
		},
		OnStop: func(context.Context) error {
			return logger.Sync()
		},
	})
}
