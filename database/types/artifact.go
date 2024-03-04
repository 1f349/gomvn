package types

import "time"

type Artifact struct {
	MvnGroup string    `json:"mvn_group"`
	Artifact string    `json:"artifact"`
	Version  string    `json:"version"`
	Modified time.Time `json:"modified"`
}
