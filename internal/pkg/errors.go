package pkg

import "errors"

var (

	//users
	ErrUsuarioNoEncontrado   = errors.New("usuario no encontrado")
	ErrUsuarioYaExiste       = errors.New("usuario ya existe")
	ErrEmailDuplicado        = errors.New("email ya esta registrado")
	ErrEmailRequerido        = errors.New("email es obligatorio")
	ErrPasswordRequerido     = errors.New("password es obligatorio")
	ErrCredencialesInvalidas = errors.New("credenciales invalidas")
	ErrUsuarioInactivo       = errors.New("usuario inactivo")
	ErrIDInvalido            = errors.New("id invalido")

	//medicos
	ErrMedicoNoEncontrado = errors.New("medico no encontrado")
	ErrMedicoInactivo     = errors.New("medico inactivo")

	ErrMatriculaDuplicada    = errors.New("matricula ya esta registrada")
	ErrMatriculaRequerida    = errors.New("matricula requerida")
	ErrEspecialidadRequerida = errors.New("especialidad requerida")

	//pacientes
	ErrPacienteNoEncontrado = errors.New("paciente no encontrado")
	ErrPacienteInactivo     = errors.New("paciente inactivo")
	ErrDNIDuplicado         = errors.New("DNI duplicado")
	ErrDNIrequerido         = errors.New("DNI requiro")

	//agenda
	//turnos
	ErrTurnoNoDisponible = errors.New("turno no disponible")

	ErrErrorPersistencia = errors.New("error interno de persistencia")
)
