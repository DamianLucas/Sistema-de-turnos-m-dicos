package dto

type CrearAgendaRequest struct {
	DiaSemana     int    `json:"dia_semana" binding:"required,min=1,max=7"`
	HoraInicio    string `json:"hora_inicio" binding:"required"`
	HoraFin       string `json:"hora_fin" binding:"required"`
	DuracionTurno int    `json:"duracion_turno" binding:"required,min=5,max=120"`
}

type ActualizarAgendaRequest struct {
	HoraInicio    string `json:"hora_inicio" binding:"required"`
	HoraFin       string `json:"hora_fin" binding:"required"`
	DuracionTurno int    `json:"duracion_turno" binding:"required,min=5,max=120"`
}
