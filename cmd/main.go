package main

import (
	"github.com/konveyor/tackle2-addon/repository"
	"github.com/konveyor/tackle2-addon/ssh"
	hub "github.com/konveyor/tackle2-hub/addon"
	"os"
	"path"
	"strings"
)

var (
	// hub integration.
	addon = hub.Addon
	// HomeDir directory.
	HomeDir   = ""
	BinDir    = ""
	SourceDir = ""
	Dir       = ""
)

func init() {
	Dir, _ = os.Getwd()
	HomeDir, _ = os.UserHomeDir()
	SourceDir = path.Join(Dir, "source")
	BinDir = path.Join(Dir, "dependencies")
}

type SoftError = hub.SoftError

//
// Data Addon data passed in the secret.
type Data struct {
	// Output directory within application bucket.
	Output string `json:"output" binding:"required"`
	// Mode options.
	Mode Mode `json:"mode"`
	// Sources list.
	Sources Sources `json:"sources"`
	// Targets list.
	Targets Targets `json:"targets"`
	// Scope options.
	Scope Scope `json:"scope"`
	// Rules options.
	Rules *Rules `json:"rules"`
}

//
// main
func main() {
	addon.Run(func() (err error) {
		//
		// Get the addon data associated with the task.
		d := &Data{}
		err = addon.DataWith(d)
		if err != nil {
			err = &SoftError{Reason: err.Error()}
			return
		}
		//
		// windup
		windup := Windup{}
		windup.Data = d
		//
		// Fetch application.
		addon.Activity("Fetching application.")
		application, err := addon.Task.Application()
		if err == nil {
			windup.application = application
		} else {
			return
		}
		//
		// Maven.
		maven := repository.Maven{
			Application: application,
			M2Dir:       "/mnt/m2",
			BinDir:      BinDir,
		}
		//
		// SSH
		agent := ssh.Agent{}
		err = agent.Start()
		if err != nil {
			return
		}
		//
		// Fetch repository.
		if !d.Mode.Binary {
			addon.Total(2)
			if application.Repository == nil {
				err = &SoftError{Reason: "Application repository not defined."}
				return
			}
			SourceDir = path.Join(
				Dir,
				strings.Split(
					path.Base(
						application.Repository.URL),
					".")[0])
			var r repository.Repository
			r, err = repository.New(SourceDir, application)
			if err != nil {
				return
			}
			err = r.Fetch()
			if err == nil {
				addon.Increment()
				windup.Mode.Repository = r
			} else {
				return
			}
			if d.Mode.WithDeps {
				err = maven.Fetch(SourceDir)
				if err != nil {
					return
				}
			}
		} else {
			if d.Mode.Artifact == "" {
				err = maven.FetchArtifact()
				if err != nil {
					return
				}
			}
		}
		//
		// Run windup.
		err = windup.Run()
		if err == nil {
			addon.Increment()
		} else {
			return
		}
		return
	})
}
