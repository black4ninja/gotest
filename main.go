package main

import (
	"context"
	"github.com/black4ninja/mi-proyecto/internal/user/domain"
	"github.com/black4ninja/mi-proyecto/pkg/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"

	oauthDelivery "github.com/black4ninja/mi-proyecto/internal/oauth/delivery"
	oauthRepo "github.com/black4ninja/mi-proyecto/internal/oauth/repository"
	oauthUseCase "github.com/black4ninja/mi-proyecto/internal/oauth/usecase"
	userDelivery "github.com/black4ninja/mi-proyecto/internal/user/delivery"
	userRepo "github.com/black4ninja/mi-proyecto/internal/user/repository"
	userUseCase "github.com/black4ninja/mi-proyecto/internal/user/usecase"
	"github.com/black4ninja/mi-proyecto/pkg/config"
	"github.com/black4ninja/mi-proyecto/pkg/middleware"

	permissionDelivery "github.com/black4ninja/mi-proyecto/internal/permission/delivery"
	permissionRepo "github.com/black4ninja/mi-proyecto/internal/permission/repository"
	permissionUseCase "github.com/black4ninja/mi-proyecto/internal/permission/usecase"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("Archivo .env no encontrado, usando variables de entorno del sistema")
	}

	// Configurar MongoDB
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	mongoDBName := getEnv("MONGO_DB", "my_database")
	mongoTimeout := 10 * time.Second

	// Conectar a MongoDB
	mongoConfig := config.MongoConfig{
		URI:      mongoURI,
		Database: mongoDBName,
		Timeout:  mongoTimeout,
	}

	mongoClient, err := config.NewMongoClient(mongoConfig)
	if err != nil {
		log.Fatalf("Error al conectar a MongoDB: %v", err)
	}

	// Configurar cierre de conexión al terminar
	setupGracefulShutdown(mongoClient)

	// Configurar modo de Gin basado en entorno
	if getEnv("ENV", "development") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// ------ COLECCIONES DE MONGODB ------
	// Colecciones existentes
	userCollection := mongoClient.Database(mongoDBName).Collection("users")
	permissionCollection := mongoClient.Database(mongoDBName).Collection("permissions")
	roleCollection := mongoClient.Database(mongoDBName).Collection("roles")
	userRoleCollection := mongoClient.Database(mongoDBName).Collection("user_roles")

	// ------ INICIALIZACIÓN DE REPOSITORIOS ------
	// Repositorios de usuario
	userRepository := userRepo.NewMongoUserRepository(userCollection)
	permissionRepository := permissionRepo.NewMongoPermissionRepository(permissionCollection)
	roleRepository := permissionRepo.NewMongoRoleRepository(roleCollection)
	userRoleRepository := permissionRepo.NewMongoUserRoleRepository(userRoleCollection, roleRepository)

	// Repositorios de OAuth
	clientCollection := config.GetCollection(mongoClient, mongoDBName, "oauth_clients")
	tokenCollection := config.GetCollection(mongoClient, mongoDBName, "oauth_tokens")
	clientRepository := oauthRepo.NewMongoClientRepository(clientCollection)
	tokenRepository := oauthRepo.NewMongoTokenRepository(tokenCollection)

	// ------ INICIALIZACIÓN DE CASOS DE USO ------
	// Caso de uso de usuario
	userService := userUseCase.NewUserUseCase(userRepository)
	permissionService := permissionUseCase.NewPermissionUseCase(permissionRepository, userRoleRepository)
	roleService := permissionUseCase.NewRoleUseCase(roleRepository, permissionRepository)
	userRoleService := permissionUseCase.NewUserRoleUseCase(userRoleRepository, roleRepository, permissionRepository)

	// Configuración de OAuth
	jwtSecret := getEnv("JWT_SECRET", "mi_secret_super_seguro")
	// Determinar el tiempo de expiración según el entorno
	var tokenExpiration time.Duration
	if getEnv("ENV", "development") == "production" {
		// En producción: 30 días
		tokenExpiration = 30 * 24 * time.Hour
	} else {
		// En desarrollo: 15 minutos
		tokenExpiration = 15 * time.Minute
	}

	// Configurar refresh token (opcional - se puede hacer también dependiente del entorno)
	var refreshExpiration time.Duration
	if getEnv("ENV", "development") == "production" {
		// En producción: 90 días (habitualmente más largo que el access token)
		refreshExpiration = 90 * 24 * time.Hour
	} else {
		// En desarrollo: 1 hora
		refreshExpiration = 1 * time.Hour
	}

	// Caso de uso de OAuth
	oauthService := oauthUseCase.NewOAuthUseCase(
		clientRepository,
		tokenRepository,
		userService,
		jwtSecret,
		tokenExpiration,
		refreshExpiration,
	)

	// ------ INICIALIZACIÓN DE MIDDLEWARES ------
	// Middleware de OAuth
	oauthMiddleware := middleware.NewOAuthMiddleware(oauthService)
	permissionMiddleware := middleware.NewPermissionMiddleware(userRoleService)

	// ------ CONFIGURACIÓN DE RUTAS ------
	// Inicializar router de Gin
	router := gin.Default()
	// Rutas base
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API Monolítica con estructura de Microservicios y OAuth 2.0",
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"mongo":  "connected",
		})
	})

	router.POST("/api/register", func(c *gin.Context) {
		var req domain.CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ValidationErrorResponse(c, err.Error())
			return
		}

		user, err := userService.CreateUser(&req)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.SuccessResponse(c, http.StatusCreated, "Usuario creado con éxito", user)
	})

	publicRoutes := router.Group("/api")
	{
		// Rutas de OAuth (públicas)
		oauthRoutes := publicRoutes.Group("/oauth")
		oauthDelivery.NewOAuthHandler(oauthRoutes, oauthService)
	}

	// Grupo de rutas para la API
	api := router.Group("/api")
	api.Use(oauthMiddleware.Protected()) // Protección aplicada solo a este grupo
	{
		// Rutas de usuarios
		userRoutes := api.Group("/users")
		userDelivery.NewUserHandler(userRoutes, userService)

		// Rutas de permisos
		permissionRoutes := api.Group("/permissions")
		permissionRoutes.Use(permissionMiddleware.RequirePermission("admin:permissions"))
		permissionDelivery.NewPermissionHandler(permissionRoutes, permissionService, roleService, userRoleService)
	}

	// ------ EJEMPLOS DE USO DEL MIDDLEWARE DE PERMISOS ------

	// Ejemplo: Proteger acceso a módulo de finanzas
	//finanzasRoutes := api.Group("/finanzas")
	//finanzasRoutes.Use(permissionMiddleware.RequireModuleAccess("finanzas"))

	// Rutas específicas con permisos más granulares
	//finanzasRoutes.GET("/reportes", permissionMiddleware.RequirePermission("finanzas:reports:read"), handleFinanzasReportes)
	//finanzasRoutes.POST("/transacciones", permissionMiddleware.RequirePermission("finanzas:transactions:write"), handleFinanzasTransacciones)

	// Ejemplo: Proteger acceso a módulo de inventario
	//inventarioRoutes := api.Group("/inventario")
	//inventarioRoutes.Use(permissionMiddleware.RequireModuleAccess("inventario"))

	// Ejemplo: Ruta que requiere cualquiera de varios permisos
	/*api.GET("/dashboard", permissionMiddleware.RequireAnyPermission(
		"admin:dashboard",
		"finanzas:dashboard",
		"inventario:dashboard",
	), handleDashboard)*/

	// Ejemplo: Ruta que requiere múltiples permisos
	/*api.POST("/importar-datos", permissionMiddleware.RequireAllPermissions(
		"admin:data:import",
		"admin:data:modify",
	), handleImportarDatos)*/

	// Obtener puerto
	port := getEnv("PORT", "3000")

	// Configurar servidor HTTP
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Iniciar el servidor en una goroutine
	go func() {
		log.Printf("Servidor iniciando en el puerto %s...\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error al iniciar el servidor: %v", err)
		}
	}()

	// Esperar señal para graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Apagando servidor...")

	// Establecer un timeout para el cierre
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error al cerrar el servidor: %v", err)
	}

	log.Println("Servidor apagado correctamente")
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

// setupGracefulShutdown configura el cierre correcto de MongoDB
func setupGracefulShutdown(client *mongo.Client) {
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Cerrando la conexión a MongoDB...")
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Error al cerrar la conexión a MongoDB: %v", err)
		}
	}()
}
