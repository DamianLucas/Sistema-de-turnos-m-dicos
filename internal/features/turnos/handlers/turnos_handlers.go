package handlers

import (
	"turnos-medicos/internal/features/turnos/dto"
	"turnos-medicos/internal/features/turnos/services"
	"turnos-medicos/internal/middleware"
	"turnos-medicos/internal/pkg"

	"github.com/gin-gonic/gin"
)

type TurnoHandler struct {
	service services.TurnoService
}

func NewTurnoHandler(service services.TurnoService) *TurnoHandler {
	return &TurnoHandler{
		service: service,
	}
}

// CrearTurno godoc
//
// @Summary Crear turno
// @Description Crea un turno disponible dentro de una agenda
// @Tags Turnos
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID de la agenda"
// @Param request body dto.CrearTurnoRequest true "Datos del turno"
// @Success 201 {object} models.Turno
// @Failure 400 {object} map[string]interface{} "Datos inválidos"
// @Failure 401 {object} map[string]interface{} "No autenticado"
// @Failure 403 {object} map[string]interface{} "Sin permisos"
// @Failure 404 {object} map[string]interface{} "Agenda no encontrada"
// @Failure 500 {object} map[string]interface{} "Error interno"
// @Router /turnos/agenda/{id}/crear [post]
func (h *TurnoHandler) CrearTurno(c *gin.Context) {

	agendaID, err := pkg.ParseInt64Param(c, "id")
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	var req dto.CrearTurnoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.BadRequest(c, "datos invalidos")
		return
	}

	turno, err := h.service.CrearTurno(c.Request.Context(), agendaID, req)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Created(c, turno, "Turno creado correctamente")
}

// ObtenerTurnoPorID godoc
// @Summary Obtener turno por ID
// @Description Obtiene un turno específico por su ID
// @Tags Turnos
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID del turno"
// @Success 200 {object} pkg.APIResponse{data=models.Turno}
// @Failure 400 {object} pkg.APIResponse
// @Failure 404 {object} pkg.APIResponse
// @Failure 401 {object} pkg.APIResponse
// @Failure 500 {object} pkg.APIResponse
// @Router /turnos/{id} [get]
func (h *TurnoHandler) ObtenerTurnoPorID(c *gin.Context) {
	turnoID, err := pkg.ParseInt64Param(c, "id")
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	//usuario autenticado
	authUserID := middleware.GetUserID(c)
	authRol := middleware.GetUserRol(c)

	turno, err := h.service.ObtenerTurnoPorID(c.Request.Context(), authUserID, authRol, turnoID)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, turno, "Turno obtenido correctamente")
}

// ListarTurnosPorMedico godoc
// @Summary Listar turnos por médico
// @Description listar turnos del médico
// @Tags Turnos
// @Security BearerAuth
// @Produce json
// @Param medicoID path int true "ID del médico"
// @Success 200 {object} pkg.APIResponse{data=[]models.Turno}
// @Failure 400 {object} pkg.APIResponse
// @Failure 404 {object} pkg.APIResponse
// @Failure 401 {object} pkg.APIResponse
// @Failure 500 {object} pkg.APIResponse
// @Router /turnos/medico/{medicoID} [get]
func (h *TurnoHandler) ListarTurnosPorMedico(c *gin.Context) {
	medicoID, err := pkg.ParseInt64Param(c, "medicoID")
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	//usuario autenticado
	authUserID := middleware.GetUserID(c)
	authRol := middleware.GetUserRol(c)

	turnos, err := h.service.ListarTurnosPorMedico(
		c.Request.Context(),
		authUserID,
		authRol,
		medicoID,
	)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}
	pkg.Success(c, turnos, "Turnos del medico obtenidos correctamente")
}

// ListarTurnosPorPaciente godoc
// @Summary Listar turnos por paciente
// @Description Obtiene todos los turnos asociados a un paciente
// @Tags Turnos
// @Security BearerAuth
// @Produce json
// @Param pacienteID path int true "ID del paciente"
// @Success 200 {object} pkg.APIResponse{data=[]models.Turno}
// @Failure 400 {object} pkg.APIResponse
// @Failure 404 {object} pkg.APIResponse
// @Failure 401 {object} pkg.APIResponse
// @Failure 500 {object} pkg.APIResponse
// @Router /turnos/paciente/{pacienteID} [get]
func (h *TurnoHandler) ListarTurnosPorPaciente(c *gin.Context) {
	pacienteID, err := pkg.ParseInt64Param(c, "pacienteID")
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	turnos, err := h.service.ListarTurnosPorPaciente(c.Request.Context(), pacienteID)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, turnos, "Turnos del paciente obtenidos correctamente")
}

// ListarTurnosDisponibles godoc
// @Summary Listar turnos disponibles
// @Description Obtiene todos los turnos disponibles de un médico
// @Tags Turnos
// @Security BearerAuth
// @Produce json
// @Param medicoID path int true "ID del médico"
// @Success 200 {object} pkg.APIResponse{data=[]dto.TurnoDisponibleResponse}
// @Failure 400 {object} pkg.APIResponse
// @Failure 404 {object} pkg.APIResponse
// @Failure 401 {object} pkg.APIResponse
// @Failure 500 {object} pkg.APIResponse
// @Router /turnos/disponibles/{medicoID} [get]
func (h *TurnoHandler) ListarTurnosDisponibles(c *gin.Context) {
	medicoID, err := pkg.ParseInt64Param(c, "medicoID")
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	turnos, err := h.service.ListarTurnosDisponibles(c.Request.Context(), medicoID)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, turnos, "Turnos disponibles obtenidos correctamente")

}

// ReservarTurno godoc
// @Summary Reservar turno
// @Description Reserva un turno disponible para un paciente
// @Tags Turnos
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID del turno"
// @Param request body dto.ReservarTurnoRequest true "Datos de la reserva"
// @Success 200 {object} pkg.APIResponse
// @Failure 400 {object} pkg.APIResponse
// @Failure 404 {object} pkg.APIResponse
// @Failure 401 {object} pkg.APIResponse
// @Failure 500 {object} pkg.APIResponse
// @Router /turnos/{id}/reservar [put]
func (h *TurnoHandler) ReservarTurno(c *gin.Context) {
	turnoID, err := pkg.ParseInt64Param(c, "id")
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	var req dto.ReservarTurnoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.BadRequest(c, "datos invalidos")
		return
	}

	err = h.service.ReservarTurno(c.Request.Context(), turnoID, req.PacienteID)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, nil, "Turno reservado correctamente")
}

// LiberarTurno godoc
// @Summary Liberar turno
// @Description Libera un turno reservado y lo vuelve a dejar disponible
// @Tags Turnos
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID del turno"
// @Success 200 {object} pkg.APIResponse
// @Failure 400 {object} pkg.APIResponse
// @Failure 404 {object} pkg.APIResponse
// @Failure 401 {object} pkg.APIResponse
// @Failure 500 {object} pkg.APIResponse
// @Router /turnos/{id}/liberar [put]
func (h *TurnoHandler) LiberarTurno(c *gin.Context) {
	turnoID, err := pkg.ParseInt64Param(c, "id")
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	err = h.service.LiberarTurno(c.Request.Context(), turnoID)
	if err != nil {
		pkg.HandleError(c, err)
		// c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	pkg.Success(c, nil, "Turno liberado correctamente")
}

// @Summary Marcar turno como atendido
// @Description Marca un turno reservado como atendido
// @Tags Turnos
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID del turno"
// @Success 200 {object} pkg.APIResponse
// @Failure 400 {object} pkg.APIResponse
// @Failure 404 {object} pkg.APIResponse
// @Failure 401 {object} pkg.APIResponse
// @Failure 500 {object} pkg.APIResponse
// @Router /turnos/{id}/atendido [put]
func (h *TurnoHandler) MarcarTurnoAtendido(c *gin.Context) {
	turnoID, err := pkg.ParseInt64Param(c, "id")
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	//usuario autenticado
	authUserID := middleware.GetUserID(c)
	authRol := middleware.GetUserRol(c)

	err = h.service.MarcarTurnoAtendido(c.Request.Context(), authUserID, authRol, turnoID)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, nil, "Turno marcado como atendido correctamente")

}

// MarcarTurnoNoAsistio godoc
// @Summary Marcar turno como no asistió
// @Description Marca un turno reservado como no asistido por el paciente
// @Tags Turnos
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID del turno"
// @Success 200 {object} pkg.APIResponse
// @Failure 400 {object} pkg.APIResponse
// @Failure 404 {object} pkg.APIResponse
// @Failure 401 {object} pkg.APIResponse
// @Failure 500 {object} pkg.APIResponse
// @Router /turnos/{id}/no-asistio [put]
func (h *TurnoHandler) MarcarTurnoNoAsistio(c *gin.Context) {
	turnoID, err := pkg.ParseInt64Param(c, "id")
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	//usuario autenticado
	authUserID := middleware.GetUserID(c)
	authRol := middleware.GetUserRol(c)

	err = h.service.MarcarTurnoNoAsistio(c.Request.Context(), authUserID, authRol, turnoID)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, nil, "Turno marcado como no asistió correctamente")
}
