package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/black4ninja/mi-proyecto/internal/oauth/domain"
)

type mongoTokenRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

// NewMongoTokenRepository crea un nuevo repositorio de tokens con MongoDB
func NewMongoTokenRepository(collection *mongo.Collection) domain.TokenRepository {
	return &mongoTokenRepository{
		collection: collection,
		timeout:    10 * time.Second,
	}
}

// Create crea un nuevo token
func (r *mongoTokenRepository) Create(token *domain.Token) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	token.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, token)
	return err
}

// GetByAccessToken obtiene un token por su access token
func (r *mongoTokenRepository) GetByAccessToken(accessToken string) (*domain.Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var token domain.Token
	err := r.collection.FindOne(ctx, bson.M{"access_token": accessToken}).Decode(&token)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("token no encontrado")
		}
		return nil, err
	}

	return &token, nil
}

// GetByRefreshToken obtiene un token por su refresh token
func (r *mongoTokenRepository) GetByRefreshToken(refreshToken string) (*domain.Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var token domain.Token
	err := r.collection.FindOne(ctx, bson.M{"refresh_token": refreshToken}).Decode(&token)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("token no encontrado")
		}
		return nil, err
	}

	return &token, nil
}

// DeleteByRefreshToken elimina un token por su refresh token
func (r *mongoTokenRepository) DeleteByRefreshToken(refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"refresh_token": refreshToken})
	return err
}

// DeleteByUserID elimina todos los tokens de un usuario
func (r *mongoTokenRepository) DeleteByUserID(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.collection.DeleteMany(ctx, bson.M{"user_id": userID})
	return err
}
