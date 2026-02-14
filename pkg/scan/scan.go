package scan

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Config struct {
	ScanimageBinary string   `hcl:"scanimage_binary"`
	ScanimageArgs   []string `hcl:"scanimage_args"`
	DestDir         string   `hcl:"dest_dir"`
}

type Scanner struct {
	ScanimageBinary string
	DestDir         string
	Logger          *zap.Logger
	Jobs            []*ScanJob
	Args []string
}

type ScannerOpts struct {
	fx.In

	Config Config
	Logger *zap.Logger
}

func NewScanner(o ScannerOpts) *Scanner {
	_ = os.Mkdir(o.Config.DestDir, os.ModePerm)
	return &Scanner{
		ScanimageBinary: o.Config.ScanimageBinary,
		DestDir:         o.Config.DestDir,
		Logger:          o.Logger.Named("scanner"),
		Args: o.Config.ScanimageArgs,

		Jobs: []*ScanJob{},
	}
}

func (s *Scanner) NewJob(ctx context.Context) *ScanJob {
	jobId := uuid.Must(uuid.NewRandom()).String()
	f := filepath.Join(s.DestDir, s.Filename(jobId))
	args := append(s.Args, "-o", f, "--format=png")
	cmd := exec.CommandContext(ctx, s.ScanimageBinary, args...)

	logger := s.Logger.With(zap.String("job_id", jobId))

	s.Logger.Info("command", zap.String("bin", cmd.Path), zap.Strings("args", cmd.Args))

	logger.Info("Created new job")

	s.Jobs = append(s.Jobs, &ScanJob{
		ID:       jobId,
		cmd:      cmd,
		filename: f,
		Logger:   logger,
	})

	return s.Jobs[len(s.Jobs)-1]
}

func (s *Scanner) Filename(jobId string) string {
	return fmt.Sprintf("%s.png", jobId)
}
