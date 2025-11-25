package model

import "time"

type Lecturer struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	LecturerID string    `json:"lecturer_id"`
	Department string    `json:"department"`
	CreatedAt  time.Time `json:"created_at"`
}

type LecturerDetail struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	FullName    string    `json:"full_name"`
	LecturerID  string    `json:"lecturer_id"`
	Department  string    `json:"department"`
	CreatedAt   time.Time `json:"created_at"`
}
