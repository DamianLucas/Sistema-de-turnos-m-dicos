package repository

import (
	"context"
	"turnos-medicos/internal/features/pacientes/models"
)

type PacienteRepository interface {
	CrearPaciente(ctx context.Context, p *models.Paciente) error
	ObtenerPacientePorID(ctx context.Context, pacienteID int64) (*models.Paciente, error)
	ObtenerPacientePorDNI(ctx context.Context, dni string) (*models.Paciente, error)
	ListarPacientesActivos(ctx context.Context) ([]*models.Paciente, error)
	DesactivarPaciente(ctx context.Context, pacienteID int64) error
	ActualizarPaciente(ctx context.Context, p *models.Paciente) error

	AsignarMedicoTratante(ctx context.Context, pacienteID, medicoID int64) error
	QuitarMedicoTratante(ctx context.Context, pacienteID int64) error
	ListarPacientesPorMedico(ctx context.Context, medicoID int64) ([]*models.Paciente, error)
}
