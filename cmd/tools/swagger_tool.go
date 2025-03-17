// cmd/tools/swagger_tool.go
package main

import (
	"fmt"
	"github.com/black4ninja/mi-proyecto/pkg/tools"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "install":
		// Instalar dependencias
		if err := tools.EnsureSwaggerDependencies(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "prepare":
		// Preparar main.go para Swagger
		if err := tools.EnsureMainHasSwaggerAnnotations(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "doc-module":
		// Documentar un módulo específico
		if len(os.Args) < 3 {
			fmt.Println("Error: Falta el nombre del módulo")
			fmt.Println("Uso: go run cmd/tools/swagger_tool.go doc-module <nombre_modulo>")
			os.Exit(1)
		}

		moduleName := os.Args[2]
		if err := tools.DocumentModule(moduleName); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "doc-all":
		// Documentar todos los módulos
		if err := tools.DocumentAllModules(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "generate":
		// Generar documentación Swagger
		if err := tools.GenerateSwaggerDocs(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "all":
		// Hacer todo de una vez
		fmt.Println("=== Instalando dependencias ===")
		if err := tools.EnsureSwaggerDependencies(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n=== Preparando main.go ===")
		if err := tools.EnsureMainHasSwaggerAnnotations(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n=== Documentando todos los módulos ===")
		if err := tools.DocumentAllModules(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n=== Generando documentación Swagger ===")
		if err := tools.GenerateSwaggerDocs(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n¡Documentación Swagger completada!")
		fmt.Println("Inicia tu servidor y accede a: http://localhost:3000/swagger/index.html")

	default:
		fmt.Printf("Comando desconocido: %s\n", command)
		showHelp()
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("Herramienta de Swagger para la API")
	fmt.Println("===================================")
	fmt.Println("Uso:")
	fmt.Println("  go run cmd/tools/swagger_tool.go <comando> [args]")
	fmt.Println()
	fmt.Println("Comandos:")
	fmt.Println("  install             - Instala dependencias necesarias")
	fmt.Println("  prepare             - Prepara main.go para Swagger")
	fmt.Println("  doc-module <nombre> - Documenta un módulo específico")
	fmt.Println("  doc-all             - Documenta todos los módulos")
	fmt.Println("  generate            - Genera la documentación Swagger")
	fmt.Println("  all                 - Ejecuta todos los pasos anteriores")
	fmt.Println()
	fmt.Println("Ejemplo:")
	fmt.Println("  go run cmd/tools/swagger_tool.go all")
}
