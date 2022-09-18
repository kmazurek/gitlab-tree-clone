package app

import (
	"errors"
	"io/fs"
	"log"
	"os"

	"github.com/zakaprov/gitlab-group-clone/infra"
	"golang.org/x/sync/errgroup"
)

type TreeClone struct {
	ErrGroup     *errgroup.Group
	GitClient    *infra.GitClient
	GitlabClient *infra.GitlabClient
}

func (tc *TreeClone) CloneGroup(groupID int, groupName string, path string) error {
	log.Println("Cloning group: " + groupName + " to path: " + path)
	path = path + "/" + groupName
	err := os.MkdirAll(path, 0755)
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
		path := path
		invalid := project.Archived || project.EmptyRepo
		if !invalid {
			tc.ErrGroup.Go(func() error {
				_, err := os.Stat(path + "/" + project.Name)
				if err != nil {
					if errors.Is(err, fs.ErrNotExist) {
						return tc.GitClient.CloneProject(path, project.Name, project.HTTPURLToRepo)

					}
					return err
				}

				return tc.GitClient.PullProject(path, project.Name)
			})
		}
	}

	return nil
}
