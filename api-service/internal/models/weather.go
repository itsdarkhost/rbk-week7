package models

import "time"

type Weather struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
}

type WeatherHistory struct {
	Id          int       `db:"id" json:"id"`
	UserId      int       `db:"user_id" json:"user_id"`
	City        string    `db:"city" json:"city"`
	Temperature float64   `db:"temperature" json:"temperature"`
	Description string    `db:"description" json:"description"`
	RequestedAt time.Time `db:"requested_at" json:"requested_at"`
}

type WeatherHistoryResponse struct {
	UserId  int              `json:"user_id"`
	City    string           `json:"city,omitempty"`
	History []WeatherHistory `json:"history"`
}
