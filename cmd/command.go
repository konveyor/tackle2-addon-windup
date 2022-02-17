package main

import (
	"bytes"
	"os/exec"
	"strings"
)

//
// Command runner.
type Command struct {
	Path    string
	Dir     string
	Options Options
	Out     bytes.Buffer
	Err     bytes.Buffer
}

//
// Run command.
func (r *Command) Run() (err error) {
	addon.Activity(
		"[CMD] Running: %s %s",
		r.Path,
		strings.Join(r.Options, " "))
	cmd := exec.Command(r.Path, r.Options...)
	cmd.Dir = r.Dir
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
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
