package terraform

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CategoryReport = "Plan Report"
)

// LogConfig is what the copy of the report that goes to the job log is built from.
type LogConfig struct {
	Enabled bool
}

// Options is what the report flags are built onto.
type Options struct {
	Destination *LogConfig
}

// NewFlags reads whether the report reaches the job log, which it does on its own
// account: the copy is what a pipeline is read with, so neither the summary file nor
// the merge request note has a say over it.
func NewFlags(opts Options) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Category: CategoryReport,
			Name:     "plan-report.log.enabled",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PLAN_REPORT_LOG_ENABLED"),
			),
			Usage:       "Write the plan report to the job log.",
			Required:    false,
			Value:       true,
			Destination: &opts.Destination.Enabled,
		},
	}
}
