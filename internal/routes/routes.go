package routes

import (
	"turnos-medicos/internal/bootstrap"
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
		users.Use(middleware.RequireRol("admin"))
		{
			users.POST("/", h.User.CrearUsuario)
			users.GET("/", h.User.ListarUsuariosActivos)
			users.GET("/:id", h.User.ObtenerUsuarioPorID)
			users.PUT("/:id", h.User.ActualizarUsuario)
			users.PATCH("/:id/desactivar", h.User.DesactivarUsuario)
		}
	}
}
