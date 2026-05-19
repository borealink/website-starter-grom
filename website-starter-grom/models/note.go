package models

import "time"

type Note struct {
	ID uint

	Title string
	Description string

	CreatedAt time.Time
	UpdatedAt time.Time
}
