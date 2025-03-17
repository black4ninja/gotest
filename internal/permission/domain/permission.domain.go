package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Permission representa un permiso individual en el sistema
// Permission representa la entidad de permission
// @Description Entidad completa de permission
type Permission struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Code        string             `json:"code" bson:"code"`     // Formato: "module:submodule:action"
	Module      string             `json:"module" bson:"module"` // Ej: "finanzas", "inventario"
	Action      string             `json:"action" bson:"action"` // Ej: "read", "write", "reports"
	Name        string             `json:"name" bson:"name"`     // Nombre para mostrar
	Description string             `json:"description" bson:"description"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// PermissionRepository define el contrato para la capa de persistencia de permisos
type PermissionRepository interface {
	GetByID(id string) (*Permission, error)
	GetByCode(code string) (*Permission, error)
	GetByModule(module string) ([]*Permission, error)
	GetAll() ([]*Permission, error)
	Create(permission *Permission) error
	Update(permission *Permission) error
	Delete(id string) error
	GetByCodesArray(codes []string) ([]*Permission, error)
}

// CreatePermissionRequest representa la solicitud para crear un permiso
// CreatePermissionRequest representa la solicitud para crear un permission
// @Description Datos necesarios para crear un permission
type CreatePermissionRequest struct {
	Code        string `json:"code" binding:"required"`
	Module      string `json:"module" binding:"required"`
	Action      string `json:"action" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdatePermissionRequest representa la solicitud para actualizar un permiso
// UpdatePermissionRequest representa la solicitud para actualizar un permission
// @Description Datos para actualizar un permission
type UpdatePermissionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// PermissionResponse representa la respuesta con datos de permisos
// PermissionResponse representa la respuesta con datos de permission
// @Description Estructura de respuesta para información de permission
type PermissionResponse struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Module      string    `json:"module"`
	Action      string    `json:"action"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PermissionUseCase define el contrato para la capa de caso de uso de permisos
type PermissionUseCase interface {
	GetPermission(id string) (*PermissionResponse, error)
	GetPermissionByCode(code string) (*PermissionResponse, error)
	GetPermissionsByModule(module string) ([]*PermissionResponse, error)
	GetAllPermissions() ([]*PermissionResponse, error)
	CreatePermission(req *CreatePermissionRequest) (*PermissionResponse, error)
	UpdatePermission(id string, req *UpdatePermissionRequest) (*PermissionResponse, error)
	DeletePermission(id string) error
	HasPermission(userID string, permissionCode string) (bool, error)
	GetPermissionsByCodesArray(codes []string) ([]*PermissionResponse, error)
}
