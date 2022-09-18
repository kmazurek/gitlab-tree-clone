package app

import (
	"log"
	"os"

	"github.com/xanzy/go-gitlab"
	"github.com/zakaprov/gitlab-group-clone/infra"
	"golang.org/x/sync/errgroup"
)

func CloneGroup(client *gitlab.Client, errGroup *errgroup.Group, groupID int, groupName string, path string) error {
	log.Println("Cloning group: " + groupName + " to path: " + path)
	path = path + "/" + groupName
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	subgroups, err := infra.ListSubgroups(client, groupID)
	if err != nil {
		return err
	}

	for _, subgroup := range subgroups {
		subgroup := subgroup
		errGroup.Go(func() error {
			return CloneGroup(client, errGroup, subgroup.ID, subgroup.Name, path)
		})
	}

	projects, err := infra.ListProjects(client, groupID)
	if err != nil {
		return err
	}
	for _, project := range projects {
		invalid := project.Archived || project.EmptyRepo
		if !invalid {
			project := project
			errGroup.Go(func() error {
				return infra.CloneProject(path, project.Name, project.SSHURLToRepo)
			})
		}
	}

	return nil
}
