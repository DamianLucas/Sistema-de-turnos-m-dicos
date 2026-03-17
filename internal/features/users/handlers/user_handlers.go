package handlers

import (
	"errors"
	"strconv"
	"turnos-medicos/internal/features/users/dto"
	"turnos-medicos/internal/features/users/services"
	"turnos-medicos/internal/pkg"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(s services.UserService) *UserHandler {
	return &UserHandler{service: s}
}

//handlers del service

// CrearUsuario
func (h *UserHandler) CrearUsuario(c *gin.Context) {

	var req dto.CrearUsuarioRequest

	//Recibir y validar el body con el dto
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.BadRequest(c, "Datos invalidos")
		return
	}

	//llamar al service y pasarle el dto
	user, err := h.service.CrearUsuario(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, pkg.ErrEmailDuplicado) {
			pkg.BadRequest(c, err.Error())
			return
		}
		pkg.InternalError(c)
		return
	}
	pkg.Created(c, user)

}

// ObtenerUsuarioPorID
func (h *UserHandler) ObtenerUsuarioPorID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		pkg.BadRequest(c, "ID invalido")
		return
	}

	userID, err := h.service.ObtenerUsuarioPorID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pkg.ErrUsuarioNoEncontrado) {
			pkg.NotFound(c, err.Error())
			return
		}
		pkg.InternalError(c)
		return
	}
	pkg.Success(c, userID)

}

// ListarUsuariosActivos
func (h *UserHandler) ListarUsuariosActivos(c *gin.Context) {
	ctx := c.Request.Context()

	usersActive, err := h.service.ListarUsuariosActivos(ctx)
	if err != nil {
		pkg.InternalError(c)
		return
	}

	pkg.Success(c, usersActive)
}

// ActualizarUsuarios
func (h *UserHandler) ActualizarUsuario(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.BadRequest(c, "ID invalido")
		return
	}

	var req dto.ActualizarUsuarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.BadRequest(c, "Datos invalidos")
		return
	}

	user, err := h.service.ActualizarUsuario(c.Request.Context(), userID, req)
	if err != nil {
		if errors.Is(err, pkg.ErrUsuarioNoEncontrado) {
			pkg.NotFound(c, err.Error())
			return
		}
		pkg.InternalError(c)
		return
	}

	pkg.Success(c, user)
}

// DesactivarUsuario
func (h *UserHandler) DesactivarUsuario(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.BadRequest(c, "ID invalido")
		return
	}

	err = h.service.DesactivarUsuario(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, pkg.ErrUsuarioNoEncontrado) {
			pkg.NotFound(c, err.Error())
			return
		}

		pkg.InternalError(c)
		return
	}

	pkg.Success(c, "Usuario desactivado")
}
