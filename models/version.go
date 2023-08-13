package models

import (
	"fmt"
	"sort"
	"strings"
)

type RemoteVersion struct {
	Version string              `json:"version"`
	Stable  bool                `json:"stable"`
	Files   []RemoteVersionFile `json:"files"`
}

func (v RemoteVersion) PrettyVersion() string {
	return strings.TrimPrefix(v.Version, "go")
}

func (v RemoteVersion) Archs() []string {
	vers := map[string]bool{}
	for _, v := range v.Files {
		if v.OS == "" || v.Arch == "" {
			continue
		}
		vers[fmt.Sprintf("%s-%s", v.OS, v.Arch)] = true
	}
	var varchs []string
	for k := range vers {
		varchs = append(varchs, k)
	}
	sort.Strings(varchs)
	return varchs
}

type RemoteVersionFile struct {
	Filename string `json:"filename"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	Sha256   string `json:"sha256"`
	Size     int    `json:"size"`
	Kind     string `json:"kind"`
}
