package dto

import "turnos-medicos/internal/features/users/models"

type CrearUsuarioRequest struct {
	Nombre   string     `json:"nombre" binding:"required"`
	Apellido string     `json:"apellido" binding:"required"`
	Email    string     `json:"email" binding:"required,email"`
	Password string     `json:"password" binding:"required,min=8"`
	Rol      models.Rol `json:"rol" binding:"required,oneof=admin medico administrativo"`
}

type ActualizarUsuarioRequest struct {
	Nombre   string     `json:"nombre"`
	Apellido string     `json:"apellido"`
	Email    string     `json:"email" binding:"omitempty,email"`
	Password string     `json:"password" binding:"omitempty,min=8"`
	Rol      models.Rol `json:"rol" binding:"omitempty,oneof=admin medico administrativo"`
}
