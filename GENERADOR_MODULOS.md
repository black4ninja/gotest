# Implementación del Generador de Módulos

Este documento explica paso a paso cómo implementar un generador de módulos para tu proyecto Go con arquitectura limpia.

## Estructura de Archivos

El generador consta de dos archivos principales:
1. `pkg/tools/module_generator.go`: Contiene la lógica y plantillas para generar módulos
2. `cmd/tools/generate_module.go`: Proporciona la interfaz de línea de comandos

## Pasos de Implementación

### 1. Configurar Permisos de Ejecución (Opcional en sistemas Unix)

```bash
chmod +x cmd/tools/generate_module.go
```

## Uso del Generador de Módulos

### Generar un Nuevo Módulo

Para crear un nuevo módulo llamado "category", ejecuta:

```bash
go run cmd/tools/generate_module.go category
```

Esto generará la siguiente estructura:

```
internal/category/
├── domain/
│   └── category.domain.go
├── repository/
│   └── mongo.category.repository.go
├── usecase/
│   └── category.usecase.go
├── delivery/
│   └── category.delivery.go
└── main_fragment.go.txt
```

### Personalizar el Módulo Generado

1. Revisa cada archivo generado y personalízalo según tus necesidades
2. Añade o modifica campos en la estructura de entidad
3. Ajusta las validaciones y la lógica de negocio

### Integrar el Módulo en la Aplicación

1. Copia el contenido de `internal/category/main_fragment.go.txt` en tu archivo `main.go`
2. Ajusta las importaciones según tu estructura de proyecto
3. Configura permisos en los scripts de inicialización si usas un sistema RBAC

## Personalización del Generador

Puedes personalizar las plantillas en `module_generator.go` para adaptarlas a tus necesidades específicas:

1. Modifica las estructuras de datos para incluir campos adicionales
2. Añade lógica de validación personalizada
3. Incluye integración con otros componentes de tu arquitectura
4. Añade generación automática de pruebas unitarias

## Ejemplo de Uso Completo

```bash
# Generar un módulo de productos
go run cmd/tools/generate_module.go product

# Personalizar entidad en domain/product.domain.go
# Añadir campos como SKU, precio, inventario, etc.

# Integrar el módulo en main.go
# Usar el fragmento generado

# Crear registros de permisos
# Añadir al script de inicialización de permisos
```

## Solución de Problemas

### Error: "template: file:1: unexpected..."

Verifica que las comillas y backticks en las plantillas estén correctamente escapados. Los caracteres `` ` `` utilizados para etiquetas de struct en Go pueden causar problemas en las plantillas.

### Error: "package ... not found"

Asegúrate de que las rutas de importación en `generate_module.go` coincidan con la estructura de tu proyecto y módulo Go.

### Otros Errores

Si enfrentas otros problemas, asegúrate de:
- Tener permisos de escritura en los directorios
- No tener conflictos con archivos existentes
- Usar la versión correcta de Go (1.16+ recomendado)