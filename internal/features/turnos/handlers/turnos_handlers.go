package handlers

import (
	"strconv"
	"turnos-medicos/internal/features/turnos/dto"
	"turnos-medicos/internal/features/turnos/services"
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

func (h *TurnoHandler) CrearTurno(c *gin.Context) {

	var req dto.CrearTurnoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.BadRequest(c, "datos invalidos")
		return
	}

	turno, err := h.service.CrearTurno(
		c.Request.Context(),
		req,
	)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Created(c, turno, "Turno creado correctamente")
}

func (h *TurnoHandler) ObtenerTurnoPorID(c *gin.Context) {
	idStr := c.Param("id")

	turnoID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.BadRequest(c, pkg.ErrIDInvalido.Error())
		return
	}
	turno, err := h.service.ObtenerTurnoPorID(c.Request.Context(), turnoID)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, turno, "Turno obtenido correctamente")
}

func (h *TurnoHandler) ListarTurnosPorMedico(c *gin.Context) {
	idStr := c.Param("medicoID")

	medicoID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.BadRequest(c, pkg.ErrIDInvalido.Error())
		return
	}

	turnos, err := h.service.ListarTurnosPorMedico(c.Request.Context(), medicoID)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, turnos, "Turnos del medico obtenidos correctamente")
}

func (h *TurnoHandler) ListarTurnosPorPaciente(c *gin.Context) {
	idStr := c.Param("pacienteID")

	pacienteID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.BadRequest(c, pkg.ErrIDInvalido.Error())
		return
	}

	turnos, err := h.service.ListarTurnosPorPaciente(c.Request.Context(), pacienteID)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, turnos, "Turnos del paciente obtenidos correctamente")
}

func (h *TurnoHandler) ListarTurnosDisponibles(c *gin.Context) {
	medicoIDStr := c.Param("medicoID")

	medicoID, err := strconv.ParseInt(medicoIDStr, 10, 64)
	if err != nil {
		pkg.BadRequest(c, pkg.ErrIDInvalido.Error())
		return
	}

	turnos, err := h.service.ListarTurnosDisponibles(c.Request.Context(), medicoID)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, turnos, "Turnos disponibles obtenidos correctamente")

}

func (h *TurnoHandler) ReservarTurno(c *gin.Context) {
	idStr := c.Param("id")

	turnoID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.BadRequest(c, pkg.ErrIDInvalido.Error())
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

func (h *TurnoHandler) LiberarTurno(c *gin.Context) {
	idStr := c.Param("id")

	turnoID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.BadRequest(c, pkg.ErrIDInvalido.Error())
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

func (h *TurnoHandler) MarcarTurnoAtendido(c *gin.Context) {
	idStr := c.Param("id")

	turnoID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.BadRequest(c, pkg.ErrIDInvalido.Error())
		return
	}

	err = h.service.MarcarTurnoAtendido(c.Request.Context(), turnoID)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, nil, "Turno marcado como atendido correctamente")

}

func (h *TurnoHandler) MarcarTurnoNoAsistio(c *gin.Context) {

	idStr := c.Param("id")

	turnoID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.BadRequest(c, pkg.ErrIDInvalido.Error())
		return
	}

	err = h.service.MarcarTurnoNoAsistio(c.Request.Context(), turnoID)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, nil, "Turno marcado como no asistió correctamente")
}
