package main

import (
	"context"
	"os"

	"github.com/tvanriel/printserver/pkg/app"
	"github.com/urfave/cli/v3"
)

func main() {
	(&cli.Command{
		Name:        "",
		Description: "Run the print and Scan server",
		Action: func(ctx context.Context, c *cli.Command) error {
			return app.Run(ctx)
		},
	}).Run(context.Background(), os.Args)
}
