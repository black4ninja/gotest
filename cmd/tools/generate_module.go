// cmd/tools/generate_module.go
package main

import (
	"fmt"
	"os"

	"github.com/black4ninja/mi-proyecto/pkg/tools"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run cmd/tools/generate_module.go <nombre_del_modulo>")
		os.Exit(1)
	}

	moduleName := os.Args[1]

	fmt.Printf("Generando módulo: %s\n", moduleName)

	if err := tools.GenerateModule(moduleName); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Módulo generado exitosamente.")
	fmt.Println("Revise los archivos generados y personalícelos según sus necesidades.")
}
