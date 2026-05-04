package services

import (
	"context"
	"fmt"
	"time"
	"turnos-medicos/internal/features/agenda/dto"
	"turnos-medicos/internal/features/agenda/models"
	agendaRepo "turnos-medicos/internal/features/agenda/repository"
	medicoRepo "turnos-medicos/internal/features/medicos/repository"
	"turnos-medicos/internal/pkg"
)

type AgendaService interface {
	CrearAgenda(ctx context.Context, medicoID int64, req dto.CrearAgendaRequest) (*models.Agenda, error)
	ObtenerAgendaPorID(ctx context.Context, agendaID int64) (*models.Agenda, error)
	ListarAgendasPorMedico(ctx context.Context, medicoID int64) ([]*models.Agenda, error)
	ActualizarAgenda(ctx context.Context, agendaID int64, req dto.ActualizarAgendaRequest) (*models.Agenda, error)
	DesactivarAgenda(ctx context.Context, id int64) error
	ActivarAgenda(ctx context.Context, id int64) error
}

type agendaService struct {
	repo       agendaRepo.AgendaRepository
	repoMedico medicoRepo.MedicoRepository
}

func NewAgendaService(repo agendaRepo.AgendaRepository, repoMedico medicoRepo.MedicoRepository) AgendaService {
	return &agendaService{
		repo:       repo,
		repoMedico: repoMedico,
	}
}

func (s *agendaService) CrearAgenda(ctx context.Context, medicoID int64, req dto.CrearAgendaRequest) (*models.Agenda, error) {
	if medicoID <= 0 {
		return nil, pkg.ErrIDInvalido
	}

	const layoutHora = "15:04" //validacion horario - convención del lenguaje

	inicio, err := time.Parse(layoutHora, req.HoraInicio)
	if err != nil {
		return nil, pkg.ErrAgendaInvalida
	}

	fin, err := time.Parse(layoutHora, req.HoraFin)
	if err != nil {
		return nil, pkg.ErrAgendaInvalida
	}

	if !fin.After(inicio) {
		return nil, pkg.ErrAgendaInvalida
	}

	medico, err := s.repoMedico.ObtenerMedicoPorID(ctx, medicoID)
	if err != nil {
		return nil, err
	}

	if !medico.Activo {
		return nil, pkg.ErrMedicoInactivo
	}

	agenda := &models.Agenda{
		MedicoID:      medicoID,
		DiaSemana:     req.DiaSemana,
		HoraInicio:    req.HoraInicio,
		HoraFin:       req.HoraFin,
		DuracionTurno: req.DuracionTurno,
	}

	err = s.repo.CrearAgenda(ctx, agenda)
	if err != nil {
		return nil, err
	}

	return agenda, nil
}

func (s *agendaService) ObtenerAgendaPorID(ctx context.Context, agendaID int64) (*models.Agenda, error) {
	if agendaID <= 0 {
		return nil, pkg.ErrIDInvalido
	}

	agenda, err := s.repo.ObtenerAgendaPorID(ctx, agendaID)
	if err != nil {
		return nil, err
	}

	return agenda, err

}

func (s *agendaService) ListarAgendasPorMedico(ctx context.Context, medicoID int64) ([]*models.Agenda, error) {
	if medicoID <= 0 {
		return nil, pkg.ErrIDInvalido
	}

	medico, err := s.repoMedico.ObtenerMedicoPorID(ctx, medicoID)
	if err != nil {
		return nil, err
	}

	if !medico.Activo {
		return nil, pkg.ErrMedicoInactivo
	}

	agendas, err := s.repo.ListarAgendasPorMedico(ctx, medicoID)
	if err != nil {
		return nil, err
	}

	if agendas == nil {
		agendas = []*models.Agenda{}
	}

	return agendas, nil
}

func (s *agendaService) ActualizarAgenda(ctx context.Context, agendaID int64, req dto.ActualizarAgendaRequest) (*models.Agenda, error) {
	if agendaID <= 0 {
		return nil, pkg.ErrIDInvalido
	}

	agendaActual, err := s.repo.ObtenerAgendaPorID(ctx, agendaID)
	if err != nil {
		return nil, err
	}

	if !agendaActual.Activo {
		return nil, pkg.ErrAgendaInactiva
	}

	const layout = "15:04"

	inicio, err := time.Parse(layout, req.HoraInicio)
	if err != nil {
		return nil, pkg.ErrAgendaInvalida
	}

	fin, err := time.Parse(layout, req.HoraFin)
	if err != nil {
		return nil, pkg.ErrAgendaInvalida
	}

	if !fin.After(inicio) {
		return nil, pkg.ErrAgendaInvalida
	}

	agendaActual.HoraInicio = req.HoraInicio
	agendaActual.HoraFin = req.HoraFin
	agendaActual.DuracionTurno = req.DuracionTurno

	err = s.repo.ActualizarAgenda(ctx, agendaActual)
	if err != nil {
		return nil, err
	}

	return agendaActual, nil
}

func (s *agendaService) DesactivarAgenda(ctx context.Context, id int64) error {
	if id <= 0 {
		return pkg.ErrIDInvalido
	}

	agenda, err := s.repo.ObtenerAgendaPorID(ctx, id)
	if err != nil {
		return err
	}

	if !agenda.Activo {
		return nil
	}

	if err := s.repo.DesactivarAgenda(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", pkg.ErrDesactivarAgenda, err)
	}

	return nil
}

func (s *agendaService) ActivarAgenda(ctx context.Context, id int64) error {
	if id <= 0 {
		return pkg.ErrIDInvalido
	}

	agenda, err := s.repo.ObtenerAgendaPorID(ctx, id)
	if err != nil {
		return err
	}

	if agenda.Activo {
		return nil
	}

	if err := s.repo.ActivarAgenda(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", pkg.ErrActivarAgenda, err)
	}

	return nil

}
