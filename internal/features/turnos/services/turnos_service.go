package services

import (
	"context"
	"errors"
	"time"

	"turnos-medicos/internal/features/turnos/dto"
	"turnos-medicos/internal/features/turnos/models"
	"turnos-medicos/internal/features/turnos/repository"
	"turnos-medicos/internal/pkg"

	agendaRepository "turnos-medicos/internal/features/agenda/repository"
	medicoRepository "turnos-medicos/internal/features/medicos/repository"
	pacienteRepository "turnos-medicos/internal/features/pacientes/repository"

	"github.com/lib/pq"
)

type TurnoService interface {
	GenerarTurnos(ctx context.Context, agendaID int64, req dto.GenerarTurnosRequest) error
	ObtenerTurnoPorID(ctx context.Context, turnoID int64) (*models.Turno, error)
	ListarTurnosPorMedico(ctx context.Context, medicoID int64) ([]*models.Turno, error)
	ListarTurnosPorPaciente(ctx context.Context, pacienteID int64) ([]*models.Turno, error)
	ListarTurnosDisponibles(ctx context.Context, medicoID int64) ([]*models.Turno, error)
	ReservarTurno(ctx context.Context, turnoID int64, pacienteID int64) error
	LiberarTurno(ctx context.Context, turnoID int64) error
	MarcarTurnoAtendido(ctx context.Context, turnoID int64) error
	MarcarTurnoNoAsistio(ctx context.Context, turnoID int64) error
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

func (s *turnoService) GenerarTurnos(ctx context.Context, agendaID int64, req dto.GenerarTurnosRequest) error {

	// Validación básica de ID.
	if agendaID <= 0 {
		return pkg.ErrIDInvalido
	}

	//obtenemos agenda
	agenda, err := s.repoAgenda.ObtenerAgendaPorID(ctx, agendaID)
	if err != nil {
		return err
	}

	//no se generan turnos desde agendas inactivas
	if !agenda.Activo {
		return pkg.ErrAgendaInactiva
	}

	//obtenemos el medico asociado a la agenda
	medico, err := s.repoMedico.ObtenerMedicoPorID(ctx, agenda.MedicoID)
	if err != nil {
		return err
	}

	//no se generan turnos para medicos inactivos
	if !medico.Activo {
		return pkg.ErrMedicoInactivo
	}

	// Parseo de fechas usando el layout estándar de Go.
	//
	// "2006-01-02" NO es una fecha arbitraria.
	// Go utiliza esta fecha específica como referencia
	// para definir formatos.
	fechaInicio, err := time.Parse("2006-01-02", req.FechaInicio)
	if err != nil {
		return pkg.ErrFechaInvalida
	}

	fechaFin, err := time.Parse("2006-01-02", req.FechaFin)
	if err != nil {
		return pkg.ErrFechaInvalida
	}

	// La fecha inicial no puede ser posterior a la final.
	if fechaInicio.After(fechaFin) {
		return pkg.ErrRangoFechasInvalido
	}

	// Limita la generación a 90 días para evitar:
	// - generación masiva accidental
	// - consumo excesivo de memoria
	// - demasiados inserts
	if fechaFin.Sub(fechaInicio).Hours() > 24*90 {
		return pkg.ErrRangoFechasExcedido
	}

	horaInicio := agenda.HoraInicio
	horaFin := agenda.HoraFin

	// Validación defensiva.
	// Aunque ya existe CHECK en PostgreSQL,
	// el dominio también valida consistencia.
	if !horaFin.After(horaInicio) {
		return pkg.ErrAgendaInvalida
	}

	// Convierte minutos enteros a Duration.
	//
	// Ejemplo:
	// 30 -> 30m
	duracion := time.Duration(agenda.DuracionTurno) * time.Minute

	// Recorre todas las fechas del rango.
	for fecha := fechaInicio; !fecha.After(fechaFin); fecha = fecha.AddDate(0, 0, 1) {

		// Filtra únicamente los días que coinciden
		// con el día configurado en agenda.
		//
		// Weekday:
		// Sunday = 0
		// Monday = 1, etc
		if fecha.Weekday() != time.Weekday(agenda.DiaSemana%7) {
			continue
		}

		// Construye la fecha/hora real del inicio del slot.
		//
		// Combina:
		// - fecha actual
		// - hora de inicio de agenda
		slotInicio := time.Date(
			fecha.Year(),
			fecha.Month(),
			fecha.Day(),
			horaInicio.Hour(),
			horaInicio.Minute(),
			0,
			0,
			time.UTC,
		)

		// Hora máxima permitida para generar slots.
		slotFinLimite := time.Date(
			fecha.Year(),
			fecha.Month(),
			fecha.Day(),
			horaFin.Hour(),
			horaFin.Minute(),
			0,
			0,
			time.UTC,
		)

		// Genera slots consecutivos.
		for slotInicio.Before(slotFinLimite) {

			// Calcula el final del slot.
			slotFin := slotInicio.Add(duracion)

			// Evita generar turnos que excedan
			// el horario configurado.
			if slotFin.After(slotFinLimite) {
				break
			}

			turno := &models.Turno{
				AgendaID:   agenda.ID,
				MedicoID:   agenda.MedicoID,
				Fecha:      fecha,
				HoraInicio: slotInicio,
				HoraFin:    slotFin,
				Estado:     models.EstadoDisponible,
			}

			err := s.repoTurno.CrearTurno(ctx, turno)
			if err != nil {

				var pqErr *pq.Error
				// PostgreSQL:
				// 23505 = unique_violation
				//
				// Si el slot ya existe:
				// - se ignora
				// - continúa la generación
				//
				// Esto hace que el proceso sea idempotente.

				if errors.As(err, &pqErr) && pqErr.Code == "23505" {
					slotInicio = slotInicio.Add(duracion)
					continue
				}

				return err
			}

			// Avanza al siguiente bloque horario.
			slotInicio = slotInicio.Add(duracion)
		}
	}

	return nil
}

func (s *turnoService) ObtenerTurnoPorID(ctx context.Context, turnoID int64) (*models.Turno, error) {
	if turnoID <= 0 {
		return nil, pkg.ErrIDInvalido
	}

	turno, err := s.repoTurno.ObtenerTurnoPorID(ctx, turnoID)
	if err != nil {
		return nil, err
	}

	return turno, nil

}

func (s *turnoService) ListarTurnosPorMedico(ctx context.Context, medicoID int64) ([]*models.Turno, error) {
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

func (s *turnoService) MarcarTurnoAtendido(ctx context.Context, turnoID int64) error {
	if turnoID <= 0 {
		return pkg.ErrIDInvalido
	}

	turno, err := s.repoTurno.ObtenerTurnoPorID(ctx, turnoID)
	if err != nil {
		return err
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

func (s *turnoService) MarcarTurnoNoAsistio(ctx context.Context, turnoID int64) error {
	if turnoID <= 0 {
		return pkg.ErrIDInvalido
	}

	turno, err := s.repoTurno.ObtenerTurnoPorID(ctx, turnoID)
	if err != nil {
		return err
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
