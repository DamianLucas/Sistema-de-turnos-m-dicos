package services

import (
	"context"
	"fmt"
	"time"
	"turnos-medicos/internal/features/agenda/dto"
	"turnos-medicos/internal/features/agenda/models"
	userModel "turnos-medicos/internal/features/users/models"

	agendaRepo "turnos-medicos/internal/features/agenda/repository"
	medicoRepo "turnos-medicos/internal/features/medicos/repository"
	"turnos-medicos/internal/pkg"
)

type AgendaService interface {
	CrearAgenda(ctx context.Context, medicoID int64, req dto.CrearAgendaRequest) (*models.Agenda, error)
	ObtenerAgendaPorID(ctx context.Context, authUserID int64, authRol userModel.Rol, agendaID int64) (*models.Agenda, error)
	ListarAgendasPorMedico(ctx context.Context, authUserID int64, authRol userModel.Rol, medicoID int64) ([]*models.Agenda, error)
	ActualizarAgenda(ctx context.Context, authUserID int64, authRol userModel.Rol, agendaID int64, req dto.ActualizarAgendaRequest) (*models.Agenda, error)
	DesactivarAgenda(ctx context.Context, authUserID int64, authRol userModel.Rol, id int64) error
	ActivarAgenda(ctx context.Context, authUserID int64, authRol userModel.Rol, id int64) error
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

	horaInicio, err := time.Parse("15:04", req.HoraInicio)
	if err != nil {
		return nil, pkg.ErrHoraInvalida
	}

	horaFin, err := time.Parse("15:04", req.HoraFin)
	if err != nil {
		return nil, pkg.ErrHoraInvalida
	}

	if !horaFin.After(horaInicio) {
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
		HoraInicio:    horaInicio,
		HoraFin:       horaFin,
		DuracionTurno: req.DuracionTurno,
	}

	if err := s.repo.CrearAgenda(ctx, agenda); err != nil {
		return nil, err
	}

	return agenda, nil
}

func (s *agendaService) ObtenerAgendaPorID(ctx context.Context, authUserID int64, authRol userModel.Rol, agendaID int64) (*models.Agenda, error) {
	if agendaID <= 0 {
		return nil, pkg.ErrIDInvalido
	}

	agenda, err := s.repo.ObtenerAgendaPorID(ctx, agendaID)
	if err != nil {
		return nil, err
	}

	// OWNERSHIP
	if authRol == userModel.RolMedico {

		medico, err := s.repoMedico.ObtenerMedicoPorID(
			ctx,
			agenda.MedicoID,
		)
		if err != nil {
			return nil, err
		}

		if medico.UserID != authUserID {
			return nil, pkg.ErrForbidden
		}
	}

	return agenda, err

}

func (s *agendaService) ListarAgendasPorMedico(ctx context.Context, authUserID int64, authRol userModel.Rol, medicoID int64) ([]*models.Agenda, error) {
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

	// OWNERSHIP
	if authRol == userModel.RolMedico {

		if medico.UserID != authUserID {
			return nil, pkg.ErrForbidden
		}
	}

	agendas, err := s.repo.ListarAgendasPorMedico(ctx, medicoID)
	if err != nil {
		return nil, err
	}

	if agendas == nil {
		return []*models.Agenda{}, nil
	}

	return agendas, nil
}

func (s *agendaService) ActualizarAgenda(ctx context.Context, authUserID int64, authRol userModel.Rol, agendaID int64, req dto.ActualizarAgendaRequest) (*models.Agenda, error) {
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

	// OWNERSHIP
	if authRol == userModel.RolMedico {

		medico, err := s.repoMedico.ObtenerMedicoPorID(
			ctx,
			agendaActual.MedicoID,
		)
		if err != nil {
			return nil, err
		}

		if medico.UserID != authUserID {
			return nil, pkg.ErrForbidden
		}
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

	agendaActual.HoraInicio = inicio
	agendaActual.HoraFin = fin
	agendaActual.DuracionTurno = req.DuracionTurno

	err = s.repo.ActualizarAgenda(ctx, agendaActual)
	if err != nil {
		return nil, err
	}

	return agendaActual, nil
}

func (s *agendaService) DesactivarAgenda(ctx context.Context, authUserID int64, authRol userModel.Rol, id int64) error {
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

	// OWNERSHIP
	if authRol == userModel.RolMedico {

		medico, err := s.repoMedico.ObtenerMedicoPorID(
			ctx,
			agenda.MedicoID,
		)
		if err != nil {
			return err
		}

		if medico.UserID != authUserID {
			return pkg.ErrForbidden
		}
	}

	if err := s.repo.DesactivarAgenda(ctx, id); err != nil {
		return fmt.Errorf("error al desactivar agenda: %w", err)
	}

	return nil
}

func (s *agendaService) ActivarAgenda(ctx context.Context, authUserID int64, authRol userModel.Rol, id int64) error {
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
	// OWNERSHIP
	if authRol == userModel.RolMedico {

		medico, err := s.repoMedico.ObtenerMedicoPorID(
			ctx,
			agenda.MedicoID,
		)
		if err != nil {
			return err
		}

		if medico.UserID != authUserID {
			return pkg.ErrForbidden
		}
	}

	if err := s.repo.ActivarAgenda(ctx, id); err != nil {
		return fmt.Errorf("error al activar la agenda: %w", err)
	}

	return nil

}
