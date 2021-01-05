package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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

func adjust(repo, path string) (root, file string) {
	root = repo
	file = path
	if repo, err := filepath.Abs(repo); err == nil {
		for {
			subdir := filepath.Join(repo, git.GitDirName)
			if x, err := os.Stat(subdir); os.IsNotExist(err) {
				basename := filepath.Base(repo)
				dirname := filepath.Dir(repo)
				if dirname == repo {
					return
				}
				repo = dirname
				path = filepath.Join(basename, path)
			} else if err == nil && x.IsDir() {
				root = repo
				file = path
				return
			}
		}
	}
	return
}

func check(commit *object.Commit, path string) (h plumbing.Hash, err error) {
	tree, err := commit.Tree()
	if err != nil {
		return
	}

	if path == "." {
		h = tree.Hash
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
		head.Hash: {},
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
				found++
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
	repo := "."
	path := "."
	if len(os.Args) > 1 {
		repo = os.Args[1]
	}
	if len(os.Args) > 2 {
		path = strings.TrimSuffix(os.Args[2], "/")
	}

	repo, path = adjust(repo, path)

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
