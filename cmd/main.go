package main

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/konveyor/tackle2-addon/repository"
	"github.com/konveyor/tackle2-addon/ssh"
	hub "github.com/konveyor/tackle2-hub/addon"
	"github.com/konveyor/tackle2-hub/nas"
)

var (
	addon     = hub.Addon
	HomeDir   = ""
	DepDir    = ""
	BinDir    = ""
	SourceDir = ""
	AppDir    = ""
	Dir       = ""
	M2Dir     = ""
	ReportDir = ""
	RuleDir   = ""
)

func init() {
	Dir, _ = os.Getwd()
	HomeDir, _ = os.UserHomeDir()
	SourceDir = path.Join(Dir, "source")
	DepDir = path.Join(Dir, "deps")
	BinDir = path.Join(Dir, "bin")
	ReportDir = path.Join(Dir, "report")
	RuleDir = path.Join(Dir, "rules")
	M2Dir = "/cache/m2"
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
		// Create directories.
		for _, dir := range []string{BinDir, M2Dir, RuleDir, ReportDir} {
			err = nas.MkDir(dir, 0755)
			if err != nil {
				return
			}
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
		// Delete report.
		mark := time.Now()
		bucket := addon.Application.Bucket(application.ID)
		err = bucket.Delete(d.Output)
		if err != nil {
			return
		}
		addon.Activity(
			"[BUCKET] Report deleted:%s duration:%v.",
			d.Output,
			time.Since(mark))
		//
		// Maven.
		maven := repository.Maven{
			Application: application,
			M2Dir:       M2Dir,
			BinDir:      DepDir,
		}
		//
		// SSH
		agent := ssh.Agent{}
		err = agent.Start()
		if err != nil {
			return
		}
		//
		// The application has modules.
		var hasModules bool
		//
		// Fetch repository.
		if !d.Mode.Binary {
			addon.Total(2)
			if application.Repository == nil {
				err = &SoftError{Reason: "Application repository not defined."}
				return
			}
			SourceDir = path.Join(
				SourceDir,
				strings.Split(
					path.Base(
						application.Repository.URL),
					".")[0])
			AppDir = path.Join(SourceDir, application.Repository.Path)
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
				hasModules, err = maven.HasModules(AppDir)
				if err != nil {
					return
				}
				if hasModules {
					err = maven.InstallArtifacts(AppDir)
					if err != nil {
						return
					}
				}
				err = maven.Fetch(AppDir)
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
		//
		// Update report.
		mark = time.Now()
		bucket = addon.Application.Bucket(application.ID)
		err = bucket.Put(ReportDir, d.Output)
		if err != nil {
			return
		}
		addon.Activity(
			"[BUCKET} Report updated:%s duration:%v.",
			d.Output,
			time.Since(mark))
		//
		// Clean up.
		if hasModules {
			err = maven.DeleteArtifacts(SourceDir)
			if err != nil {
				return
			}
		}
		return
	})
}
