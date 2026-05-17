package routes

import (
	"turnos-medicos/internal/bootstrap"
	"turnos-medicos/internal/features/users/models"
	"turnos-medicos/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, h *bootstrap.Handlers) {

	v1 := r.Group("/api/v1")

	// =========================
	// AUTH (PÚBLICO)
	// =========================

	auth := v1.Group("/auth")
	{
		auth.POST("/login", h.Auth.Login)
	}

	// =========================
	// RUTAS PRIVADAS (JWT)
	// =========================

	private := v1.Group("/")
	private.Use(middleware.RequireAuth())
	{
		// =========================
		// USERS (solo admin)
		// =========================
		users := private.Group("/users")
		users.Use(middleware.RequireRol(models.RolAdmin))
		{
			users.POST("/", h.User.CrearUsuario)
			users.GET("/", h.User.ListarUsuariosActivos)
			users.GET("/:id", h.User.ObtenerUsuarioPorID)
			users.PUT("/:id", h.User.ActualizarUsuario)
			users.PATCH("/:id/desactivar", h.User.DesactivarUsuario)
		}
	}

	// =========================
	// MEDICOS
	// =========================

	medicos := private.Group("/medicos")
	medicos.Use(middleware.RequireRol(models.RolAdmin, models.RolAdministrativo))
	{
		medicos.POST("/", h.Medico.CrearMedico)

		// handler unificado que se usa para obtener medicos por especialidad y medicos activos con los queryParams
		medicos.GET("/", h.Medico.ListarMedicos)

		medicos.GET("/matricula/:matricula", h.Medico.ObtenerMedicoPorMatricula)
		medicos.GET("/:id", h.Medico.ObtenerMedicoPorID)
		medicos.GET("/:id/pacientes", h.Medico.ListarPacientesPorMedico) //<------ aqui

		medicos.PUT("/:id", h.Medico.ActualizarMedico)
		medicos.PATCH("/:id/desactivar", h.Medico.DesactivarMedico)
		medicos.PATCH("/:id/activar", h.Medico.ActivarMedico)
	}

	// =========================
	// PACIENTES
	// =========================

	//por el momento separare en dos grupos las rutas de paciente ya que aun no implemente ownership
	pacientes := private.Group("/pacientes")

	//escritura
	pacienteWrite := pacientes.Group("/")
	pacienteWrite.Use(middleware.RequireRol(models.RolAdmin, models.RolAdministrativo))
	{
		pacienteWrite.POST("/", h.Paciente.CrearPaciente)
		pacienteWrite.PUT("/:id", h.Paciente.ActualizarPaciente)
		pacienteWrite.PATCH("/:id/desactivar", h.Paciente.DesactivarPaciente)
		pacienteWrite.PATCH("/:id/activar", h.Paciente.ActivarPaciente)
		pacienteWrite.PATCH("/:id/asignar-medico/:medicoID", h.Paciente.AsignarMedicoTratante)
		pacienteWrite.DELETE("/:id/medico", h.Paciente.QuitarMedicoTratante)
	}

	//lectura
	pacienteRead := pacientes.Group("/")
	pacienteRead.Use(middleware.RequireRol(models.RolAdmin, models.RolAdministrativo, models.RolMedico))
	{
		pacienteRead.GET("/", h.Paciente.ListarPacientesActivos)
		pacienteRead.GET("/dni/:dni", h.Paciente.ObtenerPacientePorDNI)
		pacienteRead.GET("/:id", h.Paciente.ObtenerPacientePorID)
	}

	// =========================
	// AGENDA
	// =========================

	agendas := private.Group("/agendas")
	agendas.Use(middleware.RequireRol(models.RolAdmin, models.RolAdministrativo))
	{
		agendas.POST("/medico/:medicoID", h.Agenda.CrearAgenda)

		agendas.GET("/medico/:medicoID", h.Agenda.ListarAgendasPorMedico)
		agendas.GET("/:id", h.Agenda.ObtenerAgendaPorID)

		agendas.PUT("/:id", h.Agenda.ActualizarAgenda)

		agendas.PATCH("/:id/desactivar", h.Agenda.DesactivarAgenda)
		agendas.PATCH("/:id/activar", h.Agenda.ActivarAgenda)
	}

	// =========================
	// TURNOS
	// =========================

	turnos := private.Group("/turnos")

	// generacion
	turnosGeneracion := turnos.Group("/")
	turnosGeneracion.Use(middleware.RequireRol(models.RolAdmin, models.RolAdministrativo))
	{
		turnosGeneracion.POST("/agenda/:id/crear", h.Turno.CrearTurno)
	}

	//lectura
	turnosRead := turnos.Group("/")
	turnosRead.Use(middleware.RequireRol(models.RolAdmin, models.RolAdministrativo, models.RolMedico))
	{
		turnosRead.GET("/:id", h.Turno.ObtenerTurnoPorID)

		turnosRead.GET("/medico/:medicoID", h.Turno.ListarTurnosPorMedico)

		turnosRead.GET("/paciente/:pacienteID", h.Turno.ListarTurnosPorPaciente)

		turnosRead.GET("/disponibles/:medicoID", h.Turno.ListarTurnosDisponibles)
	}

	//reservas
	turnosReserva := turnos.Group("/")
	turnosReserva.Use(middleware.RequireRol(models.RolAdmin, models.RolAdministrativo))
	{
		turnosReserva.PUT("/:id/reservar", h.Turno.ReservarTurno)
		turnosReserva.PUT("/:id/liberar", h.Turno.LiberarTurno)
	}

	//estados medicos
	turnosEstado := turnos.Group("/")
	turnosEstado.Use(middleware.RequireRol(models.RolAdmin, models.RolMedico))
	{
		turnosEstado.PUT("/:id/atendido", h.Turno.MarcarTurnoAtendido)
		turnosEstado.PUT("/:id/no-asistio", h.Turno.MarcarTurnoNoAsistio)
	}
}
