package infra

import (
	"errors"
	"log"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type GitClient struct {
	apiToken string
}

func NewGitClient(apiToken string) *GitClient {
	return &GitClient{apiToken: apiToken}
}

func (gc *GitClient) CloneProject(path string, projectName string, url string) error {
	log.Println("Cloning project: " + path + "/" + projectName)

	_, err := git.PlainClone(path+"/"+projectName, false, &git.CloneOptions{
		URL: url,
		Auth: &http.BasicAuth{
			Username: "username", // anything except an empty string
			Password: gc.apiToken,
		},
	})

	if err != nil {
		if errors.Is(err, git.ErrRepositoryAlreadyExists) {
			return nil
		}
		return err
	}

	return nil
}

func (gc *GitClient) PullProject(path string, projectName string) error {
	projectPath := path + "/" + projectName
	log.Println("Pulling project: " + projectPath)

	repo, err := git.PlainOpen(projectPath)
	if err != nil {
		return err
	}
	tree, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = tree.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: "username", // anything except an empty string
			Password: gc.apiToken,
		},
	})

	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			return nil
		}
		return err
	}

	return nil
}
