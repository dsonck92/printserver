package config

import (
	"fmt"

	"github.com/tvanriel/cloudsdk/hclconfig"
	"github.com/tvanriel/cloudsdk/http"
	"github.com/tvanriel/cloudsdk/logging"
	"github.com/tvanriel/printserver/pkg/printer"
	"github.com/tvanriel/printserver/pkg/scan"
	"go.uber.org/fx"
)

type Config struct {
	fx.Out

	Printer printer.Config        `hcl:"printer,block"`
	Scanner scan.Config           `hcl:"scanner,block"`
	Http    http.Configuration    `hcl:"http,block"`
	Logging logging.Configuration `hcl:"logging,block"`
}

func ParseConfig() (Config, error) {
	var config Config
	err := hclconfig.HclConfiguration(&config, "printserver")
	if err != nil {
		return config, fmt.Errorf("parse printserver config: %w", err)
	}

	return config, nil
}
