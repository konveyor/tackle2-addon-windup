package main

import (
	"errors"
	liberr "github.com/konveyor/controller/pkg/error"
	"github.com/konveyor/tackle2-hub/api"
	"os"
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
func newRepository(homeDir string, application *api.Application) (r Repository, err error) {
	switch strings.ToLower(application.Repository.Kind) {
	case "subversion":
		r = &Subversion{}
	default:
		r = &Git{}
	}
	r.With(homeDir, application)
	err = r.Validate()
	return
}

//
// Repository interface.
type Repository interface {
	With(homeDir string, application *api.Application)
	Fetch() (err error)
	Validate() (err error)
}

//
// SCM - source code manager.
type SCM struct {
	Application *api.Application
	HomeDir     string
}

//
// With settings.
func (r *SCM) With(homeDir string, application *api.Application) {
	r.HomeDir = homeDir
	r.Application = application
}

//
// EnsureDir creates a directory if not exists.
func (r *SCM) EnsureDir(path string, mode os.FileMode) (err error) {
	err = os.MkdirAll(path, mode)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			err = nil
		} else {
			err = liberr.Wrap(
				err,
				"path",
				path)
		}
	}
	return
}

//
// WriteKey writes the SSH key.
func (r *SCM) WriteKey(id *api.Identity) (err error) {
	if id.Key == "" {
		return
	}
	dir := pathlib.Join(HomeDir, ".ssh")
	err = r.EnsureDir(dir, 0700)
	if err != nil {
		return
	}
	path := pathlib.Join(dir, "id_git")
	_, err = os.Stat(path)
	if !errors.Is(err, os.ErrNotExist) {
		err = liberr.Wrap(os.ErrExist)
		return
	}
	f, err := os.Create(path)
	if err != nil {
		err = liberr.Wrap(
			err,
			"path",
			path)
		return
	}
	_, err = f.Write([]byte(id.Key))
	if err != nil {
		err = liberr.Wrap(
			err,
			"path",
			path)
	}
	_ = f.Close()
	return
}
