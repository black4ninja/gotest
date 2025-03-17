package domain

// GrantType representa los tipos de concesión de OAuth 2.0
const (
	GrantTypeAuthorizationCode = "authorization_code"
	GrantTypePassword          = "password"
	GrantTypeClientCredentials = "client_credentials"
	GrantTypeRefreshToken      = "refresh_token"
)

// TokenType representa los tipos de token
const (
	TokenTypeBearer = "Bearer"
)

// OAuthRequest representa la solicitud de token OAuth 2.0
type OAuthRequest struct {
	GrantType    string `json:"grant_type" binding:"required"`
	ClientID     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// OAuthResponse representa la respuesta de token OAuth 2.0
type OAuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// OAuthUseCase define el contrato para la capa de casos de uso
type OAuthUseCase interface {
	GenerateToken(req *OAuthRequest) (*OAuthResponse, error)
	ValidateToken(accessToken string) (string, map[string]interface{}, error)
	ValidateRefreshToken(refreshToken string) (*Token, error)
	RevokeToken(refreshToken string) error
}
