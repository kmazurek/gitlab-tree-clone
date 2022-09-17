package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/xanzy/go-gitlab"
	"golang.org/x/sync/errgroup"
)

const JUNI_ROOT_GROUP_ID = 7330753
const JUNI_ROOT_GROUP_NAME = "junitechnology"

func cloneGroup(client *gitlab.Client, errGroup *errgroup.Group, groupID int, groupName string, path string) error {
	log.Println("Cloning group: " + groupName + " to path: " + path)
	path = path + "/" + groupName
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	subgroups, err := listSubgroups(client, groupID)
	if err != nil {
		return err
	}

	for _, subgroup := range subgroups {
		subgroup := subgroup
		errGroup.Go(func() error {
			return cloneGroup(client, errGroup, subgroup.ID, subgroup.Name, path)
		})
	}

	projects, err := listProjects(client, groupID)
	if err != nil {
		return err
	}
	for _, project := range projects {
		invalid := project.Archived || project.EmptyRepo
		if !invalid {
			project := project
			errGroup.Go(func() error {
				return cloneProject(path, project.Name, project.SSHURLToRepo)
			})
		}
	}

	return nil
}

func cloneProject(path string, projectName string, url string) error {
	log.Println("Cloning project: " + path + "/" + projectName)

	out, err := exec.Command("git", "clone", url, path+"/"+projectName).CombinedOutput()
	if err != nil {
		output := string(out)
		if strings.Contains(output, "already exists and is not an empty directory") {
			return nil
		}

		log.Fatal(string(out))
		return err
	}

	return nil
}

func listSubgroups(client *gitlab.Client, groupId int) ([]*gitlab.Group, error) {
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

func listProjects(client *gitlab.Client, groupId int) ([]*gitlab.Project, error) {
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

type Config struct {
	GroupNames map[string]bool
	GroupIDs   map[int]bool
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

	gitlabClient, err := gitlab.NewClient(*token)
	if err != nil {
		log.Fatal(err)
		return
	}

	ctx := context.Background()
	errGroup, _ := errgroup.WithContext(ctx)

	err = cloneGroup(gitlabClient, errGroup, JUNI_ROOT_GROUP_ID, JUNI_ROOT_GROUP_NAME, ".")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = errGroup.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
