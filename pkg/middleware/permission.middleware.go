package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/black4ninja/mi-proyecto/internal/permission/domain"
)

// PermissionMiddleware es un middleware para verificar permisos
type PermissionMiddleware struct {
	userRoleUseCase domain.UserRoleUseCase
}

// NewPermissionMiddleware crea un nuevo middleware de permisos
func NewPermissionMiddleware(userRoleUseCase domain.UserRoleUseCase) *PermissionMiddleware {
	return &PermissionMiddleware{
		userRoleUseCase: userRoleUseCase,
	}
}

// RequirePermission verifica que el usuario tenga un permiso específico
func (m *PermissionMiddleware) RequirePermission(permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el ID de usuario del contexto (establecido por el middleware de autenticación)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  "No autenticado",
			})
			c.Abort()
			return
		}

		// Verificar permiso
		hasPermission, err := m.userRoleUseCase.HasPermission(userID.(string), permissionCode)
		if err != nil || !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error":  "Permiso denegado: se requiere " + permissionCode,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission verifica que el usuario tenga al menos uno de los permisos especificados
func (m *PermissionMiddleware) RequireAnyPermission(permissionCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el ID de usuario del contexto (establecido por el middleware de autenticación)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  "No autenticado",
			})
			c.Abort()
			return
		}

		// Verificar si tiene al menos uno de los permisos
		for _, permissionCode := range permissionCodes {
			hasPermission, err := m.userRoleUseCase.HasPermission(userID.(string), permissionCode)
			if err == nil && hasPermission {
				c.Next()
				return
			}
		}

		// Si no tiene ninguno de los permisos
		c.JSON(http.StatusForbidden, gin.H{
			"status": "error",
			"error":  "Permiso denegado: se requiere al menos uno de los permisos especificados",
		})
		c.Abort()
	}
}

// RequireAllPermissions verifica que el usuario tenga todos los permisos especificados
func (m *PermissionMiddleware) RequireAllPermissions(permissionCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el ID de usuario del contexto (establecido por el middleware de autenticación)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  "No autenticado",
			})
			c.Abort()
			return
		}

		// Verificar que tenga todos los permisos
		for _, permissionCode := range permissionCodes {
			hasPermission, err := m.userRoleUseCase.HasPermission(userID.(string), permissionCode)
			if err != nil || !hasPermission {
				c.JSON(http.StatusForbidden, gin.H{
					"status": "error",
					"error":  "Permiso denegado: se requieren todos los permisos especificados",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// RequireModuleAccess verifica que el usuario tenga acceso a un módulo completo
func (m *PermissionMiddleware) RequireModuleAccess(module string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el ID de usuario del contexto (establecido por el middleware de autenticación)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  "No autenticado",
			})
			c.Abort()
			return
		}

		// Verificar acceso al módulo (permisos que comienzan con "module:")
		moduleWildcard := module + ":*"
		hasPermission, err := m.userRoleUseCase.HasPermission(userID.(string), moduleWildcard)
		if err != nil || !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error":  "Permiso denegado: se requiere acceso al módulo " + module,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
