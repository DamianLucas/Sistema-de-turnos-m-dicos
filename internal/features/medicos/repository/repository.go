package repository

import (
	"context"
	modelsMedico "turnos-medicos/internal/features/medicos/models"
	modelsUser "turnos-medicos/internal/features/users/models"
)

type MedicoRepository interface {
	CrearMedico(ctx context.Context, u *modelsUser.User, m *modelsMedico.Medico) error

	ObtenerMedicoPorID(ctx context.Context, medicoId int64) (*modelsMedico.Medico, error)
	ObtenerMedicoPorMatricula(ctx context.Context, matricula string) (*modelsMedico.Medico, error)

	ListarMedicosActivos(ctx context.Context) ([]*modelsMedico.Medico, error)
	ListarMedicosPorEspecialidad(ctx context.Context, especialidad string) ([]*modelsMedico.Medico, error)

	ActualizarMedico(ctx context.Context, m *modelsMedico.Medico) error
	DesactivarMedico(ctx context.Context, medicoID int64) error
	ActivarMedico(ctx context.Context, medicoID int64) error
}
