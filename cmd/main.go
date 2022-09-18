package main

import (
	"context"
	"log"

	"github.com/alexflint/go-arg"
	"github.com/chigopher/pathlib"
	"github.com/zakaprov/gitlab-group-clone/app"
	"github.com/zakaprov/gitlab-group-clone/infra"
	"golang.org/x/sync/errgroup"
)

type args struct {
	DestinationDir string   `arg:"positional" placeholder:"DESTINATION" help:"target dir for cloning the group tree."`
	IgnoreNames    []string `arg:"--ignore-name,separate" help:"If specified, subgroups with this name will not be cloned. May be given multiple times"`
	IgnoreIDs      []int    `arg:"--ignore-id,separate" help:"If specified, subgroups with this ID will not be cloned. May be given multiple times"`
	RootID         int      `arg:"-r,--root,required" placeholder:"GROUP_ID" help:"ID of the GitLab group to use as tree root"`
	Token          string   `arg:"-t,required" help:"GitLab API access token"`
}

func (args) Description() string {
	return "Clone repositories from a GitLab group recursively."
}

func main() {
	var args args
	arg.MustParse(&args)
	ctx := context.Background()
	errGroup, _ := errgroup.WithContext(ctx)

	gc, err := infra.NewGitlabClient(args.Token)
	if err != nil {
		log.Fatal(err)
		return
	}

	treeClone := app.TreeCloner{
		ErrGroup:     errGroup,
		GitClient:    infra.NewGitClient(args.Token),
		GitlabClient: gc,
	}
	err = treeClone.CloneTree(args.RootID, pathlib.NewPath(args.DestinationDir))
	if err != nil {
		log.Fatal(err)
		return
	}

	err = errGroup.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
