package main

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func main() {
	ctx := context.TODO()
	s := memory.NewStorage()

	r, err := git.CloneContext(ctx, s, nil, &git.CloneOptions{
		URL:          "https://github.com/sjansen/mecha",
		Depth:        1,
		SingleBranch: true,
		Tags:         git.NoTags,
	})
	if err != nil {
		die(err)
	}

	ref, err := r.Head()
	if err != nil {
		die(err)
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		die(err)
	}

	fmt.Println(commit)
}
