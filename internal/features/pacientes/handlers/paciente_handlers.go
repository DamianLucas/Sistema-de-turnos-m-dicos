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

func (h *PacienteHandler) ObtenerPacientePorDNI(c *gin.Context) {
	dni := c.Param("dni")

	if dni == "" {
		pkg.BadRequest(c, pkg.ErrDNIInvalido.Error())
		return
	}

	paciente, err := h.service.ObtenerPacientePorDNI(c.Request.Context(), dni)
	if err != nil {
		if errors.Is(err, pkg.ErrPacienteNoEncontrado) {
			pkg.NotFound(c, err.Error())
			return
		}

		pkg.InternalError(c)
		return
	}

	pkg.Success(c, paciente, "Paciente obtenido por DNI correctamente")
}

func (h *PacienteHandler) ListarPacientesActivos(c *gin.Context) {
	pacientesActivos, err := h.service.ListarPacientesActivos(c.Request.Context())
	if err != nil {
		pkg.InternalError(c)
		return
	}

	pkg.Success(c, pacientesActivos, "Pacientes activos listados correctamente")
}

func (h *PacienteHandler) DesactivarPaciente(c *gin.Context) {
	idStr := c.Param("id")

	pacienteID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.BadRequest(c, pkg.ErrIDInvalido.Error())
		return
	}

	err = h.service.DesactivarPaciente(c.Request.Context(), pacienteID)
	if err != nil {
		if errors.Is(err, pkg.ErrIDInvalido) {
			pkg.BadRequest(c, err.Error())
			return
		}

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

	pkg.Success(c, nil, "Paciente desactivado correctamente")

}

func (h *PacienteHandler) ActualizarPaciente(c *gin.Context) {
	idStr := c.Param("id")

	pacienteID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.BadRequest(c, pkg.ErrIDInvalido.Error())
		return
	}

	var req dto.ActualizarPacienteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.BadRequest(c, err.Error())
		return
	}

	paciente, err := h.service.ActualizarPaciente(c.Request.Context(), pacienteID, req)
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

	pkg.Success(c, paciente, "Paciente actualizado correctamente")

}

func (h *PacienteHandler) AsignarMedicoTratante(c *gin.Context) {
	idStr := c.Param("id")
	pacienteID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.BadRequest(c, pkg.ErrIDInvalido.Error())
		return
	}

	medicoIDStr := c.Param("medicoID")
	medicoID, err := strconv.ParseInt(medicoIDStr, 10, 64)
	if err != nil {
		pkg.BadRequest(c, pkg.ErrIDInvalido.Error())
		return
	}

	err = h.service.AsignarMedicoTratante(c.Request.Context(), pacienteID, medicoID)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, nil, "Médico tratante asignado correctamente")
}

func (h *PacienteHandler) QuitarMedicoTratante(c *gin.Context) {
	idStr := c.Param("id")

	pacienteID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.BadRequest(c, pkg.ErrIDInvalido.Error())
		return
	}

	err = h.service.QuitarMedicoTratante(c.Request.Context(), pacienteID)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, nil, "Médico tratante removido correctamente")
}
