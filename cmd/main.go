package main

import (
	"context"
	"flag"
	"log"
	"strconv"
	"strings"

	"github.com/zakaprov/gitlab-group-clone/app"
	"github.com/zakaprov/gitlab-group-clone/infra"
	"golang.org/x/sync/errgroup"
)

const JUNI_ROOT_GROUP_ID = 7330753
const JUNI_ROOT_GROUP_NAME = "junitechnology"

type Config struct {
	CloneDir    string
	GroupNames  map[string]bool
	GroupIDs    map[int]bool
	RootGroupID int
}

func NewConfig(groups string, ids string) *Config {
	splitGroups := strings.Split(groups, ",")
	splitIDs := strings.Split(ids, ",")

	groupNames := make(map[string]bool)
	for _, group := range splitGroups {
		groupNames[group] = true
	}

	groupIDs := make(map[int]bool)
	for _, id := range splitIDs {
		intID, _ := strconv.Atoi(id)
		groupIDs[intID] = true
	}

	return &Config{
		GroupNames: groupNames,
		GroupIDs:   groupIDs,
	}
}

func main() {
	token := flag.String("token", "", "GitLab token to use")
	groupNames := flag.String("group-names", "", "Groups to fetch")
	groupIDs := flag.String("group-ids", "", "Groups to fetch")
	flag.Parse()

	config := NewConfig(*groupNames, *groupIDs)
	println(config)

	ctx := context.Background()
	errGroup, _ := errgroup.WithContext(ctx)

	gc, err := infra.NewGitlabClient(*token)
	if err != nil {
		log.Fatal(err)
		return
	}

	treeClone := app.TreeClone{
		ErrGroup:     errGroup,
		GitClient:    infra.NewGitClient(*token),
		GitlabClient: gc,
	}
	err = treeClone.CloneGroup(JUNI_ROOT_GROUP_ID, JUNI_ROOT_GROUP_NAME, ".")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = errGroup.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
