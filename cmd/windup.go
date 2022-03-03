package main

import (
	"github.com/konveyor/tackle2-hub/api"
	"path"
)

//
// Windup application analyzer.
type Windup struct {
	*Data
	bucket *api.Bucket
}

//
// Run windup.
func (r *Windup) Run() (err error) {
	cmd := Command{Path: "/opt/windup"}
	cmd.Options, err = r.options()
	if err != nil {
		return
	}
	err = cmd.Run()
	return
}

//
// options builds CLL options.
func (r *Windup) options() (options Options, err error) {
	options = Options{
		"--batchMode",
		"--output",
		r.bucket.Path,
	}
	err = r.Mode.AddOptions(&options)
	if err != nil {
		return
	}
	if r.Targets != nil {
		err = r.Targets.AddOptions(&options)
		if err != nil {
			return
		}
	}
	err = r.Scope.AddOptions(&options)
	if err != nil {
		return
	}
	if r.Rules != nil {
		err = r.Scope.AddOptions(&options)
		if err != nil {
			return
		}
	}
	return
}

//
// BucketPath path relative to a bucket.
type BucketPath struct {
	Bucket uint   `json:"bucket" binding:"required"`
	Path   string `json:"path" binding:"required"`
}

//
// Mode settings.
type Mode struct {
	Binary     bool        `json:"binary"`
	Artifact   *BucketPath `json:"artifact"`
	WithDeps   bool        `json:"withDeps"`
	Repository Repository
}

//
// AddOptions adds windup options.
func (r *Mode) AddOptions(options *Options) (err error) {
	if r.Repository == nil {
		return
	}
	options.add("--input", SourceDir)
	options.add("--sourceMode")
	return
}

//
// Targets list of target.
type Targets []string

//
// AddOptions add options.
func (r Targets) AddOptions(options *Options) (err error) {
	for _, target := range r {
		options.add("--target", target)
	}
	return
}

//
// Scope settings.
type Scope struct {
	WithKnown bool `json:"withKnown"`
	Packages  struct {
		Included []string `json:"included"`
		Excluded []string `json:"excluded"`
	} `json:"packages"`
}

//
// AddOptions adds windup options.
func (r *Scope) AddOptions(options *Options) (err error) {
	if len(r.Packages.Included) > 0 {
		options.add("--packages", r.Packages.Included...)
	}
	if len(r.Packages.Excluded) > 0 {
		options.add("--excludePackages", r.Packages.Excluded...)
	}
	return
}

//
// Rules settings.
type Rules struct {
	Directory *BucketPath `json:"directory"`
	Tags      struct {
		Included []string `json:"included"`
		Excluded []string `json:"excluded"`
	} `json:"tags"`
}

//
// AddOptions adds windup options.
func (r *Rules) AddOptions(options *Options) (err error) {
	bucket, bErr := addon.Bucket.Get(r.Directory.Bucket)
	if bErr != nil {
		return
	}
	options.add(
		"--userRulesDirectory",
		path.Join(
			bucket.Path,
			r.Directory.Path))
	if len(r.Tags.Included) > 0 {
		options.add("--includeTags", r.Tags.Included...)
	}
	if len(r.Tags.Excluded) > 0 {
		options.add("--excludeTags", r.Tags.Excluded...)
	}
	return
}
