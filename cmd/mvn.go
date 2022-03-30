package main

import (
	"errors"
	"github.com/clbanning/mxj"
	"github.com/konveyor/tackle2-hub/api"
	"os"
	pathlib "path"
)

const (
	SettingsFile = "settings.xml"
)

var (
	// DepsDir the maven output dir.
	DepsDir = "deps"
	// M2Dir the maven repository.
	M2Dir = "/mnt/m2"
)

func init() {
	//
	// DEPS_DIR - Maven output dir.
	if path, found := os.LookupEnv("DEPS_DIR"); found {
		DepsDir = path
	} else {
		DepsDir = pathlib.Join(Dir, DepsDir)
	}
	//
	// M2_DIR - Maven repository.
	if path, found := os.LookupEnv("M2_DIR"); found {
		M2Dir = path
	}
}

//
// Maven repository.
type Maven struct {
	Application *api.Application
	path        string
}

//
// Fetch fetches dependencies listed in the POM.
func (r *Maven) Fetch() (err error) {
	addon.Activity("[MVN] Fetch dependencies.")
	pom := pathlib.Join(SourceDir, "pom.xml")
	options := Options{
		"dependency:copy-dependencies",
		"-f",
		pom,
	}
	err = r.run(options)
	return
}

//
// FetchArtifact fetches an artifact.
func (r *Maven) FetchArtifact() (err error) {
	artifact := r.Application.Binary
	addon.Activity("[MVN] Fetch artifact %s.", artifact)
	options := Options{
		"dependency:copy",
	}
	options.addf("-Dartifact=%s", artifact)
	options.add("-Dmdep.useBaseVersion=true")
	err = r.run(options)
	return
}

//
// run executes maven.
func (r *Maven) run(options Options) (err error) {
	settings, err := r.writeSettings()
	if err != nil {
		return
	}
	err = os.MkdirAll(DepsDir, 0755)
	if err != nil {
		return
	}
	cmd := Command{Path: "/usr/bin/mvn"}
	cmd.Options = options
	cmd.Options.addf("-DoutputDirectory=%s", DepsDir)
	cmd.Options.addf("-Dmaven.repo.local=%s", M2Dir)
	cmd.Options.add("-s", settings)
	err = cmd.Run()
	return
}

//
// writeSettings writes settings file.
func (r *Maven) writeSettings() (path string, err error) {
	id, found, err := addon.Application.FindIdentity(r.Application.ID, "maven")
	if err != nil {
		return
	}
	if !found {
		return
	}
	path = pathlib.Join(Dir, SettingsFile)
	_, err = os.Stat(path)
	if !errors.Is(err, os.ErrNotExist) {
		err = os.ErrExist
		return
	}
	f, err := os.Create(path)
	if err != nil {
		return
	}
	settings := id.Settings
	settings, err = r.injectProxy(id)
	if err != nil {
		return
	}
	_, err = f.Write([]byte(settings))
	_ = f.Close()
	return
}

//
// injectProxy injects proxy settings.
func (r *Maven) injectProxy(id *api.Identity) (s string, err error) {
	m, err := mxj.NewMapXml([]byte(id.Settings))
	if err != nil {
		return
	}
	proxies, err := addon.Proxy.List()
	if err != nil {
		return
	}
	pList := []interface{}{}
	for _, p := range proxies {
		mp := mxj.Map{
			"id":       p.Kind,
			"active":   p.Enabled,
			"protocol": p.Kind,
			"host":     p.Host,
			"port":     p.Port,
		}
		if p.Identity != nil {
			pid, idErr := addon.Identity.Get(p.Identity.ID)
			if idErr != nil {
				err = idErr
				return
			}
			mp["username"] = pid.User
			mp["password"] = pid.Password
		}
		pList = append(pList, mp)
	}
	v, err := m.ValuesForPath("settings.proxies.proxy")
	if err != nil {
		return
	}
	err = m.SetValueForPath(
		mxj.Map{"proxy": append(v, pList...)},
		"settings.proxies")
	if err != nil {
		return
	}
	b, err := m.XmlIndent("", "  ")
	s = string(b)
	return
}
