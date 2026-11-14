package terraform

import (
	"context"
	"fmt"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	clientgitlab "gitlab.com/gitlab-org/api/client-go/v3"

	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

// keeps the notes per merge request, so a spec reads what each one was told.
type deployNotes struct {
	bodies map[int64]string
}

func (f *deployNotes) ListMergeRequestNotes(
	_ any,
	_ int64,
	_ *clientgitlab.ListMergeRequestNotesOptions,
	_ ...clientgitlab.RequestOptionFunc,
) ([]*clientgitlab.Note, *clientgitlab.Response, error) {
	return nil, nil, nil
}

func (f *deployNotes) CreateMergeRequestNote(
	_ any,
	mergeRequest int64,
	opt *clientgitlab.CreateMergeRequestNoteOptions,
	_ ...clientgitlab.RequestOptionFunc,
) (*clientgitlab.Note, *clientgitlab.Response, error) {
	f.bodies[mergeRequest] = *opt.Body

	return &clientgitlab.Note{ID: mergeRequest}, nil, nil
}

func (f *deployNotes) UpdateMergeRequestNote(
	_ any,
	_ int64,
	_ int64,
	_ *clientgitlab.UpdateMergeRequestNoteOptions,
	_ ...clientgitlab.RequestOptionFunc,
) (*clientgitlab.Note, *clientgitlab.Response, error) {
	ginkgo.Fail("no note is listed, so none can be updated")

	return nil, nil, nil
}

// answers the merge requests of the commit and the jobs of the pipeline, the only
// calls a deploy report makes.
type deployMergeRequests struct {
	gitlab.MergeRequestsAdapter
	merged  []*clientgitlab.BasicMergeRequest
	jobs    []*clientgitlab.Job
	commits []string
	err     error
}

func (f *deployMergeRequests) ListMergeRequestsByCommit(
	_ any,
	sha string,
	_ ...clientgitlab.RequestOptionFunc,
) ([]*clientgitlab.BasicMergeRequest, *clientgitlab.Response, error) {
	f.commits = append(f.commits, sha)

	return f.merged, &clientgitlab.Response{}, f.err
}

func (f *deployMergeRequests) ListPipelineJobs(
	_ any,
	_ int64,
	_ *clientgitlab.ListJobsOptions,
	_ ...clientgitlab.RequestOptionFunc,
) ([]*clientgitlab.Job, *clientgitlab.Response, error) {
	return f.jobs, &clientgitlab.Response{}, nil
}

var _ = ginkgo.Describe("Deploy report", func() {
	var notes *deployNotes
	var mergeRequests *deployMergeRequests

	ginkgo.BeforeEach(func() {
		notes = &deployNotes{bodies: map[int64]string{}}
		mergeRequests = &deployMergeRequests{
			// only the one merged into the default branch brought the commit in.
			merged: []*clientgitlab.BasicMergeRequest{
				{IID: 5, State: "merged", SourceBranch: "feature", TargetBranch: "main"},
				{IID: 8, State: "merged", SourceBranch: "feature", TargetBranch: "next"},
				{IID: 9, State: "opened", SourceBranch: "other", TargetBranch: "main"},
			},
			jobs: []*clientgitlab.Job{
				{Name: "plan", WebURL: "https://gitlab.test/-/jobs/1"},
				{Name: "deploy", WebURL: "https://gitlab.test/-/jobs/2"},
			},
		}
	})

	source := func(job string) *Source {
		return &Source{
			Read: func(_ context.Context, _ *Task) (Report, error) {
				return Report{Title: "Pulumi preview report"}, nil
			},
			MergeRequest: gitlab.MergeRequestReportConfig{
				Enabled:       true,
				Token:         "glpat-report",
				ApiUrl:        "https://gitlab.test/api/v4",
				ProjectId:     "group/project",
				Branch:        "main",
				DefaultBranch: "main",
				PlanJob:       "plan",
				DeployJob:     "deploy",
			},
			Notes: func(_ gitlab.MergeRequestReportConfig) (gitlab.NotesAdapter, error) { return notes, nil },
			MergeRequests: func(_ gitlab.MergeRequestConfig) (gitlab.MergeRequestsAdapter, error) {
				return mergeRequests, nil
			},
			Discriminators: func() []string { return []string{"production"} },
			Metadata: Metadata{
				Target:      "production",
				JobName:     job,
				JobUrl:      "https://gitlab.test/-/jobs/" + map[string]string{"plan": "1", "deploy": "2"}[job],
				PipelineId:  "42",
				PipelineUrl: "https://gitlab.test/-/pipelines/42",
				CommitSha:   "abc",
			},
		}
	}

	plan := func(src *Source) error {
		return DeployReportTask(NewTaskList(tests.NewPlumber().Plumber), src).Run(context.Background())
	}

	ginkgo.It("leaves the plan with the plan and the deploy job on the merge request merged into the commit", func() {
		Expect(plan(source("plan"))).To(Succeed())

		Expect(mergeRequests.commits).To(Equal([]string{"abc"}))
		Expect(notes.bodies).To(HaveLen(1))
		Expect(notes.bodies[5]).To(HavePrefix(
			"Merged into `main`. [Job `plan`](https://gitlab.test/-/jobs/1) planned the changes below for `production`, " +
				"which [the deploy job `deploy`](https://gitlab.test/-/jobs/2) applies once it is run.\n\n## Pulumi preview report\n",
		))
		Expect(notes.bodies[5]).To(HaveSuffix("<!-- gitlab-pipes:mr-report:plan:production:deploy-planned -->"))
	})

	ginkgo.It("links the pipeline once the deploy job is not in it", func() {
		mergeRequests.jobs = nil

		Expect(plan(source("plan"))).To(Succeed())

		Expect(notes.bodies[5]).To(ContainSubstring(
			"which the deploy job of [pipeline #42](https://gitlab.test/-/pipelines/42) applies once it is run.",
		))
	})

	ginkgo.DescribeTable("leaves the merge requests alone outside a pipeline of the default branch",
		func(change func(src *Source)) {
			src := source("plan")
			change(src)

			Expect(plan(src)).To(Succeed())

			Expect(mergeRequests.commits).To(BeEmpty())
			Expect(notes.bodies).To(BeEmpty())
		},
		ginkgo.Entry("a merge request pipeline", func(src *Source) { src.MergeRequest.MergeRequestIid = 7 }),
		ginkgo.Entry("another branch", func(src *Source) { src.MergeRequest.Branch = "feature" }),
		ginkgo.Entry("a tag", func(src *Source) { src.MergeRequest.Branch = "" }),
		ginkgo.Entry("the report disabled", func(src *Source) { src.MergeRequest.Enabled = false }),
	)

	ginkgo.Describe("apply", func() {
		apply := func(src *Source, err error) error {
			return src.Apply(context.Background(), fixtures.Task("apply"), func(_ context.Context) error { return err })
		}

		ginkgo.It("reports the deploy with the plan it applied", func() {
			Expect(apply(source("deploy"), nil)).To(Succeed())

			Expect(notes.bodies).To(HaveLen(1))
			Expect(notes.bodies[5]).To(Equal(
				"Deployed what was merged into `main`. [Job `deploy`](https://gitlab.test/-/jobs/2) applied the plan of " +
					"[the plan job `plan`](https://gitlab.test/-/jobs/1) for `production`.\n" +
					"\n<!-- gitlab-pipes:mr-report:deploy:production:deploy -->",
			))
		})

		ginkgo.It("reports the failed deploy and fails with the apply", func() {
			Expect(apply(source("deploy"), fmt.Errorf("apply failed"))).To(MatchError("apply failed"))

			Expect(notes.bodies[5]).To(HavePrefix(
				"Deploying what was merged into `main` failed. [Job `deploy`](https://gitlab.test/-/jobs/2) failed while applying the plan of " +
					"[the plan job `plan`](https://gitlab.test/-/jobs/1) for `production`, so it may be applied only in part.\n",
			))
		})

		ginkgo.It("never fails the deploy with its report", func() {
			mergeRequests.err = fmt.Errorf("gitlab is down")

			Expect(apply(source("deploy"), nil)).To(Succeed())
			Expect(notes.bodies).To(BeEmpty())
		})

		ginkgo.It("keeps the error of the apply once its report fails as well", func() {
			mergeRequests.err = fmt.Errorf("gitlab is down")

			Expect(apply(source("deploy"), fmt.Errorf("apply failed"))).To(MatchError("apply failed"))
		})
	})
})
