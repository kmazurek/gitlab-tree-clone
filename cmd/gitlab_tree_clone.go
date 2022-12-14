package main

import (
	"context"
	"log"

	"github.com/alexflint/go-arg"
	"github.com/chigopher/pathlib"
	"github.com/kmazurek/gitlab-tree-clone/internal/app"
	"github.com/kmazurek/gitlab-tree-clone/internal/infra"
	"golang.org/x/sync/errgroup"
)

type args struct {
	IgnoreIDs   []int    `arg:"--ignore-id,separate" help:"If specified, subgroups with this ID will not be cloned. May be given multiple times"`
	IgnoreNames []string `arg:"--ignore-name,separate" help:"If specified, subgroups with this name will not be cloned. May be given multiple times"`
	OutputDir   string   `arg:"-o,--output-dir" default:"." placeholder:"OUTPUT_DIR" help:"target dir for cloning the group tree"`
	RootID      int      `arg:"positional" placeholder:"GROUP_ID" help:"ID of the GitLab group to use as tree root"`
	Token       string   `arg:"-t,required" help:"GitLab API access token"`
}

func (args) Description() string {
	return "Clone repositories from a GitLab group recursively."
}

func main() {
	var args args
	arg.MustParse(&args)
	ctx := context.Background()

	errGroup, _ := errgroup.WithContext(ctx)
	gitClient := infra.NewGitClient(args.Token)
	gitlabClient, err := infra.NewGitlabClient(args.Token)
	if err != nil {
		log.Fatal(err)
		return
	}

	treeCloner, err := app.NewTreeCloner(gitClient, gitlabClient, errGroup, args.IgnoreIDs, args.IgnoreNames)
	if err != nil {
		log.Fatal(err)
		return
	}

	errGroup.Go(func() error {
		return treeCloner.CloneTree(args.RootID, pathlib.NewPath(args.OutputDir))
	})

	err = errGroup.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
