package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/black4ninja/mi-proyecto/internal/user/domain"
	"github.com/black4ninja/mi-proyecto/pkg/utils"
)

// UserHandler maneja las peticiones HTTP para usuarios
type UserHandler struct {
	userUseCase domain.UserUseCase
}

// NewUserHandler crea un nuevo manejador de usuarios
func NewUserHandler(router *gin.RouterGroup, useCase domain.UserUseCase) {
	handler := &UserHandler{
		userUseCase: useCase,
	}

	// Rutas públicas
	// Ninguna en este caso

	// Rutas protegidas
	router.GET("/", handler.GetAllUsers)
	router.GET("/:id", handler.GetUser)
	router.POST("/", handler.CreateUser)
	router.PUT("/:id", handler.UpdateUser)
	router.DELETE("/:id", handler.DeleteUser)
	router.PUT("/:id/archive", handler.ArchiveUser)
	router.POST("/change-password", handler.ChangePassword)
	router.GET("/me", handler.GetProfile)
}

// GetAllUsers manejador para obtener todos los usuarios
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	// Extraer todos los parámetros de consulta
	queryParams := make(map[string]string)

	// Parámetros básicos
	if status := c.Query("status"); status != "" {
		queryParams["status"] = status
	}
	if role := c.Query("role"); role != "" {
		queryParams["role"] = role
	}
	if name := c.Query("name"); name != "" {
		queryParams["name"] = name
	}
	if email := c.Query("email"); email != "" {
		queryParams["email"] = email
	}

	// Parámetros de fecha
	if createdFrom := c.Query("created_from"); createdFrom != "" {
		queryParams["created_from"] = createdFrom
	}
	if createdTo := c.Query("created_to"); createdTo != "" {
		queryParams["created_to"] = createdTo
	}
	if archivedFrom := c.Query("archived_from"); archivedFrom != "" {
		queryParams["archived_from"] = archivedFrom
	}
	if archivedTo := c.Query("archived_to"); archivedTo != "" {
		queryParams["archived_to"] = archivedTo
	}

	// Construir filtro seguro para MongoDB
	filter := utils.BuildMongoFilter(queryParams, utils.CommonUserFilterConfig)

	// Filtros de fechas como rangos
	if createdFrom := c.Query("created_from"); createdFrom != "" || c.Query("created_to") != "" {
		dateRange := utils.DateRangeFilter(c.Query("created_from"), c.Query("created_to"))
		if dateRange != nil {
			filter["created_at"] = dateRange
		}
	}

	if archivedFrom := c.Query("archived_from"); archivedFrom != "" || c.Query("archived_to") != "" {
		dateRange := utils.DateRangeFilter(c.Query("archived_from"), c.Query("archived_to"))
		if dateRange != nil {
			filter["archived_at"] = dateRange
		}
	}

	// Añadir filtros adicionales de acuerdo a la lógica de negocio
	// Por ejemplo, para usuarios no archivados cuando no se especifica estatus
	if _, hasStatus := filter["status"]; !hasStatus {
		// Si no se especificó un estatus, mostrar solo usuarios activos por defecto
		filter["status"] = utils.StatusActive
	}

	// Parámetro especial para incluir/excluir archivados
	if includeArchived := c.Query("include_archived"); includeArchived == "true" {
		// Si se solicita explícitamente incluir archivados, eliminar el filtro de status
		delete(filter, "status")
	}

	// Obtener todos los usuarios con los filtros aplicados
	users, err := h.userUseCase.GetAllUsers(filter)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuarios obtenidos con éxito", users)
}

// GetUser manejador para obtener un usuario
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userUseCase.GetUser(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuario obtenido con éxito", user)
}

// CreateUser manejador para crear un usuario
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	user, err := h.userUseCase.CreateUser(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Usuario creado con éxito", user)
}

// UpdateUser manejador para actualizar un usuario
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	user, err := h.userUseCase.UpdateUser(id, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuario actualizado con éxito", user)
}

// DeleteUser manejador para eliminar un usuario
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.userUseCase.DeleteUser(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuario eliminado con éxito", nil)
}

// ArchiveUser manejador para archivar un usuario
func (h *UserHandler) ArchiveUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.userUseCase.ArchiveUser(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuario archivado con éxito", nil)
}

// ChangePassword manejador para cambiar la contraseña
func (h *UserHandler) ChangePassword(c *gin.Context) {
	// Obtener el ID del usuario del token (middleware)
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "No autorizado")
		return
	}

	var req domain.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	if err := h.userUseCase.ChangePassword(userID.(string), &req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Contraseña cambiada con éxito", nil)
}

// GetProfile obtiene el perfil del usuario autenticado
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Obtener el ID del usuario del token (middleware)
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "No autorizado")
		return
	}

	user, err := h.userUseCase.GetUser(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Perfil obtenido con éxito", user)
}
