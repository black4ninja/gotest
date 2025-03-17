// pkg/tools/module_generator.go
// Este script genera la estructura básica de un nuevo módulo

package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// GenerateModule crea la estructura básica de un nuevo módulo
func GenerateModule(moduleName string) error {
	// Convertir a minúsculas y quitar espacios
	moduleName = strings.ToLower(strings.TrimSpace(moduleName))

	// Verificar nombre
	if moduleName == "" {
		return fmt.Errorf("el nombre del módulo no puede estar vacío")
	}

	// Rutas base
	baseDir := "internal/" + moduleName
	dirs := []string{
		baseDir + "/domain",
		baseDir + "/repository",
		baseDir + "/usecase",
		baseDir + "/delivery",
	}

	// Crear estructura de directorios
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error al crear directorio %s: %w", dir, err)
		}
		fmt.Printf("Directorio creado: %s\n", dir)
	}

	// Generar archivos
	files := map[string]string{
		baseDir + "/domain/" + moduleName + ".domain.go":               domainTemplate,
		baseDir + "/repository/mongo." + moduleName + ".repository.go": repositoryTemplate,
		baseDir + "/usecase/" + moduleName + ".usecase.go":             usecaseTemplate,
		baseDir + "/delivery/" + moduleName + ".delivery.go":           deliveryTemplate,
	}

	data := struct {
		ModuleName      string
		ModuleNameTitle string
	}{
		ModuleName:      moduleName,
		ModuleNameTitle: strings.Title(moduleName),
	}

	for file, templateContent := range files {
		if err := generateFile(file, templateContent, data); err != nil {
			return fmt.Errorf("error al generar archivo %s: %w", file, err)
		}
		fmt.Printf("Archivo generado: %s\n", file)
	}

	// Generar fragmento para main.go
	mainFragment := filepath.Join(baseDir, "main_fragment.go.txt")
	if err := generateFile(mainFragment, mainTemplate, data); err != nil {
		return fmt.Errorf("error al generar fragmento para main.go: %w", err)
	}
	fmt.Printf("\nArchivo generado: %s\n", mainFragment)
	fmt.Printf("\nFragmento para agregar a main.go creado. Revise el archivo %s\n", mainFragment)

	return nil
}

func generateFile(path, content string, data interface{}) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl, err := template.New("file").Parse(content)
	if err != nil {
		return err
	}

	return tmpl.Execute(file, data)
}

// Templates para los archivos
const domainTemplate = `package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Constantes para el estado del {{.ModuleName}}
const (
	{{.ModuleNameTitle}}StatusActive   = "active"
	{{.ModuleNameTitle}}StatusInactive = "inactive"
	{{.ModuleNameTitle}}StatusArchived = "archived"
)

// {{.ModuleNameTitle}} representa la entidad de {{.ModuleName}}
type {{.ModuleNameTitle}} struct {
	ID          primitive.ObjectID ` + "`json:\"id\" bson:\"_id,omitempty\"`" + `
	Name        string             ` + "`json:\"name\" bson:\"name\"`" + `
	Description string             ` + "`json:\"description\" bson:\"description\"`" + `
	Status      string             ` + "`json:\"status\" bson:\"status\"`" + `
	CreatedAt   time.Time          ` + "`json:\"created_at\" bson:\"created_at\"`" + `
	UpdatedAt   time.Time          ` + "`json:\"updated_at\" bson:\"updated_at\"`" + `
	ArchivedAt  *time.Time         ` + "`json:\"archived_at,omitempty\" bson:\"archived_at,omitempty\"`" + `
	// Añade aquí tus campos específicos
}

// Create{{.ModuleNameTitle}}Request representa la solicitud para crear un {{.ModuleName}}
type Create{{.ModuleNameTitle}}Request struct {
	Name        string ` + "`json:\"name\" binding:\"required\"`" + `
	Description string ` + "`json:\"description\"`" + `
	// Añade aquí tus campos específicos
}

// Update{{.ModuleNameTitle}}Request representa la solicitud para actualizar un {{.ModuleName}}
type Update{{.ModuleNameTitle}}Request struct {
	Name        string ` + "`json:\"name\"`" + `
	Description string ` + "`json:\"description\"`" + `
	Status      string ` + "`json:\"status\"`" + `
	// Añade aquí tus campos específicos
}

// {{.ModuleNameTitle}}Response representa la respuesta con datos de {{.ModuleName}}
type {{.ModuleNameTitle}}Response struct {
	ID          string     ` + "`json:\"id\"`" + `
	Name        string     ` + "`json:\"name\"`" + `
	Description string     ` + "`json:\"description\"`" + `
	Status      string     ` + "`json:\"status\"`" + `
	CreatedAt   time.Time  ` + "`json:\"created_at\"`" + `
	UpdatedAt   time.Time  ` + "`json:\"updated_at\"`" + `
	ArchivedAt  *time.Time ` + "`json:\"archived_at,omitempty\"`" + `
	// Añade aquí tus campos específicos
}

// {{.ModuleNameTitle}}Repository define el contrato para la capa de persistencia
type {{.ModuleNameTitle}}Repository interface {
	GetByID(id string) (*{{.ModuleNameTitle}}, error)
	GetAll(params map[string]interface{}) ([]*{{.ModuleNameTitle}}, error)
	Create({{.ModuleName}} *{{.ModuleNameTitle}}) error
	Update({{.ModuleName}} *{{.ModuleNameTitle}}) error
	Delete(id string) error
	Archive(id string) error
}

// {{.ModuleNameTitle}}UseCase define el contrato para la capa de casos de uso
type {{.ModuleNameTitle}}UseCase interface {
	Get{{.ModuleNameTitle}}(id string) (*{{.ModuleNameTitle}}Response, error)
	GetAll{{.ModuleNameTitle}}s(params map[string]interface{}) ([]*{{.ModuleNameTitle}}Response, error)
	Create{{.ModuleNameTitle}}(req *Create{{.ModuleNameTitle}}Request) (*{{.ModuleNameTitle}}Response, error)
	Update{{.ModuleNameTitle}}(id string, req *Update{{.ModuleNameTitle}}Request) (*{{.ModuleNameTitle}}Response, error)
	Delete{{.ModuleNameTitle}}(id string) error
	Archive{{.ModuleNameTitle}}(id string) error
}
`

const repositoryTemplate = `package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/black4ninja/mi-proyecto/internal/{{.ModuleName}}/domain"
)

type mongo{{.ModuleNameTitle}}Repository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

// NewMongo{{.ModuleNameTitle}}Repository crea un nuevo repositorio de {{.ModuleName}}s con MongoDB
func NewMongo{{.ModuleNameTitle}}Repository(collection *mongo.Collection) domain.{{.ModuleNameTitle}}Repository {
	return &mongo{{.ModuleNameTitle}}Repository{
		collection: collection,
		timeout:    10 * time.Second,
	}
}

// GetByID obtiene un {{.ModuleName}} por su ID
func (r *mongo{{.ModuleNameTitle}}Repository) GetByID(id string) (*domain.{{.ModuleNameTitle}}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var {{.ModuleName}} domain.{{.ModuleNameTitle}}
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&{{.ModuleName}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("{{.ModuleName}} no encontrado")
		}
		return nil, err
	}

	return &{{.ModuleName}}, nil
}

// GetAll obtiene todos los {{.ModuleName}}s que coincidan con los parámetros dados
func (r *mongo{{.ModuleNameTitle}}Repository) GetAll(params map[string]interface{}) ([]*domain.{{.ModuleNameTitle}}, error) {
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

	var {{.ModuleName}}s []*domain.{{.ModuleNameTitle}}
	if err := cursor.All(ctx, &{{.ModuleName}}s); err != nil {
		return nil, err
	}

	return {{.ModuleName}}s, nil
}

// Create crea un nuevo {{.ModuleName}}
func (r *mongo{{.ModuleNameTitle}}Repository) Create({{.ModuleName}} *domain.{{.ModuleNameTitle}}) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	{{.ModuleName}}.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, {{.ModuleName}})
	return err
}

// Update actualiza un {{.ModuleName}} existente
func (r *mongo{{.ModuleNameTitle}}Repository) Update({{.ModuleName}} *domain.{{.ModuleNameTitle}}) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"name":        {{.ModuleName}}.Name,
			"description": {{.ModuleName}}.Description,
			"status":      {{.ModuleName}}.Status,
			"updated_at":  time.Now(),
			// Actualiza aquí tus campos específicos
		},
	}

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": {{.ModuleName}}.ID},
		update,
	)
	return err
}

// Delete elimina un {{.ModuleName}}
func (r *mongo{{.ModuleNameTitle}}Repository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// Archive marca un {{.ModuleName}} como archivado
func (r *mongo{{.ModuleNameTitle}}Repository) Archive(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":      domain.{{.ModuleNameTitle}}StatusArchived,
			"archived_at": now,
			"updated_at":  now,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}
`

const usecaseTemplate = `package usecase

import (
	"errors"
	"time"

	"github.com/black4ninja/mi-proyecto/internal/{{.ModuleName}}/domain"
)

type {{.ModuleName}}UseCase struct {
	{{.ModuleName}}Repo domain.{{.ModuleNameTitle}}Repository
}

// New{{.ModuleNameTitle}}UseCase crea un nuevo caso de uso para {{.ModuleName}}s
func New{{.ModuleNameTitle}}UseCase({{.ModuleName}}Repo domain.{{.ModuleNameTitle}}Repository) domain.{{.ModuleNameTitle}}UseCase {
	return &{{.ModuleName}}UseCase{
		{{.ModuleName}}Repo: {{.ModuleName}}Repo,
	}
}

// Get{{.ModuleNameTitle}} obtiene un {{.ModuleName}} por su ID
func (u *{{.ModuleName}}UseCase) Get{{.ModuleNameTitle}}(id string) (*domain.{{.ModuleNameTitle}}Response, error) {
	{{.ModuleName}}, err := u.{{.ModuleName}}Repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &domain.{{.ModuleNameTitle}}Response{
		ID:          {{.ModuleName}}.ID.Hex(),
		Name:        {{.ModuleName}}.Name,
		Description: {{.ModuleName}}.Description,
		Status:      {{.ModuleName}}.Status,
		CreatedAt:   {{.ModuleName}}.CreatedAt,
		UpdatedAt:   {{.ModuleName}}.UpdatedAt,
		ArchivedAt:  {{.ModuleName}}.ArchivedAt,
		// Añade aquí tus campos específicos
	}, nil
}

// GetAll{{.ModuleNameTitle}}s obtiene todos los {{.ModuleName}}s
func (u *{{.ModuleName}}UseCase) GetAll{{.ModuleNameTitle}}s(params map[string]interface{}) ([]*domain.{{.ModuleNameTitle}}Response, error) {
	{{.ModuleName}}s, err := u.{{.ModuleName}}Repo.GetAll(params)
	if err != nil {
		return nil, err
	}

	var response []*domain.{{.ModuleNameTitle}}Response
	for _, {{.ModuleName}} := range {{.ModuleName}}s {
		response = append(response, &domain.{{.ModuleNameTitle}}Response{
			ID:          {{.ModuleName}}.ID.Hex(),
			Name:        {{.ModuleName}}.Name,
			Description: {{.ModuleName}}.Description,
			Status:      {{.ModuleName}}.Status,
			CreatedAt:   {{.ModuleName}}.CreatedAt,
			UpdatedAt:   {{.ModuleName}}.UpdatedAt,
			ArchivedAt:  {{.ModuleName}}.ArchivedAt,
			// Añade aquí tus campos específicos
		})
	}

	return response, nil
}

// Create{{.ModuleNameTitle}} crea un nuevo {{.ModuleName}}
func (u *{{.ModuleName}}UseCase) Create{{.ModuleNameTitle}}(req *domain.Create{{.ModuleNameTitle}}Request) (*domain.{{.ModuleNameTitle}}Response, error) {
	// Crear {{.ModuleName}}
	now := time.Now()
	{{.ModuleName}} := &domain.{{.ModuleNameTitle}}{
		Name:        req.Name,
		Description: req.Description,
		Status:      domain.{{.ModuleNameTitle}}StatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
		// Añade aquí tus campos específicos
	}

	if err := u.{{.ModuleName}}Repo.Create({{.ModuleName}}); err != nil {
		return nil, err
	}

	return &domain.{{.ModuleNameTitle}}Response{
		ID:          {{.ModuleName}}.ID.Hex(),
		Name:        {{.ModuleName}}.Name,
		Description: {{.ModuleName}}.Description,
		Status:      {{.ModuleName}}.Status,
		CreatedAt:   {{.ModuleName}}.CreatedAt,
		UpdatedAt:   {{.ModuleName}}.UpdatedAt,
		// Añade aquí tus campos específicos
	}, nil
}

// Update{{.ModuleNameTitle}} actualiza un {{.ModuleName}} existente
func (u *{{.ModuleName}}UseCase) Update{{.ModuleNameTitle}}(id string, req *domain.Update{{.ModuleNameTitle}}Request) (*domain.{{.ModuleNameTitle}}Response, error) {
	// Obtener {{.ModuleName}} existente
	{{.ModuleName}}, err := u.{{.ModuleName}}Repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Actualizar campos
	if req.Name != "" {
		{{.ModuleName}}.Name = req.Name
	}

	if req.Description != "" {
		{{.ModuleName}}.Description = req.Description
	}

	if req.Status != "" {
		// Validar estado
		if req.Status != domain.{{.ModuleNameTitle}}StatusActive && 
		   req.Status != domain.{{.ModuleNameTitle}}StatusInactive && 
		   req.Status != domain.{{.ModuleNameTitle}}StatusArchived {
			return nil, errors.New("estado de {{.ModuleName}} inválido")
		}
		{{.ModuleName}}.Status = req.Status
	}

	// Añade aquí tus campos específicos

	{{.ModuleName}}.UpdatedAt = time.Now()

	if err := u.{{.ModuleName}}Repo.Update({{.ModuleName}}); err != nil {
		return nil, err
	}

	return &domain.{{.ModuleNameTitle}}Response{
		ID:          {{.ModuleName}}.ID.Hex(),
		Name:        {{.ModuleName}}.Name,
		Description: {{.ModuleName}}.Description,
		Status:      {{.ModuleName}}.Status,
		CreatedAt:   {{.ModuleName}}.CreatedAt,
		UpdatedAt:   {{.ModuleName}}.UpdatedAt,
		ArchivedAt:  {{.ModuleName}}.ArchivedAt,
		// Añade aquí tus campos específicos
	}, nil
}

// Delete{{.ModuleNameTitle}} elimina un {{.ModuleName}}
func (u *{{.ModuleName}}UseCase) Delete{{.ModuleNameTitle}}(id string) error {
	return u.{{.ModuleName}}Repo.Delete(id)
}

// Archive{{.ModuleNameTitle}} archiva un {{.ModuleName}}
func (u *{{.ModuleName}}UseCase) Archive{{.ModuleNameTitle}}(id string) error {
	return u.{{.ModuleName}}Repo.Archive(id)
}
`

const deliveryTemplate = `package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/black4ninja/mi-proyecto/internal/{{.ModuleName}}/domain"
	"github.com/black4ninja/mi-proyecto/pkg/utils"
)

// {{.ModuleNameTitle}}Handler maneja las peticiones HTTP para {{.ModuleName}}s
type {{.ModuleNameTitle}}Handler struct {
	{{.ModuleName}}UseCase domain.{{.ModuleNameTitle}}UseCase
}

// New{{.ModuleNameTitle}}Handler crea un nuevo manejador de {{.ModuleName}}s
func New{{.ModuleNameTitle}}Handler(router *gin.RouterGroup, useCase domain.{{.ModuleNameTitle}}UseCase) {
	handler := &{{.ModuleNameTitle}}Handler{
		{{.ModuleName}}UseCase: useCase,
	}

	// Rutas de {{.ModuleName}}s
	router.GET("/", handler.GetAll{{.ModuleNameTitle}}s)
	router.GET("/:id", handler.Get{{.ModuleNameTitle}})
	router.POST("/", handler.Create{{.ModuleNameTitle}})
	router.PUT("/:id", handler.Update{{.ModuleNameTitle}})
	router.DELETE("/:id", handler.Delete{{.ModuleNameTitle}})
	router.PUT("/:id/archive", handler.Archive{{.ModuleNameTitle}})
}

// GetAll{{.ModuleNameTitle}}s manejador para obtener todos los {{.ModuleName}}s
func (h *{{.ModuleNameTitle}}Handler) GetAll{{.ModuleNameTitle}}s(c *gin.Context) {
	// Extraer todos los parámetros de consulta
	queryParams := make(map[string]string)

	if status := c.Query("status"); status != "" {
		queryParams["status"] = status
	}
	if name := c.Query("name"); name != "" {
		queryParams["name"] = name
	}

	// Construir filtro para MongoDB
	filter := utils.BuildMongoFilter(queryParams, utils.FilterConfig{
		"status": utils.FilterDefinition{
			AllowedValues: []string{domain.{{.ModuleNameTitle}}StatusActive, domain.{{.ModuleNameTitle}}StatusInactive, domain.{{.ModuleNameTitle}}StatusArchived},
		},
		"name": utils.FilterDefinition{
			Transformer: utils.TransformToRegex,
		},
	})

	// Si no se especificó un estado, mostrar solo {{.ModuleName}}s activos por defecto
	if _, hasStatus := filter["status"]; !hasStatus {
		filter["status"] = domain.{{.ModuleNameTitle}}StatusActive
	}

	// Obtener todos los {{.ModuleName}}s con los filtros aplicados
	{{.ModuleName}}s, err := h.{{.ModuleName}}UseCase.GetAll{{.ModuleNameTitle}}s(filter)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "{{.ModuleNameTitle}}s obtenidos con éxito", {{.ModuleName}}s)
}

// Get{{.ModuleNameTitle}} manejador para obtener un {{.ModuleName}}
func (h *{{.ModuleNameTitle}}Handler) Get{{.ModuleNameTitle}}(c *gin.Context) {
	id := c.Param("id")

	{{.ModuleName}}, err := h.{{.ModuleName}}UseCase.Get{{.ModuleNameTitle}}(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "{{.ModuleNameTitle}} obtenido con éxito", {{.ModuleName}})
}

// Create{{.ModuleNameTitle}} manejador para crear un {{.ModuleName}}
func (h *{{.ModuleNameTitle}}Handler) Create{{.ModuleNameTitle}}(c *gin.Context) {
	var req domain.Create{{.ModuleNameTitle}}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	{{.ModuleName}}, err := h.{{.ModuleName}}UseCase.Create{{.ModuleNameTitle}}(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "{{.ModuleNameTitle}} creado con éxito", {{.ModuleName}})
}

// Update{{.ModuleNameTitle}} manejador para actualizar un {{.ModuleName}}
func (h *{{.ModuleNameTitle}}Handler) Update{{.ModuleNameTitle}}(c *gin.Context) {
	id := c.Param("id")

	var req domain.Update{{.ModuleNameTitle}}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	{{.ModuleName}}, err := h.{{.ModuleName}}UseCase.Update{{.ModuleNameTitle}}(id, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "{{.ModuleNameTitle}} actualizado con éxito", {{.ModuleName}})
}

// Delete{{.ModuleNameTitle}} manejador para eliminar un {{.ModuleName}}
func (h *{{.ModuleNameTitle}}Handler) Delete{{.ModuleNameTitle}}(c *gin.Context) {
	id := c.Param("id")

	if err := h.{{.ModuleName}}UseCase.Delete{{.ModuleNameTitle}}(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "{{.ModuleNameTitle}} eliminado con éxito", nil)
}

// Archive{{.ModuleNameTitle}} manejador para archivar un {{.ModuleName}}
func (h *{{.ModuleNameTitle}}Handler) Archive{{.ModuleNameTitle}}(c *gin.Context) {
	id := c.Param("id")

	if err := h.{{.ModuleName}}UseCase.Archive{{.ModuleNameTitle}}(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "{{.ModuleNameTitle}} archivado con éxito", nil)
}
`

const mainTemplate = `// Fragmento para añadir a main.go

// En la sección de colecciones de MongoDB
{{.ModuleName}}Collection := mongoClient.Database(mongoDBName).Collection("{{.ModuleName}}s")

// En la sección de inicialización de repositorios
{{.ModuleName}}Repository := {{.ModuleName}}Repo.NewMongo{{.ModuleNameTitle}}Repository({{.ModuleName}}Collection)

// En la sección de inicialización de casos de uso
{{.ModuleName}}Service := {{.ModuleName}}UseCase.New{{.ModuleNameTitle}}UseCase({{.ModuleName}}Repository)

// En la sección de configuración de rutas
// Rutas de {{.ModuleName}}s (protegidas con OAuth)
{{.ModuleName}}Routes := api.Group("/{{.ModuleName}}s")
{{.ModuleName}}Routes.Use(permissionMiddleware.RequirePermission("{{.ModuleName}}s:access")) // Opcional: middleware de permisos
{{.ModuleName}}Delivery.New{{.ModuleNameTitle}}Handler({{.ModuleName}}Routes, {{.ModuleName}}Service)

// También puedes añadir permisos en scripts/init_permissions_and_admin.go:
// createDefaultPermission(permissionService, "{{.ModuleName}}s:access", "{{.ModuleName}}s", "access", "Acceso a {{.ModuleName}}s", "Permite acceso básico al módulo de {{.ModuleName}}s")
// createDefaultPermission(permissionService, "{{.ModuleName}}s:read", "{{.ModuleName}}s", "read", "Ver {{.ModuleName}}s", "Permite ver {{.ModuleName}}s")
// createDefaultPermission(permissionService, "{{.ModuleName}}s:write", "{{.ModuleName}}s", "write", "Gestionar {{.ModuleName}}s", "Permite crear y modificar {{.ModuleName}}s")
// createDefaultPermission(permissionService, "{{.ModuleName}}s:delete", "{{.ModuleName}}s", "delete", "Eliminar {{.ModuleName}}s", "Permite eliminar {{.ModuleName}}s")
`
