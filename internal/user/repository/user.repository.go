package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/black4ninja/mi-proyecto/internal/user/domain"
)

type mongoUserRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

// NewMongoUserRepository crea un nuevo repositorio de usuarios con MongoDB
func NewMongoUserRepository(collection *mongo.Collection) domain.UserRepository {
	return &mongoUserRepository{
		collection: collection,
		timeout:    10 * time.Second,
	}
}

// GetByID obtiene un usuario por su ID
func (r *mongoUserRepository) GetByID(id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user domain.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail obtiene un usuario por su email
func (r *mongoUserRepository) GetByEmail(email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	return &user, nil
}

// GetAll obtiene todos los usuarios que coincidan con los parámetros dados
func (r *mongoUserRepository) GetAll(params map[string]interface{}) ([]*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	// Construir filtro
	filter := bson.M{}
	for key, value := range params {
		filter[key] = value
	}

	opts := options.Find()
	opts.SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// Create crea un nuevo usuario
func (r *mongoUserRepository) Create(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	user.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// Update actualiza un usuario existente
func (r *mongoUserRepository) Update(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"email":      user.Email,
			"status":     user.Status,
			"role":       user.Role,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		update,
	)
	return err
}

// Delete elimina un usuario
func (r *mongoUserRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// Archive marca un usuario como archivado
func (r *mongoUserRepository) Archive(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":      domain.UserStatusArchived,
			"archived_at": now,
			"updated_at":  now,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

// UpdateRefreshToken actualiza el token de refresco de un usuario
func (r *mongoUserRepository) UpdateRefreshToken(userID string, refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"refresh_token": refreshToken,
			"updated_at":    time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

// GetByRefreshToken obtiene un usuario por su token de refresco
func (r *mongoUserRepository) GetByRefreshToken(refreshToken string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"refresh_token": refreshToken}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("token de refresco inválido")
		}
		return nil, err
	}

	return &user, nil
}
