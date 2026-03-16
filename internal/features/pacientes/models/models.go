package models

import "time"

type Paciente struct {
	ID              int64     `json:"id"`
	Nombre          string    `json:"nombre"`
	Apellido        string    `json:"apellido"`
	DNI             string    `json:"dni"`
	Email           string    `json:"email"`
	Telefono        string    `json:"telefono"`
	FechaNacimiento time.Time `json:"fecha_nacimiento"`
	Direccion       string    `json:"direccion"`
	ObraSocial      string    `json:"obra_social"`
	MedicoTratante  *int64    `json:"medico_tratante_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
