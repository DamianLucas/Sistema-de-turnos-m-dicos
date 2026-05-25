package services

import (
	"context"
	"time"

	"turnos-medicos/internal/features/turnos/dto"
	"turnos-medicos/internal/features/turnos/models"
	"turnos-medicos/internal/features/turnos/repository"
	userModel "turnos-medicos/internal/features/users/models"
	"turnos-medicos/internal/pkg"

	agendaRepository "turnos-medicos/internal/features/agenda/repository"
	medicoRepository "turnos-medicos/internal/features/medicos/repository"
	pacienteRepository "turnos-medicos/internal/features/pacientes/repository"
)

type TurnoService interface {
	CrearTurno(ctx context.Context, req dto.CrearTurnoRequest) (*models.Turno, error)
	ObtenerTurnoPorID(ctx context.Context, authUserID int64, authRol userModel.Rol, turnoID int64) (*models.Turno, error)
	ListarTurnosPorMedico(ctx context.Context, authUserID int64, authRol userModel.Rol, medicoID int64) ([]*models.Turno, error)
	ListarTurnosPorPaciente(ctx context.Context, pacienteID int64) ([]*models.Turno, error)
	ListarTurnosDisponibles(ctx context.Context, medicoID int64) ([]*models.Turno, error)
	ReservarTurno(ctx context.Context, turnoID int64, pacienteID int64) error
	LiberarTurno(ctx context.Context, turnoID int64) error
	MarcarTurnoAtendido(ctx context.Context, authUserID int64, authRol userModel.Rol, turnoID int64) error
	MarcarTurnoNoAsistio(ctx context.Context, authUserID int64, authRol userModel.Rol, turnoID int64) error
}

type turnoService struct {
	repoTurno    repository.TurnoRepository
	repoAgenda   agendaRepository.AgendaRepository
	repoMedico   medicoRepository.MedicoRepository
	repoPaciente pacienteRepository.PacienteRepository
}

func NewTurnoService(
	repoTurno repository.TurnoRepository,
	repoAgenda agendaRepository.AgendaRepository,
	repoMedico medicoRepository.MedicoRepository,
	repoPaciente pacienteRepository.PacienteRepository,
) TurnoService {
	return &turnoService{
		repoTurno:    repoTurno,
		repoAgenda:   repoAgenda,
		repoMedico:   repoMedico,
		repoPaciente: repoPaciente,
	}
}

func (s *turnoService) CrearTurno(ctx context.Context, req dto.CrearTurnoRequest) (*models.Turno, error) {

	if req.AgendaID <= 0 {
		return nil, pkg.ErrIDInvalido
	}

	// AGENDA
	agenda, err := s.repoAgenda.ObtenerAgendaPorID(ctx, req.AgendaID)
	if err != nil {
		return nil, err
	}

	if !agenda.Activo {
		return nil, pkg.ErrAgendaInactiva
	}

	// MEDICO
	medico, err := s.repoMedico.ObtenerMedicoPorID(ctx, agenda.MedicoID)
	if err != nil {
		return nil, err
	}

	if !medico.Activo {
		return nil, pkg.ErrMedicoInactivo
	}

	// FECHA
	fecha, err := time.Parse("2006-01-02", req.Fecha)
	if err != nil {
		return nil, pkg.ErrFechaInvalida
	}

	// HORA
	horaInicio, err := time.Parse("15:04", req.HoraInicio)
	if err != nil {
		return nil, pkg.ErrHoraInvalida
	}

	// VALIDAR DIA SEMANA
	if fecha.Weekday() != time.Weekday(agenda.DiaSemana%7) {
		return nil, pkg.ErrAgendaInvalida
	}

	// CONSTRUIR DATETIME REAL
	slotInicio := time.Date(
		fecha.Year(), fecha.Month(), fecha.Day(),
		horaInicio.Hour(), horaInicio.Minute(),
		0, 0, time.UTC,
	)

	agendaInicio := time.Date(
		fecha.Year(), fecha.Month(), fecha.Day(),
		agenda.HoraInicio.Hour(), agenda.HoraInicio.Minute(),
		0, 0, time.UTC,
	)

	agendaFin := time.Date(
		fecha.Year(), fecha.Month(), fecha.Day(),
		agenda.HoraFin.Hour(), agenda.HoraFin.Minute(),
		0, 0, time.UTC,
	)

	if slotInicio.Before(agendaInicio) || !slotInicio.Before(agendaFin) {
		return nil, pkg.ErrHorarioFueraAgenda
	}

	duracion := time.Duration(agenda.DuracionTurno) * time.Minute
	slotFin := slotInicio.Add(duracion)

	if slotFin.After(agendaFin) {
		return nil, pkg.ErrHorarioFueraAgenda
	}

	if slotInicio.Before(time.Now()) {
		return nil, pkg.ErrTurnoExpirado
	}

	turno := &models.Turno{
		AgendaID:   agenda.ID,
		MedicoID:   agenda.MedicoID,
		PacienteID: nil,
		Fecha:      fecha,
		HoraInicio: slotInicio,
		HoraFin:    slotFin,
		Estado:     models.EstadoDisponible,
	}

	if err := s.repoTurno.CrearTurno(ctx, turno); err != nil {
		return nil, err
	}

	return turno, nil
}

func (s *turnoService) ObtenerTurnoPorID(ctx context.Context, authUserID int64, authRol userModel.Rol, turnoID int64) (*models.Turno, error) {
	if turnoID <= 0 {
		return nil, pkg.ErrIDInvalido
	}

	turno, err := s.repoTurno.ObtenerTurnoPorID(ctx, turnoID)
	if err != nil {
		return nil, err
	}

	// OWNERSHIP
	if authRol == userModel.RolMedico {
		medico, err := s.repoMedico.ObtenerMedicoPorID(ctx, turno.MedicoID)
		if err != nil {
			return nil, err
		}

		if medico.UserID != authUserID {
			return nil, pkg.ErrForbidden
		}
	}
	return turno, nil

}

func (s *turnoService) ListarTurnosPorMedico(ctx context.Context, authUserID int64, authRol userModel.Rol, medicoID int64) ([]*models.Turno, error) {
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

	turnos, err := s.repoTurno.ListarTurnosPorMedico(ctx, medicoID)
	if err != nil {
		return nil, err
	}

	if turnos == nil {
		return []*models.Turno{}, nil
	}

	return turnos, nil

}

func (s *turnoService) ListarTurnosPorPaciente(ctx context.Context, pacienteID int64) ([]*models.Turno, error) {
	if pacienteID <= 0 {
		return nil, pkg.ErrIDInvalido
	}

	_, err := s.repoPaciente.ObtenerPacientePorID(ctx, pacienteID)
	if err != nil {
		return nil, err
	}

	turnos, err := s.repoTurno.ListarTurnosPorPaciente(ctx, pacienteID)
	if err != nil {
		return nil, err
	}

	if turnos == nil {
		return []*models.Turno{}, nil
	}

	return turnos, nil

}

func (s *turnoService) ListarTurnosDisponibles(ctx context.Context, medicoID int64) ([]*models.Turno, error) {
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

	turnos, err := s.repoTurno.ListarTurnosDisponibles(ctx, medicoID)
	if err != nil {
		return nil, err
	}

	if turnos == nil {
		return []*models.Turno{}, nil
	}

	return turnos, nil

}

func (s *turnoService) ReservarTurno(ctx context.Context, turnoID int64, pacienteID int64) error {

	if turnoID <= 0 || pacienteID <= 0 {
		return pkg.ErrIDInvalido
	}

	turno, err := s.repoTurno.ObtenerTurnoPorID(ctx, turnoID)
	if err != nil {
		return err
	}

	// Solo se pueden reservar turnos disponibles.
	if turno.Estado != models.EstadoDisponible {
		return pkg.ErrTurnoNoDisponible
	}

	paciente, err := s.repoPaciente.ObtenerPacientePorID(ctx, pacienteID)
	if err != nil {
		return err
	}

	// Solo pacientes activos pueden reservar.
	if !paciente.Activo {
		return pkg.ErrPacienteInactivo
	}

	// Evita reservar turnos pasados
	ahora := time.Now()

	fechaTurno := time.Date(
		turno.Fecha.Year(),
		turno.Fecha.Month(),
		turno.Fecha.Day(),
		turno.HoraInicio.Hour(),
		turno.HoraInicio.Minute(),
		0,
		0,
		time.UTC,
	)

	if fechaTurno.Before(ahora) {
		return pkg.ErrTurnoExpirado
	}

	if err := s.repoTurno.ReservarTurno(ctx, turnoID, pacienteID); err != nil {
		return err
	}

	return nil
}

func (s *turnoService) LiberarTurno(ctx context.Context, turnoID int64) error {
	if turnoID <= 0 {
		return pkg.ErrIDInvalido
	}

	turno, err := s.repoTurno.ObtenerTurnoPorID(ctx, turnoID)
	if err != nil {
		return err
	}

	//solo se liberan turnos reservados
	if turno.Estado != models.EstadoReservado {
		return pkg.ErrTurnoNoReservado
	}

	// Debe existir un paciente asignado.
	if turno.PacienteID == nil {
		return pkg.ErrTurnoSinPaciente
	}

	if err := s.repoTurno.LiberarTurno(ctx, turnoID); err != nil {
		return err
	}

	return nil

}

func (s *turnoService) MarcarTurnoAtendido(ctx context.Context, authUserID int64, authRol userModel.Rol, turnoID int64) error {
	if turnoID <= 0 {
		return pkg.ErrIDInvalido
	}

	turno, err := s.repoTurno.ObtenerTurnoPorID(ctx, turnoID)
	if err != nil {
		return err
	}

	// OWNERSHIP
	if authRol == userModel.RolMedico {

		medico, err := s.repoMedico.ObtenerMedicoPorID(ctx, turno.MedicoID)
		if err != nil {
			return err
		}

		if medico.UserID != authUserID {
			return pkg.ErrForbidden
		}
	}

	// Solo un turno reservado puede marcarse como atendido.
	if turno.Estado != models.EstadoReservado {
		return pkg.ErrEstadoTurnoInvalido
	}

	// Validación defensiva.
	if turno.PacienteID == nil {
		return pkg.ErrTurnoSinPaciente
	}

	if err := s.repoTurno.MarcarTurnoAtendido(ctx, turnoID); err != nil {
		return err
	}

	return nil
}

func (s *turnoService) MarcarTurnoNoAsistio(ctx context.Context, authUserID int64, authRol userModel.Rol, turnoID int64) error {
	if turnoID <= 0 {
		return pkg.ErrIDInvalido
	}

	turno, err := s.repoTurno.ObtenerTurnoPorID(ctx, turnoID)
	if err != nil {
		return err
	}

	// OWNERSHIP
	if authRol == userModel.RolMedico {

		medico, err := s.repoMedico.ObtenerMedicoPorID(ctx, turno.MedicoID)
		if err != nil {
			return err
		}

		if medico.UserID != authUserID {
			return pkg.ErrForbidden
		}
	}

	// Solo un turno reservado puede marcarse como no asistió
	if turno.Estado != models.EstadoReservado {
		return pkg.ErrEstadoTurnoInvalido
	}

	// Validación defensiva.
	if turno.PacienteID == nil {
		return pkg.ErrTurnoSinPaciente
	}

	if err := s.repoTurno.MarcarTurnoNoAsistio(ctx, turnoID); err != nil {
		return err
	}

	return nil

}
