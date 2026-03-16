package models

import "time"

type Medico struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Matricula    string    `json:"matricula"`
	Especialidad string    `json:"especialidad"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
