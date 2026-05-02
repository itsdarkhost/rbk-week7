package models

type City struct {
	Id     int    `db:"id" json:"id"`
	UserId int    `db:"user_id" json:"user_id"`
	Name   string `db:"name" json:"name"`
}
