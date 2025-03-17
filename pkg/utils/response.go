package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response estructura estándar para respuestas JSON
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse envía una respuesta exitosa
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// ErrorResponse envía una respuesta de error
func ErrorResponse(c *gin.Context, statusCode int, errorMsg string) {
	c.JSON(statusCode, Response{
		Status: "error",
		Error:  errorMsg,
	})
}

// ValidationErrorResponse envía respuesta de errores de validación
func ValidationErrorResponse(c *gin.Context, errorMsg string) {
	ErrorResponse(c, http.StatusBadRequest, errorMsg)
}

// NotFoundResponse envía respuesta para recursos no encontrados
func NotFoundResponse(c *gin.Context, resourceName string) {
	ErrorResponse(c, http.StatusNotFound, resourceName+" no encontrado")
}

// InternalErrorResponse envía respuesta para errores internos
func InternalErrorResponse(c *gin.Context) {
	ErrorResponse(c, http.StatusInternalServerError, "Error interno del servidor")
}
