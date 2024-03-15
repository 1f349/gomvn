package database

import "strings"

func (a Artifact) GetPath() string {
	return strings.Replace(a.MvnGroup, ".", "/", -1) + "/" + a.Artifact + "/" + a.Version
}
