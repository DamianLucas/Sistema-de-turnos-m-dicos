package handlers

import (
	"errors"
	"strconv"
	"turnos-medicos/internal/features/medicos/dto"
	"turnos-medicos/internal/features/medicos/models"
	"turnos-medicos/internal/features/medicos/services"
	pacienteService "turnos-medicos/internal/features/pacientes/services"

	"turnos-medicos/internal/pkg"

	"github.com/gin-gonic/gin"
)

// IMPLEMENTAR
type MedicoHandler struct {
	service         services.MedicoService
	pacienteService pacienteService.PacienteService
}

func NewMedicoHandler(s services.MedicoService, ps pacienteService.PacienteService) *MedicoHandler {
	return &MedicoHandler{
		service:         s,
		pacienteService: ps,
	}
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

	pkg.Created(c, medico, "Medico creado correctamente")
}

func (h *MedicoHandler) ObtenerMedicoPorID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		pkg.BadRequest(c, "ID invalido")
		return
	}

	medico, err := h.service.ObtenerMedicoPorID(c.Request.Context(), id)

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
	pkg.Success(c, medico, "Medico obtenido correctamente")
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
	pkg.Success(c, medico, "Medico obtenido por matricula correctamente")
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

	pkg.Success(c, medicos, "Medicos listados correctamente")
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

	pkg.Success(c, medico, "Medico actualizado correctamente")
}

func (h *MedicoHandler) DesactivarMedico(c *gin.Context) {
	idStr := c.Param("id")

	medicoID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.BadRequest(c, "ID inválido")
		return
	}

	err = h.service.DesactivarMedico(c.Request.Context(), medicoID)
	if err != nil {
		if errors.Is(err, pkg.ErrIDInvalido) {
			pkg.BadRequest(c, err.Error())
			return
		}

		if errors.Is(err, pkg.ErrMedicoNoEncontrado) {
			pkg.NotFound(c, err.Error())
			return
		}

		if errors.Is(err, pkg.ErrMedicoInactivo) {
			pkg.BadRequest(c, err.Error())
			return
		}

		pkg.InternalError(c)
		return
	}

	pkg.Success(c, nil, "Médico desactivado correctamente")
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

	pkg.Success(c, nil, "medico activado correctamente")
}

func (h *MedicoHandler) ListarPacientesPorMedico(c *gin.Context) {
	idStr := c.Param("id")
	medicoID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.BadRequest(c, pkg.ErrIDInvalido.Error())
		return
	}

	// Llamamos al servicio de pacientes porque lo que queremos obtener son pacientes
	pacientes, err := h.pacienteService.ListarPacientesPorMedico(c.Request.Context(), medicoID)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, pacientes, "Pacientes listados correctamente")
}
