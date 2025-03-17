package usecase

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/black4ninja/mi-proyecto/internal/user/domain"
)

type userUseCase struct {
	userRepo domain.UserRepository
}

// NewUserUseCase crea un nuevo caso de uso para usuarios
func NewUserUseCase(userRepo domain.UserRepository) domain.UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

// GetUser obtiene un usuario por su ID
func (u *userUseCase) GetUser(id string) (*domain.UserResponse, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &domain.UserResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		Name:      user.Name,
		Status:    user.Status,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetUserByEmail obtiene un usuario por su email
func (u *userUseCase) GetUserByEmail(email string) (*domain.User, error) {
	return u.userRepo.GetByEmail(email)
}

// GetAllUsers obtiene todos los usuarios
func (u *userUseCase) GetAllUsers(params map[string]interface{}) ([]*domain.UserResponse, error) {
	users, err := u.userRepo.GetAll(params)
	if err != nil {
		return nil, err
	}

	var response []*domain.UserResponse
	for _, user := range users {
		response = append(response, &domain.UserResponse{
			ID:        user.ID.Hex(),
			Email:     user.Email,
			Name:      user.Name,
			Status:    user.Status,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return response, nil
}

// CreateUser crea un nuevo usuario
func (u *userUseCase) CreateUser(req *domain.CreateUserRequest) (*domain.UserResponse, error) {
	// Verificar si el email ya existe
	existingUser, err := u.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("el email ya está registrado")
	}

	// Hashear contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Establecer rol por defecto si no se proporciona
	role := req.Role
	if role == "" {
		role = "user"
	}

	// Crear usuario
	now := time.Now()
	user := &domain.User{
		Email:     req.Email,
		Name:      req.Name,
		Password:  string(hashedPassword),
		Status:    domain.UserStatusActive,
		Role:      role,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &domain.UserResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		Name:      user.Name,
		Status:    user.Status,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// UpdateUser actualiza un usuario existente
func (u *userUseCase) UpdateUser(id string, req *domain.UpdateUserRequest) (*domain.UserResponse, error) {
	// Obtener usuario existente
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Verificar si se intenta cambiar el email y si ya existe
	if req.Email != "" && req.Email != user.Email {
		existingUser, err := u.userRepo.GetByEmail(req.Email)
		if err == nil && existingUser != nil {
			return nil, errors.New("el email ya está registrado")
		}
		user.Email = req.Email
	}

	// Actualizar campos
	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Status != "" {
		user.Status = req.Status
	}

	if req.Role != "" {
		user.Role = req.Role
	}

	user.UpdatedAt = time.Now()

	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &domain.UserResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		Name:      user.Name,
		Status:    user.Status,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// DeleteUser elimina un usuario
func (u *userUseCase) DeleteUser(id string) error {
	return u.userRepo.Delete(id)
}

// ArchiveUser archiva un usuario
func (u *userUseCase) ArchiveUser(id string) error {
	return u.userRepo.Archive(id)
}

// ChangePassword cambia la contraseña de un usuario
func (u *userUseCase) ChangePassword(userID string, req *domain.ChangePasswordRequest) error {
	// Obtener usuario
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Verificar contraseña antigua
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("contraseña antigua incorrecta")
	}

	// Hashear nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Actualizar contraseña
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	return u.userRepo.Update(user)
}

// ValidateCredentials valida las credenciales de un usuario
func (u *userUseCase) ValidateCredentials(email string, password string) (*domain.User, error) {
	// Buscar usuario
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Verificar si el usuario está activo
	if user.Status != domain.UserStatusActive {
		return nil, errors.New("usuario inactivo")
	}

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	return user, nil
}

// UpdateRefreshToken actualiza el token de refresco de un usuario
func (u *userUseCase) UpdateRefreshToken(userID string, refreshToken string) error {
	return u.userRepo.UpdateRefreshToken(userID, refreshToken)
}

// GetUserByRefreshToken obtiene un usuario por su token de refresco
func (u *userUseCase) GetUserByRefreshToken(refreshToken string) (*domain.User, error) {
	return u.userRepo.GetByRefreshToken(refreshToken)
}
