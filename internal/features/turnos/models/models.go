package models

import "time"

type Turno struct {
	ID         int64     `json:"id"`
	AgendaID   int64     `json:"agenda_id"`
	MedicoID   int64     `json:"medico_id"`
	PacienteID *int64    `json:"paciente_id"`
	Fecha      time.Time `json:"fecha"`
	HoraInicio string    `json:"hora_inicio"`
	HoraFin    string    `json:"hora_fin"`
	Estado     string    `json:"estado"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
