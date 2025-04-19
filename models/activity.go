package models

import "time"

type Activity struct {
	Type string
	CreatedAt time.Time
	UpdatedAt time.Time
}