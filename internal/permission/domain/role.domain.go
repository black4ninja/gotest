package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Role representa un rol que agrupa múltiples permisos
type Role struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Permissions []string           `json:"permissions" bson:"permissions"` // Lista de códigos de permisos
	IsSystem    bool               `json:"is_system" bson:"is_system"`     // Indica si es un rol de sistema (no modificable)
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// UserRole representa la asignación de roles y permisos a un usuario
type UserRole struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      string             `json:"user_id" bson:"user_id"`
	Roles       []string           `json:"roles" bson:"roles"`             // IDs de roles asignados
	Permissions []string           `json:"permissions" bson:"permissions"` // Permisos específicos adicionales
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// RoleRepository define el contrato para la capa de persistencia de roles
type RoleRepository interface {
	GetByID(id string) (*Role, error)
	GetByName(name string) (*Role, error)
	GetAll() ([]*Role, error)
	Create(role *Role) error
	Update(role *Role) error
	Delete(id string) error
	AddPermission(roleID string, permissionCode string) error
	RemovePermission(roleID string, permissionCode string) error
}

// UserRoleRepository define el contrato para la capa de persistencia de asignaciones usuario-rol
type UserRoleRepository interface {
	GetByUserID(userID string) (*UserRole, error)
	Create(userRole *UserRole) error
	Update(userRole *UserRole) error
	Delete(id string) error
	AddRole(userID string, roleID string) error
	RemoveRole(userID string, roleID string) error
	AddPermission(userID string, permissionCode string) error
	RemovePermission(userID string, permissionCode string) error
	GetUserPermissions(userID string) ([]string, error) // Devuelve todos los permisos de un usuario (roles + específicos)
}

// CreateRoleRequest representa la solicitud para crear un rol
type CreateRoleRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"` // Lista de códigos de permisos
}

// UpdateRoleRequest representa la solicitud para actualizar un rol
type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// AssignRoleRequest representa la solicitud para asignar un rol a un usuario
type AssignRoleRequest struct {
	UserID string `json:"user_id" binding:"required"`
	RoleID string `json:"role_id" binding:"required"`
}

// AssignPermissionRequest representa la solicitud para asignar un permiso a un usuario
type AssignPermissionRequest struct {
	UserID         string `json:"user_id" binding:"required"`
	PermissionCode string `json:"permission_code" binding:"required"`
}

// RoleResponse representa la respuesta con datos de roles
type RoleResponse struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Permissions []*PermissionResponse `json:"permissions"`
	IsSystem    bool                  `json:"is_system"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

// UserRoleResponse representa la respuesta con datos de asignaciones usuario-rol
type UserRoleResponse struct {
	ID          string                `json:"id"`
	UserID      string                `json:"user_id"`
	Roles       []*RoleResponse       `json:"roles"`
	Permissions []*PermissionResponse `json:"permissions"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

// RoleUseCase define el contrato para la capa de caso de uso de roles
type RoleUseCase interface {
	GetRole(id string) (*RoleResponse, error)
	GetRoleByName(name string) (*RoleResponse, error)
	GetAllRoles() ([]*RoleResponse, error)
	CreateRole(req *CreateRoleRequest) (*RoleResponse, error)
	UpdateRole(id string, req *UpdateRoleRequest) (*RoleResponse, error)
	DeleteRole(id string) error
	AddPermissionToRole(roleID string, permissionCode string) error
	RemovePermissionFromRole(roleID string, permissionCode string) error
}

// UserRoleUseCase define el contrato para la capa de caso de uso de asignaciones usuario-rol
type UserRoleUseCase interface {
	GetUserRoles(userID string) (*UserRoleResponse, error)
	AssignRoleToUser(req *AssignRoleRequest) error
	RemoveRoleFromUser(req *AssignRoleRequest) error
	AssignPermissionToUser(req *AssignPermissionRequest) error
	RemovePermissionFromUser(req *AssignPermissionRequest) error
	GetUserPermissions(userID string) ([]string, error)
	HasPermission(userID string, permissionCode string) (bool, error)
}
