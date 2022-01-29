// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coding

import (
	"crypto/subtle"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jenkins-x/go-scm/scm"
)

type webhookService struct {
	client *wrapper
}

func (s *webhookService) Parse(req *http.Request, fn scm.SecretFunc) (scm.Webhook, error) {
	data, err := ioutil.ReadAll(
		io.LimitReader(req.Body, 10000000),
	)
	if err != nil {
		return nil, err
	}

	var hook scm.Webhook
	event := req.Header.Get("X-Coding-Service-Hook-Event")
	hookId := req.Header.Get("X-Coding-Service-Hook-Id")
	// hookAction := req.Header.Get("X-Coding-Service-Hook-Action")
	// delivery := req.Header.Get("X-Coding-Delivery")

	switch event {
	case "GIT_PUSHED":
		hook, err = parsePushHook(data, hookId)
	case "Issue Hook":
		return nil, scm.UnknownWebhook{Event: event}
	case "GIT_MR_CREATED":
	case "GIT_MR_UPDATED":
	case "GIT_MR_MERGED":
	case "GIT_MR_CLOSED":
		hook, err = parsePullRequestHook(data, event, hookId)
	default:
		return nil, scm.UnknownWebhook{Event: event}
	}
	if err != nil {
		return nil, err
	}

	// get the gitlab shared token to verify the payload
	// authenticity. If no key is provided, no validation
	// is performed.
	token, err := fn(hook)
	if err != nil {
		return hook, err
	} else if token == "" {
		return hook, nil
	}

	if subtle.ConstantTimeCompare([]byte(req.Header.Get("X-Gitlab-Token")), []byte(token)) == 0 {
		return hook, scm.ErrSignatureInvalid
	}

	return hook, nil
}

func parsePushHook(data []byte, hookId string) (scm.Webhook, error) {
	src := new(pushHook)
	err := json.Unmarshal(data, src)
	if err != nil {
		return nil, err
	}

	hook := convertPushHook(src)
	hook.GUID = hookId

	return hook, nil
}

func parsePullRequestHook(data []byte, action string, hookId string) (scm.Webhook, error) {
	src := new(pullRequestHook)
	err := json.Unmarshal(data, src)
	if err != nil {
		return nil, err
	}

	hook := convertPullRequestHook(src, action)
	hook.GUID = hookId

	return hook, nil

}

func convertPushHook(src *pushHook) *scm.PushHook {
	repo := *convertRepositoryHook(&src.Repository, &src.Project)
	dst := &scm.PushHook{
		Ref:     src.Ref,
		Repo:    repo,
		Before:  src.Before,
		After:   src.After,
		Compare: src.Compare,
		Commit: scm.Commit{
			Sha:     src.HeadCommit.Id,
			Message: src.HeadCommit.Message, // NOTE this is set below
			Author: scm.Signature{
				Login: src.HeadCommit.Author.Username,
				Name:  src.HeadCommit.Author.Name,
				Email: src.HeadCommit.Author.Email,
			},
			Committer: scm.Signature{
				Login: src.HeadCommit.Commiter.Username,
				Name:  src.HeadCommit.Commiter.Name,
				Email: src.HeadCommit.Commiter.Email,
			},
			Link: src.HeadCommit.Url, // NOTE this is set below
		},
		Sender: scm.User{
			Login:  src.Sender.Login,
			Name:   src.Sender.Name,
			Avatar: src.Sender.AvatarUrl,
		},
	}
	// if len(src.Commits) > 0 {
	// 	// get the last commit (most recent)
	// 	dst.Commit.Message = src.Commits[len(src.Commits)-1].Message
	// 	dst.Commit.Link = src.Commits[len(src.Commits)-1].Url
	// }
	return dst
}

func convertPullRequestHook(src *pullRequestHook, action string) *scm.PullRequestHook {
	mr := src.MergeRequest

	pr := scm.PullRequest{
		Number: mr.Number,
		Title:  mr.Title,
		Body:   mr.Body,
		State:  codingStateToSCMState(mr.State),
		//	Sha:    sha,
		//	Ref:    ,
		Base: scm.PullRequestBranch{
			Ref:  mr.Base.Ref,
			Sha:  mr.Base.Sha,
			Repo: *convertRepositoryHook(&mr.Base.Repo, &src.Project),
		},
		Head: scm.PullRequestBranch{
			Ref:  mr.Head.Ref,
			Sha:  mr.Head.Sha,
			Repo: *convertRepositoryHook(&mr.Head.Repo, &src.Project),
		},
		Source: mr.Head.Repo.DefaultBranch,
		Target: mr.Base.Repo.DefaultBranch,
		// Fork:     ,
		Link:     mr.HtmlUrl,
		Merged:   mr.Merged,
		MergeSha: mr.MergeCommitSha,
		// Created   : src.ObjectAttributes.CreatedAt,
		// Updated  : src.ObjectAttributes.UpdatedAt, // 2017-12-10 17:01:11 UTC
		Author: scm.User{
			Login:  mr.User.Login,
			Name:   mr.User.Name,
			Email:  "", // TODO how do we get the pull request author email?
			Avatar: mr.User.AvatarUrl,
		},
	}

	pr.Closed = pr.State == "closed"

	return &scm.PullRequestHook{
		Action:      convertAction(action),
		PullRequest: pr,
		Repo:        *convertRepositoryHook(&src.Repository, &src.Project),
		Sender: scm.User{
			Login:  src.Sender.Login,
			Name:   src.Sender.Name,
			Email:  "", // TODO how do we get the pull request author email?
			Avatar: src.Sender.AvatarUrl,
		},
	}
}

func convertRepositoryHook(from *repository, project *hookProject) *scm.Repository {
	names := strings.Split(from.FullName, "/")
	name := names[len(names)-1]
	if len(names) > 2 {
		name = strings.Join(names[1:], "/")
	}
	return &scm.Repository{
		ID:        strconv.Itoa(int(from.Id)),
		Namespace: names[0],
		Name:      name,
		FullName:  from.FullName,
		Clone:     from.CloneUrl,
		CloneSSH:  from.SshUrl,
		Link:      from.HtmlUrl,
		Branch:    from.DefaultBranch,
		Created:   time.Unix(from.CreateAt/1000, 0),
		Updated:   time.Unix(from.UpdatedAt/1000, 0),
		Private:   from.Private, // TODO how do we correctly set Private vs Public?
	}
}

func convertAction(src string) (action scm.Action) {
	switch src {
	case "GIT_MR_CREATED":
		return scm.ActionCreate
	case "GIT_MR_UPDATED":
		return scm.ActionUpdate
	case "GIT_MR_MERGED":
		return scm.ActionMerge
	case "GIT_MR_CLOSED":
		return scm.ActionClose
	default:
		return
	}
}

type (
	userInfo struct {
		Id         int64  `json:"id"`
		Login      string `json:"login"`
		AvatarUrl  string `json:"avatar_url"`
		Url        string `json:"url"`
		HtmlUrl    string `json:"html_url"`
		Name       string `json:"name"`
		NamePinYin string `json:"name_pinyin"`
	}

	repository struct {
		Id            int64    `json:"id"`
		Name          string   `json:"name"`
		FullName      string   `json:"full_name"` //用户/项目/仓库名
		Owner         userInfo `json:"owner"`
		Private       bool     `json:"private"`
		HtmlUrl       string   `json:"html_url"`
		Description   string   `json:"description"`
		Fork          bool     `json:"fork"`
		CreateAt      int64    `json:"created_at"`
		UpdatedAt     int64    `json:"updated_at"`
		CloneUrl      string   `json:"clone_url"`
		SshUrl        string   `json:"ssh_url"`
		DefaultBranch string   `json:"default_branch"`
		VcsType       string   `json:"vcs_type"`
	}

	hookProject struct {
		Id          int64  `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Url         string `json:"url"`
	}

	hookTeam struct {
		Id           int64  `json:"id"`
		Domain       string `json:"domain"`
		Name         string `json:"name"`
		NamePinYin   string `json:"name_pinyin"`
		Introduction string `json:"introduction"`
		Avatar       string `json:"avatar"`
		Url          string `json:"url"`
	}

	simpleUserInfo struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	commit struct {
		Id        string         `json:"id"`
		Message   string         `json:"message"`
		Timestamp int64          `json:"timestamp"`
		Url       string         `json:"url"`
		Author    simpleUserInfo `json:"author"`
		Commiter  simpleUserInfo `json:"committer"`
		Added     []string
		Removed   []string
		Modified  []string
	}

	pushHook struct {
		Event      string         `json:"event"`
		EventName  string         `json:"eventName"` //事件中文名
		Ref        string         `json:"ref"`
		Before     string         `json:"before"`  //推送之前旧的 sha 值
		After      string         `json:"after"`   //推送之后新的 sha 值
		Created    bool           `json:"created"` //是否新增
		Deleted    bool           `json:"deleted"` //是否删除
		Compare    string         `json:"compare"` //对比地址
		Commits    []commit       `json:"commits"`
		HeadCommit commit         `json:"head_commit"`
		Pusher     simpleUserInfo `json:"pusher"`
		Repository repository     `json:"repository"`
		Sender     userInfo       `json:"sender"`
		Project    hookProject    `json:"project"`
		Team       hookTeam       `json:"team"`
	}

	hookBranch struct {
		Ref  string     `json:"ref"`
		Sha  string     `json:"sha"`
		User userInfo   `json:"user"`
		Repo repository `json:"repo"`
	}

	hookMergeRequest struct {
		Id             int        `json:"id"`
		HtmlUrl        string     `json:"html_url"`
		PatchUrl       string     `json:"patch_url"`
		DiffUrl        string     `json:"diff_url"`
		Number         int        `json:"number"`
		State          string     `json:"state"`
		Title          string     `json:"title"`
		Body           string     `json:"body"`
		User           userInfo   `json:"user"`
		CreateAt       int64      `json:"created_at"`
		UpdateAt       int64      `json:"updated_at"`
		MergeCommitSha string     `json:"merge_commit_sha"`
		Merged         bool       `json:"merged"`
		Comments       int        `json:"comments"`
		Commits        int        `json:"commits"`
		Additions      int        `json:"additions"`
		Deletions      int        `json:"deletions"`
		ChangedFiles   int        `json:"changed_files"`
		Head           hookBranch `json:"head"`
		Base           hookBranch `json:"base"`
	}

	pullRequestHook struct {
		Event        string           `json:"event"`
		EventName    string           `json:"eventName"` //事件中文名
		MergeRequest hookMergeRequest `json:"mergeRequest"`
		Repository   repository       `json:"repository"`
		Sender       userInfo         `json:"sender"`
		Project      hookProject      `json:"project"`
		Team         hookTeam         `json:"team"`
	}
)
