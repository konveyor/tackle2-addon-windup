package main

import (
	"os"
	pathlib "path"
)

//
// Windup application analyzer.
type Windup struct {
	*Data
	bucket string
}

//
// Run windup.
func (r *Windup) Run() (err error) {
	err = os.RemoveAll(r.output())
	if err != nil {
		return
	}
	err = os.Mkdir(r.output(), 0777)
	if err != nil {
		return
	}
	cmd := Command{Path: "/opt/windup"}
	cmd.Options, err = r.options()
	if err != nil {
		return
	}
	err = cmd.Run()
	return
}

//
// output returns output directory.
func (r *Windup) output() string {
	return pathlib.Join(
		r.bucket,
		r.Output)
}

//
// options builds CLL options.
func (r *Windup) options() (options Options, err error) {
	options = Options{
		"--batchMode",
		"--output",
		r.output(),
	}
	err = r.Mode.AddOptions(&options)
	if err != nil {
		return
	}
	if r.Sources != nil {
		err = r.Sources.AddOptions(&options)
		if err != nil {
			return
		}
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
// Mode settings.
type Mode struct {
	Binary     bool   `json:"binary"`
	Artifact   string `json:"artifact"`
	WithDeps   bool   `json:"withDeps"`
	Repository Repository
}

//
// AddOptions adds windup options.
func (r *Mode) AddOptions(options *Options) (err error) {
	if r.Binary {
		binDir := pathlib.Join(addon.Task.Bucket(), r.Artifact)
		options.add("--input", binDir)
	} else {
		options.add("--input", SourceDir)
		options.add("--sourceMode")
	}

	return
}

//
// Sources list of sources.
type Sources []string

//
// AddOptions add options.
func (r Sources) AddOptions(options *Options) (err error) {
	for _, source := range r {
		options.add("--source", source)
	}
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
	Directory string `json:"directory" binding:"required"`
	Tags      struct {
		Included []string `json:"included"`
		Excluded []string `json:"excluded"`
	} `json:"tags"`
}

//
// AddOptions adds windup options.
func (r *Rules) AddOptions(options *Options) (err error) {
	options.add(
		"--userRulesDirectory",
		pathlib.Join(
			addon.Task.Bucket(),
			r.Directory))
	if len(r.Tags.Included) > 0 {
		options.add("--includeTags", r.Tags.Included...)
	}
	if len(r.Tags.Excluded) > 0 {
		options.add("--excludeTags", r.Tags.Excluded...)
	}
	return
}
