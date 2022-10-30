package pipe

import (
	"github.com/urfave/cli/v2"
)

//revive:disable:line-length-limit

const (
	category_gitlab          = "Gitlab"
	category_gitlab_pipeline = "Pipeline"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category:    category_gitlab,
		Name:        "gl.api_url",
		Usage:       "Gitlab API URL of the instance.",
		Required:    true,
		EnvVars:     []string{"CI_API_V4_URL"},
		Value:       "",
		Destination: &TL.Pipe.Gitlab.ApiUrl,
	},

	&cli.StringFlag{
		Category:    category_gitlab,
		Name:        "gl.token",
		Usage:       "Token for Gitlab API authentication.",
		Required:    true,
		EnvVars:     []string{"GL_TOKEN"},
		Value:       "",
		Destination: &TL.Pipe.Gitlab.Token,
	},

	&cli.StringFlag{
		Category:    category_gitlab,
		Name:        "gl.job_token",
		Usage:       "Job token coming from the build job.",
		Required:    false,
		EnvVars:     []string{"CI_JOB_TOKEN"},
		Value:       "",
		Destination: &TL.Pipe.Gitlab.JobToken,
	},

	&cli.StringFlag{
		Category:    category_gitlab_pipeline,
		Name:        "gl_pipeline.project_id",
		Usage:       "Parent project id.",
		Required:    true,
		EnvVars:     []string{"CI_PROJECT_ID"},
		Value:       "",
		Destination: &TL.Pipe.Gitlab.ParentProjectId,
	},

	&cli.StringFlag{
		Category:    category_gitlab_pipeline,
		Name:        "gl_pipeline.parent_pipeline_id",
		Usage:       "Pipeline id of the parent pipeline.",
		Required:    true,
		EnvVars:     []string{"PARENT_PIPELINE_ID"},
		Value:       "",
		Destination: &TL.Pipe.Gitlab.ParentPipelineId,
	},

	&cli.StringFlag{
		Category:    category_gitlab_pipeline,
		Name:        "gl_pipeline.download_artifacts",
		Usage:       "Names of the jobs that yield artifacts from the parent job.",
		Required:    true,
		EnvVars:     []string{"PARENT_DOWNLOAD_ARTIFACTS"},
		Destination: &TL.Pipe.Gitlab.DownloadArtifacts,
	},
}
