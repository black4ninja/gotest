package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/black4ninja/mi-proyecto/internal/oauth/domain"
	"github.com/black4ninja/mi-proyecto/pkg/utils"
)

// OAuthHandler maneja las peticiones HTTP para OAuth
type OAuthHandler struct {
	oauthUseCase domain.OAuthUseCase
}

// NewOAuthHandler crea un nuevo manejador de OAuth
func NewOAuthHandler(router *gin.RouterGroup, useCase domain.OAuthUseCase) {
	handler := &OAuthHandler{
		oauthUseCase: useCase,
	}

	// Rutas OAuth
	router.POST("/token", handler.GenerateToken)
	router.POST("/revoke", handler.RevokeToken)
}

// GenerateToken manejador para generar tokens OAuth
func (h *OAuthHandler) GenerateToken(c *gin.Context) {
	var req domain.OAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.oauthUseCase.GenerateToken(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, token)
}

// RevokeToken manejador para revocar tokens
func (h *OAuthHandler) RevokeToken(c *gin.Context) {
	type RevokeRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	var req RevokeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.oauthUseCase.RevokeToken(req.RefreshToken); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Token revocado con éxito", nil)
}
