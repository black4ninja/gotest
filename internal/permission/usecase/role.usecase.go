package usecase

import (
	"errors"
	"time"

	"github.com/black4ninja/mi-proyecto/internal/permission/domain"
)

type roleUseCase struct {
	roleRepo       domain.RoleRepository
	permissionRepo domain.PermissionRepository
}

// NewRoleUseCase crea un nuevo caso de uso para roles
func NewRoleUseCase(roleRepo domain.RoleRepository, permissionRepo domain.PermissionRepository) domain.RoleUseCase {
	return &roleUseCase{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

// GetRole obtiene un rol por su ID
func (u *roleUseCase) GetRole(id string) (*domain.RoleResponse, error) {
	role, err := u.roleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Obtener los permisos asociados
	permissions, err := u.permissionRepo.GetByCodesArray(role.Permissions)
	if err != nil {
		return nil, err
	}

	// Convertir permisos al formato de respuesta
	var permissionsResponse []*domain.PermissionResponse
	for _, p := range permissions {
		permissionsResponse = append(permissionsResponse, &domain.PermissionResponse{
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

	return &domain.RoleResponse{
		ID:          role.ID.Hex(),
		Name:        role.Name,
		Description: role.Description,
		Permissions: permissionsResponse,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}, nil
}

// GetRoleByName obtiene un rol por su nombre
func (u *roleUseCase) GetRoleByName(name string) (*domain.RoleResponse, error) {
	role, err := u.roleRepo.GetByName(name)
	if err != nil {
		return nil, err
	}

	// Obtener los permisos asociados
	permissions, err := u.permissionRepo.GetByCodesArray(role.Permissions)
	if err != nil {
		return nil, err
	}

	// Convertir permisos al formato de respuesta
	var permissionsResponse []*domain.PermissionResponse
	for _, p := range permissions {
		permissionsResponse = append(permissionsResponse, &domain.PermissionResponse{
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

	return &domain.RoleResponse{
		ID:          role.ID.Hex(),
		Name:        role.Name,
		Description: role.Description,
		Permissions: permissionsResponse,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}, nil
}

// GetAllRoles obtiene todos los roles
func (u *roleUseCase) GetAllRoles() ([]*domain.RoleResponse, error) {
	roles, err := u.roleRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var response []*domain.RoleResponse

	// Para cada rol, obtener sus permisos
	for _, role := range roles {
		// Obtener los permisos asociados
		permissions, err := u.permissionRepo.GetByCodesArray(role.Permissions)
		if err != nil {
			continue // Ignorar errores y seguir con el siguiente rol
		}

		// Convertir permisos al formato de respuesta
		var permissionsResponse []*domain.PermissionResponse
		for _, p := range permissions {
			permissionsResponse = append(permissionsResponse, &domain.PermissionResponse{
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

		response = append(response, &domain.RoleResponse{
			ID:          role.ID.Hex(),
			Name:        role.Name,
			Description: role.Description,
			Permissions: permissionsResponse,
			IsSystem:    role.IsSystem,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})
	}

	return response, nil
}

// CreateRole crea un nuevo rol
func (u *roleUseCase) CreateRole(req *domain.CreateRoleRequest) (*domain.RoleResponse, error) {
	// Verificar que no exista un rol con el mismo nombre
	existingRole, err := u.roleRepo.GetByName(req.Name)
	if err == nil && existingRole != nil {
		return nil, errors.New("ya existe un rol con este nombre")
	}

	// Verificar que los permisos existan
	for _, pCode := range req.Permissions {
		_, err := u.permissionRepo.GetByCode(pCode)
		if err != nil {
			return nil, errors.New("permiso no válido: " + pCode)
		}
	}

	// Crear rol
	now := time.Now()
	role := &domain.Role{
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permissions,
		IsSystem:    false, // No es un rol de sistema
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = u.roleRepo.Create(role)
	if err != nil {
		return nil, err
	}

	// Obtener los permisos para la respuesta
	permissions, err := u.permissionRepo.GetByCodesArray(role.Permissions)
	if err != nil {
		// Si hay error al obtener permisos, devolvemos el rol sin permisos
		return &domain.RoleResponse{
			ID:          role.ID.Hex(),
			Name:        role.Name,
			Description: role.Description,
			Permissions: []*domain.PermissionResponse{},
			IsSystem:    role.IsSystem,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		}, nil
	}

	// Convertir permisos al formato de respuesta
	var permissionsResponse []*domain.PermissionResponse
	for _, p := range permissions {
		permissionsResponse = append(permissionsResponse, &domain.PermissionResponse{
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

	return &domain.RoleResponse{
		ID:          role.ID.Hex(),
		Name:        role.Name,
		Description: role.Description,
		Permissions: permissionsResponse,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}, nil
}

// UpdateRole actualiza un rol existente
func (u *roleUseCase) UpdateRole(id string, req *domain.UpdateRoleRequest) (*domain.RoleResponse, error) {
	// Obtener rol existente
	role, err := u.roleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Verificar que no sea un rol de sistema
	if role.IsSystem {
		return nil, errors.New("no se puede modificar un rol de sistema")
	}

	// Actualizar campos
	if req.Name != "" && req.Name != role.Name {
		// Verificar que no exista otro rol con el nuevo nombre
		existingRole, err := u.roleRepo.GetByName(req.Name)
		if err == nil && existingRole != nil && existingRole.ID.Hex() != id {
			return nil, errors.New("ya existe un rol con este nombre")
		}

		role.Name = req.Name
	}

	if req.Description != "" {
		role.Description = req.Description
	}

	role.UpdatedAt = time.Now()

	// Guardar cambios
	err = u.roleRepo.Update(role)
	if err != nil {
		return nil, err
	}

	// Obtener los permisos para la respuesta
	permissions, err := u.permissionRepo.GetByCodesArray(role.Permissions)
	if err != nil {
		// Si hay error al obtener permisos, devolvemos el rol sin permisos
		return &domain.RoleResponse{
			ID:          role.ID.Hex(),
			Name:        role.Name,
			Description: role.Description,
			Permissions: []*domain.PermissionResponse{},
			IsSystem:    role.IsSystem,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		}, nil
	}

	// Convertir permisos al formato de respuesta
	var permissionsResponse []*domain.PermissionResponse
	for _, p := range permissions {
		permissionsResponse = append(permissionsResponse, &domain.PermissionResponse{
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

	return &domain.RoleResponse{
		ID:          role.ID.Hex(),
		Name:        role.Name,
		Description: role.Description,
		Permissions: permissionsResponse,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}, nil
}

// DeleteRole elimina un rol
func (u *roleUseCase) DeleteRole(id string) error {
	return u.roleRepo.Delete(id)
}

// AddPermissionToRole añade un permiso a un rol
func (u *roleUseCase) AddPermissionToRole(roleID string, permissionCode string) error {
	// Verificar que el permiso exista
	_, err := u.permissionRepo.GetByCode(permissionCode)
	if err != nil {
		return errors.New("permiso no válido: " + permissionCode)
	}

	return u.roleRepo.AddPermission(roleID, permissionCode)
}

// RemovePermissionFromRole elimina un permiso de un rol
func (u *roleUseCase) RemovePermissionFromRole(roleID string, permissionCode string) error {
	return u.roleRepo.RemovePermission(roleID, permissionCode)
}
