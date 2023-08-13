package models

import "strings"

type InstalledVersion struct {
	Path    string
	RawName string
	Status  string
	Ok      bool
}

func (i InstalledVersion) Version() string {
	return strings.TrimPrefix(i.RawName, "go")
}
