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
	ErrListarUsuariosActivos = errors.New("error al listar usuarios activos")
	ErrDesactivarUsuario     = errors.New("error desactivando usuario")

	//medicos
	ErrMedicoNoEncontrado   = errors.New("medico no encontrado")
	ErrMedicoInactivo       = errors.New("medico inactivo")
	ErrListarMedicosActivos = errors.New("error al listar medicos activos")
	ErrDesactivarMedico     = errors.New("error desactivando medico")

	ErrMatriculaDuplicada    = errors.New("matricula ya esta registrada")
	ErrMatriculaRequerida    = errors.New("matricula requerida")
	ErrEspecialidadRequerida = errors.New("especialidad requerida")

	//pacientes
	ErrPacienteNoEncontrado     = errors.New("paciente no encontrado")
	ErrPacienteInactivo         = errors.New("paciente inactivo")
	ErrDNIDuplicado             = errors.New("DNI duplicado")
	ErrDNIrequerido             = errors.New("DNI requiro")
	ErrDNIInvalido              = errors.New("DNI invalido")
	ErrListarPacientesActivos   = errors.New("error al listar pacientes activos")
	ErrDesactivarPaciente       = errors.New("error desactivando paciente")
	ErrActualizarPaciente       = errors.New("error al actualizar paciente")
	ErrAsignarMedicoPaciente    = errors.New("error al asignar medico tratante")
	ErrQuitarMedicoPaciente     = errors.New("error al quitar medico tratante")
	ErrListarPacientesPorMedico = errors.New("error al listar pacientes por medico")

	//agenda
	//turnos
	ErrTurnoNoDisponible = errors.New("turno no disponible")

	ErrErrorPersistencia = errors.New("error interno de persistencia")
)
