package main

import (
	"github.com/konveyor/tackle2-hub/api"
)

//
// Windup application analyzer.
type Windup struct {
	repository Repository
	bucket     *api.Bucket
	packages   []string
	targets    []string
}

//
// Run windup.
func (r *Windup) Run() (err error) {
	cmd := Command{Path: "/opt/windup"}
	cmd.Options = r.options()
	err = cmd.Run()
	if err != nil {
		return
	}

	return
}

//
// options builds CLL options.
func (r *Windup) options() (options Options) {
	options = Options{
		"--batchMode",
		"--output",
		r.bucket.Path,
	}
	options.add("--target", r.targets...)
	options.add("--input", r.repository.Path())
	if r.repository != nil {
		options.add("--sourceMode")
	}
	if len(r.packages) > 0 {
		options.add("--packages", r.packages...)
	}
	return
}
