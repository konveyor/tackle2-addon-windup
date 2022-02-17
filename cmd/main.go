package main

import (
	"errors"
	hub "github.com/konveyor/tackle2-hub/addon"
	"github.com/konveyor/tackle2-hub/api"
	"os"
	"strings"
	"time"
)

var (
	// hub integration.
	addon = hub.Addon
)

const (
	DefaultTarget = "cloud-readiness"
)

//
// Checkpoint mapping.
type Checkpoint = map[string]time.Duration

//
// Artifact uploaded.
type Artifact struct {
	Bucket uint   `json:"bucket" binding:"required"`
	Path   string `json:"path" binding:"required"`
}

//
// Data Addon data passed in the secret.
type Data struct {
	Application  uint       `json:"application" binding:"required"`
	Binary       bool       `json:"binary"`
	Dependencies bool       `json:"dependencies"`
	Targets      []string   `json:"targets"`
	Packages     []string   `json:"packages"`
	Artifact     *Artifact  `json:"artifact"`
	Checkpoint   Checkpoint `json:"checkpoint"`
}

//
// validate settings.
// Default settings not specified.
func (d *Data) validate() (err error) {
	if d.Application == 0 {
		err = errors.New("Application not specified.")
		return
	}
	if len(d.Targets) == 0 {
		d.Targets = []string{DefaultTarget}
	}

	return
}

//
// main
func main() {
	addon.Run(func() (err error) {
		windup := Windup{}
		//
		// Get the addon data associated with the task.
		d := &Data{}
		err = addon.DataWith(d)
		if err != nil {
			return
		}
		//
		// Debugging.
		checkpoint(d, "started")
		defer func() {
			checkpoint(d, "done")
		}()

		addon.Activity("Working in: %s", cwd())
		//
		// Validate the addon data.
		err = d.validate()
		if err != nil {
			return
		}
		//
		// Fetch application.
		addon.Activity("Fetching application.")
		application, err := addon.Application.Get(d.Application)
		if err != nil {
			return
		}
		//
		// Fetch repository.
		if !d.Binary {
			addon.Total(2)
			if application.Repository == nil {
				err = errors.New("Application repository not defined.")
				return
			}
			var r Repository
			r, err = newRepository(application.Repository)
			if err != nil {
				return
			}
			err = r.Fetch("git")
			if err == nil {
				addon.Increment()
				windup.repository = r
			} else {
				return
			}
		}
		//
		// Create the bucket.
		addon.Activity("Ensure bucket (Windup).")
		bucket, err := ensureBucket(d)
		if err == nil {
			addon.Activity("Using bucket id:%d at:%s.", bucket.ID, bucket.Path)
			windup.bucket = bucket
		} else {
			return
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

//
// Checkpoint for container inspection.
func checkpoint(d *Data, name string) {
	duration, found := d.Checkpoint[name]
	if !found {
		return
	}
	m := time.Minute * duration
	addon.Activity(
		"[Debug] paused at: %s for: %v",
		strings.ToUpper(name),
		m)
	time.Sleep(m)
	addon.Activity("[Debug] Resumed.")
}

//
// ensureBucket to store windup report.
func ensureBucket(d *Data) (bucket *api.Bucket, err error) {
	bucket, err = addon.Bucket.Ensure(d.Application, "Windup")
	if err != nil {
		return
	}
	err = addon.Bucket.Purge(bucket)
	return
}

//
// cwd
func cwd() (path string) {
	path, _ = os.Getwd()
	return
}
