package middleware

import (
	"turnos-medicos/internal/features/users/models"

	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) int64 {
	return c.MustGet("userID").(int64)
}

func GetUserRol(c *gin.Context) models.Rol {
	return c.MustGet("rol").(models.Rol)
}
