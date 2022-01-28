// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coding

import (
	"context"
	"errors"
	"strconv"

	"github.com/jenkins-x/go-scm/scm"
)

type pullService struct {
	client *wrapper
}

type createMergeRequest struct {
	apiRequest
	DepotId    int    `json:"DepotId"`
	Title      string `json:"Title"`
	Content    string `json:"Content"`
	SrcBranch  string `json:"SrcBranch"`
	DestBranch string `json:"DestBranch"`
}

type mergeRequestInfo struct {
	Describe     string `json:"Describe"`
	Status       string `json:"Status"`
	Title        string `json:"Title"`
	TargetBranch string `json:"TargetBranch"`
	SourceBranch string `json:"SourceBranch"`
}

type createMergeRequestResponse struct {
	Response struct {
		apiResponse
		MergeInfo *struct {
			ProjectId      int `json:"ProjectId"`
			DepotId        int `json:"DepotId"`
			MergeRequestId int `json:"MergeRequestId"`
		}
	} `json:"Response"`
}

type getMergeRequestInfoRequest struct {
	apiRequest
	DepotId int `json:"DepotId"`
	MergeId int `json:"MergeId"`
}

type getMergeRequestInfoResponse struct {
	Response struct {
		apiResponse
		MergeRequestInfo *mergeRequestInfo `json:"MergeRequestInfo"`
	} `json:"Response"`
}

type mergeRequest struct {
	apiRequest
	DepotId           int    `json:"DepotId"`
	MergeId           int    `json:"MergeId"`
	Message           string `json:"Message"`
	IsDelSourceBranch bool   `json:"IsDelSourceBranch"`
	IsFastForward     bool   `json:"IsFastForward"`
	Squash            bool   `json:"Squash"`
}

type mergeResponse struct {
	Response struct {
		apiResponse
	} `json:"Response"`
}

type closeMergeRequest struct {
	apiRequest
	DepotId int `json:"DepotId"`
	MergeId int `json:"MergeId"`
}

func (s *pullService) Find(ctx context.Context, repo string, number int) (*scm.PullRequest, *scm.Response, error) {
	repoInfo, _, err := s.client.Repositories.Find(ctx, repo)
	if err != nil {
		return nil, nil, err
	}

	in := new(getMergeRequestInfoRequest)
	in.Action = "DescribeMergeRequest"
	did, err := strconv.Atoi(repoInfo.ID)
	if err != nil {
		return nil, nil, err
	}
	in.DepotId = did
	in.MergeId = number

	out := new(getMergeRequestInfoResponse)
	res, err := s.client.do(ctx, "POST", "", in, out)
	if err != nil {
		return nil, nil, err
	}

	if out.Response.Error != nil {
		return nil, nil, errors.New(out.Response.Error.Message)
	}

	return nil, res, err
}

func (s *pullService) FindComment(ctx context.Context, repo string, index, id int) (*scm.Comment, *scm.Response, error) {
	// path := fmt.Sprintf("api/v4/projects/%s/merge_requests/%d/notes/%d", encode(repo), index, id)
	// out := new(issueComment)
	// res, err := s.client.do(ctx, "GET", path, nil, out)
	// return convertIssueComment(out), res, err
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) List(ctx context.Context, repo string, opts scm.PullRequestListOptions) ([]*scm.PullRequest, *scm.Response, error) {
	// path := fmt.Sprintf("api/v4/projects/%s/merge_requests?%s", encode(repo), encodePullRequestListOptions(opts))
	// out := []*pr{}
	// res, err := s.client.do(ctx, "GET", path, nil, &out)
	// if err != nil {
	// 	return nil, res, err
	// }
	// convRepos, convRes, err := s.convertPullRequestList(ctx, out)
	// if err != nil {
	// 	return nil, convRes, err
	// }
	// return convRepos, res, nil
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) ListChanges(ctx context.Context, repo string, number int, opts scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	// path := fmt.Sprintf("api/v4/projects/%s/merge_requests/%d/changes?%s", encode(repo), number, encodeListOptions(opts))
	// out := new(changes)
	// res, err := s.client.do(ctx, "GET", path, nil, &out)
	// return convertChangeList(out.Changes), res, err
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) ListComments(ctx context.Context, repo string, index int, opts scm.ListOptions) ([]*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) ListLabels(ctx context.Context, repo string, number int, opts scm.ListOptions) ([]*scm.Label, *scm.Response, error) {
	// mr, _, err := s.Find(ctx, repo, number)
	// if err != nil {
	// 	return nil, nil, err
	// }

	// return mr.Labels, nil, nil
	return nil, nil, nil
}

func (s *pullService) ListEvents(ctx context.Context, repo string, index int, opts scm.ListOptions) ([]*scm.ListedIssueEvent, *scm.Response, error) {
	// path := fmt.Sprintf("api/v4/projects/%s/merge_requests/%d/resource_label_events?%s", encode(repo), index, encodeListOptions(opts))
	// out := []*labelEvent{}
	// res, err := s.client.do(ctx, "GET", path, nil, &out)
	// return convertLabelEvents(out), res, err
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) AddLabel(ctx context.Context, repo string, number int, label string) (*scm.Response, error) {
	return nil, nil
}

func (s *pullService) DeleteLabel(ctx context.Context, repo string, number int, label string) (*scm.Response, error) {
	return nil, nil
}

// func (s *pullService) setLabels(ctx context.Context, repo string, number int, labelsStr string, operation string) (*scm.Response, error) {
// 	// in := url.Values{}
// 	// in.Set(operation, labelsStr)
// 	// path := fmt.Sprintf("api/v4/projects/%s/merge_requests/%d?%s", encode(repo), number, in.Encode())

// 	// return s.client.do(ctx, "PUT", path, nil, nil)
// 	return nil, scm.ErrNotSupported
// }

func (s *pullService) CreateComment(ctx context.Context, repo string, index int, input *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) DeleteComment(ctx context.Context, repo string, index, id int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *pullService) EditComment(ctx context.Context, repo string, number int, id int, input *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) Merge(ctx context.Context, repo string, number int, options *scm.PullRequestMergeOptions) (*scm.Response, error) {
	repoInfo, _, err := s.client.Repositories.Find(ctx, repo)
	if err != nil {
		return nil, err
	}

	in := new(mergeRequest)
	in.Action = "ModifyMergeMR"
	did, err := strconv.Atoi(repoInfo.ID)
	if err != nil {
		return nil, err
	}
	in.DepotId = did
	in.MergeId = number
	in.IsDelSourceBranch = options.DeleteSourceBranch
	in.IsFastForward = false
	in.Message = options.CommitTitle
	in.Squash = false

	out := new(mergeResponse)

	res, err := s.client.do(ctx, "POST", "", in, out)

	if err != nil {
		return nil, err
	}

	if out.Response.Error != nil {
		return nil, errors.New(out.Response.Error.Message)
	}

	return res, err
}

func (s *pullService) Close(ctx context.Context, repo string, number int) (*scm.Response, error) {

	repoInfo, _, err := s.client.Repositories.Find(ctx, repo)
	if err != nil {
		return nil, err
	}

	in := new(closeMergeRequest)
	in.Action = "ModifyCloseMR"
	did, err := strconv.Atoi(repoInfo.ID)
	if err != nil {
		return nil, err
	}
	in.DepotId = did
	in.MergeId = number

	out := new(mergeResponse)

	res, err := s.client.do(ctx, "POST", "", in, out)

	if err != nil {
		return nil, err
	}

	if out.Response.Error != nil {
		return nil, errors.New(out.Response.Error.Message)
	}

	return res, err
}

func (s *pullService) Reopen(ctx context.Context, repo string, number int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *pullService) AssignIssue(ctx context.Context, repo string, number int, logins []string) (*scm.Response, error) {
	// pr, _, err := s.Find(ctx, repo, number)
	// if err != nil {
	// 	return nil, err
	// }

	// allAssignees := map[int]struct{}{}
	// for _, assignee := range pr.Assignees {
	// 	allAssignees[assignee.ID] = struct{}{}
	// }
	// for _, l := range logins {
	// 	u, _, err := s.client.Users.FindLogin(ctx, l)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	allAssignees[u.ID] = struct{}{}
	// }

	// var assigneeIDs []int
	// for i := range allAssignees {
	// 	assigneeIDs = append(assigneeIDs, i)
	// }

	// return s.setAssignees(ctx, repo, number, assigneeIDs)
	return nil, scm.ErrNotFound
}

// func (s *pullService) setAssignees(ctx context.Context, repo string, number int, ids []int) (*scm.Response, error) {
// 	if len(ids) == 0 {
// 		ids = append(ids, 0)
// 	}
// 	in := &updateMergeRequestOptions{
// 		AssigneeIDs: ids,
// 	}
// 	path := fmt.Sprintf("api/v4/projects/%s/merge_requests/%d", encode(repo), number)

// 	return s.client.do(ctx, "PUT", path, in, nil)
// }

func (s *pullService) UnassignIssue(ctx context.Context, repo string, number int, logins []string) (*scm.Response, error) {
	// pr, _, err := s.Find(ctx, repo, number)
	// if err != nil {
	// 	return nil, err
	// }
	// var assignees []int
	// for _, assignee := range pr.Assignees {
	// 	shouldKeep := true
	// 	for _, l := range logins {
	// 		if assignee.Login == l {
	// 			shouldKeep = false
	// 		}
	// 	}
	// 	if shouldKeep {
	// 		assignees = append(assignees, assignee.ID)
	// 	}
	// }

	// return s.setAssignees(ctx, repo, number, assignees)
	return nil, scm.ErrNotFound
}

func (s *pullService) RequestReview(ctx context.Context, repo string, number int, logins []string) (*scm.Response, error) {
	return nil, scm.ErrNotFound
}

func (s *pullService) UnrequestReview(ctx context.Context, repo string, number int, logins []string) (*scm.Response, error) {
	return nil, scm.ErrNotFound
}

func (s *pullService) Create(ctx context.Context, repo string, input *scm.PullRequestInput) (*scm.PullRequest, *scm.Response, error) {

	rep, _, err := s.client.Repositories.Find(ctx, repo)
	if err != nil {
		return nil, nil, err
	}

	in := new(createMergeRequest)
	in.Action = "CreateGitMergeReq"
	did, _ := strconv.Atoi(rep.ID)
	in.DepotId = did
	in.Title = input.Title
	in.Content = input.Body
	in.SrcBranch = input.Head
	in.DestBranch = input.Base

	out := new(createMergeRequestResponse)
	res, err := s.client.do(ctx, "POST", "", in, out)
	if err != nil {
		return nil, res, err
	}
	if out.Response.Error != nil {
		return nil, res, errors.New(out.Response.Error.Message)
	}

	prInfo := &scm.PullRequest{
		Number: out.Response.MergeInfo.MergeRequestId,
		Title:  input.Title,
		Body:   input.Body,
		Source: input.Head,
		Target: input.Base,
	}

	return prInfo, res, nil

}

func (s *pullService) Update(ctx context.Context, repo string, number int, input *scm.PullRequestInput) (*scm.PullRequest, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) SetMilestone(ctx context.Context, repo string, prID int, number int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *pullService) ClearMilestone(ctx context.Context, repo string, prID int) (*scm.Response, error) {
	// zeroVal := 0
	// updateOpts := &updateMergeRequestOptions{
	// 	MilestoneID: &zeroVal,
	// }
	// _, res, err := s.updateMergeRequestField(ctx, repo, prID, updateOpts)
	// return res, err
	return nil, scm.ErrNotSupported
}

// type updateMergeRequestOptions struct {
// 	Title              *string `json:"title,omitempty"`
// 	Description        *string `json:"description,omitempty"`
// 	TargetBranch       *string `json:"target_branch,omitempty"`
// 	AssigneeID         *int    `json:"assignee_id,omitempty"`
// 	AssigneeIDs        []int   `json:"assignee_ids,omitempty"`
// 	Labels             *string `json:"labels,omitempty"`
// 	MilestoneID        *int    `json:"milestone_id,omitempty"`
// 	StateEvent         *string `json:"state_event,omitempty"`
// 	RemoveSourceBranch *bool   `json:"remove_source_branch,omitempty"`
// 	Squash             *bool   `json:"squash,omitempty"`
// 	DiscussionLocked   *bool   `json:"discussion_locked,omitempty"`
// 	AllowCollaboration *bool   `json:"allow_collaboration,omitempty"`
// }

// func (s *pullService) updateMergeRequestField(ctx context.Context, repo string, number int, input *updateMergeRequestOptions) (*scm.PullRequest, *scm.Response, error) {
// 	path := fmt.Sprintf("api/v4/projects/%s/merge_requests/%d", encode(repo), number)

// 	out := new(pr)
// 	res, err := s.client.do(ctx, "PUT", path, input, out)
// 	if err != nil {
// 		return nil, res, err
// 	}
// 	convRepo, convRes, err := s.convertPullRequest(ctx, out)
// 	if err != nil {
// 		return nil, convRes, err
// 	}
// 	return convRepo, res, nil
// }

// type pr struct {
// 	Number          int       `json:"iid"`
// 	Sha             string    `json:"sha"`
// 	Title           string    `json:"title"`
// 	Desc            string    `json:"description"`
// 	State           string    `json:"state"`
// 	SourceProjectID int       `json:"source_project_id"`
// 	TargetProjectID int       `json:"target_project_id"`
// 	Labels          []*string `json:"labels"`
// 	Link            string    `json:"web_url"`
// 	WIP             bool      `json:"work_in_progress"`
// 	Author          user      `json:"author"`
// 	MergeStatus     string    `json:"merge_status"`
// 	SourceBranch    string    `json:"source_branch"`
// 	TargetBranch    string    `json:"target_branch"`
// 	Created         time.Time `json:"created_at"`
// 	Updated         time.Time `json:"updated_at"`
// 	Closed          time.Time
// 	DiffRefs        struct {
// 		BaseSHA string `json:"base_sha"`
// 		HeadSHA string `json:"head_sha"`
// 	} `json:"diff_refs"`
// 	Assignee  *user   `json:"assignee"`
// 	Assignees []*user `json:"assignees"`
// }

// type changes struct {
// 	Changes []*change
// }

// type change struct {
// 	OldPath string `json:"old_path"`
// 	NewPath string `json:"new_path"`
// 	Added   bool   `json:"new_file"`
// 	Renamed bool   `json:"renamed_file"`
// 	Deleted bool   `json:"deleted_file"`
// 	Diff    string `json:"diff"`
// }

// type prInput struct {
// 	Title        string `json:"title"`
// 	Description  string `json:"description"`
// 	SourceBranch string `json:"source_branch"`
// 	TargetBranch string `json:"target_branch"`
// }

// type pullRequestMergeRequest struct {
// 	CommitMessage             string `json:"merge_commit_message,omitempty"`
// 	SquashCommitMessage       string `json:"squash_commit_message,omitempty"`
// 	Squash                    string `json:"squash,omitempty"`
// 	RemoveSourceBranch        string `json:"should_remove_source_branch,omitempty"`
// 	SHA                       string `json:"sha,omitempty"`
// 	MergeWhenPipelineSucceeds string `json:"merge_when_pipeline_succeeds,omitempty"`
// }

func (s *pullService) convertPullRequestList(ctx context.Context, from []*mergeRequestInfo) ([]*scm.PullRequest, *scm.Response, error) {
	to := []*scm.PullRequest{}
	for _, v := range from {
		converted, res, err := s.convertPullRequest(ctx, v)
		if err != nil {
			return nil, res, err
		}
		to = append(to, converted)
	}
	return to, nil, nil
}

func (s *pullService) convertPullRequest(ctx context.Context, from *mergeRequestInfo) (*scm.PullRequest, *scm.Response, error) {

	return &scm.PullRequest{
		Title: from.Title,
		Body:  from.Describe,
		State: codingStateToSCMState(from.Status),
		// Labels:         convertPullRequestLabels(from.Labels),
		// Sha:            from.Sha,
		// Ref:            fmt.Sprintf("refs/merge-requests/%d/head", from.Number),
		Source: from.SourceBranch,
		Target: from.TargetBranch,
		// Link:           from.Link,
		// Draft:          from.WIP,
		// Closed:         from.State != "opened",
		// Merged:         from.State == "merged",
		// Mergeable:      scm.ToMergeableState(from.MergeStatus) == scm.MergeableStateMergeable,
		// MergeableState: scm.ToMergeableState(from.MergeStatus),
		// Author:         *convertUser(&from.Author),
		// Assignees:      assignees,
		Head: scm.PullRequestBranch{
			Ref: from.SourceBranch,
			// Sha:  headSHA,
			// Repo: *headRepo,
		},
		Base: scm.PullRequestBranch{
			Ref: from.TargetBranch,
			// Sha:  from.DiffRefs.BaseSHA,
			// Repo: *baseRepo,
		},
		// Created: from.Created,
		// Updated: from.Updated,
		// Fork:    sourceRepo.PathNamespace,
	}, nil, nil
}

// func (s *pullService) getSourceFork(ctx context.Context, from *pr) (repository, error) {
// 	path := fmt.Sprintf("api/v4/projects/%d", from.SourceProjectID)
// 	sourceRepo := repository{}
// 	_, err := s.client.do(ctx, "GET", path, nil, &sourceRepo)
// 	if err != nil {
// 		return repository{}, err
// 	}
// 	return sourceRepo, nil
// }

// func convertPullRequestLabels(from []*string) []*scm.Label {
// 	var labels []*scm.Label
// 	for _, label := range from {
// 		l := *label
// 		labels = append(labels, &scm.Label{
// 			Name: l,
// 		})
// 	}
// 	return labels
// }

// func convertChangeList(from []*change) []*scm.Change {
// 	to := []*scm.Change{}
// 	for _, v := range from {
// 		to = append(to, convertChange(v))
// 	}
// 	return to
// }

// func convertChange(from *change) *scm.Change {
// 	to := &scm.Change{
// 		Path:         from.NewPath,
// 		PreviousPath: from.OldPath,
// 		Added:        from.Added,
// 		Deleted:      from.Deleted,
// 		Renamed:      from.Renamed,
// 		Patch:        from.Diff,
// 	}
// 	if to.Path == "" {
// 		to.Path = from.OldPath
// 	}
// 	return to
// }

func codingStateToSCMState(glState string) string {

	// 	CANMERGE	状态可自动合并
	// ACCEPTED	状态已接受
	// CANNOTMERGE	状态不可自动合并
	// REFUSED	状态已拒绝(关闭)
	// CANCEL	取消
	// MERGING	正在合并中
	// ABNORMAL	状态异常

	switch glState {
	case "CANMERGE":
		return "mergeable"
	case "CANNOTMERGE":
		return "conflict"
	case "MERGING":
		return "cannot_be_merged"
	case "CANCEL":
	case "REFUSED":
		return "closed"
	default:
		return "closed"
	}

	return "closed"
}
