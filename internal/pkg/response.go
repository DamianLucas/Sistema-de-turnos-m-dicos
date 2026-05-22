package pkg

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func JSON(c *gin.Context, status int, payload ApiResponse) {
	c.JSON(status, payload)
}

//helpers

func Success(c *gin.Context, data interface{}, message string) {
	JSON(c, http.StatusOK, ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Created(c *gin.Context, data interface{}, message string) {
	JSON(c, http.StatusCreated, ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func BadRequest(c *gin.Context, err string) {
	JSON(c, http.StatusBadRequest, ApiResponse{
		Success: false,
		Error:   err,
	})
}

func NotFound(c *gin.Context, err string) {
	JSON(c, http.StatusNotFound, ApiResponse{
		Success: false,
		Error:   err,
	})
}

func InternalError(c *gin.Context) {
	JSON(c, http.StatusInternalServerError, ApiResponse{
		Success: false,
		Error:   "error interno del servidor",
	})
}

func Unauthorized(c *gin.Context, err string) {
	JSON(c, http.StatusUnauthorized, ApiResponse{
		Success: false,
		Error:   err,
	})
}

func Forbidden(c *gin.Context, err string) {
	JSON(c, http.StatusForbidden, ApiResponse{
		Success: false,
		Error:   err,
	})
}

func HandleError(c *gin.Context, err error) {
	var httpErr HTTPError

	//buscamos si el error implementa el error
	if errors.As(err, &httpErr) {

		JSON(c, httpErr.StatusCode(), ApiResponse{
			Success: false,
			Error:   httpErr.Error(),
		})

		return
	}

	// Error inesperado
	log.Printf(
		"[ERROR] method=%s path=%s error=%v", c.Request.Method, c.Request.URL.Path, err,
	)

	InternalError(c)
}
