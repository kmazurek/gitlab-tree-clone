package infra

import "github.com/xanzy/go-gitlab"

type GitlabClient struct {
	client *gitlab.Client
}

func NewGitlabClient(apiToken string) (*GitlabClient, error) {
	client, err := gitlab.NewClient(apiToken)
	if err != nil {
		return nil, err
	}
	return &GitlabClient{client: client}, nil
}

func (gc *GitlabClient) GetGroup(groupId int) (*gitlab.Group, error) {
	group, _, err := gc.client.Groups.GetGroup(groupId, nil)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (gc *GitlabClient) ListSubgroups(groupId int) ([]*gitlab.Group, error) {
	opt := &gitlab.ListDescendantGroupsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 20,
			Page:    1,
		},
	}

	var result []*gitlab.Group

	for {
		groups, resp, err := gc.client.Groups.ListDescendantGroups(groupId, opt)
		if err != nil {
			return nil, err
		}
		result = append(result, groups...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return result, nil
}

func (gc *GitlabClient) ListProjects(groupId int) ([]*gitlab.Project, error) {
	opt := &gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 20,
			Page:    1,
		},
	}

	var result []*gitlab.Project

	for {
		projects, resp, err := gc.client.Groups.ListGroupProjects(groupId, opt)
		if err != nil {
			return nil, err
		}
		result = append(result, projects...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return result, nil
}
