package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/black4ninja/mi-proyecto/internal/permission/domain"
)

type mongoPermissionRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

// NewMongoPermissionRepository crea un nuevo repositorio de permisos con MongoDB
func NewMongoPermissionRepository(collection *mongo.Collection) domain.PermissionRepository {
	return &mongoPermissionRepository{
		collection: collection,
		timeout:    10 * time.Second,
	}
}

// GetByID obtiene un permiso por su ID
func (r *mongoPermissionRepository) GetByID(id string) (*domain.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var permission domain.Permission
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&permission)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("permiso no encontrado")
		}
		return nil, err
	}

	return &permission, nil
}

// GetByCode obtiene un permiso por su código
func (r *mongoPermissionRepository) GetByCode(code string) (*domain.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var permission domain.Permission
	err := r.collection.FindOne(ctx, bson.M{"code": code}).Decode(&permission)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("permiso no encontrado")
		}
		return nil, err
	}

	return &permission, nil
}

// GetByModule obtiene los permisos para un módulo específico
func (r *mongoPermissionRepository) GetByModule(module string) ([]*domain.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"module": module})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var permissions []*domain.Permission
	if err := cursor.All(ctx, &permissions); err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetAll obtiene todos los permisos
func (r *mongoPermissionRepository) GetAll() ([]*domain.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	opts := options.Find().SetSort(bson.M{"module": 1, "code": 1})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var permissions []*domain.Permission
	if err := cursor.All(ctx, &permissions); err != nil {
		return nil, err
	}

	return permissions, nil
}

// Create crea un nuevo permiso
func (r *mongoPermissionRepository) Create(permission *domain.Permission) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	// Verificar si ya existe un permiso con el mismo código
	count, err := r.collection.CountDocuments(ctx, bson.M{"code": permission.Code})
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("ya existe un permiso con este código")
	}

	// Crear el permiso
	permission.ID = primitive.NewObjectID()
	_, err = r.collection.InsertOne(ctx, permission)

	return err
}

// Update actualiza un permiso existente
func (r *mongoPermissionRepository) Update(permission *domain.Permission) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"name":        permission.Name,
			"description": permission.Description,
			"updated_at":  time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": permission.ID}, update)

	return err
}

// Delete elimina un permiso
func (r *mongoPermissionRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})

	return err
}

// GetByCodesArray obtiene permisos por array de códigos
func (r *mongoPermissionRepository) GetByCodesArray(codes []string) ([]*domain.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"code": bson.M{"$in": codes}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var permissions []*domain.Permission
	if err := cursor.All(ctx, &permissions); err != nil {
		return nil, err
	}

	return permissions, nil
}
