package gitlab

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	clientgitlab "gitlab.com/gitlab-org/api/client-go/v3"
)

type fakeMergeRequests struct {
	listed  []*clientgitlab.ListProjectMergeRequestsOptions
	created []*clientgitlab.CreateMergeRequestOptions
	updated map[int64]*clientgitlab.UpdateMergeRequestOptions

	open      []*clientgitlab.BasicMergeRequest
	listErr   error
	createErr error
	updateErr error
}

var _ MergeRequestsAdapter = (*fakeMergeRequests)(nil)

// answers the listing the way the API filters it, so only merge requests in the
// requested state ever reach the upsert.
func (f *fakeMergeRequests) ListProjectMergeRequests(
	_ any,
	opt *clientgitlab.ListProjectMergeRequestsOptions,
	_ ...clientgitlab.RequestOptionFunc,
) ([]*clientgitlab.BasicMergeRequest, *clientgitlab.Response, error) {
	f.listed = append(f.listed, opt)
	if f.listErr != nil {
		return nil, nil, f.listErr
	}

	listed := []*clientgitlab.BasicMergeRequest{}
	for _, mergeRequest := range f.open {
		if opt.State == nil || mergeRequest.State == *opt.State {
			listed = append(listed, mergeRequest)
		}
	}

	return listed, &clientgitlab.Response{}, nil
}

func (f *fakeMergeRequests) CreateMergeRequest(
	_ any,
	opt *clientgitlab.CreateMergeRequestOptions,
	_ ...clientgitlab.RequestOptionFunc,
) (*clientgitlab.MergeRequest, *clientgitlab.Response, error) {
	f.created = append(f.created, opt)
	if f.createErr != nil {
		return nil, nil, f.createErr
	}

	return &clientgitlab.MergeRequest{
		BasicMergeRequest: clientgitlab.BasicMergeRequest{IID: 99, WebURL: "https://gitlab.test/-/merge_requests/99"},
	}, nil, nil
}

func (f *fakeMergeRequests) UpdateMergeRequest(
	_ any,
	mergeRequest int64,
	opt *clientgitlab.UpdateMergeRequestOptions,
	_ ...clientgitlab.RequestOptionFunc,
) (*clientgitlab.MergeRequest, *clientgitlab.Response, error) {
	f.updated[mergeRequest] = opt
	if f.updateErr != nil {
		return nil, nil, f.updateErr
	}

	return &clientgitlab.MergeRequest{
		BasicMergeRequest: clientgitlab.BasicMergeRequest{IID: mergeRequest},
	}, nil, nil
}

var _ = Describe("Merge request upsert", func() {
	config := MergeRequestConfig{
		ProjectId:    "3",
		SourceBranch: "sync/source-generate",
		TargetBranch: "next",
		Title:        "chore: sync generated sources",
	}

	var mergeRequests *fakeMergeRequests

	BeforeEach(func() {
		mergeRequests = &fakeMergeRequests{updated: map[int64]*clientgitlab.UpdateMergeRequestOptions{}}
	})

	mergeRequest := func(iid int64, state string) *clientgitlab.BasicMergeRequest {
		return &clientgitlab.BasicMergeRequest{IID: iid, State: state}
	}

	upsert := func() (*MergeRequestResult, error) {
		return UpsertMergeRequest(context.Background(), mergeRequests, config, "sync body")
	}

	It("lists the open merge requests of the branch pair", func() {
		_, err := upsert()
		Expect(err).NotTo(HaveOccurred())
		Expect(mergeRequests.listed).To(HaveLen(1))
		Expect(*mergeRequests.listed[0].State).To(Equal("opened"))
		Expect(*mergeRequests.listed[0].SourceBranch).To(Equal(config.SourceBranch))
		Expect(*mergeRequests.listed[0].TargetBranch).To(Equal(config.TargetBranch))
	})

	It("creates a merge request for the branch pair when none exists", func() {
		result, err := upsert()
		Expect(err).NotTo(HaveOccurred())
		Expect(result.MergeRequestIid).To(Equal(int64(99)))
		Expect(result.WebUrl).To(Equal("https://gitlab.test/-/merge_requests/99"))
		Expect(result.Created).To(BeTrue())
		Expect(result.Action()).To(Equal("created"))
		Expect(mergeRequests.created).To(HaveLen(1))
		Expect(*mergeRequests.created[0].Title).To(Equal(config.Title))
		Expect(*mergeRequests.created[0].Description).To(Equal("sync body"))
		Expect(*mergeRequests.created[0].SourceBranch).To(Equal(config.SourceBranch))
		Expect(*mergeRequests.created[0].TargetBranch).To(Equal(config.TargetBranch))
		Expect(mergeRequests.updated).To(BeEmpty())
	})

	It("updates the description of the open merge request", func() {
		mergeRequests.open = []*clientgitlab.BasicMergeRequest{mergeRequest(7, "opened")}

		result, err := upsert()
		Expect(err).NotTo(HaveOccurred())
		Expect(result.MergeRequestIid).To(Equal(int64(7)))
		Expect(result.Created).To(BeFalse())
		Expect(result.Action()).To(Equal("updated"))
		Expect(mergeRequests.updated).To(HaveKey(int64(7)))
		Expect(*mergeRequests.updated[7].Description).To(Equal("sync body"))
		Expect(mergeRequests.created).To(BeEmpty())
	})

	DescribeTable(
		"recreates instead of updating",
		func(state string) {
			mergeRequests.open = []*clientgitlab.BasicMergeRequest{mergeRequest(7, state)}

			result, err := upsert()
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Created).To(BeTrue())
			Expect(mergeRequests.updated).To(BeEmpty())
		},
		Entry("a merged prior merge request", "merged"),
		Entry("a closed prior merge request", "closed"),
	)

	It("surfaces a listing failure", func() {
		mergeRequests.listErr = fmt.Errorf("boom")

		_, err := upsert()
		Expect(err).To(MatchError(ContainSubstring("list GitLab merge requests")))
		Expect(mergeRequests.created).To(BeEmpty())
	})

	It("surfaces a create failure", func() {
		mergeRequests.createErr = fmt.Errorf("boom")

		_, err := upsert()
		Expect(err).To(MatchError(ContainSubstring("create GitLab merge request")))
	})

	It("surfaces an update failure", func() {
		mergeRequests.open = []*clientgitlab.BasicMergeRequest{mergeRequest(7, "opened")}
		mergeRequests.updateErr = fmt.Errorf("boom")

		_, err := upsert()
		Expect(err).To(MatchError(ContainSubstring("update GitLab merge request")))
	})
})
