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
