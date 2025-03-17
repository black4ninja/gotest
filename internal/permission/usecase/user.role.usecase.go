package usecase

import (
	"errors"

	"github.com/black4ninja/mi-proyecto/internal/permission/domain"
)

type userRoleUseCase struct {
	userRoleRepo   domain.UserRoleRepository
	roleRepo       domain.RoleRepository
	permissionRepo domain.PermissionRepository
}

// NewUserRoleUseCase crea un nuevo caso de uso para asignaciones usuario-rol
func NewUserRoleUseCase(
	userRoleRepo domain.UserRoleRepository,
	roleRepo domain.RoleRepository,
	permissionRepo domain.PermissionRepository,
) domain.UserRoleUseCase {
	return &userRoleUseCase{
		userRoleRepo:   userRoleRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

// GetUserRoles obtiene los roles y permisos asignados a un usuario
func (u *userRoleUseCase) GetUserRoles(userID string) (*domain.UserRoleResponse, error) {
	// Obtener asignación de usuario
	userRole, err := u.userRoleRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Obtener roles
	var roles []*domain.RoleResponse
	for _, roleID := range userRole.Roles {
		role, err := u.roleRepo.GetByID(roleID)
		if err != nil {
			continue // Ignorar roles que no existan
		}

		// Obtener permisos del rol
		permissions, err := u.permissionRepo.GetByCodesArray(role.Permissions)
		if err != nil {
			// Si hay error, añadir el rol sin permisos
			roles = append(roles, &domain.RoleResponse{
				ID:          role.ID.Hex(),
				Name:        role.Name,
				Description: role.Description,
				Permissions: []*domain.PermissionResponse{},
				IsSystem:    role.IsSystem,
				CreatedAt:   role.CreatedAt,
				UpdatedAt:   role.UpdatedAt,
			})
			continue
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

		roles = append(roles, &domain.RoleResponse{
			ID:          role.ID.Hex(),
			Name:        role.Name,
			Description: role.Description,
			Permissions: permissionsResponse,
			IsSystem:    role.IsSystem,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})
	}

	// Obtener permisos específicos del usuario
	permissions, err := u.permissionRepo.GetByCodesArray(userRole.Permissions)
	if err != nil {
		// Si hay error, devolver sin permisos específicos
		return &domain.UserRoleResponse{
			ID:          userRole.ID.Hex(),
			UserID:      userRole.UserID,
			Roles:       roles,
			Permissions: []*domain.PermissionResponse{},
			CreatedAt:   userRole.CreatedAt,
			UpdatedAt:   userRole.UpdatedAt,
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

	return &domain.UserRoleResponse{
		ID:          userRole.ID.Hex(),
		UserID:      userRole.UserID,
		Roles:       roles,
		Permissions: permissionsResponse,
		CreatedAt:   userRole.CreatedAt,
		UpdatedAt:   userRole.UpdatedAt,
	}, nil
}

// AssignRoleToUser asigna un rol a un usuario
func (u *userRoleUseCase) AssignRoleToUser(req *domain.AssignRoleRequest) error {
	// Verificar que el rol exista
	_, err := u.roleRepo.GetByID(req.RoleID)
	if err != nil {
		return errors.New("rol no válido")
	}

	return u.userRoleRepo.AddRole(req.UserID, req.RoleID)
}

// RemoveRoleFromUser elimina un rol de un usuario
func (u *userRoleUseCase) RemoveRoleFromUser(req *domain.AssignRoleRequest) error {
	return u.userRoleRepo.RemoveRole(req.UserID, req.RoleID)
}

// AssignPermissionToUser asigna un permiso específico a un usuario
func (u *userRoleUseCase) AssignPermissionToUser(req *domain.AssignPermissionRequest) error {
	// Verificar que el permiso exista
	_, err := u.permissionRepo.GetByCode(req.PermissionCode)
	if err != nil {
		return errors.New("permiso no válido")
	}

	return u.userRoleRepo.AddPermission(req.UserID, req.PermissionCode)
}

// RemovePermissionFromUser elimina un permiso específico de un usuario
func (u *userRoleUseCase) RemovePermissionFromUser(req *domain.AssignPermissionRequest) error {
	return u.userRoleRepo.RemovePermission(req.UserID, req.PermissionCode)
}

// GetUserPermissions obtiene todos los permisos de un usuario
func (u *userRoleUseCase) GetUserPermissions(userID string) ([]string, error) {
	return u.userRoleRepo.GetUserPermissions(userID)
}

// HasPermission verifica si un usuario tiene un permiso específico
func (u *userRoleUseCase) HasPermission(userID string, permissionCode string) (bool, error) {
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
