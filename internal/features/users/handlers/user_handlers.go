package handlers

import (
	"errors"
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
			pkg.BadRequest(c, pkg.ErrEmailDuplicado.Error())
			return
		}
		pkg.InternalError(c)
		return
	}
	pkg.Created(c, user, "usuario creado correctamente")

}

func (h *UserHandler) ObtenerUsuarioPorID(c *gin.Context) {
	id, err := pkg.ParseInt64Param(c, "id")
	if err != nil {
		pkg.BadRequest(c, "id invalido")
		return
	}

	userID, err := h.service.ObtenerUsuarioPorID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pkg.ErrUsuarioNoEncontrado) {
			pkg.BadRequest(c, pkg.ErrUsuarioNoEncontrado.Error())
			return
		}
		pkg.InternalError(c)
		return
	}
	pkg.Success(c, userID, "usuario obtenido correctamente")

}

func (h *UserHandler) ListarUsuariosActivos(c *gin.Context) {
	ctx := c.Request.Context()

	usersActive, err := h.service.ListarUsuariosActivos(ctx)
	if err != nil {
		pkg.InternalError(c)
		return
	}

	pkg.Success(c, usersActive, "usuarios activos listados correctamente")
}

func (h *UserHandler) ActualizarUsuario(c *gin.Context) {
	userID, err := pkg.ParseInt64Param(c, "id")
	if err != nil {
		pkg.BadRequest(c, "id invalido")
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
			pkg.BadRequest(c, pkg.ErrUsuarioNoEncontrado.Error())
			return
		}
		pkg.InternalError(c)
		return
	}

	pkg.Success(c, user, "usuario actualizado correctamente")
}

func (h *UserHandler) DesactivarUsuario(c *gin.Context) {
	userID, err := pkg.ParseInt64Param(c, "id")
	if err != nil {
		pkg.BadRequest(c, "id invalido")
		return
	}

	err = h.service.DesactivarUsuario(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, pkg.ErrUsuarioNoEncontrado) {
			pkg.BadRequest(c, pkg.ErrUsuarioNoEncontrado.Error())
			return
		}

		pkg.InternalError(c)
		return
	}

	pkg.Success(c, nil, "Usuario desactivado")
}
