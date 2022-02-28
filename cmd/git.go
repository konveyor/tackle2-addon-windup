package main

import (
	"errors"
	"github.com/konveyor/tackle2-hub/api"
	"net/url"
	"os"
	pathlib "path"
)

//
// Git repository.
type Git struct {
	application *api.Application
}

//
// Fetch clones the repository.
func (r *Git) Fetch() (err error) {
	repository := r.application.Repository
	addon.Activity("[GIT] Cloning: %s", repository.URL)
	_ = os.RemoveAll(SourceDir)
	id, hasCreds, err := addon.Application.FindIdentity(r.application.ID, "git")
	if err != nil {
		return
	}
	err = r.writeConfig()
	if err != nil {
		return
	}
	if hasCreds {
		err = r.writeCreds(repository.URL, id)
		if err != nil {
			return
		}
	}
	cmd := Command{Path: "/usr/bin/git"}
	cmd.Options.add("clone", repository.URL, SourceDir)
	err = cmd.Run()
	return
}

//
// writeConfig writes config file.
func (r *Git) writeConfig() (err error) {
	path := pathlib.Join(HomeDir, ".gitconfig")
	_, err = os.Stat(path)
	if !errors.Is(err, os.ErrNotExist) {
		err = os.ErrExist
		return
	}
	f, err := os.Create(path)
	if err != nil {
		return
	}
	_, err = f.Write([]byte("[credential]\n"))
	_, err = f.Write([]byte("helper = store\n"))
	return
}

//
// writeCreds writes credentials (store) file.
func (r *Git) writeCreds(u string, id *api.Identity) (err error) {
	path := pathlib.Join(HomeDir, ".git-credentials")
	_, err = os.Stat(path)
	if !errors.Is(err, os.ErrNotExist) {
		err = os.ErrExist
		return
	}
	f, err := os.Create(path)
	if err != nil {
		return
	}
	parsed, err := url.Parse(u)
	if err != nil {
		return
	}
	entry := parsed.Scheme
	entry += "://"
	if id.User != "" {
		entry += id.User
		entry += ":"
	}
	if id.Password != "" {
		entry += id.Password
		entry += "@"
	}
	entry += parsed.Host
	_, err = f.Write([]byte(entry))
	_ = f.Close()
	return
}
