# Checklist para Crear un Nuevo Módulo

- Para una implementación paso a paso ver [EJEMPLO MÓDULO](EJEMPLO_MODULO)
- Para crear un módulo de manera automática se puede usar [GENERADOR MÓDULOS](GENERADOR_MODULOS)

## 1. Planificación y Diseño
- [ ] Definir el propósito y responsabilidades del módulo
- [ ] Identificar las entidades principales y sus relaciones
- [ ] Definir los casos de uso y funcionalidades (CRUD básico, lógica específica)
- [ ] Planificar los endpoints de API necesarios
- [ ] Identificar dependencias con otros módulos

## 2. Crear la Estructura de Carpetas

```
internal/
└── mi_modulo/              # Nombre del nuevo módulo
    ├── domain/             # Definiciones y contratos
    ├── repository/         # Implementación de persistencia
    ├── usecase/            # Implementación de lógica de negocio
    └── delivery/           # Controladores HTTP/Handlers
```

## 3. Implementar el Dominio (`domain/`)

- [ ] Crear el archivo principal `mi_modulo.domain.go`
    - [ ] Definir las estructuras/entidades principales
    - [ ] Definir constantes y tipos
    - [ ] Definir interfaces de repositorio (contratos de persistencia)
    - [ ] Definir interfaces de casos de uso (contratos de lógica de negocio)
    - [ ] Definir estructuras de solicitud/respuesta

## 4. Implementar el Repositorio (`repository/`)

- [ ] Crear el archivo `mongo.mi_modulo.repository.go`
    - [ ] Implementar la estructura del repositorio
    - [ ] Implementar el constructor `NewMongoMiModuloRepository`
    - [ ] Implementar todos los métodos definidos en la interfaz del dominio
    - [ ] Manejar errores y conversiones de tipos (MongoDB → Dominio)
    - [ ] Implementar validaciones necesarias

## 5. Implementar los Casos de Uso (`usecase/`)

- [ ] Crear el archivo `mi_modulo.usecase.go`
    - [ ] Implementar la estructura del caso de uso
    - [ ] Implementar el constructor `NewMiModuloUseCase`
    - [ ] Implementar todos los métodos definidos en la interfaz del dominio
    - [ ] Aplicar lógica de negocio y validaciones
    - [ ] Manejar errores y respuestas

## 6. Implementar el Delivery (Handlers HTTP) (`delivery/`)

- [ ] Crear el archivo `mi_modulo.delivery.go`
    - [ ] Implementar la estructura del handler
    - [ ] Implementar el constructor `NewMiModuloHandler`
    - [ ] Definir las rutas en el constructor
    - [ ] Implementar manejadores para cada ruta
    - [ ] Validar entradas y controlar errores
    - [ ] Convertir respuestas al formato estándar

## 7. Integrar con el Sistema Principal

- [ ] Inicializar el Repositorio en `main.go`
  ```go
  miModuloCollection := mongoClient.Database(mongoDBName).Collection("mi_modulo")
  miModuloRepository := repository.NewMongoMiModuloRepository(miModuloCollection)
  ```

- [ ] Inicializar el Caso de Uso en `main.go`
  ```go
  miModuloService := usecase.NewMiModuloUseCase(miModuloRepository)
  ```

- [ ] Configurar rutas y middleware en `main.go`
  ```go
  miModuloRoutes := api.Group("/mi-modulo")
  miModuloRoutes.Use(permissionMiddleware.RequirePermission("mi_modulo:access"))
  delivery.NewMiModuloHandler(miModuloRoutes, miModuloService)
  ```

## 8. Implementar Pruebas

- [ ] Pruebas unitarias para el repositorio
- [ ] Pruebas unitarias para los casos de uso
- [ ] Pruebas para los handlers HTTP
- [ ] Pruebas de integración (opcional)

## 9. Documentación y Finalización

- [ ] Documentar los endpoints de API (formatos de solicitud/respuesta)
- [ ] Actualizar el README del proyecto con la nueva funcionalidad
- [ ] Revisar el código para asegurar que sigue los estándares del proyecto
- [ ] Asegurar que todos los tests pasen