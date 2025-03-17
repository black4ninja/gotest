package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Constantes para el estado del usuario
const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
	UserStatusArchived = "archived"
)

// User representa la entidad de usuario
// @Description Entidad completa de usuario
type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty" example:"60f1e5e5e5e5e5e5e5e5e5e5"`  // ID único del usuario
	Email        string             `json:"email" bson:"email" example:"usuario@example.com"`            // Email del usuario
	Name         string             `json:"name" bson:"name" example:"Juan Pérez"`                       // Nombre completo del usuario
	Password     string             `json:"-" bson:"password"`                                           // Contraseña hasheada (no incluida en JSON)
	Status       string             `json:"status" bson:"status" example:"active"`                       // Estado: active, inactive, archived
	Role         string             `json:"role" bson:"role" example:"user"`                             // Rol del usuario
	RefreshToken string             `json:"-" bson:"refresh_token,omitempty"`                            // Token de refresco (no incluido en JSON)
	CreatedAt    time.Time          `json:"created_at" bson:"created_at" example:"2023-07-10T15:04:05Z"` // Fecha de creación
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at" example:"2023-07-10T15:04:05Z"` // Fecha de última actualización
	ArchivedAt   *time.Time         `json:"archived_at,omitempty" bson:"archived_at,omitempty"`          // Fecha de archivado (si aplica)
}

// CreateUserRequest representa la solicitud para crear un usuario
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
}

// UpdateUserRequest representa la solicitud para actualizar un usuario
type UpdateUserRequest struct {
	Name   string `json:"name"`
	Email  string `json:"email" binding:"omitempty,email"`
	Status string `json:"status"`
	Role   string `json:"role"`
}

// ChangePasswordRequest representa la solicitud para cambiar contraseña
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// UserResponse representa la respuesta con datos de usuario
// @Description Estructura de respuesta para información de usuario
type UserResponse struct {
	ID        string    `json:"id" example:"60f1e5e5e5e5e5e5e5e5e5e5"`     // ID único del usuario
	Email     string    `json:"email" example:"usuario@example.com"`       // Email del usuario
	Name      string    `json:"name" example:"Juan Pérez"`                 // Nombre completo del usuario
	Status    string    `json:"status" example:"active"`                   // Estado: active, inactive, archived
	Role      string    `json:"role" example:"user"`                       // Rol del usuario
	CreatedAt time.Time `json:"created_at" example:"2023-07-10T15:04:05Z"` // Fecha de creación
	UpdatedAt time.Time `json:"updated_at" example:"2023-07-10T15:04:05Z"` // Fecha de última actualización
}

// UserRepository define el contrato para la capa de persistencia
type UserRepository interface {
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetAll(params map[string]interface{}) ([]*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id string) error
	Archive(id string) error
	UpdateRefreshToken(userID string, refreshToken string) error
	GetByRefreshToken(refreshToken string) (*User, error)
}

// UserUseCase define el contrato para la capa de casos de uso
type UserUseCase interface {
	GetUser(id string) (*UserResponse, error)
	GetUserByEmail(email string) (*User, error)
	GetAllUsers(params map[string]interface{}) ([]*UserResponse, error)
	CreateUser(req *CreateUserRequest) (*UserResponse, error)
	UpdateUser(id string, req *UpdateUserRequest) (*UserResponse, error)
	DeleteUser(id string) error
	ArchiveUser(id string) error
	ChangePassword(userID string, req *ChangePasswordRequest) error
	ValidateCredentials(email string, password string) (*User, error)
	UpdateRefreshToken(userID string, refreshToken string) error
	GetUserByRefreshToken(refreshToken string) (*User, error)
}
