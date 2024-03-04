// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"database/sql"
)

type Artifact struct {
	MvnGroup string `json:"mvn_group"`
	Artifact string `json:"artifact"`
	Version  string `json:"version"`
	Modified string `json:"modified"`
}

type Path struct {
	UserID    sql.NullInt64 `json:"user_id"`
	Path      string        `json:"path"`
	Deploy    sql.NullInt64 `json:"deploy"`
	CreatedAt sql.NullTime  `json:"created_at"`
	UpdatedAt sql.NullTime  `json:"updated_at"`
}

type User struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	Admin     sql.NullInt64 `json:"admin"`
	TokenHash string        `json:"token_hash"`
	CreatedAt sql.NullTime  `json:"created_at"`
	UpdatedAt sql.NullTime  `json:"updated_at"`
}