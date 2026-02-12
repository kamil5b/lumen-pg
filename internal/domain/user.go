package domain

import "time"

type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
}

type CreateUserInput struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
