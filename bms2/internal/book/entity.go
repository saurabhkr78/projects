package book

import "time"

type Book struct {
	ID          int64
	Title       string
	Author      string
	ISBN        string
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
