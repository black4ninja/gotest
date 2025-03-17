package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/black4ninja/mi-proyecto/internal/permission/domain"
)

type mongoUserRoleRepository struct {
	collection *mongo.Collection
	roleRepo   domain.RoleRepository
	timeout    time.Duration
}

// NewMongoUserRoleRepository crea un nuevo repositorio de asignaciones usuario-rol con MongoDB
func NewMongoUserRoleRepository(collection *mongo.Collection, roleRepo domain.RoleRepository) domain.UserRoleRepository {
	return &mongoUserRoleRepository{
		collection: collection,
		roleRepo:   roleRepo,
		timeout:    10 * time.Second,
	}
}

// GetByUserID obtiene las asignaciones de rol para un usuario
func (r *mongoUserRoleRepository) GetByUserID(userID string) (*domain.UserRole, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var userRole domain.UserRole
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&userRole)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Si no existe, creamos uno nuevo con roles y permisos vacíos
			userRole = domain.UserRole{
				ID:          primitive.NewObjectID(),
				UserID:      userID,
				Roles:       []string{},
				Permissions: []string{},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			_, err = r.collection.InsertOne(ctx, userRole)
			if err != nil {
				return nil, err
			}
			return &userRole, nil
		}
		return nil, err
	}

	return &userRole, nil
}

// Create crea una nueva asignación usuario-rol
func (r *mongoUserRoleRepository) Create(userRole *domain.UserRole) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	// Verificar si ya existe una asignación para este usuario
	count, err := r.collection.CountDocuments(ctx, bson.M{"user_id": userRole.UserID})
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("ya existe una asignación para este usuario")
	}

	// Crear la asignación
	userRole.ID = primitive.NewObjectID()
	_, err = r.collection.InsertOne(ctx, userRole)

	return err
}

// Update actualiza una asignación usuario-rol existente
func (r *mongoUserRoleRepository) Update(userRole *domain.UserRole) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"roles":       userRole.Roles,
			"permissions": userRole.Permissions,
			"updated_at":  time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": userRole.ID}, update)

	return err
}

// Delete elimina una asignación usuario-rol
func (r *mongoUserRoleRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})

	return err
}

// AddRole añade un rol a un usuario
func (r *mongoUserRoleRepository) AddRole(userID string, roleID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	// Verificar que el rol existe
	_, err := r.roleRepo.GetByID(roleID)
	if err != nil {
		return err
	}

	// Obtener o crear asignación de usuario
	userRole, err := r.GetByUserID(userID)
	if err != nil {
		return err
	}

	// Verificar si el rol ya está asignado
	for _, rid := range userRole.Roles {
		if rid == roleID {
			return errors.New("el rol ya está asignado a este usuario")
		}
	}

	// Añadir el rol
	update := bson.M{
		"$push": bson.M{
			"roles": roleID,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"user_id": userID}, update)

	return err
}

// RemoveRole elimina un rol de un usuario
func (r *mongoUserRoleRepository) RemoveRole(userID string, roleID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	update := bson.M{
		"$pull": bson.M{
			"roles": roleID,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"user_id": userID}, update)

	return err
}

// AddPermission añade un permiso específico a un usuario
func (r *mongoUserRoleRepository) AddPermission(userID string, permissionCode string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	// Obtener o crear asignación de usuario
	userRole, err := r.GetByUserID(userID)
	if err != nil {
		return err
	}

	// Verificar si el permiso ya está asignado
	for _, p := range userRole.Permissions {
		if p == permissionCode {
			return errors.New("el permiso ya está asignado a este usuario")
		}
	}

	// Añadir el permiso
	update := bson.M{
		"$push": bson.M{
			"permissions": permissionCode,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"user_id": userID}, update)

	return err
}

// RemovePermission elimina un permiso específico de un usuario
func (r *mongoUserRoleRepository) RemovePermission(userID string, permissionCode string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	update := bson.M{
		"$pull": bson.M{
			"permissions": permissionCode,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"user_id": userID}, update)

	return err
}

// GetUserPermissions obtiene todos los permisos de un usuario (combinando los de sus roles y los específicos)
func (r *mongoUserRoleRepository) GetUserPermissions(userID string) ([]string, error) {
	_, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	// Obtener asignación de usuario
	userRole, err := r.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Conjunto para almacenar permisos únicos
	permissionsSet := make(map[string]bool)

	// Añadir permisos específicos del usuario
	for _, p := range userRole.Permissions {
		permissionsSet[p] = true
	}

	// Añadir permisos de cada rol
	for _, roleID := range userRole.Roles {
		role, err := r.roleRepo.GetByID(roleID)
		if err != nil {
			continue // Ignorar roles que no existan
		}

		for _, p := range role.Permissions {
			permissionsSet[p] = true
		}
	}

	// Convertir conjunto a slice
	permissions := make([]string, 0, len(permissionsSet))
	for p := range permissionsSet {
		permissions = append(permissions, p)
	}

	return permissions, nil
}
