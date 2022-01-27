// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coding

import (
	"context"
	"errors"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/driver/internal/null"
)

type userService struct {
	client *wrapper
}

type findUserResponse struct {
	Response struct {
		apiResponse
		User user `json:"User"`
	}
}

func (s *userService) CreateToken(context.Context, string, string) (*scm.UserToken, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *userService) DeleteToken(context.Context, int64) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *userService) Find(ctx context.Context) (*scm.User, *scm.Response, error) {
	out := new(findUserResponse)
	body := apiRequest{
		Action: "DescribeCodingCurrentUser",
	}
	res, err := s.client.do(ctx, "POST", "", &body, out)
	if out.Response.Error != nil {
		return nil, res, errors.New(out.Response.Error.Message)
	}
	res.ID = out.Response.RequestId
	return convertUser(&out.Response.User), res, err
}

func (s *userService) FindLogin(ctx context.Context, login string) (*scm.User, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *userService) FindEmail(ctx context.Context) (string, *scm.Response, error) {
	user, res, err := s.Find(ctx)
	return user.Email, res, err
}

func (s *userService) ListInvitations(ctx context.Context) ([]*scm.Invitation, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *userService) AcceptInvitation(ctx context.Context, invitationID int64) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

type user struct {
	ID              int         `json:"Id"`
	Status          int         `json:"Status`
	Email           null.String `json:"Email"`
	GlobalKey       string      `json:"GlobalKey"`
	Avatar          string      `json:"Avatar"`
	Name            string      `json:"Name"`
	NamePinYin      string      `json:"NamePinYin"`
	Phone           null.String `json:"Phone"`
	PhoneValidation int         `json:"PhoneValidation"`
	EmailValidation int         `json:"EmailValidation"`
	PhoneRegionCode null.String `json:"PhoneRegionCode"`
	TeamId          int         `json:"TeamId"`
}

// type repositoryInvitation struct {
// 	ID      int64      `json:"id,omitempty"`
// 	Repo    repository `json:"repository,omitempty"`
// 	Invitee user       `json:"invitee,omitempty"`
// 	Inviter user       `json:"inviter,omitempty"`

// 	Permissions string    `json:"permissions,omitempty"`
// 	CreatedAt   time.Time `json:"created_at,omitempty"`
// 	URL         string    `json:"url,omitempty"`
// 	HTMLURL     string    `json:"html_url,omitempty"`
// }

// func convertRepositoryInvitation(from *repositoryInvitation) *scm.Invitation {
// 	return &scm.Invitation{
// 		ID:          from.ID,
// 		Repo:        convertRepository(&from.Repo),
// 		Invitee:     convertUser(&from.Invitee),
// 		Inviter:     convertUser(&from.Inviter),
// 		Permissions: from.Permissions,
// 		Link:        from.URL,
// 		Created:     from.CreatedAt,
// 	}
// }

func convertUser(from *user) *scm.User {
	if from == nil {
		return nil
	}
	return &scm.User{
		ID:     from.ID,
		Avatar: from.Avatar,
		Email:  from.Email.String,
		Login:  from.NamePinYin,
		Name:   from.Name,
	}
}
