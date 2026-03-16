package routes

import (
	"turnos-medicos/internal/bootstrap"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, h *bootstrap.Handlers) {

	// v1 := r.Group("/api/v1")

	// // rutas públicas (sin RequiereAuth)
	// auth := v1.Group("/auth")
	// {
	// 	auth.POST("/login", auth.Handlers.Login)
	// }

	// //rutas privadas
	// users := v1.Group("/users")
	// users.Use(middleware.RequireAuth())
	// {
	// 	users.POST("/", users.)
	// }
}
