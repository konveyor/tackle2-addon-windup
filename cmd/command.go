package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

//
// Command runner.
type Command struct {
	Path    string
	Dir     string
	Options Options
}

//
// Run executes the command.
// The command and output are both reported in
// task Report.Activity.
func (r *Command) Run() (err error) {
	addon.Activity(
		"[CMD] Running: %s %s",
		r.Path,
		strings.Join(r.Options, " "))
	cmd := exec.Command(r.Path, r.Options...)
	cmd.Dir = r.Dir
	b, err := cmd.CombinedOutput()
	exitErr := &exec.ExitError{}
	if errors.As(err, &exitErr) {
		output := string(b)
		for _, line := range strings.Split(output, "\n") {
			addon.Activity(
				"> %s",
				line)
		}
	}

	return
}

//
// RunSilent executes the command.
// Nothing reported in task Report.Activity.
func (r *Command) RunSilent() (err error) {
	cmd := exec.Command(r.Path, r.Options...)
	cmd.Dir = r.Dir
	err = cmd.Run()
	return
}

//
// Options are CLI options.
type Options []string

//
// add
func (a *Options) add(option string, s ...string) {
	*a = append(*a, option)
	*a = append(*a, s...)
}

//
// add
func (a *Options) addf(option string, x ...interface{}) {
	*a = append(*a, fmt.Sprintf(option, x...))
}
