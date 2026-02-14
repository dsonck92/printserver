package app

import (
	"context"

	"github.com/tvanriel/cloudsdk/http"
	"github.com/tvanriel/cloudsdk/logging"
	"github.com/tvanriel/printserver/pkg/config"
	"github.com/tvanriel/printserver/pkg/printer"
	"github.com/tvanriel/printserver/pkg/scan"
	"github.com/tvanriel/printserver/pkg/web"
	"go.uber.org/fx"
)

func Run(_ context.Context) error {
	app := fx.New(
		logging.Module,
		logging.FXLogger(),

		scan.Module,
		printer.Module,
		web.Module,
		http.Module,
		fx.Provide(config.ParseConfig),
		fx.Invoke(func(_ *http.Http) {
		}),
	)

	app.Run()

	return nil
}
