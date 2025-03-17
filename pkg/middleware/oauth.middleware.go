package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/black4ninja/mi-proyecto/internal/oauth/domain"
	"github.com/black4ninja/mi-proyecto/pkg/utils"
)

// OAuthMiddleware maneja la autenticación a nivel de middleware
type OAuthMiddleware struct {
	oauthUseCase domain.OAuthUseCase
}

// NewOAuthMiddleware crea un nuevo middleware de OAuth
func NewOAuthMiddleware(oauthUseCase domain.OAuthUseCase) *OAuthMiddleware {
	return &OAuthMiddleware{
		oauthUseCase: oauthUseCase,
	}
}

// Protected protege rutas verificando el token OAuth
func (m *OAuthMiddleware) Protected() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el header de autorización
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "No autorizado: token no proporcionado")
			c.Abort()
			return
		}

		// Verificar el formato del token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Formato de token inválido")
			c.Abort()
			return
		}

		// Obtener el token
		accessToken := parts[1]

		// Validar el token
		userID, claims, err := m.oauthUseCase.ValidateToken(accessToken)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		// Almacenar el userID en el contexto
		c.Set("userID", userID)

		// Almacenar claims en el contexto
		if claims != nil {
			for key, value := range claims {
				c.Set(key, value)
			}
		}

		c.Next()
	}
}

// RequireScope verifica que el token tenga el scope requerido
func (m *OAuthMiddleware) RequireScope(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar si hay scopes en el contexto
		scopes, exists := c.Get("scopes")
		if !exists {
			utils.ErrorResponse(c, http.StatusForbidden, "Acceso denegado: no se encontraron scopes")
			c.Abort()
			return
		}

		// Verificar si el scope requerido está presente
		scopeList, ok := scopes.([]string)
		if !ok || !contains(scopeList, scope) {
			utils.ErrorResponse(c, http.StatusForbidden, "Acceso denegado: scope requerido "+scope)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole verifica que el usuario tenga el rol requerido
func (m *OAuthMiddleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar si hay rol en el contexto
		userRole, exists := c.Get("role")
		if !exists {
			utils.ErrorResponse(c, http.StatusForbidden, "Acceso denegado: no se encontró rol")
			c.Abort()
			return
		}

		// Verificar si el rol es el requerido
		if userRole.(string) != role {
			utils.ErrorResponse(c, http.StatusForbidden, "Acceso denegado: rol requerido "+role)
			c.Abort()
			return
		}

		c.Next()
	}
}

// contains verifica si un slice contiene un elemento
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
