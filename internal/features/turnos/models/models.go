package models

import "time"

type EstadoTurno string

const (
	EstadoDisponible EstadoTurno = "disponible"
	EstadoReservado  EstadoTurno = "reservado"
	EstadoAtendido   EstadoTurno = "atendido"
	EstadoNoAsistio  EstadoTurno = "no_asistio"
)

type Turno struct {
	ID         int64       `json:"id"`
	AgendaID   int64       `json:"agenda_id"`
	MedicoID   int64       `json:"medico_id"`
	PacienteID *int64      `json:"paciente_id,omitempty"`
	Fecha      time.Time   `json:"fecha"`
	HoraInicio time.Time   `json:"hora_inicio"`
	HoraFin    time.Time   `json:"hora_fin"`
	Estado     EstadoTurno `json:"estado"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}
