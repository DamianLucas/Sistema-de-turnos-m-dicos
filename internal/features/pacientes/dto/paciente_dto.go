package dto

import "time"

type CrearPacienteRequest struct {
	Nombre          string    `json:"nombre" binding:"required"`
	Apellido        string    `json:"apellido" binding:"required"`
	DNI             string    `json:"dni" binding:"required"`
	FechaNacimiento time.Time `json:"fecha_nacimiento" binding:"required"`
	Telefono        string    `json:"telefono"`
	Email           string    `json:"email" binding:"omitempty,email"`
	Direccion       string    `json:"direccion"`
	ObraSocial      string    `json:"obra_social"`
}

type ActualizarPacienteRequest struct {
	Nombre     string `json:"nombre"`
	Apellido   string `json:"apellido"`
	Telefono   string `json:"telefono"`
	Email      string `json:"email" binding:"omitempty,email"`
	Direccion  string `json:"direccion"`
	ObraSocial string `json:"obra_social"`
}
