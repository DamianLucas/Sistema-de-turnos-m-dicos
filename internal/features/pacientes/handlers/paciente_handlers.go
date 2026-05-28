package handlers

import (
	"errors"
	"turnos-medicos/internal/features/pacientes/dto"
	"turnos-medicos/internal/features/pacientes/services"
	"turnos-medicos/internal/middleware"
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

	id, err := pkg.ParseInt64Param(c, "id")
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	// usuario autenticado
	authUserID := middleware.GetUserID(c)
	authRol := middleware.GetUserRol(c)

	paciente, err := h.service.ObtenerPacientePorID(
		c.Request.Context(),
		authUserID,
		authRol,
		id,
	)
	if err != nil {
		pkg.HandleError(c, err)
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
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, paciente, "Paciente obtenido por DNI correctamente")
}

func (h *PacienteHandler) ListarPacientesActivos(c *gin.Context) {

	authUserID := middleware.GetUserID(c)
	authRol := middleware.GetUserRol(c)

	pacientesActivos, err := h.service.ListarPacientesActivos(c.Request.Context(), authUserID, authRol)
	if err != nil {
		pkg.InternalError(c)
		return
	}

	pkg.Success(c, pacientesActivos, "Pacientes activos listados correctamente")
}

func (h *PacienteHandler) DesactivarPaciente(c *gin.Context) {
	pacienteID, err := pkg.ParseInt64Param(c, "id")
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

func (h *PacienteHandler) ActivarPaciente(c *gin.Context) {
	pacienteID, err := pkg.ParseInt64Param(c, "id")
	if err != nil {
		pkg.BadRequest(c, pkg.ErrIDInvalido.Error())
		return
	}

	err = h.service.ActivarPaciente(c.Request.Context(), pacienteID)
	if err != nil {
		if errors.Is(err, pkg.ErrIDInvalido) {
			pkg.BadRequest(c, err.Error())
			return
		}

		if errors.Is(err, pkg.ErrPacienteNoEncontrado) {
			pkg.NotFound(c, err.Error())
			return
		}

		if errors.Is(err, pkg.ErrPacienteYaActivo) {
			pkg.BadRequest(c, err.Error())
			return
		}

		pkg.InternalError(c)
		return
	}

	pkg.Success(c, nil, "Paciente activado correctamente")

}

func (h *PacienteHandler) ActualizarPaciente(c *gin.Context) {
	pacienteID, err := pkg.ParseInt64Param(c, "id")
	if err != nil {
		pkg.BadRequest(c, pkg.ErrIDInvalido.Error())
		return
	}

	var req dto.ActualizarPacienteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.BadRequest(c, err.Error())
		return
	}

	// usuario autenticado
	authUserID := middleware.GetUserID(c)
	authRol := middleware.GetUserRol(c)

	paciente, err := h.service.ActualizarPaciente(c.Request.Context(), authUserID, authRol, pacienteID, req)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, paciente, "Paciente actualizado correctamente")

}

func (h *PacienteHandler) AsignarMedicoTratante(c *gin.Context) {
	pacienteID, err := pkg.ParseInt64Param(c, "id")
	if err != nil {
		pkg.BadRequest(c, pkg.ErrIDInvalido.Error())
		return
	}

	medicoID, err := pkg.ParseInt64Param(c, "medicoID")
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
	pacienteID, err := pkg.ParseInt64Param(c, "id")
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
