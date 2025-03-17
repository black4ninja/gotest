package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/black4ninja/mi-proyecto/internal/permission/domain"
	"github.com/black4ninja/mi-proyecto/pkg/utils"
)

// PermissionHandler maneja las peticiones HTTP para permisos
type PermissionHandler struct {
	permissionUC domain.PermissionUseCase
	roleUC       domain.RoleUseCase
	userRoleUC   domain.UserRoleUseCase
}

// NewPermissionHandler crea un nuevo manejador de permisos
func NewPermissionHandler(
	router *gin.RouterGroup,
	permissionUC domain.PermissionUseCase,
	roleUC domain.RoleUseCase,
	userRoleUC domain.UserRoleUseCase,
) {
	handler := &PermissionHandler{
		permissionUC: permissionUC,
		roleUC:       roleUC,
		userRoleUC:   userRoleUC,
	}

	// Rutas de permisos
	permissions := router.Group("/permissions")
	{
		permissions.GET("/", handler.GetAllPermissions)
		permissions.GET("/:id", handler.GetPermission)
		permissions.GET("/code/:code", handler.GetPermissionByCode)
		permissions.GET("/module/:module", handler.GetPermissionsByModule)
		permissions.POST("/", handler.CreatePermission)
		permissions.PUT("/:id", handler.UpdatePermission)
		permissions.DELETE("/:id", handler.DeletePermission)
	}

	// Rutas de roles
	roles := router.Group("/roles")
	{
		roles.GET("/", handler.GetAllRoles)
		roles.GET("/:id", handler.GetRole)
		roles.GET("/name/:name", handler.GetRoleByName)
		roles.POST("/", handler.CreateRole)
		roles.PUT("/:id", handler.UpdateRole)
		roles.DELETE("/:id", handler.DeleteRole)
		roles.POST("/:id/permissions", handler.AddPermissionToRole)
		roles.DELETE("/:id/permissions/:permissionCode", handler.RemovePermissionFromRole)
	}

	// Rutas de asignación usuario-rol
	userRoles := router.Group("/user-roles")
	{
		userRoles.GET("/:userID", handler.GetUserRoles)
		userRoles.POST("/assign-role", handler.AssignRoleToUser)
		userRoles.DELETE("/remove-role", handler.RemoveRoleFromUser)
		userRoles.POST("/assign-permission", handler.AssignPermissionToUser)
		userRoles.DELETE("/remove-permission", handler.RemovePermissionFromUser)
		userRoles.GET("/:userID/permissions", handler.GetUserPermissions)
		userRoles.GET("/:userID/has-permission/:permissionCode", handler.CheckUserPermission)
	}
}

// GetAllPermissions manejador para obtener todos los permisos
// @Summary Obtener todos los permissions
// @Description Obtiene una lista de todos los permissions con filtrado opcional
// @Tags permissions
// @Accept json
// @Produce json
// @Param status query string false "Estado del permission (active, inactive, archived)"
// @Param name query string false "Nombre del permission (búsqueda parcial)"
// @Success 200 {object} utils.Response{data=[]domain.PermissionResponse} "Lista de permissions"
// @Failure 500 {object} utils.Response "Error interno"
// @Router /permissions [get]
// @Security BearerAuth
func (h *PermissionHandler) GetAllPermissions(c *gin.Context) {
	permissions, err := h.permissionUC.GetAllPermissions()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permisos obtenidos con éxito", permissions)
}

// GetPermission manejador para obtener un permiso por ID
// @Summary Obtener un permission
// @Description Obtiene un permission por su ID
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path string true "ID del permission"
// @Success 200 {object} utils.Response{data=domain.PermissionResponse} "Permission obtenido"
// @Failure 404 {object} utils.Response "No encontrado"
// @Failure 500 {object} utils.Response "Error interno"
// @Router /permissions/{id} [get]
// @Security BearerAuth
func (h *PermissionHandler) GetPermission(c *gin.Context) {
	id := c.Param("id")

	permission, err := h.permissionUC.GetPermission(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permiso obtenido con éxito", permission)
}

// GetPermissionByCode manejador para obtener un permiso por código
func (h *PermissionHandler) GetPermissionByCode(c *gin.Context) {
	code := c.Param("code")

	permission, err := h.permissionUC.GetPermissionByCode(code)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permiso obtenido con éxito", permission)
}

// GetPermissionsByModule manejador para obtener permisos por módulo
func (h *PermissionHandler) GetPermissionsByModule(c *gin.Context) {
	module := c.Param("module")

	permissions, err := h.permissionUC.GetPermissionsByModule(module)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permisos obtenidos con éxito", permissions)
}

// CreatePermission manejador para crear un permiso
// @Summary Crear un permission
// @Description Crea un nuevo permission
// @Tags permissions
// @Accept json
// @Produce json
// @Param permission body domain.CreatePermissionRequest true "Datos del permission"
// @Success 201 {object} utils.Response{data=domain.PermissionResponse} "Permission creado"
// @Failure 400 {object} utils.Response "Datos inválidos"
// @Failure 500 {object} utils.Response "Error interno"
// @Router /permissions [post]
// @Security BearerAuth
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req domain.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	permission, err := h.permissionUC.CreatePermission(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Permiso creado con éxito", permission)
}

// UpdatePermission manejador para actualizar un permiso
// @Summary Actualizar un permission
// @Description Actualiza un permission existente
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path string true "ID del permission"
// @Param permission body domain.UpdatePermissionRequest true "Datos a actualizar"
// @Success 200 {object} utils.Response{data=domain.PermissionResponse} "Permission actualizado"
// @Failure 400 {object} utils.Response "Datos inválidos"
// @Failure 404 {object} utils.Response "No encontrado"
// @Failure 500 {object} utils.Response "Error interno"
// @Router /permissions/{id} [put]
// @Security BearerAuth
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id := c.Param("id")

	var req domain.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	permission, err := h.permissionUC.UpdatePermission(id, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permiso actualizado con éxito", permission)
}

// DeletePermission manejador para eliminar un permiso
// @Summary Eliminar un permission
// @Description Elimina un permission por su ID
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path string true "ID del permission"
// @Success 200 {object} utils.Response "Permission eliminado"
// @Failure 404 {object} utils.Response "No encontrado"
// @Failure 500 {object} utils.Response "Error interno"
// @Router /permissions/{id} [delete]
// @Security BearerAuth
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id := c.Param("id")

	err := h.permissionUC.DeletePermission(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permiso eliminado con éxito", nil)
}

// GetAllRoles manejador para obtener todos los roles
func (h *PermissionHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.roleUC.GetAllRoles()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Roles obtenidos con éxito", roles)
}

// GetRole manejador para obtener un rol por ID
func (h *PermissionHandler) GetRole(c *gin.Context) {
	id := c.Param("id")

	role, err := h.roleUC.GetRole(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Rol obtenido con éxito", role)
}

// GetRoleByName manejador para obtener un rol por nombre
func (h *PermissionHandler) GetRoleByName(c *gin.Context) {
	name := c.Param("name")

	role, err := h.roleUC.GetRoleByName(name)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Rol obtenido con éxito", role)
}

// CreateRole manejador para crear un rol
func (h *PermissionHandler) CreateRole(c *gin.Context) {
	var req domain.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	role, err := h.roleUC.CreateRole(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Rol creado con éxito", role)
}

// UpdateRole manejador para actualizar un rol
func (h *PermissionHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")

	var req domain.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	role, err := h.roleUC.UpdateRole(id, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Rol actualizado con éxito", role)
}

// DeleteRole manejador para eliminar un rol
func (h *PermissionHandler) DeleteRole(c *gin.Context) {
	id := c.Param("id")

	err := h.roleUC.DeleteRole(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Rol eliminado con éxito", nil)
}

// AddPermissionToRole manejador para añadir un permiso a un rol
func (h *PermissionHandler) AddPermissionToRole(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		PermissionCode string `json:"permission_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	err := h.roleUC.AddPermissionToRole(id, req.PermissionCode)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permiso añadido al rol con éxito", nil)
}

// RemovePermissionFromRole manejador para eliminar un permiso de un rol
func (h *PermissionHandler) RemovePermissionFromRole(c *gin.Context) {
	id := c.Param("id")
	permissionCode := c.Param("permissionCode")

	err := h.roleUC.RemovePermissionFromRole(id, permissionCode)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permiso eliminado del rol con éxito", nil)
}

// GetUserRoles manejador para obtener los roles de un usuario
func (h *PermissionHandler) GetUserRoles(c *gin.Context) {
	userID := c.Param("userID")

	userRoles, err := h.userRoleUC.GetUserRoles(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Roles de usuario obtenidos con éxito", userRoles)
}

// AssignRoleToUser manejador para asignar un rol a un usuario
func (h *PermissionHandler) AssignRoleToUser(c *gin.Context) {
	var req domain.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	err := h.userRoleUC.AssignRoleToUser(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Rol asignado al usuario con éxito", nil)
}

// RemoveRoleFromUser manejador para eliminar un rol de un usuario
func (h *PermissionHandler) RemoveRoleFromUser(c *gin.Context) {
	var req domain.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	err := h.userRoleUC.RemoveRoleFromUser(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Rol eliminado del usuario con éxito", nil)
}

// AssignPermissionToUser manejador para asignar un permiso a un usuario
func (h *PermissionHandler) AssignPermissionToUser(c *gin.Context) {
	var req domain.AssignPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	err := h.userRoleUC.AssignPermissionToUser(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permiso asignado al usuario con éxito", nil)
}

// RemovePermissionFromUser manejador para eliminar un permiso de un usuario
func (h *PermissionHandler) RemovePermissionFromUser(c *gin.Context) {
	var req domain.AssignPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	err := h.userRoleUC.RemovePermissionFromUser(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permiso eliminado del usuario con éxito", nil)
}

// GetUserPermissions manejador para obtener los permisos de un usuario
func (h *PermissionHandler) GetUserPermissions(c *gin.Context) {
	userID := c.Param("userID")

	permissions, err := h.userRoleUC.GetUserPermissions(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permisos de usuario obtenidos con éxito", permissions)
}

// CheckUserPermission manejador para verificar si un usuario tiene un permiso
func (h *PermissionHandler) CheckUserPermission(c *gin.Context) {
	userID := c.Param("userID")
	permissionCode := c.Param("permissionCode")

	hasPermission, err := h.userRoleUC.HasPermission(userID, permissionCode)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Verificación de permiso completada", gin.H{
		"has_permission": hasPermission,
	})
}
