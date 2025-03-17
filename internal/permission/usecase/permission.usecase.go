package usecase

import (
	"errors"
	"time"

	"github.com/black4ninja/mi-proyecto/internal/permission/domain"
)

type permissionUseCase struct {
	permissionRepo domain.PermissionRepository
	userRoleRepo   domain.UserRoleRepository
}

// NewPermissionUseCase crea un nuevo caso de uso para permisos
func NewPermissionUseCase(permissionRepo domain.PermissionRepository, userRoleRepo domain.UserRoleRepository) domain.PermissionUseCase {
	return &permissionUseCase{
		permissionRepo: permissionRepo,
		userRoleRepo:   userRoleRepo,
	}
}

// GetPermission obtiene un permiso por su ID
func (u *permissionUseCase) GetPermission(id string) (*domain.PermissionResponse, error) {
	permission, err := u.permissionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &domain.PermissionResponse{
		ID:          permission.ID.Hex(),
		Code:        permission.Code,
		Module:      permission.Module,
		Action:      permission.Action,
		Name:        permission.Name,
		Description: permission.Description,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	}, nil
}

// GetPermissionByCode obtiene un permiso por su código
func (u *permissionUseCase) GetPermissionByCode(code string) (*domain.PermissionResponse, error) {
	permission, err := u.permissionRepo.GetByCode(code)
	if err != nil {
		return nil, err
	}

	return &domain.PermissionResponse{
		ID:          permission.ID.Hex(),
		Code:        permission.Code,
		Module:      permission.Module,
		Action:      permission.Action,
		Name:        permission.Name,
		Description: permission.Description,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	}, nil
}

// GetPermissionsByModule obtiene los permisos para un módulo específico
func (u *permissionUseCase) GetPermissionsByModule(module string) ([]*domain.PermissionResponse, error) {
	permissions, err := u.permissionRepo.GetByModule(module)
	if err != nil {
		return nil, err
	}

	var response []*domain.PermissionResponse
	for _, p := range permissions {
		response = append(response, &domain.PermissionResponse{
			ID:          p.ID.Hex(),
			Code:        p.Code,
			Module:      p.Module,
			Action:      p.Action,
			Name:        p.Name,
			Description: p.Description,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	return response, nil
}

// GetAllPermissions obtiene todos los permisos
func (u *permissionUseCase) GetAllPermissions() ([]*domain.PermissionResponse, error) {
	permissions, err := u.permissionRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var response []*domain.PermissionResponse
	for _, p := range permissions {
		response = append(response, &domain.PermissionResponse{
			ID:          p.ID.Hex(),
			Code:        p.Code,
			Module:      p.Module,
			Action:      p.Action,
			Name:        p.Name,
			Description: p.Description,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	return response, nil
}

// CreatePermission crea un nuevo permiso
func (u *permissionUseCase) CreatePermission(req *domain.CreatePermissionRequest) (*domain.PermissionResponse, error) {
	// Validar que el código sea único
	existingPermission, err := u.permissionRepo.GetByCode(req.Code)
	if err == nil && existingPermission != nil {
		return nil, errors.New("ya existe un permiso con este código")
	}

	// Crear permiso
	now := time.Now()
	permission := &domain.Permission{
		Code:        req.Code,
		Module:      req.Module,
		Action:      req.Action,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = u.permissionRepo.Create(permission)
	if err != nil {
		return nil, err
	}

	return &domain.PermissionResponse{
		ID:          permission.ID.Hex(),
		Code:        permission.Code,
		Module:      permission.Module,
		Action:      permission.Action,
		Name:        permission.Name,
		Description: permission.Description,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	}, nil
}

// UpdatePermission actualiza un permiso existente
func (u *permissionUseCase) UpdatePermission(id string, req *domain.UpdatePermissionRequest) (*domain.PermissionResponse, error) {
	// Obtener permiso existente
	permission, err := u.permissionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Actualizar campos
	if req.Name != "" {
		permission.Name = req.Name
	}

	if req.Description != "" {
		permission.Description = req.Description
	}

	permission.UpdatedAt = time.Now()

	// Guardar cambios
	err = u.permissionRepo.Update(permission)
	if err != nil {
		return nil, err
	}

	return &domain.PermissionResponse{
		ID:          permission.ID.Hex(),
		Code:        permission.Code,
		Module:      permission.Module,
		Action:      permission.Action,
		Name:        permission.Name,
		Description: permission.Description,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	}, nil
}

// DeletePermission elimina un permiso
func (u *permissionUseCase) DeletePermission(id string) error {
	return u.permissionRepo.Delete(id)
}

// HasPermission verifica si un usuario tiene un permiso específico
func (u *permissionUseCase) HasPermission(userID string, permissionCode string) (bool, error) {
	// Obtener todos los permisos del usuario
	permissions, err := u.userRoleRepo.GetUserPermissions(userID)
	if err != nil {
		return false, err
	}

	// Verificar si el permiso específico está en la lista
	for _, p := range permissions {
		if p == permissionCode {
			return true, nil
		}

		// Comprobar permisos comodín (por ejemplo, "module:*" o "module:submodule:*")
		if isWildcardMatch(p, permissionCode) {
			return true, nil
		}
	}

	return false, nil
}

// GetPermissionsByCodesArray obtiene permisos por array de códigos
func (u *permissionUseCase) GetPermissionsByCodesArray(codes []string) ([]*domain.PermissionResponse, error) {
	permissions, err := u.permissionRepo.GetByCodesArray(codes)
	if err != nil {
		return nil, err
	}

	var response []*domain.PermissionResponse
	for _, p := range permissions {
		response = append(response, &domain.PermissionResponse{
			ID:          p.ID.Hex(),
			Code:        p.Code,
			Module:      p.Module,
			Action:      p.Action,
			Name:        p.Name,
			Description: p.Description,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	return response, nil
}

// isWildcardMatch verifica si un permiso coincide con un comodín
// Por ejemplo, "module:*" coincidiría con "module:action"
func isWildcardMatch(pattern, permissionCode string) bool {
	// Si el patrón termina en *, es un comodín
	if len(pattern) > 2 && pattern[len(pattern)-1] == '*' {
		// Quitar el * del final
		prefix := pattern[:len(pattern)-1]

		// Si el permiso comienza con el prefijo, es una coincidencia
		return len(permissionCode) >= len(prefix) && permissionCode[:len(prefix)] == prefix
	}

	return false
}
