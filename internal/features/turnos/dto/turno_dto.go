package dto

type CrearTurnoRequest struct {
	AgendaID   int64  `json:"agenda_id" binding:"required"`
	Fecha      string `json:"fecha" binding:"required"`
	HoraInicio string `json:"hora_inicio" binding:"required"`
}

// para handler
type ReservarTurnoRequest struct {
	PacienteID int64 `json:"paciente_id" binding:"required"`
}
