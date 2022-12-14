package infra

import (
	"errors"
	"log"

	"github.com/chigopher/pathlib"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/kmazurek/gitlab-tree-clone/internal/domain"
)

// implements domain.GitClient
var _ domain.GitClient = (*GitClient)(nil)

type GitClient struct {
	apiToken string
}

func NewGitClient(apiToken string) *GitClient {
	return &GitClient{apiToken: apiToken}
}

func (gc *GitClient) CloneRepo(url string, path *pathlib.Path) error {
	log.Println("Cloning project:", path.String())

	_, err := git.PlainClone(path.String(), false, &git.CloneOptions{
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
