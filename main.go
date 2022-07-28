package main

import (
	"context"
	"net/http"
	"os"

	"github.com/HDRUK/dar-integration-sandbox/api"
	"github.com/gorilla/mux"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(mux.NewRouter),
		fx.Provide(api.ProvideDatabase),
		fx.Invoke(api.LoadVariables),
		fx.Invoke(api.NewServer),
		fx.Invoke(registerHooks),
		api.LoggerFXModule,
	).Run()
}

func registerHooks(lifecycle fx.Lifecycle, mux *mux.Router, logger *zap.SugaredLogger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go http.ListenAndServe(":"+os.Getenv("PORT"), mux)

			return nil
		},
		OnStop: func(context.Context) error {
			return logger.Sync()
		},
	})
}
