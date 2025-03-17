// cmd/scripts/init_permissions_and_admin.go
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// Dominios
	permDomain "github.com/black4ninja/mi-proyecto/internal/permission/domain"
	userDomain "github.com/black4ninja/mi-proyecto/internal/user/domain"

	// Repositorios
	permRepo "github.com/black4ninja/mi-proyecto/internal/permission/repository"
	userRepo "github.com/black4ninja/mi-proyecto/internal/user/repository"

	// Casos de uso
	permUseCase "github.com/black4ninja/mi-proyecto/internal/permission/usecase"
	userUseCase "github.com/black4ninja/mi-proyecto/internal/user/usecase"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("Archivo .env no encontrado, usando variables de entorno del sistema")
	}

	// Configurar MongoDB
	mongoURI := getEnvPerms("MONGO_URI", "mongodb://localhost:27017")
	mongoDBName := getEnvPerms("MONGO_DB", "my_database")
	mongoTimeout := 10 * time.Second

	// Conectar a MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), mongoTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error al conectar a MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Inicializar repositorios
	permissionCollection := client.Database(mongoDBName).Collection("permissions")
	roleCollection := client.Database(mongoDBName).Collection("roles")
	userCollection := client.Database(mongoDBName).Collection("users")
	userRoleCollection := client.Database(mongoDBName).Collection("user_roles")

	permissionRepository := permRepo.NewMongoPermissionRepository(permissionCollection)
	roleRepository := permRepo.NewMongoRoleRepository(roleCollection)
	userRepository := userRepo.NewMongoUserRepository(userCollection)
	userRoleRepository := permRepo.NewMongoUserRoleRepository(userRoleCollection, roleRepository)

	// Inicializar casos de uso
	permissionService := permUseCase.NewPermissionUseCase(permissionRepository, userRoleRepository)
	roleService := permUseCase.NewRoleUseCase(roleRepository, permissionRepository)
	userService := userUseCase.NewUserUseCase(userRepository)
	userRoleService := permUseCase.NewUserRoleUseCase(userRoleRepository, roleRepository, permissionRepository)

	// Inicializar permisos y roles
	log.Println("Iniciando creación de permisos y roles predeterminados...")
	initializeDefaultPermissionsAndRoles(permissionService, roleService)
	log.Println("Permisos y roles predeterminados creados correctamente")

	// Crear usuario administrador
	log.Println("Iniciando creación de usuario administrador predeterminado...")
	createDefaultAdminUser(userService, roleService, userRoleService)
	log.Println("Usuario administrador predeterminado creado correctamente")
}

// initializeDefaultPermissionsAndRoles crea permisos y roles predeterminados
func initializeDefaultPermissionsAndRoles(
	permissionService permDomain.PermissionUseCase,
	roleService permDomain.RoleUseCase,
) {
	// Crear permisos administrativos
	createDefaultPermission(permissionService, "admin:permissions", "admin", "permissions", "Administrar permisos", "Permite administrar permisos y roles")
	createDefaultPermission(permissionService, "admin:users", "admin", "users", "Administrar usuarios", "Permite administrar usuarios")
	createDefaultPermission(permissionService, "admin:dashboard", "admin", "dashboard", "Dashboard administrativo", "Acceso al dashboard administrativo")
	createDefaultPermission(permissionService, "admin:data:import", "admin", "data:import", "Importar datos", "Permite importar datos")
	createDefaultPermission(permissionService, "admin:data:modify", "admin", "data:modify", "Modificar datos", "Permite modificar datos del sistema")

	// Crear permisos de módulo financiero
	createDefaultPermission(permissionService, "finanzas:read", "finanzas", "read", "Ver finanzas", "Acceso de lectura al módulo financiero")
	createDefaultPermission(permissionService, "finanzas:write", "finanzas", "write", "Editar finanzas", "Permite crear y editar datos financieros")
	createDefaultPermission(permissionService, "finanzas:reports:read", "finanzas", "reports:read", "Ver reportes financieros", "Acceso a reportes financieros")
	createDefaultPermission(permissionService, "finanzas:reports:export", "finanzas", "reports:export", "Exportar reportes", "Permite exportar reportes financieros")
	createDefaultPermission(permissionService, "finanzas:transactions:write", "finanzas", "transactions:write", "Crear transacciones", "Permite crear transacciones financieras")
	createDefaultPermission(permissionService, "finanzas:dashboard", "finanzas", "dashboard", "Dashboard financiero", "Acceso al dashboard financiero")

	// Crear permisos de módulo inventario
	createDefaultPermission(permissionService, "inventario:read", "inventario", "read", "Ver inventario", "Acceso de lectura al módulo de inventario")
	createDefaultPermission(permissionService, "inventario:write", "inventario", "write", "Editar inventario", "Permite crear y editar elementos del inventario")
	createDefaultPermission(permissionService, "inventario:reports", "inventario", "reports", "Reportes de inventario", "Acceso a reportes de inventario")
	createDefaultPermission(permissionService, "inventario:dashboard", "inventario", "dashboard", "Dashboard inventario", "Acceso al dashboard de inventario")

	// Crear roles predeterminados

	// Rol de administrador
	adminPerms := []string{
		"admin:permissions",
		"admin:users",
		"admin:dashboard",
		"admin:data:import",
		"admin:data:modify",
		"finanzas:read",
		"finanzas:write",
		"finanzas:reports:read",
		"finanzas:reports:export",
		"finanzas:transactions:write",
		"finanzas:dashboard",
		"inventario:read",
		"inventario:write",
		"inventario:reports",
		"inventario:dashboard",
	}

	createDefaultRole(roleService, "Administrador", "Acceso completo al sistema", adminPerms)

	// Rol de gerente financiero
	finanzasPerms := []string{
		"finanzas:read",
		"finanzas:write",
		"finanzas:reports:read",
		"finanzas:reports:export",
		"finanzas:transactions:write",
		"finanzas:dashboard",
	}

	createDefaultRole(roleService, "Gerente Financiero", "Gestión del módulo financiero", finanzasPerms)

	// Rol de analista financiero
	analistaPerms := []string{
		"finanzas:read",
		"finanzas:reports:read",
		"finanzas:dashboard",
	}

	createDefaultRole(roleService, "Analista Financiero", "Visualización de datos financieros", analistaPerms)

	// Rol de gerente de inventario
	inventarioPerms := []string{
		"inventario:read",
		"inventario:write",
		"inventario:reports",
		"inventario:dashboard",
	}

	createDefaultRole(roleService, "Gerente de Inventario", "Gestión del inventario", inventarioPerms)
}

// createDefaultPermission crea un permiso si no existe
func createDefaultPermission(
	permissionService permDomain.PermissionUseCase,
	code, module, action, name, description string,
) {
	// Verificar si ya existe
	_, err := permissionService.GetPermissionByCode(code)
	if err == nil {
		return // Ya existe, no hacer nada
	}

	// Crear permiso
	permissionService.CreatePermission(&permDomain.CreatePermissionRequest{
		Code:        code,
		Module:      module,
		Action:      action,
		Name:        name,
		Description: description,
	})
}

// createDefaultRole crea un rol si no existe
func createDefaultRole(
	roleService permDomain.RoleUseCase,
	name, description string,
	permissions []string,
) {
	// Verificar si ya existe
	_, err := roleService.GetRoleByName(name)
	if err == nil {
		return // Ya existe, no hacer nada
	}

	// Crear rol
	roleService.CreateRole(&permDomain.CreateRoleRequest{
		Name:        name,
		Description: description,
		Permissions: permissions,
	})
}

// createDefaultAdminUser crea un usuario administrador si no existe
func createDefaultAdminUser(
	userService userDomain.UserUseCase,
	roleService permDomain.RoleUseCase,
	userRoleService permDomain.UserRoleUseCase,
) {
	// Configuración del usuario admin predeterminado
	adminEmail := getEnv("DEFAULT_ADMIN_EMAIL", "admin@sistema.com")
	adminPassword := getEnv("DEFAULT_ADMIN_PASSWORD", "AdminPass123!")

	// Verificar si ya existe
	existingUser, err := userService.GetUserByEmail(adminEmail)
	if err == nil {
		log.Println("El usuario admin ya existe, verificando permisos...")

		// Si existe, verificar si tiene el rol admin y asignarlo si no lo tiene
		userRoles, err := userRoleService.GetUserRoles(existingUser.ID)
		if err == nil && len(userRoles.Roles) > 0 {
			log.Println("El usuario admin ya tiene roles asignados")
			return
		}

		// Buscar el rol de Administrador
		adminRole, err := roleService.GetRoleByName("Administrador")
		if err != nil {
			log.Printf("Error al buscar rol de administrador: %v", err)
			return
		}

		// Asignar el rol al usuario existente
		err = userRoleService.AssignRoleToUser(&permDomain.AssignRoleRequest{
			UserID: existingUser.ID,
			RoleID: adminRole.ID,
		})
		if err != nil {
			log.Printf("Error al asignar rol admin: %v", err)
		} else {
			log.Println("Rol de administrador asignado correctamente al usuario existente")
		}

		return
	}

	// Crear el usuario admin
	adminUser, err := userService.CreateUser(&userDomain.CreateUserRequest{
		Name:     "Administrador del Sistema",
		Email:    adminEmail,
		Password: adminPassword,
		Role:     "admin",
	})
	if err != nil {
		log.Printf("Error al crear usuario admin: %v", err)
		return
	}

	// Buscar el rol de Administrador
	adminRole, err := roleService.GetRoleByName("Administrador")
	if err != nil {
		log.Printf("Error al buscar rol de administrador: %v", err)
		return
	}

	// Asignar el rol al usuario
	err = userRoleService.AssignRoleToUser(&permDomain.AssignRoleRequest{
		UserID: adminUser.ID,
		RoleID: adminRole.ID,
	})
	if err != nil {
		log.Printf("Error al asignar rol admin: %v", err)
		return
	}

	log.Printf("Usuario administrador creado con éxito: %s", adminEmail)
	log.Printf("Contraseña: %s (cámbiala después de iniciar sesión)", adminPassword)
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnvPerms(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}
