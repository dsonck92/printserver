package printer

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Config struct {
	Name string `hcl:"name,label"`

	LPBinary string   `hcl:"lp_binary"`
	LPArgs   []string `hcl:"lp_args"`
	DestDir  string   `hcl:"dest_dir"`
}

type PrinterOpts struct {
	fx.In

	Logger *zap.Logger

	Config []Config
}

type Printer struct {
	Logger *zap.Logger

	LPBinary string
	LPArgs   []string
	DestDir  string

	Jobs []*PrintJob
}

type Printers map[string]*Printer

func NewPrinter(o PrinterOpts) Printers {
	printers := make(Printers, len(o.Config))
	for _, c := range o.Config {
		_ = os.Mkdir(c.DestDir, os.ModePerm)

		printers[c.Name] = &Printer{
			Logger:   o.Logger.Named("printer").Named(c.Name),
			LPBinary: c.LPBinary,
			LPArgs:   c.LPArgs,
			DestDir:  c.DestDir,
			Jobs:     []*PrintJob{},
		}
	}

	return printers
}

func (p *Printer) NewJob(ctx context.Context, filename string, r io.Reader) (*PrintJob, error) {
	jobId := uuid.Must(uuid.NewRandom()).String()
	jobDir := filepath.Join(p.DestDir, jobId)
	err := os.Mkdir(jobDir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("create job dir: %w", err)
	}

	fullFilename := filepath.Join(jobDir, filename)
	f, err := os.Create(fullFilename)
	if err != nil {
		return nil, fmt.Errorf("create output file: %w", err)
	}

	defer f.Close()
	_, err = io.Copy(f, r)
	if err != nil {
		return nil, fmt.Errorf("save to-be-printed-file on disk: %w", err)
	}

	args := append(p.LPArgs, fullFilename)

	cmd := exec.CommandContext(ctx, p.LPBinary, args...)

	p.Logger.Info("command", zap.String("bin", cmd.Path), zap.Strings("args", cmd.Args))

	logger := p.Logger.With(zap.String("job_id", jobId))
	logger.Info("Created new job")

	p.Jobs = append(p.Jobs, &PrintJob{
		ID:       jobId,
		cmd:      cmd,
		filename: fullFilename,
		Logger:   logger,
	})

	return p.Jobs[len(p.Jobs)-1], nil
}
