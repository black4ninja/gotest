// pkg/tools/swagger_documenter.go
package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// DocumentModule añade comentarios Swagger a un módulo específico
func DocumentModule(moduleName string) error {
	// Convertir a minúsculas y quitar espacios
	moduleName = strings.ToLower(strings.TrimSpace(moduleName))

	// Verificar nombre
	if moduleName == "" {
		return fmt.Errorf("el nombre del módulo no puede estar vacío")
	}

	// Ruta al archivo del handler
	handlerPath := filepath.Join("internal", moduleName, "delivery", moduleName+".delivery.go")

	// Verificar que el archivo exista
	if _, err := os.Stat(handlerPath); os.IsNotExist(err) {
		return fmt.Errorf("el archivo %s no existe", handlerPath)
	}

	// Leer el archivo
	content, err := os.ReadFile(handlerPath)
	if err != nil {
		return fmt.Errorf("error al leer el archivo %s: %w", handlerPath, err)
	}

	// Añadir comentarios Swagger a los métodos
	updatedContent, err := addSwaggerComments(string(content), moduleName)
	if err != nil {
		return fmt.Errorf("error al añadir comentarios Swagger: %w", err)
	}

	// Escribir el archivo actualizado
	err = os.WriteFile(handlerPath, []byte(updatedContent), 0644)
	if err != nil {
		return fmt.Errorf("error al escribir el archivo %s: %w", handlerPath, err)
	}

	fmt.Printf("Comentarios Swagger añadidos al módulo %s\n", moduleName)

	// También documentar el dominio
	err = documentDomain(moduleName)
	if err != nil {
		fmt.Printf("Advertencia: No se pudo documentar completamente el dominio: %v\n", err)
	}

	return nil
}

// DocumentAllModules añade comentarios Swagger a todos los módulos
func DocumentAllModules() error {
	// Buscar todos los módulos
	modules, err := findModules()
	if err != nil {
		return fmt.Errorf("error al buscar módulos: %w", err)
	}

	fmt.Printf("Se encontraron %d módulos\n", len(modules))

	// Documentar cada módulo
	for _, module := range modules {
		fmt.Printf("Documentando módulo: %s\n", module)
		err := DocumentModule(module)
		if err != nil {
			fmt.Printf("Error al documentar el módulo %s: %v\n", module, err)
		}
	}

	return nil
}

// findModules encuentra todos los módulos en la carpeta internal
func findModules() ([]string, error) {
	var modules []string

	dirs, err := os.ReadDir("internal")
	if err != nil {
		return nil, fmt.Errorf("error al leer el directorio internal: %w", err)
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			// Verificar si es un módulo (tiene subcarpeta domain)
			domainPath := filepath.Join("internal", dir.Name(), "domain")
			if _, err := os.Stat(domainPath); !os.IsNotExist(err) {
				modules = append(modules, dir.Name())
			}
		}
	}

	return modules, nil
}

// addSwaggerComments añade comentarios Swagger a los métodos del handler
func addSwaggerComments(content, moduleName string) (string, error) {
	// Convertir primera letra a mayúscula para el nombre del handler
	moduleTitle := strings.Title(moduleName)

	// Crear expresiones regulares para encontrar métodos
	methodPatterns := []struct {
		regex   *regexp.Regexp
		comment string
	}{
		{
			// GetAll[ModuleName]s
			regexp.MustCompile(`func \(h \*` + moduleTitle + `Handler\) GetAll` + moduleTitle + `s\(c \*gin\.Context\)`),
			fmt.Sprintf(`// @Summary Obtener todos los %ss
// @Description Obtiene una lista de todos los %ss con filtrado opcional
// @Tags %ss
// @Accept json
// @Produce json
// @Param status query string false "Estado del %s (active, inactive, archived)"
// @Param name query string false "Nombre del %s (búsqueda parcial)"
// @Success 200 {object} utils.Response{data=[]domain.%sResponse} "Lista de %ss"
// @Failure 500 {object} utils.Response "Error interno"
// @Router /%ss [get]
// @Security BearerAuth`, moduleName, moduleName, moduleName, moduleName, moduleName, moduleTitle, moduleName, moduleName),
		},
		{
			// Get[ModuleName]
			regexp.MustCompile(`func \(h \*` + moduleTitle + `Handler\) Get` + moduleTitle + `\(c \*gin\.Context\)`),
			fmt.Sprintf(`// @Summary Obtener un %s
// @Description Obtiene un %s por su ID
// @Tags %ss
// @Accept json
// @Produce json
// @Param id path string true "ID del %s"
// @Success 200 {object} utils.Response{data=domain.%sResponse} "%s obtenido"
// @Failure 404 {object} utils.Response "No encontrado"
// @Failure 500 {object} utils.Response "Error interno"
// @Router /%ss/{id} [get]
// @Security BearerAuth`, moduleName, moduleName, moduleName, moduleName, moduleTitle, moduleTitle, moduleName),
		},
		{
			// Create[ModuleName]
			regexp.MustCompile(`func \(h \*` + moduleTitle + `Handler\) Create` + moduleTitle + `\(c \*gin\.Context\)`),
			fmt.Sprintf(`// @Summary Crear un %s
// @Description Crea un nuevo %s
// @Tags %ss
// @Accept json
// @Produce json
// @Param %s body domain.Create%sRequest true "Datos del %s"
// @Success 201 {object} utils.Response{data=domain.%sResponse} "%s creado"
// @Failure 400 {object} utils.Response "Datos inválidos"
// @Failure 500 {object} utils.Response "Error interno"
// @Router /%ss [post]
// @Security BearerAuth`, moduleName, moduleName, moduleName, moduleName, moduleTitle, moduleName, moduleTitle, moduleTitle, moduleName),
		},
		{
			// Update[ModuleName]
			regexp.MustCompile(`func \(h \*` + moduleTitle + `Handler\) Update` + moduleTitle + `\(c \*gin\.Context\)`),
			fmt.Sprintf(`// @Summary Actualizar un %s
// @Description Actualiza un %s existente
// @Tags %ss
// @Accept json
// @Produce json
// @Param id path string true "ID del %s"
// @Param %s body domain.Update%sRequest true "Datos a actualizar"
// @Success 200 {object} utils.Response{data=domain.%sResponse} "%s actualizado"
// @Failure 400 {object} utils.Response "Datos inválidos"
// @Failure 404 {object} utils.Response "No encontrado"
// @Failure 500 {object} utils.Response "Error interno"
// @Router /%ss/{id} [put]
// @Security BearerAuth`, moduleName, moduleName, moduleName, moduleName, moduleName, moduleTitle, moduleTitle, moduleTitle, moduleName),
		},
		{
			// Delete[ModuleName]
			regexp.MustCompile(`func \(h \*` + moduleTitle + `Handler\) Delete` + moduleTitle + `\(c \*gin\.Context\)`),
			fmt.Sprintf(`// @Summary Eliminar un %s
// @Description Elimina un %s por su ID
// @Tags %ss
// @Accept json
// @Produce json
// @Param id path string true "ID del %s"
// @Success 200 {object} utils.Response "%s eliminado"
// @Failure 404 {object} utils.Response "No encontrado"
// @Failure 500 {object} utils.Response "Error interno"
// @Router /%ss/{id} [delete]
// @Security BearerAuth`, moduleName, moduleName, moduleName, moduleName, moduleTitle, moduleName),
		},
		{
			// Archive[ModuleName]
			regexp.MustCompile(`func \(h \*` + moduleTitle + `Handler\) Archive` + moduleTitle + `\(c \*gin\.Context\)`),
			fmt.Sprintf(`// @Summary Archivar un %s
// @Description Archiva un %s por su ID
// @Tags %ss
// @Accept json
// @Produce json
// @Param id path string true "ID del %s"
// @Success 200 {object} utils.Response "%s archivado"
// @Failure 404 {object} utils.Response "No encontrado"
// @Failure 500 {object} utils.Response "Error interno"
// @Router /%ss/{id}/archive [put]
// @Security BearerAuth`, moduleName, moduleName, moduleName, moduleName, moduleTitle, moduleName),
		},
	}

	// Verificar si ya tiene comentarios Swagger
	if strings.Contains(content, "// @Summary") {
		fmt.Printf("ADVERTENCIA: El módulo %s ya parece tener comentarios Swagger\n", moduleName)
		return content, nil
	}

	// Añadir comentarios a los métodos
	newContent := content
	for _, pattern := range methodPatterns {
		// Solo añadir si encuentra el método
		if pattern.regex.MatchString(newContent) {
			newContent = pattern.regex.ReplaceAllString(newContent, pattern.comment+"\n$0")
		}
	}

	return newContent, nil
}

// documentDomain añade comentarios al dominio
func documentDomain(moduleName string) error {
	domainPath := filepath.Join("internal", moduleName, "domain", moduleName+".domain.go")

	// Verificar que el archivo exista
	if _, err := os.Stat(domainPath); os.IsNotExist(err) {
		return fmt.Errorf("el archivo %s no existe", domainPath)
	}

	// Leer el archivo
	content, err := os.ReadFile(domainPath)
	if err != nil {
		return fmt.Errorf("error al leer el archivo %s: %w", domainPath, err)
	}

	// Verificar si ya tiene comentarios Swagger
	if strings.Contains(string(content), "@Description") {
		return nil // Ya está documentado
	}

	// Convertir primera letra a mayúscula para el nombre del dominio
	moduleTitle := strings.Title(moduleName)

	// Patrones para estructuras principales
	entityRegex := regexp.MustCompile(`type ` + moduleTitle + ` struct {`)
	createReqRegex := regexp.MustCompile(`type Create` + moduleTitle + `Request struct {`)
	updateReqRegex := regexp.MustCompile(`type Update` + moduleTitle + `Request struct {`)
	responseRegex := regexp.MustCompile(`type ` + moduleTitle + `Response struct {`)

	// Comentarios a añadir
	entityComment := fmt.Sprintf("// %s representa la entidad de %s\n// @Description Entidad completa de %s", moduleTitle, moduleName, moduleName)
	createReqComment := fmt.Sprintf("// Create%sRequest representa la solicitud para crear un %s\n// @Description Datos necesarios para crear un %s", moduleTitle, moduleName, moduleName)
	updateReqComment := fmt.Sprintf("// Update%sRequest representa la solicitud para actualizar un %s\n// @Description Datos para actualizar un %s", moduleTitle, moduleName, moduleName)
	responseComment := fmt.Sprintf("// %sResponse representa la respuesta con datos de %s\n// @Description Estructura de respuesta para información de %s", moduleTitle, moduleName, moduleName)

	// Aplicar reemplazos
	contentStr := string(content)
	contentStr = entityRegex.ReplaceAllString(contentStr, entityComment+"\n$0")
	contentStr = createReqRegex.ReplaceAllString(contentStr, createReqComment+"\n$0")
	contentStr = updateReqRegex.ReplaceAllString(contentStr, updateReqComment+"\n$0")
	contentStr = responseRegex.ReplaceAllString(contentStr, responseComment+"\n$0")

	// Guardar el archivo modificado
	err = os.WriteFile(domainPath, []byte(contentStr), 0644)
	if err != nil {
		return fmt.Errorf("error al escribir el archivo %s: %w", domainPath, err)
	}

	fmt.Printf("Comentarios Swagger añadidos al dominio del módulo %s\n", moduleName)
	return nil
}

// GenerateSwaggerDocs ejecuta swag init para generar la documentación
func GenerateSwaggerDocs() error {
	// Comprobar que el archivo main.go tiene las anotaciones necesarias
	err := EnsureMainHasSwaggerAnnotations()
	if err != nil {
		return fmt.Errorf("error al verificar las anotaciones en main.go: %w", err)
	}

	// Ejecutar swag init
	fmt.Println("Ejecutando swag init para generar la documentación...")
	cmd := RunCommand("swag", "init", "--generalInfo", "main.go", "--parseDependency", "--output", "./docs")
	if cmd.Err != nil {
		return fmt.Errorf("error al ejecutar swag init: %w", cmd.Err)
	}

	fmt.Println("Documentación Swagger generada correctamente")
	fmt.Println("Accede a la documentación en: http://localhost:3000/swagger/index.html")
	return nil
}

// EnsureMainHasSwaggerAnnotations verifica y añade anotaciones Swagger en main.go
func EnsureMainHasSwaggerAnnotations() error {
	mainPath := "main.go"

	// Verificar que el archivo exista
	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		return fmt.Errorf("el archivo %s no existe", mainPath)
	}

	// Leer el archivo
	content, err := os.ReadFile(mainPath)
	if err != nil {
		return fmt.Errorf("error al leer el archivo %s: %w", mainPath, err)
	}

	// Verificar si ya tiene anotaciones Swagger
	if strings.Contains(string(content), "@title") {
		return nil // Ya tiene anotaciones
	}

	// Comentarios Swagger para añadir al inicio del archivo
	swaggerAnnotations := `// @title API Monolítica con Arquitectura de Microservicios
// @version 1.0
// @description API REST monolítica con arquitectura interna de microservicios usando Go, Gin y MongoDB
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.tu-compania.com/support
// @contact.email soporte@tu-compania.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /api
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Ingresa tu token con el formato: Bearer {token}

`

	// Encontrar la línea "package main"
	packageLineRegex := regexp.MustCompile(`(?m)^package main`)
	if !packageLineRegex.MatchString(string(content)) {
		return fmt.Errorf("no se encontró la línea 'package main' en el archivo")
	}

	// Añadir anotaciones antes de "package main"
	updatedContent := packageLineRegex.ReplaceAllString(string(content), swaggerAnnotations+"package main")

	// Verificar que ginSwagger está importado
	if !strings.Contains(updatedContent, "swaggo/gin-swagger") {
		// Buscar la sección de importaciones
		importRegex := regexp.MustCompile(`import \(([\s\S]*?)\)`)
		imports := importRegex.FindStringSubmatch(updatedContent)
		if len(imports) < 2 {
			fmt.Println("ADVERTENCIA: No se pudo encontrar el bloque de importaciones")
		} else {
			// Añadir importaciones de Swagger
			swaggerImports := `
	// Swagger
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/black4ninja/mi-proyecto/docs" // Documentación generada automáticamente`

			// Reemplazar el bloque de importaciones
			updatedImports := imports[1] + swaggerImports
			updatedContent = importRegex.ReplaceAllString(updatedContent, "import ("+updatedImports+")")
		}
	}

	// Verificar que la ruta de Swagger está configurada
	if !strings.Contains(updatedContent, "ginSwagger.WrapHandler") {
		// Buscar dónde añadir la ruta
		routerRegex := regexp.MustCompile(`router\s*:=\s*gin\.Default\(\)`)
		if routerRegex.MatchString(updatedContent) {
			// Añadir después de la inicialización del router
			swaggerRoute := `
	// Configurar Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))`

			// No reemplazar directamente para no afectar código existente
			// Solo añadir una nota al final del archivo
			updatedContent += "\n// TODO: Añadir esta línea después de la inicialización del router:" + swaggerRoute + "\n"
		}
	}

	// Guardar el archivo modificado
	err = os.WriteFile(mainPath, []byte(updatedContent), 0644)
	if err != nil {
		return fmt.Errorf("error al escribir el archivo %s: %w", mainPath, err)
	}

	fmt.Println("Anotaciones Swagger añadidas a main.go")
	return nil
}

// EnsureSwaggerDependencies asegura que las dependencias necesarias estén instaladas
func EnsureSwaggerDependencies() error {
	fmt.Println("Verificando/instalando dependencias de Swagger...")

	// Verificar/instalar swag
	_, err := exec.LookPath("swag")
	if err != nil {
		fmt.Println("Instalando swag...")
		cmd := RunCommand("go", "install", "github.com/swaggo/swag/cmd/swag@latest")
		if cmd.Err != nil {
			return fmt.Errorf("error al instalar swag: %w", cmd.Err)
		}
	}

	// Instalar bibliotecas de Go necesarias
	dependencies := []string{
		"github.com/swaggo/swag/cmd/swag@latest",
		"github.com/swaggo/gin-swagger@latest",
		"github.com/swaggo/files@latest",
	}

	for _, dep := range dependencies {
		fmt.Printf("Instalando %s...\n", dep)
		cmd := RunCommand("go", "get", "-u", dep)
		if cmd.Err != nil {
			fmt.Printf("ADVERTENCIA: Error al instalar %s: %v\n", dep, cmd.Err)
		}
	}

	fmt.Println("Dependencias instaladas/verificadas")
	return nil
}

// CommandResult contiene el resultado de ejecutar un comando
type CommandResult struct {
	Output string
	Err    error
}

// RunCommand ejecuta un comando y devuelve su salida y error
func RunCommand(name string, args ...string) CommandResult {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return CommandResult{
		Output: string(output),
		Err:    err,
	}
}
