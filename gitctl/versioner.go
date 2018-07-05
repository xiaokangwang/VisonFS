package gitctl

import "gopkg.in/src-d/go-git.v4"

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
	wc.AddGlob("*.*")
	wc.Commit("New File Version", &git.CommitOptions{})
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
