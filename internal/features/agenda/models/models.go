package models

import "time"

type Agenda struct {
	ID            int64     `json:"id"`
	MedicoID      int64     `json:"medico_id"`
	DiaSemana     int       `json:"dia_semana"`
	HoraInicio    time.Time `json:"hora_inicio"`
	HoraFin       time.Time `json:"hora_fin"`
	DuracionTurno int       `json:"duracion_turno"`
	Activo        bool      `json:"activo"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
