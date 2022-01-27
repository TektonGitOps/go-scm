package coding

import (
	"context"
	"errors"

	"github.com/jenkins-x/go-scm/scm"
)

type organizationService struct {
	client *wrapper
}

type findProjectRequest struct {
	apiRequest
	ProjectName string `json:"ProjectName"`
}

type findProjectResponse struct {
	Response *struct {
		apiResponse
		Project projectItem `json:"Project"`
	}
}

type listProjectsRequest struct {
	apiRequest
	UserId int `json:"userId"`
}

type projectItem struct {
	Name        string `json:"Name"`
	Id          int    `json:"Id"`
	Type        int    `json:"Type"`
	DisplayName string `json:"DisplayName"`
	Icon        string `json:"Icon"`
	Description string `json:"Description"`
	CreatedAt   int64  `json:"CreatedAt"`
	MaxMember   int    `json:"MaxMember"`
	TeamId      int    `json:"TeamId"`
	UserOwnerId int    `json:"UserOwnerId"`
	IsDemo      bool   `json:"IsDemo"`
	Archived    bool   `json:"Archived"`
	StartDate   int64  `json:"StartDate"`
	UpdatedAt   int64  `json:"UpdatedAt"`
	TeamOwnerId int    `json:"TeamOwnerId"`
	EndDate     int64  `json:"EndDate"`
	Status      int    `json:"Status"`
}

type projectListResponse struct {
	Response *struct {
		apiResponse
		ProjectList []*projectItem `json:"ProjectList"`
	}
}

type listProjectMembersRequest struct {
	apiRequest
	PageNumber int `json:"PageNumber"`
	PageSize   int `json:"PageSize"`
	ProjectId  int `json:"ProjectId"`
	//	RoleId     null.Int `json:"RoleId"`
}

type projectMemberItem struct {
	Id              int    `json:"Id"`
	TeamId          int    `json:"TeamId"`
	Name            string `json:"Name"`
	NamePinYin      string `json:"NamePinYin"`
	Avatar          string `json:"Avatar"`
	Email           string `json:"Email"`
	Phone           string `json:"Phone"`
	EmailValidation int    `json:"EmailValidation"`
	PhoneValidation int    `json:"PhoneValidation"`
	Status          int    `json:"Status"`
	GlobalKey       string `json:"GlobalKey"`
	Roles           []*struct {
		RoleType     string `json:"RoleType"`
		RoleId       int    `json:"RoleId"`
		RoleTypeName string `json:"RoleTypeName"`
	}
}

type listProjectMemberResponse struct {
	Response struct {
		apiResponse
		Data struct {
			PageNumber     int                  `json:"PageNumber"`
			PageSize       int                  `json:"PageSize"`
			TotalCount     int                  `json:"TotalCount"`
			ProjectMembers []*projectMemberItem `json:"ProjectMembers"`
		} `json:"Data"`
	} `json:"Response"`
}

func (s *organizationService) Create(context.Context, *scm.OrganizationInput) (*scm.Organization, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *organizationService) Delete(context.Context, string) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *organizationService) IsMember(ctx context.Context, org string, user string) (bool, *scm.Response, error) {
	return false, nil, scm.ErrNotSupported
}

func (s *organizationService) Find(ctx context.Context, name string) (*scm.Organization, *scm.Response, error) {
	body := findProjectRequest{
		apiRequest: apiRequest{
			Action: "DescribeProjectByName",
		},
		ProjectName: name,
	}

	out := new(findProjectResponse)
	res, err := s.client.do(ctx, "POST", "", &body, out)

	if out.Response.Error != nil {
		return nil, res, errors.New(out.Response.Error.Message)
	}

	return convertOrganization(&out.Response.Project), res, err
}

func (s *organizationService) List(ctx context.Context, opts scm.ListOptions) ([]*scm.Organization, *scm.Response, error) {
	user, res, err := s.client.Users.Find(ctx)
	if err != nil {
		return nil, nil, err
	}
	body := listProjectsRequest{
		apiRequest: apiRequest{
			Action: "DescribeUserProjects",
		},
		UserId: user.ID,
	}

	out := new(projectListResponse)
	res, err = s.client.do(ctx, "POST", "", &body, out)

	if out.Response.Error != nil {
		return nil, res, errors.New(out.Response.Error.Message)
	}

	return convertOrganizationList(out.Response.ProjectList), res, err
}

func (s *organizationService) ListOrgMembers(ctx context.Context, org string, ops scm.ListOptions) ([]*scm.TeamMember, *scm.Response, error) {
	orgRes, _, err := s.Find(ctx, org)
	if err != nil {
		return nil, nil, err
	}

	body := listProjectMembersRequest{
		apiRequest: apiRequest{
			Action: "DescribeProjectMembers",
		},
		ProjectId:  orgRes.ID,
		PageNumber: 1,
		PageSize:   1000,
	}

	out := new(listProjectMemberResponse)
	res, err := s.client.do(ctx, "POST", "", &body, out)

	if out.Response.Error != nil {
		return nil, res, errors.New(out.Response.Error.Message)
	}

	return convertTeamMembers(out.Response.Data.ProjectMembers), res, err
}

func (s *organizationService) ListTeamMembers(ctx context.Context, id int, role string, opts scm.ListOptions) ([]*scm.TeamMember, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *organizationService) ListPendingInvitations(ctx context.Context, org string, opts scm.ListOptions) ([]*scm.OrganizationPendingInvite, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *organizationService) ListMemberships(ctx context.Context, opts scm.ListOptions) ([]*scm.Membership, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func convertOrganizationList(from []*projectItem) []*scm.Organization {
	to := []*scm.Organization{}
	for _, v := range from {
		to = append(to, convertOrganization(v))
	}
	return to
}

func (s *organizationService) AcceptOrganizationInvitation(ctx context.Context, org string) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *organizationService) IsAdmin(ctx context.Context, org string, user string) (bool, *scm.Response, error) {
	return false, nil, scm.ErrNotSupported
}

func (s *organizationService) ListTeams(ctx context.Context, org string, opts scm.ListOptions) ([]*scm.Team, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func convertOrganization(from *projectItem) *scm.Organization {
	return &scm.Organization{
		ID:     from.Id,
		Name:   from.Name,
		Avatar: from.Icon,
		Permissions: scm.Permissions{
			MembersCreateInternal: false,
			MembersCreatePublic:   false,
			MembersCreatePrivate:  false,
		},
	}
}

func convertTeamMembers(from []*projectMemberItem) []*scm.TeamMember {
	to := []*scm.TeamMember{}
	for _, v := range from {
		member := convertTeamMember(v)
		if member != nil {
			to = append(to, member)
		}
	}
	return to
}

func convertTeamMember(from *projectMemberItem) *scm.TeamMember {
	if from == nil {
		return nil
	}
	return &scm.TeamMember{
		Login: from.Email,
	}
}
