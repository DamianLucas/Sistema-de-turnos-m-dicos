package models

import (
	"time"
)

/*
si bien los campos como Nombre, Apellido , Email y Activo estan en la tabla user los agregue
para no tener que hacer JOIN cada vez que necesite mostrar un médico completo.
*/
type Medico struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Nombre       string    `json:"nombre"`
	Apellido     string    `json:"apellido"`
	Email        string    `json:"email"`
	Matricula    string    `json:"matricula"`
	Especialidad string    `json:"especialidad"`
	Activo       bool      `json:"activo"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
