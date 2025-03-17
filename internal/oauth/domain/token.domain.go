package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Token representa un token OAuth 2.0
type Token struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AccessToken      string             `json:"access_token" bson:"access_token"`
	RefreshToken     string             `json:"refresh_token" bson:"refresh_token"`
	UserID           string             `json:"user_id" bson:"user_id,omitempty"`
	ClientID         string             `json:"client_id" bson:"client_id"`
	Scopes           []string           `json:"scopes" bson:"scopes"`
	ExpiresAt        time.Time          `json:"expires_at" bson:"expires_at"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	RefreshExpiresAt time.Time          `json:"refresh_expires_at" bson:"refresh_expires_at"`
}

// TokenRepository define el contrato para la capa de persistencia
type TokenRepository interface {
	Create(token *Token) error
	GetByAccessToken(accessToken string) (*Token, error)
	GetByRefreshToken(refreshToken string) (*Token, error)
	DeleteByRefreshToken(refreshToken string) error
	DeleteByUserID(userID string) error
}
