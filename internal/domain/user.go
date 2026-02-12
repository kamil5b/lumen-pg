package domain

import "time"

// User represents a PostgreSQL user/role in the system
type User struct {
	Username  string
	CreatedAt time.Time
}

// LoginInput represents user login credentials
type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
