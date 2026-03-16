package models

import "time"

type Agenda struct {
	ID            int64     `json:"id"`
	MedicoID      int64     `json:"medico_id"`
	DiaSemana     int64     `json:"dia_semana"`
	HoraInicio    string    `json:"hora_inicio"`
	HoraFin       string    `json:"hora_fin"`
	DuracionTurno int64     `json:"duracion_turno"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
