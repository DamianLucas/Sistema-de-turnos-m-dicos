package middleware

import (
	"strings"
	"turnos-medicos/internal/utils"

	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Obtener Autorizacion
		authHeader := c.GetHeader("Authorization")

		//verificar que exista
		if authHeader == "" {
			utils.Unauthorized(c, "No autorizado")
			c.Abort()
			return
		}

		//verificar formato: "Bearer token"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(c, "No autorizado")
			c.Abort()
			return
		}

		tokenString := parts[1]
		//validar token
		claims, err := utils.ValidarToken(tokenString)
		if err != nil {
			utils.Unauthorized(c, "No autorizado")
			c.Abort()
			return
		}

		//guardar claims en el contexto
		c.Set("userID", claims.UserID)
		c.Set("rol", claims.Rol)

		c.Next()
	}
}
