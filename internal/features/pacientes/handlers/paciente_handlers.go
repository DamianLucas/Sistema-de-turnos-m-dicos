package handlers

import (
	"errors"
	"strconv"
	"turnos-medicos/internal/features/pacientes/dto"
	"turnos-medicos/internal/features/pacientes/services"
	"turnos-medicos/internal/pkg"

	"github.com/gin-gonic/gin"
)

// IMPLEMENTAR
type PacienteHandler struct {
	service services.PacienteService
}

func NewPacienteHandler(s services.PacienteService) *PacienteHandler {
	return &PacienteHandler{service: s}
}

func (h *PacienteHandler) CrearPaciente(c *gin.Context) {

	var req dto.CrearPacienteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.BadRequest(c, "datos invalidos")
		return
	}

	paciente, err := h.service.CrearPaciente(c.Request.Context(), req)
	if err != nil {

		if errors.Is(err, pkg.ErrDNIDuplicado) {
			pkg.BadRequest(c, err.Error())
			return
		}
		if errors.Is(err, pkg.ErrEmailDuplicado) {
			pkg.BadRequest(c, pkg.ErrEmailDuplicado.Error())
			return
		}

		pkg.InternalError(c)
		return
	}

	pkg.Created(c, paciente, "Paciente creado correctamente")
}

func (h *PacienteHandler) ObtenerPacientePorID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.BadRequest(c, "ID invalido")
		return
	}

	paciente, err := h.service.ObtenerPacientePorID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pkg.ErrPacienteNoEncontrado) {
			pkg.NotFound(c, err.Error())
			return
		}
		if errors.Is(err, pkg.ErrPacienteInactivo) {
			pkg.BadRequest(c, err.Error())
			return
		}
		pkg.InternalError(c)
		return
	}

	pkg.Success(c, paciente, "Paciente obtenido correctamente")
}
