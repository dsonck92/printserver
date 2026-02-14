package web

import (
	"context"
	"sort"

	echo "github.com/labstack/echo/v4"
	"github.com/tvanriel/cloudsdk/http"
	"github.com/tvanriel/printserver/pkg/printer"
	"github.com/tvanriel/printserver/pkg/scan"
	"github.com/tvanriel/printserver/pkg/web/assets"
	"github.com/tvanriel/printserver/pkg/web/views"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var _ http.RouteGroup = (*Controller)(nil)

type ControllerOpts struct {
	fx.In

	Logger   *zap.Logger
	Printers printer.Printers
	Scanner  *scan.Scanner
}

func NewController(o ControllerOpts) *Controller {
	return &Controller{
		Logger:   o.Logger.Named("web"),
		Printers: o.Printers,
		Scanner:  o.Scanner,
	}
}

type Controller struct {
	Logger   *zap.Logger
	Printers printer.Printers
	Scanner  *scan.Scanner
}

// ApiGroup implements http.RouteGroup.
func (c *Controller) ApiGroup() string {
	return ""
}

// Handler implements http.RouteGroup.
func (c *Controller) Handler(g *echo.Group) {
	g.GET("", c.index)
	g.GET("bootstrap.min.css", c.Asset(assets.Bootstrap, "text/css"))
	g.GET("htmx.min.js", c.Asset(assets.HTMX, "application/javascript"))
	g.POST("print/:printer", c.print)
	g.POST("scan", c.scan)
	g.GET("scan/:id", c.scanid)
	g.GET("print/:printer/:id", c.printid)
	g.GET("scanimage/:id", c.scanimage)
	g.GET("scanjobs", c.scanjobs)
	g.GET("printjobs/:printer", c.printjobs)
}

func (c *Controller) Asset(x []byte, s string) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", s)
		_, _ = c.Response().Writer.Write(x)
		return nil
	}
}

// Version implements http.RouteGroup.
func (c *Controller) Version() string {
	return ""
}

func (c *Controller) print(ctx echo.Context) error {
	pr := c.Printers[ctx.Param("printer")]

	f, h, err := ctx.Request().FormFile("file")
	if err != nil {
		c.Logger.Error("read form file", zap.Error(err))
		return nil
	}
	job, err := pr.NewJob(context.Background(), h.Filename, f)
	if err != nil {
		c.Logger.Error("submit job", zap.Error(err))
	}

	go job.Run()

	return c.printjobs(ctx)
}

func (c *Controller) printid(ctx echo.Context) error {
	pr := c.Printers[ctx.Param("printer")]

	for _, j := range pr.Jobs {
		if j.ID == ctx.Param("id") {
			views.PrintJob(j).Render(ctx.Request().Context(), ctx.Response().Writer)
		}
	}

	return nil
}

func (c *Controller) scanid(ctx echo.Context) error {
	for _, j := range c.Scanner.Jobs {
		if j.ID == ctx.Param("id") {
			views.ScanJob(j).Render(ctx.Request().Context(), ctx.Response().Writer)
		}
	}

	return nil
}

func (c *Controller) scanimage(ctx echo.Context) error {
	for _, j := range c.Scanner.Jobs {
		if j.ID == ctx.Param("id") {
			ctx.File(j.Filename())
		}
	}

	return nil
}

func (c *Controller) printjobs(ctx echo.Context) error {
	pr := c.Printers[ctx.Param("printer")]

	if err := views.PrintJobs(pr.Jobs).Render(ctx.Request().Context(), ctx.Response().Writer); err != nil {
		c.Logger.Error("write printjobs response", zap.Error(err))
	}

	return nil
}

func (c *Controller) scanjobs(ctx echo.Context) error {
	if err := views.ScanJobs(c.Scanner.Jobs).Render(ctx.Request().Context(), ctx.Response().Writer); err != nil {
		c.Logger.Error("write scanjobs response", zap.Error(err))
	}

	return nil
}

func (c *Controller) scan(ctx echo.Context) error {
	j := c.Scanner.NewJob(context.Background())
	go j.Run()

	return c.scanjobs(ctx)
}

func (c *Controller) index(ctx echo.Context) error {
	printers := make([]string, 0, len(c.Printers))
	for name := range c.Printers {
		printers = append(printers, name)
	}

	sort.Strings(printers)

	if err := views.Index(printers).Render(ctx.Request().Context(), ctx.Response().Writer); err != nil {
		c.Logger.Error("write index response", zap.Error(err))
	}

	return nil
}
