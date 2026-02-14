package printer

import (
	"os/exec"

	"go.uber.org/zap"
)

type PrintJob struct {
	cmd      *exec.Cmd
	ID       string
	filename string
	State    State
	Logger   *zap.Logger
	Output   []byte
}

type State int

const (
	StateCreated State = iota
	StateStarted
	StateSucceeded
	StateFailed
)

func (j *PrintJob) Run() error {
	j.Logger.Info("Start job")

	j.State = StateStarted
	var err error

	j.Output, err = j.cmd.CombinedOutput()
	if err != nil {
		j.State = StateFailed

		j.Logger.Error("Job failed")

		return err
	}

	j.State = StateSucceeded

	j.Logger.Info("Job succeeded")

	return nil
}

func (j *PrintJob) Filename() string {
	return j.filename
}
