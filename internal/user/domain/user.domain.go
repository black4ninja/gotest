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
type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email        string             `json:"email" bson:"email"`
	Name         string             `json:"name" bson:"name"`
	Password     string             `json:"-" bson:"password"`
	Status       string             `json:"status" bson:"status"`
	Role         string             `json:"role" bson:"role"`
	RefreshToken string             `json:"-" bson:"refresh_token,omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	ArchivedAt   *time.Time         `json:"archived_at,omitempty" bson:"archived_at,omitempty"`
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
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
