package main

import (
	"errors"
	hub "github.com/konveyor/tackle2-hub/addon"
	"github.com/konveyor/tackle2-hub/api"
	"os"
)

var (
	// hub integration.
	addon = hub.Addon
	// HomeDir directory.
	HomeDir = ""
	Dir     = ""
)

func init() {
	HomeDir, _ = os.UserHomeDir()
	Dir, _ = os.Getwd()
}

//
// Data Addon data passed in the secret.
type Data struct {
	Application uint `json:"application" binding:"required"`
	//
	Mode    Mode    `json:"mode"`
	Targets Targets `json:"targets"`
	Scope   Scope   `json:"scope"`
	Rules   *Rules  `json:"rules"`
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
			return
		}
		windup := Windup{}
		windup.Data = d
		//
		// Fetch application.
		addon.Activity("Fetching application.")
		application, err := addon.Application.Get(d.Application)
		if err != nil {
			return
		}
		//
		// Fetch repository.
		if !d.Mode.Binary {
			addon.Total(2)
			if application.Repository == nil {
				err = errors.New("Application repository not defined.")
				return
			}
			var r Repository
			r, err = newRepository(application)
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
				mvn := Maven{}
				mvn.application = application
				err = mvn.Fetch()
				if err != nil {
					return
				}
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
// ensureBucket to store windup report.
func ensureBucket(d *Data) (bucket *api.Bucket, err error) {
	bucket, err = addon.Bucket.Ensure(d.Application, "Windup")
	if err != nil {
		return
	}
	err = addon.Bucket.Purge(bucket)
	return
}
