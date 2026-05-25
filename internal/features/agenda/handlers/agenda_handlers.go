package handlers

import (
	"turnos-medicos/internal/features/agenda/dto"
	"turnos-medicos/internal/features/agenda/services"
	"turnos-medicos/internal/middleware"
	"turnos-medicos/internal/pkg"

	"github.com/gin-gonic/gin"
)

type AgendaHandler struct {
	service services.AgendaService
}

func NewAgendaHandler(s services.AgendaService) *AgendaHandler {
	return &AgendaHandler{service: s}
}

func (h *AgendaHandler) CrearAgenda(c *gin.Context) {
	medicoID, err := pkg.ParseInt64Param(c, "medicoID")
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	var req dto.CrearAgendaRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.BadRequest(c, "Datos inválidos")
		return
	}

	agenda, err := h.service.CrearAgenda(c.Request.Context(), medicoID, req)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Created(c, agenda, "Agenda creada correctamente")
}

func (h *AgendaHandler) ObtenerAgendaPorID(c *gin.Context) {
	agendaID, err := pkg.ParseInt64Param(c, "id")
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	//usuario autenticado
	authUserID := middleware.GetUserID(c)
	authRol := middleware.GetUserRol(c)

	agenda, err := h.service.ObtenerAgendaPorID(c.Request.Context(), authUserID, authRol, agendaID)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, agenda, "Agenda obtenida correctamente")
}

func (h *AgendaHandler) ListarAgendasPorMedico(c *gin.Context) {
	medicoID, err := pkg.ParseInt64Param(c, "medicoID")
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	//usuario autenticado
	authUserID := middleware.GetUserID(c)
	authRol := middleware.GetUserRol(c)

	listarAgenda, err := h.service.ListarAgendasPorMedico(c.Request.Context(), authUserID, authRol, medicoID)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, listarAgenda, "Agendas listadas correctamente")

}

func (h *AgendaHandler) ActualizarAgenda(c *gin.Context) {
	agendaID, err := pkg.ParseInt64Param(c, "id")
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	var req dto.ActualizarAgendaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.BadRequest(c, err.Error())
		return
	}

	//usuario autenticado
	authUserID := middleware.GetUserID(c)
	authRol := middleware.GetUserRol(c)

	agenda, err := h.service.ActualizarAgenda(c.Request.Context(), authUserID, authRol, agendaID, req)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, agenda, "Agenda actualizada correctamente")
}

func (h *AgendaHandler) DesactivarAgenda(c *gin.Context) {
	agendaID, err := pkg.ParseInt64Param(c, "id")
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	authUserID := middleware.GetUserID(c)
	authRol := middleware.GetUserRol(c)

	err = h.service.DesactivarAgenda(c.Request.Context(), authUserID, authRol, agendaID)

	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, nil, "Agenda desactivada correctamente")

}

func (h *AgendaHandler) ActivarAgenda(c *gin.Context) {
	agendaID, err := pkg.ParseInt64Param(c, "id")
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	authUserID := middleware.GetUserID(c)
	authRol := middleware.GetUserRol(c)

	err = h.service.ActivarAgenda(c.Request.Context(), authUserID, authRol, agendaID)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	pkg.Success(c, nil, "Agenda activada correctamente")
}
