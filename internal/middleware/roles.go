package middleware

import (
	"turnos-medicos/internal/features/users/models"
	"turnos-medicos/internal/utils"

	"github.com/gin-gonic/gin"
)

func RequireRol(rolesPermitidos ...models.Rol) gin.HandlerFunc {
	return func(c *gin.Context) {

		rolCtx, exist := c.Get("rol")
		if !exist {
			utils.Unauthorized(c, "No autorizado")
			c.Abort()
			return
		}

		rolUsuario, ok := rolCtx.(models.Rol)
		if !ok {
			utils.Unauthorized(c, "Rol invalido")
			c.Abort()
			return
		}

		// verificar si el rol del usuario está permitido
		for _, rol := range rolesPermitidos {
			if rolUsuario == rol {
				c.Next()
				return
			}
		}

		// si no coincide ningún rol
		utils.Forbidden(c, "No tiene permisos para acceder a este recurso")
		c.Abort()
	}
}
