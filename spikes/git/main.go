package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type byCommitter []*object.Commit

func (a byCommitter) Len() int      { return len(a) }
func (a byCommitter) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byCommitter) Less(i, j int) bool {
	return a[i].Committer.When.Before(a[j].Committer.When)
}

func die(err error) {
	fmt.Fprintln(os.Stdout, "FATAL:", err)
	os.Exit(1)
}

func check(commit *object.Commit, path string) (h plumbing.Hash, err error) {
	tree, err := commit.Tree()
	if err != nil {
		return
	}

	subtree, err := tree.Tree(path)
	if err != nil {
		return
	}

	h = subtree.Hash
	return
}

func open(path string) (c *object.Commit, err error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return
	}

	ref, err := r.Head()
	if err != nil {
		return
	}

	c, err = r.CommitObject(ref.Hash())
	return
}

func search(head *object.Commit, path string, target plumbing.Hash) (*object.Commit, error) {
	seen := map[plumbing.Hash]struct{}{
		head.Hash: struct{}{},
	}
	stack := []*object.Commit{head}
	leafs := []*object.Commit{}

	for {
		n := len(stack)
		if n < 1 {
			break
		}

		c := stack[n-1]
		stack = stack[:n-1]
		found := 0

		err := c.Parents().ForEach(func(parent *object.Commit) error {
			if _, ok := seen[parent.Hash]; ok {
				return nil
			}
			if h, err := check(parent, path); err != nil {
				if err != object.ErrDirectoryNotFound && err != object.ErrFileNotFound {
					return err
				}
			} else if h == target {
				seen[parent.Hash] = struct{}{}
				stack = append(stack, parent)
				found += 1
			}
			return nil
		})
		if err != nil {
			return nil, nil
		}
		if found == 0 {
			leafs = append(leafs, c)
		}
	}

	sort.Sort(byCommitter(leafs))
	return leafs[0], nil
}

func main() {
	repo := os.Args[1]
	path := os.Args[2]
	path = strings.TrimSuffix(path, "/")

	head, err := open(repo)
	if err != nil {
		die(err)
	}

	target, err := check(head, path)
	if err != nil {
		die(err)
	}

	oldest, err := search(head, path, target)
	if err != nil {
		die(err)
	}

	fmt.Print(oldest)
}
