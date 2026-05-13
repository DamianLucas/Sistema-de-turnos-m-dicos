package repository

import (
	"context"
	"turnos-medicos/internal/features/turnos/models"
)

type TurnoRepository interface {
	CrearTurno(ctx context.Context, turno *models.Turno) error

	ObtenerTurnoPorID(ctx context.Context, turnoID int64) (*models.Turno, error)

	ListarTurnosPorMedico(ctx context.Context, medicoID int64) ([]*models.Turno, error)
	ListarTurnosPorPaciente(ctx context.Context, pacienteID int64) ([]*models.Turno, error)
	ListarTurnosDisponibles(ctx context.Context, medicoID int64) ([]*models.Turno, error)

	ReservarTurno(ctx context.Context, turnoID int64, pacienteID int64) error
	LiberarTurno(ctx context.Context, turnoID int64) error

	MarcarTurnoAtendido(ctx context.Context, turnoID int64) error
	MarcarTurnoNoAsistio(ctx context.Context, turnoID int64) error
}
