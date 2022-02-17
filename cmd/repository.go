package main

import (
	"errors"
	"github.com/konveyor/tackle2-hub/api"
	"os"
	"path"
	"strings"
)

//
// Factory.
func newRepository(r *api.Repository) (rp Repository, err error) {
	kind := r.Kind
	if kind == "" {
		if strings.HasSuffix(r.URL, ".git") {
			kind = "git"
		} else {
			kind = "svn"
		}
	}
	switch kind {
	case "git":
		rp = &Git{Repository: r}
	case "svn":
	case "mvn":
	default:
		err = errors.New("unknown kind")
	}

	return
}

//
// Repository interface.
type Repository interface {
	Fetch(path string) (err error)
	Path() string
}

//
// Git repository.
type Git struct {
	*api.Repository
	path string
}

//
// Fetch the repository.
func (r *Git) Fetch(path string) (err error) {
	addon.Activity("[GIT] Cloning: %s", r.URL)
	r.path = path
	_ = os.RemoveAll(r.path)
	cmd := Command{Path: "/usr/bin/git"}
	cmd.Options.add("clone", r.URL, path)
	err = cmd.Run()
	if err != nil {
		return
	}

	return
}

//
// Path to the fetched repository.
func (r *Git) Path() string {
	return path.Join(cwd(), r.path)
}
