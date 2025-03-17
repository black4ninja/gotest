package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Client representa un cliente OAuth 2.0
type Client struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ClientID     string             `json:"client_id" bson:"client_id"`
	ClientSecret string             `json:"client_secret" bson:"client_secret"`
	Name         string             `json:"name" bson:"name"`
	RedirectURIs []string           `json:"redirect_uris" bson:"redirect_uris"`
	GrantTypes   []string           `json:"grant_types" bson:"grant_types"`
	Scopes       []string           `json:"scopes" bson:"scopes"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

// ClientRepository define el contrato para la capa de persistencia
type ClientRepository interface {
	GetByClientID(clientID string) (*Client, error)
	ValidateClient(clientID, clientSecret string) (*Client, error)
	Create(client *Client) error
	Update(client *Client) error
	Delete(id string) error
}
