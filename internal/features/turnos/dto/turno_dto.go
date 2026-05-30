package dto

import "time"

type CrearTurnoRequest struct {
	AgendaID   int64  `json:"agenda_id" binding:"required"`
	Fecha      string `json:"fecha" binding:"required"`
	HoraInicio string `json:"hora_inicio" binding:"required"`
}

// response para desacoplar
type TurnoDisponibleResponse struct {
	ID           int64     `json:"id"`
	Fecha        time.Time `json:"fecha"`
	HoraInicio   string    `json:"hora_inicio"`
	MedicoID     int64     `json:"medico_id"`
	MedicoNombre string    `json:"medico_nombre"`
	Especialidad string    `json:"especialidad"`
}

// para handler
type ReservarTurnoRequest struct {
	PacienteID int64 `json:"paciente_id" binding:"required"`
}
