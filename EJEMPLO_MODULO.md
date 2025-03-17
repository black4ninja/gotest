# Implementación Guiada de un Módulo: Productos

Este documento proporciona un ejemplo detallado de cómo implementar un nuevo módulo en la arquitectura actual, siguiendo la guía de NUEVO_MODULO.md. Usaremos un módulo de "Productos" como ejemplo.

## 1. Estructura de Carpetas

Primero, creamos la estructura de carpetas para nuestro módulo:

```
internal/
└── product/              # Módulo de productos
    ├── domain/           # Definiciones y contratos
    ├── repository/       # Implementación de persistencia 
    ├── usecase/          # Implementación de lógica de negocio
    └── delivery/         # Controladores HTTP/Handlers
```

## 2. Implementación del Dominio (Domain)

El archivo `internal/product/domain/product.domain.go` define nuestras entidades y contratos:

```go
package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Constantes para el estado del producto
const (
	ProductStatusActive   = "active"
	ProductStatusInactive = "inactive"
	ProductStatusArchived = "archived"
)

// Product representa la entidad de producto
type Product struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price"`
	Stock       int                `json:"stock" bson:"stock"`
	Category    string             `json:"category" bson:"category"`
	Status      string             `json:"status" bson:"status"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	ArchivedAt  *time.Time         `json:"archived_at,omitempty" bson:"archived_at,omitempty"`
}

// CreateProductRequest representa la solicitud para crear un producto
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gte=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	Category    string  `json:"category" binding:"required"`
}

// UpdateProductRequest representa la solicitud para actualizar un producto
type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"omitempty,gte=0"`
	Stock       int     `json:"stock" binding:"omitempty,gte=0"`
	Category    string  `json:"category"`
	Status      string  `json:"status"`
}

// ProductResponse representa la respuesta con datos de producto
type ProductResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Stock       int        `json:"stock"`
	Category    string     `json:"category"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ArchivedAt  *time.Time `json:"archived_at,omitempty"`
}

// ProductRepository define el contrato para la capa de persistencia
type ProductRepository interface {
	GetByID(id string) (*Product, error)
	GetAll(params map[string]interface{}) ([]*Product, error)
	Create(product *Product) error
	Update(product *Product) error
	Delete(id string) error
	Archive(id string) error
}

// ProductUseCase define el contrato para la capa de casos de uso
type ProductUseCase interface {
	GetProduct(id string) (*ProductResponse, error)
	GetAllProducts(params map[string]interface{}) ([]*ProductResponse, error)
	CreateProduct(req *CreateProductRequest) (*ProductResponse, error)
	UpdateProduct(id string, req *UpdateProductRequest) (*ProductResponse, error)
	DeleteProduct(id string) error
	ArchiveProduct(id string) error
}
```

### Explicación del Dominio

1. **Constantes**: Definimos estados posibles para los productos.
2. **Entidad Product**: La estructura principal que representa un producto en el sistema.
3. **DTOs de Solicitud/Respuesta**: Estructuras específicas para entrada y salida de datos.
    - `CreateProductRequest`: Para crear nuevos productos
    - `UpdateProductRequest`: Para actualizar productos existentes
    - `ProductResponse`: Para devolver datos de productos al cliente
4. **Interfaces de Contrato**:
    - `ProductRepository`: Define qué operaciones de persistencia se necesitan
    - `ProductUseCase`: Define la API de lógica de negocio

## 3. Implementación del Repositorio (Repository)

El archivo `internal/product/repository/mongo.product.repository.go` implementa el acceso a datos:

```go
package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/black4ninja/mi-proyecto/internal/product/domain"
)

type mongoProductRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

// NewMongoProductRepository crea un nuevo repositorio de productos con MongoDB
func NewMongoProductRepository(collection *mongo.Collection) domain.ProductRepository {
	return &mongoProductRepository{
		collection: collection,
		timeout:    10 * time.Second,
	}
}

// GetByID obtiene un producto por su ID
func (r *mongoProductRepository) GetByID(id string) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var product domain.Product
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("producto no encontrado")
		}
		return nil, err
	}

	return &product, nil
}

// GetAll obtiene todos los productos que coincidan con los parámetros dados
func (r *mongoProductRepository) GetAll(params map[string]interface{}) ([]*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	// Construir filtro
	filter := bson.M{}
	for key, value := range params {
		filter[key] = value
	}

	opts := options.Find()
	opts.SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

// Create crea un nuevo producto
func (r *mongoProductRepository) Create(product *domain.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	product.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, product)
	return err
}

// Update actualiza un producto existente
func (r *mongoProductRepository) Update(product *domain.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"stock":       product.Stock,
			"category":    product.Category,
			"status":      product.Status,
			"updated_at":  time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": product.ID},
		update,
	)
	return err
}

// Delete elimina un producto
func (r *mongoProductRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// Archive marca un producto como archivado
func (r *mongoProductRepository) Archive(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":      domain.ProductStatusArchived,
			"archived_at": now,
			"updated_at":  now,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}
```

### Explicación del Repositorio

1. **Estructura**: `mongoProductRepository` implementa la interfaz `ProductRepository`.
2. **Constructor**: `NewMongoProductRepository` crea una nueva instancia con la colección MongoDB.
3. **Métodos de acceso a datos**:
    - Cada método implementa una operación específica del contrato
    - Se usan contextos con timeout para seguridad
    - Se convierten IDs de string a ObjectID cuando sea necesario
    - Se manejan errores específicos de MongoDB

## 4. Implementación de Casos de Uso (UseCase)

El archivo `internal/product/usecase/product.usecase.go` implementa la lógica de negocio:

```go
package usecase

import (
	"errors"
	"time"

	"github.com/black4ninja/mi-proyecto/internal/product/domain"
)

type productUseCase struct {
	productRepo domain.ProductRepository
}

// NewProductUseCase crea un nuevo caso de uso para productos
func NewProductUseCase(productRepo domain.ProductRepository) domain.ProductUseCase {
	return &productUseCase{
		productRepo: productRepo,
	}
}

// GetProduct obtiene un producto por su ID
func (u *productUseCase) GetProduct(id string) (*domain.ProductResponse, error) {
	product, err := u.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &domain.ProductResponse{
		ID:          product.ID.Hex(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    product.Category,
		Status:      product.Status,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		ArchivedAt:  product.ArchivedAt,
	}, nil
}

// GetAllProducts obtiene todos los productos
func (u *productUseCase) GetAllProducts(params map[string]interface{}) ([]*domain.ProductResponse, error) {
	products, err := u.productRepo.GetAll(params)
	if err != nil {
		return nil, err
	}

	var response []*domain.ProductResponse
	for _, product := range products {
		response = append(response, &domain.ProductResponse{
			ID:          product.ID.Hex(),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Category:    product.Category,
			Status:      product.Status,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
			ArchivedAt:  product.ArchivedAt,
		})
	}

	return response, nil
}

// CreateProduct crea un nuevo producto
func (u *productUseCase) CreateProduct(req *domain.CreateProductRequest) (*domain.ProductResponse, error) {
	// Crear producto
	now := time.Now()
	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		Status:      domain.ProductStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := u.productRepo.Create(product); err != nil {
		return nil, err
	}

	return &domain.ProductResponse{
		ID:          product.ID.Hex(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    product.Category,
		Status:      product.Status,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}

// UpdateProduct actualiza un producto existente
func (u *productUseCase) UpdateProduct(id string, req *domain.UpdateProductRequest) (*domain.ProductResponse, error) {
	// Obtener producto existente
	product, err := u.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Actualizar campos
	if req.Name != "" {
		product.Name = req.Name
	}

	if req.Description != "" {
		product.Description = req.Description
	}

	if req.Price >= 0 {
		product.Price = req.Price
	}

	if req.Stock >= 0 {
		product.Stock = req.Stock
	}

	if req.Category != "" {
		product.Category = req.Category
	}

	if req.Status != "" {
		// Validar estado
		if req.Status != domain.ProductStatusActive && 
		   req.Status != domain.ProductStatusInactive && 
		   req.Status != domain.ProductStatusArchived {
			return nil, errors.New("estado de producto inválido")
		}
		product.Status = req.Status
	}

	product.UpdatedAt = time.Now()

	if err := u.productRepo.Update(product); err != nil {
		return nil, err
	}

	return &domain.ProductResponse{
		ID:          product.ID.Hex(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    product.Category,
		Status:      product.Status,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		ArchivedAt:  product.ArchivedAt,
	}, nil
}

// DeleteProduct elimina un producto
func (u *productUseCase) DeleteProduct(id string) error {
	return u.productRepo.Delete(id)
}

// ArchiveProduct archiva un producto
func (u *productUseCase) ArchiveProduct(id string) error {
	return u.productRepo.Archive(id)
}
```

### Explicación de los Casos de Uso

1. **Estructura**: `productUseCase` implementa la interfaz `ProductUseCase`.
2. **Constructor**: `NewProductUseCase` recibe las dependencias necesarias (el repositorio).
3. **Lógica de negocio**:
    - Los métodos aplican la lógica específica de cada operación
    - Se realizan validaciones adicionales
    - Se transforman entidades del dominio a DTOs de respuesta
    - Se delegan operaciones de persistencia al repositorio

## 5. Implementación del Delivery (HTTP Handlers)

El archivo `internal/product/delivery/product.delivery.go` implementa los handlers HTTP:

```go
package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/black4ninja/mi-proyecto/internal/product/domain"
	"github.com/black4ninja/mi-proyecto/pkg/utils"
)

// ProductHandler maneja las peticiones HTTP para productos
type ProductHandler struct {
	productUseCase domain.ProductUseCase
}

// NewProductHandler crea un nuevo manejador de productos
func NewProductHandler(router *gin.RouterGroup, useCase domain.ProductUseCase) {
	handler := &ProductHandler{
		productUseCase: useCase,
	}

	// Rutas de productos
	router.GET("/", handler.GetAllProducts)
	router.GET("/:id", handler.GetProduct)
	router.POST("/", handler.CreateProduct)
	router.PUT("/:id", handler.UpdateProduct)
	router.DELETE("/:id", handler.DeleteProduct)
	router.PUT("/:id/archive", handler.ArchiveProduct)
}

// GetAllProducts manejador para obtener todos los productos
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	// Extraer todos los parámetros de consulta
	queryParams := make(map[string]string)

	if status := c.Query("status"); status != "" {
		queryParams["status"] = status
	}
	if category := c.Query("category"); category != "" {
		queryParams["category"] = category
	}
	if name := c.Query("name"); name != "" {
		queryParams["name"] = name
	}

	// Construir filtro para MongoDB
	filter := utils.BuildMongoFilter(queryParams, utils.FilterConfig{
		"status": utils.FilterDefinition{
			AllowedValues: []string{domain.ProductStatusActive, domain.ProductStatusInactive, domain.ProductStatusArchived},
		},
		"category": utils.FilterDefinition{},
		"name": utils.FilterDefinition{
			Transformer: utils.TransformToRegex,
		},
	})

	// Si no se especificó un estado, mostrar solo productos activos por defecto
	if _, hasStatus := filter["status"]; !hasStatus {
		filter["status"] = domain.ProductStatusActive
	}

	// Obtener todos los productos con los filtros aplicados
	products, err := h.productUseCase.GetAllProducts(filter)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Productos obtenidos con éxito", products)
}

// GetProduct manejador para obtener un producto
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")

	product, err := h.productUseCase.GetProduct(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Producto obtenido con éxito", product)
}

// CreateProduct manejador para crear un producto
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req domain.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	product, err := h.productUseCase.CreateProduct(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Producto creado con éxito", product)
}

// UpdateProduct manejador para actualizar un producto
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var req domain.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	product, err := h.productUseCase.UpdateProduct(id, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Producto actualizado con éxito", product)
}

// DeleteProduct manejador para eliminar un producto
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if err := h.productUseCase.DeleteProduct(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Producto eliminado con éxito", nil)
}

// ArchiveProduct manejador para archivar un producto
func (h *ProductHandler) ArchiveProduct(c *gin.Context) {
	id := c.Param("id")

	if err := h.productUseCase.ArchiveProduct(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Producto archivado con éxito", nil)
}
```

### Explicación del Delivery

1. **Estructura**: `ProductHandler` maneja las peticiones HTTP para productos.
2. **Constructor**: `NewProductHandler` configura las rutas y dependencias.
3. **Definición de rutas**: Se definen todos los endpoints necesarios.
4. **Handlers HTTP**:
    - Reciben solicitudes HTTP mediante Gin
    - Extraen y validan parámetros o cuerpos JSON
    - Llaman a los métodos apropiados del caso de uso
    - Transforman la respuesta al formato estándar de la API

## 6. Integración con el Sistema Principal

Los cambios necesarios en `main.go`:

```go
// En la sección de colecciones de MongoDB
productCollection := mongoClient.Database(mongoDBName).Collection("products")

// En la sección de inicialización de repositorios
productRepository := productRepo.NewMongoProductRepository(productCollection)

// En la sección de inicialización de casos de uso
productService := productUseCase.NewProductUseCase(productRepository)

// En la sección de configuración de rutas
// Rutas de productos (protegidas con OAuth)
productRoutes := api.Group("/products")
productRoutes.Use(permissionMiddleware.RequirePermission("products:access")) // Opcional: añadir middleware de permisos
productDelivery.NewProductHandler(productRoutes, productService)

// También puedes añadir rutas específicas con diferentes permisos
// Por ejemplo:
api.GET("/catalog", permissionMiddleware.RequirePermission("products:read"), func(c *gin.Context) {
    // Lógica para mostrar catálogo público
})
```

### Explicación de la Integración

1. **Inicialización de la colección MongoDB** para productos.
2. **Creación del repositorio** pasando la colección MongoDB.
3. **Creación del caso de uso** pasando el repositorio.
4. **Configuración de rutas** utilizando grupos de Gin.
5. **Aplicación de middleware de permisos** para proteger las rutas.

## 7. Configuración de Permisos

En `scripts/init_permissions_and_admin.go`:

```go
// Crear permisos para el nuevo módulo
createDefaultPermission(permissionService, "products:access", "products", "access", "Acceso a productos", "Permite acceso básico al módulo de productos")
createDefaultPermission(permissionService, "products:read", "products", "read", "Ver productos", "Permite ver productos")
createDefaultPermission(permissionService, "products:write", "products", "write", "Gestionar productos", "Permite crear y modificar productos")
createDefaultPermission(permissionService, "products:delete", "products", "delete", "Eliminar productos", "Permite eliminar productos")

// Asignar permisos a roles existentes
roleService.AddPermissionToRole(adminRoleID, "products:access")
roleService.AddPermissionToRole(adminRoleID, "products:read")
roleService.AddPermissionToRole(adminRoleID, "products:write")
roleService.AddPermissionToRole(adminRoleID, "products:delete")
```

### Explicación de Permisos

1. **Creación de permisos** específicos para el módulo de productos.
2. **Asignación de permisos** a roles existentes (como el rol de administrador).

## Resumen

Este ejemplo completo muestra cómo crear un nuevo módulo siguiendo los principios de Clean Architecture:

1. **Domain**: Define contratos e interfaces
2. **Repository**: Implementa acceso a datos
3. **UseCase**: Implementa lógica de negocio
4. **Delivery**: Implementa handlers HTTP
5. **Integración**: Conecta el módulo con el sistema principal
6. **Permisos**: Configura el control de acceso al módulo

Siguiendo este patrón, puedes crear cualquier nuevo módulo de manera consistente y mantenible.