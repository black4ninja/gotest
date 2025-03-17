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

type mongoRoleRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

// NewMongoRoleRepository crea un nuevo repositorio de roles con MongoDB
func NewMongoRoleRepository(collection *mongo.Collection) domain.RoleRepository {
	return &mongoRoleRepository{
		collection: collection,
		timeout:    10 * time.Second,
	}
}

// GetByID obtiene un rol por su ID
func (r *mongoRoleRepository) GetByID(id string) (*domain.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var role domain.Role
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("rol no encontrado")
		}
		return nil, err
	}

	return &role, nil
}

// GetByName obtiene un rol por su nombre
func (r *mongoRoleRepository) GetByName(name string) (*domain.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var role domain.Role
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("rol no encontrado")
		}
		return nil, err
	}

	return &role, nil
}

// GetAll obtiene todos los roles
func (r *mongoRoleRepository) GetAll() ([]*domain.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	opts := options.Find().SetSort(bson.M{"name": 1})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var roles []*domain.Role
	if err := cursor.All(ctx, &roles); err != nil {
		return nil, err
	}

	return roles, nil
}

// Create crea un nuevo rol
func (r *mongoRoleRepository) Create(role *domain.Role) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	// Verificar si ya existe un rol con el mismo nombre
	count, err := r.collection.CountDocuments(ctx, bson.M{"name": role.Name})
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("ya existe un rol con este nombre")
	}

	// Crear el rol
	role.ID = primitive.NewObjectID()
	_, err = r.collection.InsertOne(ctx, role)

	return err
}

// Update actualiza un rol existente
func (r *mongoRoleRepository) Update(role *domain.Role) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	// Verificar si es un rol de sistema
	var existingRole domain.Role
	err := r.collection.FindOne(ctx, bson.M{"_id": role.ID}).Decode(&existingRole)
	if err != nil {
		return err
	}

	if existingRole.IsSystem {
		return errors.New("no se puede modificar un rol de sistema")
	}

	update := bson.M{
		"$set": bson.M{
			"name":        role.Name,
			"description": role.Description,
			"updated_at":  time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": role.ID}, update)

	return err
}

// Delete elimina un rol
func (r *mongoRoleRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Verificar si es un rol de sistema
	var role domain.Role
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&role)
	if err != nil {
		return err
	}

	if role.IsSystem {
		return errors.New("no se puede eliminar un rol de sistema")
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})

	return err
}

// AddPermission añade un permiso a un rol
func (r *mongoRoleRepository) AddPermission(roleID string, permissionCode string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		return err
	}

	// Verificar si es un rol de sistema
	var role domain.Role
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&role)
	if err != nil {
		return err
	}

	// Verificar si el permiso ya está en el rol
	for _, perm := range role.Permissions {
		if perm == permissionCode {
			return errors.New("el permiso ya está asignado a este rol")
		}
	}

	update := bson.M{
		"$push": bson.M{
			"permissions": permissionCode,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)

	return err
}

// RemovePermission elimina un permiso de un rol
func (r *mongoRoleRepository) RemovePermission(roleID string, permissionCode string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		return err
	}

	// Verificar si es un rol de sistema
	var role domain.Role
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&role)
	if err != nil {
		return err
	}

	if role.IsSystem {
		return errors.New("no se puede modificar un rol de sistema")
	}

	update := bson.M{
		"$pull": bson.M{
			"permissions": permissionCode,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)

	return err
}
