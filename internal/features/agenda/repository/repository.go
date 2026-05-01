package repository

import (
	"context"
	"turnos-medicos/internal/features/agenda/models"
)

type AgendaRepository interface {
	CrearAgenda(ctx context.Context, agenda *models.Agenda) error
	ObtenerAgendaPorID(ctx context.Context, agendaID int64) (*models.Agenda, error)
	ListarAgendasPorMedico(ctx context.Context, medicoID int64) ([]*models.Agenda, error)
	ActualizarAgenda(ctx context.Context, agenda *models.Agenda) error
	DesactivarAgenda(ctx context.Context, id int64) error
	ActivarAgenda(ctx context.Context, id int64) error
}
