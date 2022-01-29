// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package coding implements a Coding client.
package coding

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"code.gitea.io/sdk/gitea"
	"github.com/jenkins-x/go-scm/scm"
)

type apiRequest struct {
	Action string `json:"Action"`
}

type apiResponse struct {
	RequestId string `json:"RequestId"`
	Error     *struct {
		Message string `json:"Message"`
		Code    string `json:"Code"`
	} `json:"Error"`
}

// wraper wraps the Client to provide high level helper functions
// for making http requests and unmarshaling the response.
type wrapper struct {
	*scm.Client
	token string
}

// NewWebHookService creates a new instance of the webhook service without the rest of the client
// func NewWebHookService() scm.WebhookService {
// 	return &webhookService{nil}
// }

// New returns a new Gitea API client without a token set
func New(uri string) (*scm.Client, error) {
	return NewWithToken(uri, "")
}

// NewDefault return a new Coding API Client
func NewDefault() *scm.Client {
	client, _ := New("https://e.coding.net/open-api")
	return client
}

// NewDefaultWithToken
func NewDefaultWithToken(token string) *scm.Client {
	client, _ := NewWithToken("https://e.coding.net/open-api", token)
	return client
}

// NewWithToken returns a new Coding API client with the token set.
func NewWithToken(uri string, token string) (*scm.Client, error) {
	base, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(base.Path, "open-api") && !strings.HasSuffix(base.Path, "open-api/") {
		if !strings.HasSuffix(base.Path, "/") {
			base.Path = base.Path + "/"
		}
		base.Path = base.Path + "open-api"
	}

	client := &wrapper{Client: new(scm.Client)}
	client.token = token

	if err != nil {
		return nil, err
	}
	client.BaseURL = base
	// initialize services
	client.Driver = scm.DriverCoding
	// client.Contents = &contentService{client}
	// client.Git = &gitService{client}
	// client.Issues = &issueService{client}
	// client.Milestones = &milestoneService{client}
	client.Organizations = &organizationService{client}
	client.PullRequests = &pullService{client}
	client.Repositories = &repositoryService{client}
	// client.Reviews = &reviewService{client}
	// client.Releases = &releaseService{client}
	client.Users = &userService{client}
	client.Webhooks = &webhookService{client}
	return client.Client, nil
}

// do wraps the Client.Do function by creating the Request and
// unmarshalling the response.
func (c *wrapper) do(ctx context.Context, method, path string, in, out interface{}) (*scm.Response, error) {
	req := &scm.Request{
		Method: method,
		Path:   path,
	}
	// if we are posting or putting data, we need to
	// write it to the body of the request.
	if in != nil {
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(in)
		req.Header = map[string][]string{
			"Content-Type": {"application/json"},
		}
		req.Body = buf
	}

	// add token header
	if c.token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("token %s", c.token))
	}

	// execute the http request
	res, err := c.Client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// if an error is encountered, unmarshal and return the
	// error response.
	if res.Status > 300 {
		return res, errors.New(
			http.StatusText(res.Status),
		)
	}

	if out == nil {
		return res, nil
	}

	// if raw output is expected, copy to the provided
	// buffer and exit.
	if w, ok := out.(io.Writer); ok {
		io.Copy(w, res.Body)
		return res, nil
	}

	// if a json response is expected, parse and return
	// the json response.
	return res, json.NewDecoder(res.Body).Decode(out)
}

// toSCMResponse creates a new Response for the provided
// http.Response. r must not be nil.
func toSCMResponse(r *gitea.Response) *scm.Response {
	if r == nil {
		return nil
	}
	res := &scm.Response{
		Status: r.StatusCode,
		Header: r.Header,
		Body:   r.Body,
	}
	res.PopulatePageValues()
	return res
}

func getInfoFromGitUrl(gitUrl string) (string, string, string, string, error) {
	urlParts, err := url.Parse(gitUrl)
	if err != nil {
		return "", "", "", "", err
	}

	urlParts.Path = strings.TrimSuffix(urlParts.Path, ".git")
	urlParts.Path = strings.TrimSuffix(urlParts.Path, "/")
	urlParts.Path = strings.TrimPrefix(urlParts.Path, "/")

	names := strings.Split(urlParts.Path, "/")

	owner := ""
	project := ""
	name := ""
	fullName := urlParts.Path
	if len(names) == 3 {
		owner = names[0]
		project = names[1]
		name = scm.Join(project, names[2])
	} else if len(names) == 2 {
		owner = names[0]
		project = names[0]
		name = scm.Join(project, names[2])
	} else if len(names) == 1 {
		owner = names[0]
		project = names[0]
		name = names[0]
	} else {
		return "", "", "", "", errors.New("invalid names")
	}
	return owner, project, name, fullName, nil
}

// func toGiteaListOptions(in scm.ListOptions) gitea.ListOptions {
// 	return gitea.ListOptions{
// 		Page:     in.Page,
// 		PageSize: in.Size,
// 	}
// }
