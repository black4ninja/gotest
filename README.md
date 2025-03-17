# API Monolítica con Arquitectura de Microservicios

Este proyecto implementa una API REST monolítica con una arquitectura interna que simula microservicios, utilizando Go, Gin y MongoDB. El diseño sigue los principios de Clean Architecture y está estructurado para facilitar el mantenimiento y el crecimiento.

## Características

- **Autenticación OAuth 2.0**: Sistema completo de autenticación con soporte para diferentes tipos de grants
- **Gestión de permisos RBAC**: Sistema de roles y permisos detallado
- **Arquitectura limpia**: Separación de responsabilidades (dominio, casos de uso, repositorios, delivery)
- **Desarrollo rápido**: Configuración para desarrollo en tiempo real con Air
- **MongoDB**: Persistencia basada en MongoDB
- **Middleware de seguridad**: Protección de rutas y validación de permisos

## Estructura del Proyecto

```
proyecto-raiz/
├── main.go                                 # Punto de entrada principal
├── go.mod
├── .air.toml                               # Configuración para desarrollo en tiempo real
├── .env                                    # Variables de entorno
├── internal/                               # Código de aplicación organizado por dominios
│   ├── user/                               # Módulo de usuarios
│   │   ├── domain/                         # Definiciones y contratos
│   │   ├── repository/                     # Implementación de persistencia
│   │   ├── usecase/                        # Implementación de lógica de negocio
│   │   └── delivery/                       # Controladores HTTP/Handlers
│   ├── oauth/                              # Módulo de OAuth 2.0
│   │   ├── domain/
│   │   ├── repository/
│   │   ├── usecase/
│   │   └── delivery/
│   └── permission/                         # Módulo de permisos y roles
│       ├── domain/
│       ├── repository/
│       ├── usecase/
│       └── delivery/
├── pkg/                                    # Código compartido entre módulos
│   ├── config/                             # Configuraciones
│   ├── middleware/                         # Middlewares
│   └── utils/                              # Utilidades
└── scripts/                                # Scripts de utilidad
    ├── init_oauth_client.go                # Crea cliente OAuth inicial
    └── init_permissions_and_admin.go       # Crea permisos y admin inicial
```

## Arquitectura

Cada módulo sigue la siguiente estructura interna basada en Clean Architecture:

1. **Domain**: Entidades de negocio e interfaces para casos de uso y repositorios
2. **Repository**: Implementaciones de persistencia (adaptadores MongoDB)
3. **UseCase**: Implementación de lógica de negocio
4. **Delivery**: Controladores HTTP (utilizando Gin)

## Requisitos

- Go 1.18 o superior
- MongoDB 4.4 o superior
- Air (para desarrollo en tiempo real)

## Instalación

1. Clonar el repositorio:
   ```bash
   git clone https://github.com/usuario/mi-proyecto.git
   cd mi-proyecto
   ```

2. Instalar dependencias:
   ```bash
   go mod download
   ```

3. Instalar Air para desarrollo en tiempo real:
   ```bash
   go install github.com/cosmtrek/air@latest
   ```

4. Configurar variables de entorno (copia el archivo .env.example):
   ```bash
   cp .env.example .env
   # Edita .env con tu configuración
   ```

## Configuración

Edita el archivo `.env` para configurar:

```
# Servidor
PORT=3000
ENV=development

# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DB=my_database

# OAuth
JWT_SECRET=your_secret_key_here
TOKEN_EXP=7200  # Tiempo de expiración del token en segundos

# Admin predeterminado (para scripts de inicialización)
DEFAULT_ADMIN_EMAIL=admin@ejemplo.com
DEFAULT_ADMIN_PASSWORD=adminPass123!
```

## Inicialización de Datos

Antes de utilizar el sistema, puedes ejecutar los scripts de inicialización:

1. Inicializar cliente OAuth:
   ```bash
   go run scripts/init_oauth_client.go
   ```

2. Inicializar permisos y usuario administrador:
   ```bash
   go run scripts/init_permissions_and_admin.go
   ```

## Ejecución

### Desarrollo (con Air para recarga en tiempo real)

```bash
air
```

### Producción

```bash
go build -o main .
./main
```

## API Endpoints

### Autenticación (OAuth 2.0)

- **POST /api/oauth/token**: Genera un token de acceso
    - Grant types: `password`, `client_credentials`, `refresh_token`
- **POST /api/oauth/revoke**: Revoca un token de acceso

### Usuarios

- **GET /api/users**: Lista todos los usuarios (protegido)
- **GET /api/users/:id**: Obtiene un usuario por su ID (protegido)
- **POST /api/users**: Crea un nuevo usuario (protegido)
- **PUT /api/users/:id**: Actualiza un usuario existente (protegido)
- **DELETE /api/users/:id**: Elimina un usuario (protegido)
- **PUT /api/users/:id/archive**: Archiva un usuario (protegido)

### Permisos y Roles

- **GET /api/permissions/permissions**: Lista todos los permisos (protegido)
- **GET /api/permissions/roles**: Lista todos los roles (protegido)
- **POST /api/permissions/roles/:id/permissions**: Asigna un permiso a un rol (protegido)
- **POST /api/permissions/user-roles/assign-role**: Asigna un rol a un usuario (protegido)

## Creación de un Nuevo Módulo

Para crear un nuevo módulo, sigue el checklist proporcionado en el archivo [NUEVO_MODULO.md](./NUEVO_MODULO.md).

## Desarrollo con Air

Air permite la recarga automática de tu aplicación cuando detecta cambios en los archivos. El archivo `.air.toml` ya está configurado para observar cambios en archivos `.go`, `.tpl`, `.tmpl` y `.html`.

Para modificar qué archivos observa o ignora, edita estas secciones en `.air.toml`:

```toml
[build]
# Archivos a observar para cambios
include_ext = ["go", "tpl", "tmpl", "html"]
# Archivos a ignorar
exclude_dir = ["assets", "tmp", "vendor"]
```

## Pruebas

Ejecutar todas las pruebas:

```bash
go test ./...
```

## Licencia

[MIT](LICENSE)