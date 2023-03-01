package main

import (
	"bufio"
	"github.com/konveyor/tackle2-addon/command"
	"github.com/konveyor/tackle2-addon/repository"
	"github.com/konveyor/tackle2-hub/api"
	"github.com/konveyor/tackle2-hub/nas"
	"os"
	pathlib "path"
)

//
// Windup application analyzer.
type Windup struct {
	application *api.Application
	*Data
}

//
// Run windup.
func (r *Windup) Run() (err error) {
	cmd := command.Command{Path: "/opt/windup"}
	cmd.Options, err = r.options()
	if err != nil {
		return
	}
	err = cmd.Run()
	if err != nil {
		r.reportLog()
	}

	return
}

//
// reportLog reports the log content.
func (r *Windup) reportLog() {
	path := pathlib.Join(
		HomeDir,
		".mta",
		"log",
		"mta.log")
	f, err := os.Open(path)
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		addon.Activity(">> %s\n", scanner.Text())
	}
	_ = f.Close()
}

//
// options builds CLL options.
func (r *Windup) options() (options command.Options, err error) {
	options = command.Options{
		"--batchMode",
		"--output",
		ReportDir,
	}
	err = r.maven(&options)
	if err != nil {
		return
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
		err = r.Rules.AddOptions(&options)
		if err != nil {
			return
		}
	}
	return
}

//
// maven add --input for maven artifacts.
func (r *Windup) maven(options *command.Options) (err error) {
	found, err := nas.HasDir(DepDir)
	if found {
		options.Add("--input", DepDir)
	}
	return
}

//
// Mode settings.
type Mode struct {
	Binary     bool   `json:"binary"`
	Artifact   string `json:"artifact"`
	WithDeps   bool   `json:"withDeps"`
	Diva       bool   `json:"diva"`
	CSV        bool   `json:"csv"`
	Repository repository.Repository
}

//
// AddOptions adds windup options.
func (r *Mode) AddOptions(options *command.Options) (err error) {
	if r.Binary {
		if r.Artifact != "" {
			bucket := addon.Bucket()
			err = bucket.Get(r.Artifact, BinDir)
			if err != nil {
				return
			}
			options.Add("--input", BinDir)
		}
	} else {
		options.Add("--input", AppDir)
	}
	if r.Diva {
		options.Add("--enableTransactionAnalysis")
	}
	if r.CSV {
		options.Add("--exportCSV")
	}

	return
}

//
// Sources list of sources.
type Sources []string

//
// AddOptions add options.
func (r Sources) AddOptions(options *command.Options) (err error) {
	for _, source := range r {
		options.Add("--source", source)
	}
	return
}

//
// Targets list of target.
type Targets []string

//
// AddOptions add options.
func (r Targets) AddOptions(options *command.Options) (err error) {
	for _, target := range r {
		options.Add("--target", target)
	}
	return
}

//
// Scope settings.
type Scope struct {
	WithKnown bool `json:"withKnown"`
	Packages  struct {
		Included []string `json:"included,omitempty"`
		Excluded []string `json:"excluded,omitempty"`
	} `json:"packages"`
}

//
// AddOptions adds windup options.
func (r *Scope) AddOptions(options *command.Options) (err error) {
	if r.WithKnown {
		options.Add("--analyzeKnownLibraries")
	}
	if len(r.Packages.Included) > 0 {
		options.Add("--packages", r.Packages.Included...)
	}
	if len(r.Packages.Excluded) > 0 {
		options.Add("--excludePackages", r.Packages.Excluded...)
	}
	return
}

//
// Rules settings.
type Rules struct {
	Path string `json:"path" binding:"required"`
	Tags struct {
		Included []string `json:"included,omitempty"`
		Excluded []string `json:"excluded,omitempty"`
	} `json:"tags"`
}

//
// AddOptions adds windup options.
func (r *Rules) AddOptions(options *command.Options) (err error) {
	options.Add(
		"--userRulesDirectory",
		RuleDir)
	bucket := addon.Bucket()
	err = bucket.Get(r.Path, RuleDir)
	if err != nil {
		return
	}
	if len(r.Tags.Included) > 0 {
		options.Add("--includeTags", r.Tags.Included...)
	}
	if len(r.Tags.Excluded) > 0 {
		options.Add("--excludeTags", r.Tags.Excluded...)
	}
	return
}
