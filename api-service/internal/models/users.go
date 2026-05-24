package models

import "time"

type User struct {
	Id           int        `db:"id" json:"id"`
	Username     string     `db:"username" json:"username"`
	Email        string     `db:"email" json:"email"`
	PasswordHash string     `db:"password_hash" json:"-"`
	Role         string     `db:"role" json:"role"`
	DeletedAt    *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
