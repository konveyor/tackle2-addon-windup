package main

import (
	"errors"
	"github.com/konveyor/tackle2-hub/api"
	pathlib "path"
	"strings"
)

var (
	// SourceDir repository path.
	SourceDir = "source"
)

func init() {
	SourceDir = pathlib.Join(Dir, SourceDir)
}

//
// Factory.
func newRepository(a *api.Application) (rp Repository, err error) {
	kind := a.Repository.Kind
	if kind == "" {
		if strings.HasSuffix(a.Repository.URL, ".git") {
			kind = "git"
		} else {
			kind = "svn"
		}
	}
	switch kind {
	case "git":
		rp = &Git{application: a}
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
	Fetch() (err error)
}
