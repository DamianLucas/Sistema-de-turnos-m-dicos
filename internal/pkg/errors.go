package pkg

import "errors"

var (
	ErrUsuarioNoEncontrado   = errors.New("usuario no encontrado")
	ErrUsuarioYaExiste       = errors.New("usuario ya existe")
	ErrEmailDuplicado        = errors.New("email ya esta registrado")
	ErrEmailRequerido        = errors.New("email es obligatorio")
	ErrPasswordRequerido     = errors.New("password es obligatorio")
	ErrCredencialesInvalidas = errors.New("credenciales invalidas")
	ErrUsuarioInactivo       = errors.New("usuario inactivo")
	ErrTurnoNoDisponible     = errors.New("turno no disponible")
	ErrIDInvalido            = errors.New("id invalido")

	ErrMedicoNoEncontrado = errors.New("medico no encontrado")
	ErrMedicoInactivo     = errors.New("medico inctivo")

	ErrMatriculaDuplicada    = errors.New("matricula ya esta registrada")
	ErrMatriculaRequerida    = errors.New("matricula requerida")
	ErrEspecialidadRequerida = errors.New("especialidad requerida")

	ErrErrorPersistencia = errors.New("error interno de persistencia")
)
