package domain

import (
	"github.com/chigopher/pathlib"
	"github.com/xanzy/go-gitlab"
)

type GitClient interface {
	CloneRepo(url string, path *pathlib.Path) error
}

type GitlabClient interface {
	GetGroup(groupId int) (*gitlab.Group, error)
	ListProjects(groupId int) ([]*gitlab.Project, error)
	ListSubgroups(groupId int) ([]*gitlab.Group, error)
}