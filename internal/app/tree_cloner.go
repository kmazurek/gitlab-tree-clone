package app

import (
	"log"

	"github.com/chigopher/pathlib"
	"github.com/kmazurek/gitlab-tree-clone/internal/domain"
	"github.com/kmazurek/gitlab-tree-clone/internal/util"
	"golang.org/x/sync/errgroup"
)

type TreeCloner struct {
	errGroup     *errgroup.Group
	gitClient    domain.GitClient
	gitlabClient domain.GitlabClient
	ignoreIDs    map[int]bool
	ignoreNames  map[string]bool
}

func NewTreeCloner(
	gitClient domain.GitClient,
	gitlabClient domain.GitlabClient,
	errGroup *errgroup.Group,
	ignoreIDs []int,
	ignoreNames []string) (*TreeCloner, error) {
	var idMap = make(map[int]bool)
	for _, id := range ignoreIDs {
		idMap[id] = true
	}
	var nameMap = make(map[string]bool)
	for _, name := range ignoreNames {
		nameMap[name] = true
	}

	return &TreeCloner{
		errGroup:     errGroup,
		gitClient:    gitClient,
		gitlabClient: gitlabClient,
		ignoreIDs:    idMap,
		ignoreNames:  nameMap,
	}, nil
}

func (tc *TreeCloner) CloneTree(groupID int, path *pathlib.Path) error {
	group, err := tc.gitlabClient.GetGroup(groupID)
	if err != nil {
		return err
	}
	return tc.cloneGroup(group.ID, group.Name, path)
}

func (tc *TreeCloner) cloneGroup(groupID int, groupName string, path *pathlib.Path) error {
	log.Println("Cloning group:", groupName, "to path:", path.String())
	groupPath := path.Join(groupName)

	err := groupPath.MkdirAll()
	if err != nil {
		return err
	}

	subgroups, err := tc.gitlabClient.ListSubgroups(groupID)
	if err != nil {
		return err
	}

	for _, subgroup := range subgroups {
		subgroup := subgroup
		tc.errGroup.Go(func() error {
			if util.MapContains(tc.ignoreIDs, subgroup.ID) || util.MapContains(tc.ignoreNames, subgroup.Name) {
				return nil
			}
			return tc.cloneGroup(subgroup.ID, subgroup.Name, groupPath)
		})
	}

	projects, err := tc.gitlabClient.ListProjects(groupID)
	if err != nil {
		return err
	}
	for _, project := range projects {
		project := project
		invalid := project.Archived || project.EmptyRepo
		if !invalid {
			tc.errGroup.Go(func() error {
				return tc.gitClient.CloneRepo(project.HTTPURLToRepo, groupPath.Join(project.Name))
			})
		}
	}

	return nil
}
