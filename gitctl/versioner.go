package gitctl

import (
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type Gitctl struct {
	path    string
	current *git.Repository
}

func NewGitctl(path string) *Gitctl {
	return &Gitctl{path: path}
}

func (gic *Gitctl) NewVerison() {
	gic.ensureCurrent()
	wc, err := gic.current.Worktree()
	if err != nil {
		panic(err)
	}
	wc.AddGlob("autocommit/*")
	var author object.Signature
	author.Name = "Auto Commiter"
	author.Email = "stub@stub.kkdev.org"
	author.When = time.Now()
	_, err = wc.Commit("New File Version", &git.CommitOptions{Author: author})
	if err != nil {
		panic(err)
	}
}

func (gic *Gitctl) ensureCurrent() {
	if gic.current == nil {
		var err error
		gic.current, err = git.PlainOpen(gic.path)
		if err != nil {
			panic(err)
		}
	}
}
