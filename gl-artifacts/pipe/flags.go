package pipe

import (
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "gl.token",
		Usage:       "Token for gitlab api authentication.",
		Required:    true,
		EnvVars:     []string{"GL_TOKEN", "BANDAID_API_TOKEN"},
		Value:       "https://gitlab.kilic.dev",
		Destination: &Pipe.Gitlab.Token,
	},

	&cli.StringFlag{
		Name:        "gl.job_token",
		Usage:       "Job token coming from the build job.",
		Required:    false,
		EnvVars:     []string{"CI_JOB_TOKEN"},
		Value:       "https://gitlab.kilic.dev",
		Destination: &Pipe.Gitlab.JobToken,
	},

	&cli.StringFlag{
		Name:        "gl_pipeline.project_id",
		Usage:       "Parent project id.",
		Required:    true,
		EnvVars:     []string{"CI_PROJECT_ID"},
		Value:       "",
		Destination: &Pipe.Gitlab.ParentProjectId,
	},

	&cli.StringFlag{
		Name:        "gl_pipeline.parent_pipeline_id",
		Usage:       "Pipeline id of the parent pipeline.",
		Required:    true,
		EnvVars:     []string{"PARENT_PIPELINE_ID"},
		Value:       "",
		Destination: &Pipe.Gitlab.ParentPipelineId,
	},

	&cli.StringSliceFlag{
		Name:        "gl_pipeline.download_artifacts",
		Usage:       "Names of the jobs that yield artifacts from the parent job.",
		Required:    true,
		EnvVars:     []string{"PARENT_DOWNLOAD_ARTIFACTS"},
		Destination: &Pipe.Gitlab.DownloadArtifacts,
	},
}
