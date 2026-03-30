package dto

type CrearMedicoRequest struct {
	Nombre       string `json:"nombre" binding:"required"`
	Apellido     string `json:"apellido" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=8"`
	Matricula    string `json:"matricula" binding:"required"`
	Especialidad string `json:"especialidad" binding:"required"`
}

type ActualizarMedicoRequest struct {
	Nombre       string `json:"nombre"`
	Apellido     string `json:"apellido"`
	Email        string `json:"email" binding:"omitempty,email"`
	Especialidad string `json:"especialidad"`
}
