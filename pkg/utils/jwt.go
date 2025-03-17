package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims es la estructura de claims para el JWT
type Claims struct {
	UserID string   `json:"user_id"`
	Role   string   `json:"role"`
	Scopes []string `json:"scopes"`
	jwt.RegisteredClaims
}

// GenerateJWT genera un nuevo token JWT
func GenerateJWT(userID, role string, scopes []string, secret string, expiration time.Duration) (string, error) {
	// Preparar claims
	claims := &Claims{
		UserID: userID,
		Role:   role,
		Scopes: scopes,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Crear token con claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar token con clave secreta
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT valida un token JWT y retorna los claims
func ValidateJWT(tokenString, secret string) (string, map[string]interface{}, error) {
	// Parsear token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verificar que el método de firma sea HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inválido")
		}

		// Retornar clave secreta
		return []byte(secret), nil
	})

	if err != nil {
		return "", nil, err
	}

	// Verificar que el token sea válido
	if !token.Valid {
		return "", nil, errors.New("token inválido")
	}

	// Obtener claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil, errors.New("no se pudieron obtener los claims")
	}

	// Extraer userID
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", nil, errors.New("user_id no encontrado en claims")
	}

	// Convertir claims a map
	claimsMap := make(map[string]interface{})
	for key, value := range claims {
		claimsMap[key] = value
	}

	return userID, claimsMap, nil
}
