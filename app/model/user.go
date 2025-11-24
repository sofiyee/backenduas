package model

import "time"

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"password,omitempty"`      // INPUT DARI USER
	PasswordHash string    `json:"-"`                       // TIDAK DIKIRIM KE CLIENT
	FullName     string    `json:"full_name"`
	RoleID       string    `json:"role_id"`
	RoleName     string    `json:"role_name"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}


