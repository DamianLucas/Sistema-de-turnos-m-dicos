package pkg

import (
	"errors"
)

var (

	// =========================================================
	// GENERALES
	// =========================================================

	ErrIDInvalido    = NewBadRequestError("id invalido")
	ErrHoraInvalida  = NewBadRequestError("hora invalida")
	ErrFechaInvalida = NewBadRequestError("fecha invalida")

	ErrErrorPersistencia = errors.New("error interno de persistencia")

	ErrForbidden = NewForbiddenError("no autorizado para acceder a este recurso")

	// =========================================================
	// USERS
	// =========================================================

	// NOT FOUND
	ErrUsuarioNoEncontrado = NewNotFoundError(
		"usuario no encontrado",
	)

	// BAD REQUEST
	ErrUsuarioInactivo = NewBadRequestError(
		"usuario inactivo",
	)

	ErrEmailRequerido = NewBadRequestError(
		"email es obligatorio",
	)

	ErrPasswordRequerido = NewBadRequestError(
		"password es obligatorio",
	)

	// CONFLICT
	ErrUsuarioYaExiste = NewConflictError(
		"usuario ya existe",
	)

	ErrEmailDuplicado = NewConflictError(
		"email ya esta registrado",
	)

	// UNAUTHORIZED
	ErrCredencialesInvalidas = NewUnauthorizedError(
		"credenciales invalidas",
	)

	// =========================================================
	// MEDICOS
	// =========================================================

	// NOT FOUND
	ErrMedicoNoEncontrado = NewNotFoundError(
		"medico no encontrado",
	)

	// BAD REQUEST
	ErrMedicoInactivo = NewBadRequestError(
		"medico inactivo",
	)

	ErrMatriculaRequerida = NewBadRequestError(
		"matricula requerida",
	)

	ErrEspecialidadRequerida = NewBadRequestError(
		"especialidad requerida",
	)

	// CONFLICT
	ErrMatriculaDuplicada = NewConflictError(
		"matricula ya esta registrada",
	)

	// =========================================================
	// PACIENTES
	// =========================================================

	// NOT FOUND
	ErrPacienteNoEncontrado = NewNotFoundError(
		"paciente no encontrado",
	)

	// BAD REQUEST
	ErrPacienteInactivo = NewBadRequestError(
		"paciente inactivo",
	)

	ErrPacienteYaActivo = NewBadRequestError(
		"paciente ya activo",
	)

	ErrDNIInvalido = NewBadRequestError(
		"DNI invalido",
	)

	ErrDNIrequerido = NewBadRequestError(
		"DNI requerido",
	)

	// CONFLICT
	ErrDNIDuplicado = NewConflictError(
		"DNI duplicado",
	)

	// =========================================================
	// AGENDA
	// =========================================================

	// NOT FOUND
	ErrAgendaNoEncontrada = NewNotFoundError(
		"agenda no encontrada",
	)

	// BAD REQUEST
	ErrAgendaInvalida = NewBadRequestError(
		"datos de agenda invalidos",
	)

	ErrAgendaInactiva = NewBadRequestError(
		"agenda inactiva",
	)

	// CONFLICT
	ErrAgendaDuplicada = NewConflictError(
		"agenda ya existe para ese dia",
	)

	// =========================================================
	// TURNOS
	// =========================================================

	// NOT FOUND
	ErrTurnoNoEncontrado = NewNotFoundError(
		"turno no encontrado",
	)

	// BAD REQUEST
	ErrTurnoNoDisponible = NewBadRequestError(
		"turno no disponible",
	)

	ErrTurnoExpirado = NewBadRequestError(
		"turno expirado",
	)

	ErrTurnoNoReservado = NewBadRequestError(
		"turno no reservado",
	)

	ErrTurnoSinPaciente = NewBadRequestError(
		"turno sin paciente asignado",
	)

	ErrEstadoTurnoInvalido = NewBadRequestError(
		"estado de turno invalido",
	)

	ErrHorarioFueraAgenda = NewBadRequestError(
		"horario fuera de agenda",
	)

	// CONFLICT
	ErrTurnoDuplicado = NewConflictError(
		"ya existe un turno para ese horario",
	)
)
