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
	LPBinary string   `hcl:"lp_binary"`
	LPArgs   []string `hcl:"lp_args"`
	DestDir  string   `hcl:"dest_dir"`
}

type PrinterOpts struct {
	fx.In

	Logger *zap.Logger

	Config Config
}

type Printer struct {
	Logger *zap.Logger

	Config Config
	Jobs   []*PrintJob
}

func NewPrinter(o PrinterOpts) *Printer {

	_ = os.Mkdir(o.Config.DestDir, os.ModePerm)

	return &Printer{
		Logger: o.Logger.Named("printer"),
		Config: o.Config,

		Jobs: []*PrintJob{},
	}
}

func (p *Printer) NewJob(ctx context.Context, filename string, r io.Reader) (*PrintJob, error) {
	jobId := uuid.Must(uuid.NewRandom()).String()
	jobDir := filepath.Join(p.Config.DestDir, jobId)
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

	args := append(p.Config.LPArgs, fullFilename)

	cmd := exec.CommandContext(ctx, p.Config.LPBinary, args...)

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
