package usecase

import (
	"errors"
	"strings"
	"time"

	"github.com/black4ninja/mi-proyecto/internal/oauth/domain"
	userDomain "github.com/black4ninja/mi-proyecto/internal/user/domain"
	"github.com/black4ninja/mi-proyecto/pkg/utils"
)

type oauthUseCase struct {
	clientRepo domain.ClientRepository
	tokenRepo  domain.TokenRepository
	userUC     userDomain.UserUseCase
	jwtSecret  string
	tokenExp   time.Duration
	refreshExp time.Duration
}

// NewOAuthUseCase crea un nuevo caso de uso para OAuth
func NewOAuthUseCase(
	clientRepo domain.ClientRepository,
	tokenRepo domain.TokenRepository,
	userUC userDomain.UserUseCase,
	jwtSecret string,
	tokenExp time.Duration,
	refreshExp time.Duration,
) domain.OAuthUseCase {
	return &oauthUseCase{
		clientRepo: clientRepo,
		tokenRepo:  tokenRepo,
		userUC:     userUC,
		jwtSecret:  jwtSecret,
		tokenExp:   tokenExp,
		refreshExp: refreshExp,
	}
}

// GenerateToken genera un token OAuth 2.0
func (u *oauthUseCase) GenerateToken(req *domain.OAuthRequest) (*domain.OAuthResponse, error) {
	// Validar cliente
	client, err := u.clientRepo.ValidateClient(req.ClientID, req.ClientSecret)
	if err != nil {
		return nil, err
	}

	// Verificar si el tipo de concesión es válido para este cliente
	if !contains(client.GrantTypes, req.GrantType) {
		return nil, errors.New("tipo de concesión no permitido para este cliente")
	}

	// Verificar scopes
	var scopes []string
	if req.Scope != "" {
		requestedScopes := strings.Split(req.Scope, " ")
		for _, s := range requestedScopes {
			if contains(client.Scopes, s) {
				scopes = append(scopes, s)
			}
		}
	}

	// Si no se proporcionaron scopes válidos, usar los scopes por defecto del cliente
	if len(scopes) == 0 {
		scopes = client.Scopes
	}

	// Generar tokens según el tipo de concesión
	switch req.GrantType {
	case domain.GrantTypePassword:
		return u.handlePasswordGrant(req, client, scopes)
	case domain.GrantTypeRefreshToken:
		return u.handleRefreshTokenGrant(req, client, scopes)
	case domain.GrantTypeClientCredentials:
		return u.handleClientCredentialsGrant(client, scopes)
	default:
		return nil, errors.New("tipo de concesión no implementado")
	}
}

// handlePasswordGrant maneja la concesión de tipo password
func (u *oauthUseCase) handlePasswordGrant(req *domain.OAuthRequest, client *domain.Client, scopes []string) (*domain.OAuthResponse, error) {
	// Validar que se proporcionaron username y password
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("nombre de usuario y contraseña requeridos")
	}

	// Validar credenciales del usuario
	user, err := u.userUC.ValidateCredentials(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	// Generar tokens
	accessToken, err := utils.GenerateJWT(user.ID.Hex(), user.Role, scopes, u.jwtSecret, u.tokenExp)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRandomToken(32)
	if err != nil {
		return nil, err
	}

	// Guardar token en la base de datos
	expiresAt := time.Now().Add(u.tokenExp)
	refreshExpiresAt := time.Now().Add(u.refreshExp)
	token := &domain.Token{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		UserID:           user.ID.Hex(),
		ClientID:         client.ClientID,
		Scopes:           scopes,
		ExpiresAt:        expiresAt,
		RefreshExpiresAt: refreshExpiresAt,
		CreatedAt:        time.Now(),
	}

	if err := u.tokenRepo.Create(token); err != nil {
		return nil, err
	}

	// Actualizar refresh token del usuario
	if err := u.userUC.UpdateRefreshToken(user.ID.Hex(), refreshToken); err != nil {
		return nil, err
	}

	// Preparar respuesta
	return &domain.OAuthResponse{
		AccessToken:  accessToken,
		TokenType:    domain.TokenTypeBearer,
		ExpiresIn:    int(u.tokenExp.Seconds()),
		RefreshToken: refreshToken,
		Scope:        strings.Join(scopes, " "),
	}, nil
}

// handleRefreshTokenGrant maneja la concesión de tipo refresh_token
// En internal/oauth/usecase/oauth_usecase.go, actualiza la función handleRefreshTokenGrant

func (u *oauthUseCase) handleRefreshTokenGrant(req *domain.OAuthRequest, client *domain.Client, scopes []string) (*domain.OAuthResponse, error) {
	// Validar que se proporcionó un refresh token
	if req.RefreshToken == "" {
		return nil, errors.New("refresh token requerido")
	}

	// Usar la nueva función de validación
	oldToken, err := u.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Verificar que el token pertenezca al mismo cliente
	if oldToken.ClientID != client.ClientID {
		return nil, errors.New("refresh token no válido para este cliente")
	}

	// Si no se proporcionaron scopes, usar los del token anterior
	if len(scopes) == 0 {
		scopes = oldToken.Scopes
	}

	// Generar nuevos tokens
	accessToken, err := utils.GenerateJWT(oldToken.UserID, "", scopes, u.jwtSecret, u.tokenExp)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRandomToken(32)
	if err != nil {
		return nil, err
	}

	// Eliminar token antiguo
	if err := u.tokenRepo.DeleteByRefreshToken(req.RefreshToken); err != nil {
		return nil, err
	}

	// Guardar nuevo token con fechas de expiración configuradas
	accessExpiresAt := time.Now().Add(u.tokenExp)
	refreshExpiresAt := time.Now().Add(u.refreshExp)

	token := &domain.Token{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		UserID:           oldToken.UserID,
		ClientID:         client.ClientID,
		Scopes:           scopes,
		ExpiresAt:        accessExpiresAt,
		RefreshExpiresAt: refreshExpiresAt,
		CreatedAt:        time.Now(),
	}

	if err := u.tokenRepo.Create(token); err != nil {
		return nil, err
	}

	// Actualizar refresh token del usuario si hay un usuario asociado
	if oldToken.UserID != "" {
		if err := u.userUC.UpdateRefreshToken(oldToken.UserID, refreshToken); err != nil {
			return nil, err
		}
	}

	// Preparar respuesta
	return &domain.OAuthResponse{
		AccessToken:  accessToken,
		TokenType:    domain.TokenTypeBearer,
		ExpiresIn:    int(u.tokenExp.Seconds()),
		RefreshToken: refreshToken,
		Scope:        strings.Join(scopes, " "),
	}, nil
}

// handleClientCredentialsGrant maneja la concesión de tipo client_credentials
func (u *oauthUseCase) handleClientCredentialsGrant(client *domain.Client, scopes []string) (*domain.OAuthResponse, error) {
	// Generar access token para el cliente (sin usuario asociado)
	accessToken, err := utils.GenerateJWT("", "client", scopes, u.jwtSecret, u.tokenExp)
	if err != nil {
		return nil, err
	}

	// No se genera refresh token para client credentials
	expiresAt := time.Now().Add(u.tokenExp)
	token := &domain.Token{
		AccessToken: accessToken,
		ClientID:    client.ClientID,
		Scopes:      scopes,
		ExpiresAt:   expiresAt,
		CreatedAt:   time.Now(),
	}

	if err := u.tokenRepo.Create(token); err != nil {
		return nil, err
	}

	// Preparar respuesta
	return &domain.OAuthResponse{
		AccessToken: accessToken,
		TokenType:   domain.TokenTypeBearer,
		ExpiresIn:   int(u.tokenExp.Seconds()),
		Scope:       strings.Join(scopes, " "),
	}, nil
}

// ValidateToken valida un token de acceso
func (u *oauthUseCase) ValidateToken(accessToken string) (string, map[string]interface{}, error) {
	// Verificar que el token exista en la base de datos
	token, err := u.tokenRepo.GetByAccessToken(accessToken)
	if err != nil {
		return "", nil, errors.New("token inválido")
	}

	// Verificar que el token no haya expirado
	if time.Now().After(token.ExpiresAt) {
		return "", nil, errors.New("token expirado")
	}

	// Verificar y decodificar JWT
	userID, claims, err := utils.ValidateJWT(accessToken, u.jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return userID, claims, nil
}

// ValidateRefreshToken valida un token de refresco y retorna el token si es válido
func (u *oauthUseCase) ValidateRefreshToken(refreshToken string) (*domain.Token, error) {
	// Buscar token en la base de datos
	token, err := u.tokenRepo.GetByRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("refresh token inválido o no encontrado")
	}

	// Verificar que el token no haya expirado
	if token.RefreshExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token expirado")
	}

	// Si hay un usuario asociado, verificar que esté activo
	if token.UserID != "" {
		user, err := u.userUC.GetUser(token.UserID)
		if err != nil {
			return nil, errors.New("usuario no encontrado")
		}

		if user.Status != userDomain.UserStatusActive {
			return nil, errors.New("usuario inactivo")
		}
	}

	return token, nil
}

// RevokeToken revoca un token de refresco
func (u *oauthUseCase) RevokeToken(refreshToken string) error {
	// Eliminar token
	if err := u.tokenRepo.DeleteByRefreshToken(refreshToken); err != nil {
		return err
	}

	// Intentar obtener usuario por refresh token
	user, err := u.userUC.GetUserByRefreshToken(refreshToken)
	if err == nil {
		// Limpiar refresh token del usuario
		_ = u.userUC.UpdateRefreshToken(user.ID.Hex(), "")
	}

	return nil
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
