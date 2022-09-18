package app

import (
	"log"

	"github.com/chigopher/pathlib"
	"github.com/zakaprov/gitlab-group-clone/infra"
	"golang.org/x/sync/errgroup"
)

type TreeClone struct {
	ErrGroup     *errgroup.Group
	GitClient    *infra.GitClient
	GitlabClient *infra.GitlabClient
}

func (tc *TreeClone) CloneGroup(groupID int, groupName string, path *pathlib.Path) error {
	log.Println("Cloning group: " + groupName + " to path: " + path.String())
	groupPath := path.Join(groupName)
	err := groupPath.MkdirAll()
	if err != nil {
		return err
	}

	subgroups, err := tc.GitlabClient.ListSubgroups(groupID)
	if err != nil {
		return err
	}

	for _, subgroup := range subgroups {
		subgroup := subgroup
		tc.ErrGroup.Go(func() error {
			return tc.CloneGroup(subgroup.ID, subgroup.Name, path)
		})
	}

	projects, err := tc.GitlabClient.ListProjects(groupID)
	if err != nil {
		return err
	}
	for _, project := range projects {
		project := project
		invalid := project.Archived || project.EmptyRepo
		if !invalid {
			tc.ErrGroup.Go(func() error {
				return tc.GitClient.CloneProject(groupPath.Join(project.Name), project.HTTPURLToRepo)
			})
		}
	}

	return nil
}
