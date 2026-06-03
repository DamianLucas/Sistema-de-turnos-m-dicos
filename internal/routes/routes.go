package routes

import (
	"turnos-medicos/internal/bootstrap"
	"turnos-medicos/internal/features/users/models"
	"turnos-medicos/internal/middleware"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "turnos-medicos/docs"
)

func SetupRoutes(r *gin.Engine, h *bootstrap.Handlers) {

	r.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

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

	medicosAdmin := medicos.Group("/")
	medicosAdmin.Use(middleware.RequireRol(models.RolAdmin, models.RolAdministrativo))
	{
		medicosAdmin.POST("/", h.Medico.CrearMedico)
		// handler unificado que se usa para obtener medicos por especialidad y medicos activos con los queryParams
		medicosAdmin.GET("/", h.Medico.ListarMedicos)
		medicosAdmin.GET("/matricula/:matricula", h.Medico.ObtenerMedicoPorMatricula)
		medicosAdmin.GET("/:id", h.Medico.ObtenerMedicoPorID)
		medicosAdmin.PUT("/:id", h.Medico.ActualizarMedico)
		medicosAdmin.PATCH("/:id/desactivar", h.Medico.DesactivarMedico)
		medicosAdmin.PATCH("/:id/activar", h.Medico.ActivarMedico)
	}
	//ownership
	medicosOwner := medicos.Group("/")
	medicosOwner.Use(middleware.RequireRol(models.RolAdmin, models.RolAdministrativo, models.RolMedico))
	{
		medicosOwner.GET("/:id/pacientes", h.Medico.ListarPacientesPorMedico)
	}

	// =========================
	// PACIENTES
	// =========================

	pacientes := private.Group("/pacientes")

	pacienteAdmin := pacientes.Group("/")
	pacienteAdmin.Use(middleware.RequireRol(models.RolAdmin, models.RolAdministrativo))
	{
		pacienteAdmin.POST("/", h.Paciente.CrearPaciente)
		pacienteAdmin.PATCH("/:id/asignar-medico/:medicoID", h.Paciente.AsignarMedicoTratante)
		pacienteAdmin.DELETE("/:id/medico", h.Paciente.QuitarMedicoTratante)
		pacienteAdmin.GET("/dni/:dni", h.Paciente.ObtenerPacientePorDNI)
		pacienteAdmin.PATCH("/:id/desactivar", h.Paciente.DesactivarPaciente)
		pacienteAdmin.PATCH("/:id/activar", h.Paciente.ActivarPaciente)
	}
	//ownership
	pacienteOwner := pacientes.Group("/")
	pacienteOwner.Use(middleware.RequireRol(models.RolAdmin, models.RolAdministrativo, models.RolMedico))
	{
		pacienteOwner.GET("/", h.Paciente.ListarPacientesActivos)
		pacienteOwner.GET("/:id", h.Paciente.ObtenerPacientePorID)
		pacienteOwner.PUT("/:id", h.Paciente.ActualizarPaciente)
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

	turnosAdmin := turnos.Group("/")
	turnosAdmin.Use(middleware.RequireRol(models.RolAdmin, models.RolAdministrativo))
	{
		turnosAdmin.POST("/agenda/:id/crear", h.Turno.CrearTurno)
		turnosAdmin.GET("/paciente/:pacienteID", h.Turno.ListarTurnosPorPaciente)
		turnosAdmin.GET("/disponibles/:medicoID", h.Turno.ListarTurnosDisponibles)
		turnosAdmin.PUT("/:id/reservar", h.Turno.ReservarTurno)
		turnosAdmin.PUT("/:id/liberar", h.Turno.LiberarTurno)
		turnosAdmin.PUT("/:id/atendido", h.Turno.MarcarTurnoAtendido)
		turnosAdmin.PUT("/:id/no-asistio", h.Turno.MarcarTurnoNoAsistio)
	}
	//Ownership
	turnosRead := turnos.Group("/")
	turnosRead.Use(middleware.RequireRol(models.RolAdmin, models.RolAdministrativo, models.RolMedico))
	{
		turnosRead.GET("/:id", h.Turno.ObtenerTurnoPorID)
		turnosRead.GET("/medico/:medicoID", h.Turno.ListarTurnosPorMedico)
	}
}
