package handler

import (
	"errors"
	"turnos-medicos/internal/features/auth/dto"
	"turnos-medicos/internal/features/auth/service"
	"turnos-medicos/internal/pkg"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// Login godoc
//
// @Summary Iniciar sesión
// @Description Autentica un usuario y devuelve un JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Credenciales de acceso"
// @Success 200 {object} map[string]interface{} "Login correcto"
// @Failure 400 {object} map[string]interface{} "Request inválido"
// @Failure 401 {object} map[string]interface{} "Credenciales inválidas"
// @Failure 500 {object} map[string]interface{} "Error interno"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.BadRequest(c, err.Error())
		return
	}

	resp, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, pkg.ErrCredencialesInvalidas) {
			pkg.Unauthorized(c, "Email o Password invalido")
			return
		}
		pkg.InternalError(c)
		return
	}

	pkg.Success(c, resp, "Login correcto")
}
