package handlers

import (
	"errors"
	"strconv"
	"turnos-medicos/internal/features/medicos/dto"
	"turnos-medicos/internal/features/medicos/models"
	"turnos-medicos/internal/features/medicos/services"
	"turnos-medicos/internal/pkg"

	"github.com/gin-gonic/gin"
)

// IMPLEMENTAR
type MedicoHandler struct {
	service services.MedicoService
}

func NewMedicoHandler(s services.MedicoService) *MedicoHandler {
	return &MedicoHandler{service: s}
}

func (h *MedicoHandler) CrearMedico(c *gin.Context) {

	var req dto.CrearMedicoRequest

	// Recibir y validar el body
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.BadRequest(c, "datos invalidos")
		return
	}

	medico, err := h.service.CrearMedico(c.Request.Context(), req)
	if err != nil {

		if errors.Is(err, pkg.ErrEmailDuplicado) {
			pkg.BadRequest(c, pkg.ErrEmailDuplicado.Error())
			return
		}

		if errors.Is(err, pkg.ErrMatriculaDuplicada) {
			pkg.BadRequest(c, pkg.ErrMatriculaDuplicada.Error())
			return
		}

		pkg.InternalError(c)
		// c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	pkg.Created(c, medico)
}

func (h *MedicoHandler) ObtenerMedicoPorID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		pkg.BadRequest(c, "ID invalido")
		return
	}

	medicoID, err := h.service.ObtenerMedicoPorID(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, pkg.ErrMedicoNoEncontrado) {
			pkg.NotFound(c, err.Error())
			return
		}
		if errors.Is(err, pkg.ErrMedicoInactivo) {
			pkg.NotFound(c, err.Error())
			return
		}
		pkg.InternalError(c)
		return
	}
	pkg.Success(c, medicoID)
}

func (h *MedicoHandler) ObtenerMedicoPorMatricula(c *gin.Context) {
	matricula := c.Param("matricula")
	if matricula == "" {
		pkg.BadRequest(c, "matricula es requerida")
		return
	}

	medico, err := h.service.ObtenerMedicoPorMatricula(c.Request.Context(), matricula)
	if err != nil {
		if errors.Is(err, pkg.ErrMedicoNoEncontrado) {
			pkg.NotFound(c, pkg.ErrMedicoNoEncontrado.Error())
			return
		}

		pkg.InternalError(c)
		return
	}
	pkg.Success(c, medico)
}

func (h *MedicoHandler) ListarMedicos(c *gin.Context) {

	especialidad := c.Query("especialidad")

	var (
		medicos []*models.Medico
		err     error
	)

	if especialidad != "" {
		medicos, err = h.service.ListarMedicosPorEspecialidad(
			c.Request.Context(),
			especialidad,
		)
	} else {
		medicos, err = h.service.ListarMedicosActivos(
			c.Request.Context(),
		)
	}

	if err != nil {
		pkg.InternalError(c)
		return
	}

	if medicos == nil {
		medicos = []*models.Medico{}
	}

	pkg.Success(c, medicos)
}

func (h *MedicoHandler) ActualizarMedico(c *gin.Context) {
	idStr := c.Param("id")
	medicoID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.BadRequest(c, "ID invalido")
		return
	}

	var req dto.ActualizarMedicoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.BadRequest(c, "Datos invalidos")
		return
	}

	medico, err := h.service.ActualizarMedico(c.Request.Context(), medicoID, req)
	if err != nil {
		if errors.Is(err, pkg.ErrMedicoNoEncontrado) {
			pkg.NotFound(c, pkg.ErrMedicoNoEncontrado.Error())
			return
		}
		if errors.Is(err, pkg.ErrMedicoInactivo) {
			pkg.NotFound(c, pkg.ErrMedicoInactivo.Error())
			return
		}
		pkg.InternalError(c)
		return
	}

	pkg.Success(c, medico)
}

func (h *MedicoHandler) DesactivarMedico(c *gin.Context) {
	idStr := c.Param("id")
	medicoID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.BadRequest(c, "ID invalido")
		return
	}

	err = h.service.DesactivarMedico(c.Request.Context(), medicoID)
	if err != nil {
		if errors.Is(err, pkg.ErrMedicoNoEncontrado) {
			pkg.NotFound(c, pkg.ErrMedicoNoEncontrado.Error())
			return
		}

		pkg.InternalError(c)
		return
	}

	pkg.Success(c, "Medico desactivado correctamente")
}

func (h *MedicoHandler) ActivarMedico(c *gin.Context) {
	idStr := c.Param("id")
	medicoID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.BadRequest(c, "ID invalido")
		return
	}

	err = h.service.ActivarMedico(c.Request.Context(), medicoID)
	if err != nil {

		if errors.Is(err, pkg.ErrMedicoNoEncontrado) {
			pkg.NotFound(c, pkg.ErrMedicoNoEncontrado.Error())
			return
		}

		pkg.InternalError(c)
		return
	}

	pkg.Success(c, "medico activado correctamente")
}
