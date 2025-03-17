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

type mongoClientRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

// NewMongoClientRepository crea un nuevo repositorio de clientes con MongoDB
func NewMongoClientRepository(collection *mongo.Collection) domain.ClientRepository {
	return &mongoClientRepository{
		collection: collection,
		timeout:    10 * time.Second,
	}
}

// GetByClientID obtiene un cliente por su ID de cliente
func (r *mongoClientRepository) GetByClientID(clientID string) (*domain.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var client domain.Client
	err := r.collection.FindOne(ctx, bson.M{"client_id": clientID}).Decode(&client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("cliente no encontrado")
		}
		return nil, err
	}

	return &client, nil
}

// ValidateClient valida las credenciales de un cliente
func (r *mongoClientRepository) ValidateClient(clientID, clientSecret string) (*domain.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var client domain.Client
	err := r.collection.FindOne(ctx, bson.M{
		"client_id":     clientID,
		"client_secret": clientSecret,
	}).Decode(&client)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("credenciales de cliente inválidas")
		}
		return nil, err
	}

	return &client, nil
}

// Create crea un nuevo cliente
func (r *mongoClientRepository) Create(client *domain.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	client.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, client)
	return err
}

// Update actualiza un cliente existente
func (r *mongoClientRepository) Update(client *domain.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"name":          client.Name,
			"redirect_uris": client.RedirectURIs,
			"grant_types":   client.GrantTypes,
			"scopes":        client.Scopes,
			"updated_at":    time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": client.ID},
		update,
	)
	return err
}

// Delete elimina un cliente
func (r *mongoClientRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
