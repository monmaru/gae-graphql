package model

import (
	"time"
)

type User struct {
	ID    string `json:"id" datastore:"-"`
	Name  string `json:"name"`
	EMail string `json:"email"`
}

type Blog struct {
	ID        string    `json:"id" datastore:"-"`
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
}

type BlogList struct {
	Nodes      []Blog `json:"nodes"`
	TotalCount int    `json:"totalCount"`
}
