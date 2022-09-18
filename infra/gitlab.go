package infra

import "github.com/xanzy/go-gitlab"

func GetGroup(client *gitlab.Client, groupId int) (*gitlab.Group, error) {
	group, _, err := client.Groups.GetGroup(groupId, nil)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func ListSubgroups(client *gitlab.Client, groupId int) ([]*gitlab.Group, error) {
	opt := &gitlab.ListDescendantGroupsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 20,
			Page:    1,
		},
	}

	var result []*gitlab.Group

	for {
		groups, resp, err := client.Groups.ListDescendantGroups(groupId, opt)
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

func ListProjects(client *gitlab.Client, groupId int) ([]*gitlab.Project, error) {
	opt := &gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 20,
			Page:    1,
		},
	}

	var result []*gitlab.Project

	for {
		projects, resp, err := client.Groups.ListGroupProjects(groupId, opt)
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
