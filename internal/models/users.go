package models

import "time"

type User struct {
	Id        int        `db:"id" json:"id"`
	Username  string     `db:"username" json:"username"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
