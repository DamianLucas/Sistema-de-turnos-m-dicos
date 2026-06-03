package handlers

import (
	"errors"
	"turnos-medicos/internal/features/medicos/dto"
	"turnos-medicos/internal/features/medicos/models"
	"turnos-medicos/internal/features/medicos/services"
	pacienteService "turnos-medicos/internal/features/pacientes/services"
	"turnos-medicos/internal/middleware"

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

// CrearMedico godoc
//
// @Summary Crear médico
// @Description Crea un nuevo médico en el sistema
// @Tags Medicos
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CrearMedicoRequest true "Datos del médico"
// @Success 201 {object} map[string]interface{} "Médico creado correctamente"
// @Failure 400 {object} map[string]interface{} "Datos inválidos, email o matrícula duplicada"
// @Failure 401 {object} map[string]interface{} "No autenticado"
// @Failure 403 {object} map[string]interface{} "Sin permisos"
// @Failure 500 {object} map[string]interface{} "Error interno"
// @Router /medicos [post]
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

// ObtenerMedicoPorID godoc
//
// @Summary Obtener médico por ID
// @Description Obtiene un médico por su ID
// @Tags Medicos
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID del médico"
// @Success 200 {object} models.Medico
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 401 {object} map[string]interface{} "No autenticado"
// @Failure 403 {object} map[string]interface{} "Sin permisos"
// @Failure 404 {object} map[string]interface{} "Médico no encontrado"
// @Failure 500 {object} map[string]interface{} "Error interno"
// @Router /medicos/{id} [get]
func (h *MedicoHandler) ObtenerMedicoPorID(c *gin.Context) {
	id, err := pkg.ParseInt64Param(c, "id")

	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	// usuario autenticado
	authUserID := middleware.GetUserID(c)
	authRol := middleware.GetUserRol(c)

	medico, err := h.service.ObtenerMedicoPorID(c.Request.Context(), authUserID, authRol, id)
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
		pkg.HandleError(c, pkg.ErrMatriculaRequerida)
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

// ListarMedicos godoc
//
// @Summary Listar médicos
// @Description Lista todos los médicos activos o filtra por especialidad
// @Tags Medicos
// @Security BearerAuth
// @Produce json
// @Param especialidad query string false "Filtrar por especialidad"
// @Success 200 {array} models.Medico
// @Failure 401 {object} map[string]interface{} "No autenticado"
// @Failure 403 {object} map[string]interface{} "Sin permisos"
// @Failure 500 {object} map[string]interface{} "Error interno"
// @Router /medicos [get]
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
	medicoID, err := pkg.ParseInt64Param(c, "id")
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	var req dto.ActualizarMedicoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.HandleError(c, err)
		return
	}

	authUserID := middleware.GetUserID(c)
	authRol := middleware.GetUserRol(c)

	medico, err := h.service.ActualizarMedico(c.Request.Context(), authUserID, authRol, medicoID, req)
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
	medicoID, err := pkg.ParseInt64Param(c, "id")
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
	medicoID, err := pkg.ParseInt64Param(c, "id")
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
	medicoID, err := pkg.ParseInt64Param(c, "id")
	if err != nil {
		pkg.BadRequest(c, pkg.ErrIDInvalido.Error())
		return
	}

	//usuario autenticado
	authUserID := middleware.GetUserID(c)
	authRol := middleware.GetUserRol(c)

	// Llamamos al servicio de pacientes porque lo que queremos obtener son pacientes
	pacientes, err := h.pacienteService.ListarPacientesPorMedico(c.Request.Context(), authUserID, authRol, medicoID)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, pacientes, "Pacientes listados correctamente")
}
