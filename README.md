# Proyecto base

## Estructura del proyecto
```
proyecto-raiz/
├── main.go                                 # Punto de entrada principal
├── go.mod
├── .env
├── internal/
│   ├── user/                               # Módulo de usuarios
│   │   ├── domain/
│   │   │   └── user.go                     # Entidad y contratos de usuario
│   │   ├── repository/
│   │   │   └── mongo_repository.go         # Repositorio de usuarios con MongoDB
│   │   ├── usecase/
│   │   │   └── user_usecase.go             # Casos de uso de usuarios
│   │   └── delivery/
│   │       └── http_handler.go             # Handlers HTTP para usuarios
│   └── oauth/                              # Módulo de OAuth 2.0
│       ├── domain/
│       │   ├── oauth.go                    # Entidades y contratos de OAuth
│       │   ├── client.go                   # Entidad de cliente OAuth
│       │   └── token.go                    # Entidad de token OAuth
│       ├── repository/
│       │   ├── mongo_client_repository.go  # Repositorio de clientes OAuth
│       │   └── mongo_token_repository.go   # Repositorio de tokens OAuth
│       ├── usecase/
│       │   └── oauth_usecase.go            # Casos de uso de OAuth
│       └── delivery/
│           └── http_handler.go             # Handlers HTTP para OAuth
└── pkg/
    ├── config/                             # Configuraciones
    │   └── mongo.go                        # Configuración de MongoDB
    ├── middleware/                         # Middlewares
    │   └── oauth_middleware.go             # Middleware de OAuth
    └── utils/                              # Utilidades
        ├── crypto.go                       # Funciones de criptografía
        ├── jwt.go                          # Funciones para JWT
        └── response.go                     # Utilidades para respuestas HTTP
```