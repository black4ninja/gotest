# Recomendaciones para Trabajar con esta Arquitectura

## Buenas Prácticas

### 1. Definición Clara de Responsabilidades

- **Domain**: Define qué hace el módulo, no cómo lo hace.
    - Mantén las interfaces claras y enfocadas en un propósito.
    - No introduzcas dependencias de frameworks o tecnologías específicas.

- **Repository**: Implementa el acceso a datos sin lógica de negocio.
    - Preocúpate solo por cómo guardar y recuperar datos.
    - Maneja la conversión entre estructuras de datos y entidades de dominio.

- **UseCase**: Implementa la lógica de negocio pura.
    - Aquí van las validaciones, reglas de negocio y flujos de trabajo.
    - Coordina las operaciones entre múltiples repositorios si es necesario.

- **Delivery**: Gestiona la comunicación con el cliente.
    - Traduce las solicitudes HTTP a llamadas de casos de uso.
    - Formatea las respuestas según las necesidades del cliente.

### 2. Gestión de Errores

- Utiliza tipos de error personalizados para situaciones específicas del dominio.
- Mantén mensajes de error consistentes y útiles para el usuario final.
- En el delivery, traduce los errores a códigos HTTP apropiados.

### 3. Validaciones

- Las validaciones básicas de formato deben estar en las estructuras de request (binding).
- Las validaciones de negocio deben estar en los casos de uso.
- No dupliques validaciones entre capas.

### 4. Testing

- Prueba cada capa de forma independiente:
    - **Repository**: Pruebas contra una BD de prueba o con mocks.
    - **UseCase**: Pruebas unitarias con repositorios mockeados.
    - **Delivery**: Pruebas de integración con endpoints completos.

### 5. Documentación

- Documenta el propósito de cada módulo y sus interfaces principales.
- Usa comentarios para explicar decisiones de diseño no obvias.
- Mantén ejemplos de uso en la documentación.

## Tips para el Desarrollo

### 1. Trabajo con Air

```bash
# Iniciar el servidor con recarga automática
air

# Ejecutar con opciones específicas
air -c .air.toml
```

- Si Air no detecta cambios correctamente, puedes forzar un reinicio tocando el archivo main.go.
- Considera configurar diferentes archivos `.air.toml` para diferentes entornos.

### 2. Depuración

- Usa el paquete `log` para depuración básica durante el desarrollo.
- Para depuración más detallada, considera usar Delve:
  ```bash
  go install github.com/go-delve/delve/cmd/dlv@latest
  dlv debug main.go
  ```
- Implementa logs estructurados para entornos de producción.

### 3. Migraciones de Base de Datos

Para proyectos más complejos, considera implementar un sistema de migraciones:

```go
// En pkg/database/migrations.go
func RunMigrations(db *mongo.Database) error {
    // Implementación de migraciones
}
```

### 4. Scripts Útiles

Algunos scripts que podrías necesitar:

- **Generación de módulos**: Un script que genere la estructura básica de un nuevo módulo.
- **Validación de código**: Un script que ejecute linters y herramientas de análisis estático.
- **Generación de documentación**: Un script para generar documentación a partir de comentarios.

## Evolución de la Arquitectura

A medida que el proyecto crezca, considera:

### 1. Comunicación entre Módulos

- Define bien las dependencias entre módulos.
- Considera usar un patrón mediador o eventos para comunicación desacoplada.
- Evita referencias circulares entre módulos.

### 2. Servicios Compartidos

- Identifica funcionalidades transversales (logging, métricas, notificaciones).
- Impleméntalas como servicios independientes que pueden ser inyectados.

### 3. Transición a Microservicios

Si eventualmente necesitas migrar a microservicios:

- La separación por módulos facilitará la transición.
- Considera empezar con una arquitectura de monolito modular (como la actual).
- Luego puedes extraer módulos individuales como servicios independientes.
- Implementa API Gateways y Service Discovery cuando sea necesario.

## Solución de Problemas Comunes

### Error: "panic: runtime error: invalid memory address or nil pointer dereference"

- Verifica que todas las dependencias necesarias estén inicializadas correctamente.
- Asegúrate de que las inyecciones de dependencias sean correctas.
- Comprueba que las conexiones a bases de datos estén establecidas.

### Error: "context deadline exceeded"

- Revisa los timeouts de las operaciones de MongoDB.
- Asegúrate de que las operaciones de base de datos no sean demasiado pesadas.
- Verifica la conectividad y latencia de la red.

### Error: "cannot unmarshal type"

- Verifica las estructuras de datos para serialización/deserialización.
- Asegúrate de que los campos JSON coincidan con las estructuras Go.
- Comprueba los tipos de datos (string vs int, etc.).

### Error: "Cannot use 'X' (type primitive.ObjectID) as the type string"

- Este error ocurre cuando intentas usar un ObjectID de MongoDB directamente como string.
- La solución es usar el método `.Hex()` para convertir el ID a su representación en string:
  ```go
  // Incorrecto
  userID := user.ID           // Tipo: primitive.ObjectID

  // Correcto
  userID := user.ID.Hex()     // Tipo: string
  ```
- Este error es común cuando trabajas con IDs entre diferentes módulos, ya que algunos esperan strings y otros trabajan con primitive.ObjectID.

### Error: "parsing time: extra text"

- Ocurre cuando intentas parsear una fecha con formato incorrecto.
- Verifica el formato de las fechas que estás procesando.
- Usa el formato correcto de time.Parse:
  ```go
  // Formato RFC3339
  time.Parse(time.RFC3339, "2023-01-01T12:00:00Z")
  
  // Otros formatos personalizados
  time.Parse("2006-01-02", "2023-01-01")
  ```

## Optimización de Rendimiento

### 1. Índices en MongoDB

Crea índices para mejorar el rendimiento de las consultas:

```go
// En una función de inicialización
func createIndices(collection *mongo.Collection) {
    // Índice simple
    collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
        Keys: bson.D{{Key: "email", Value: 1}},
        Options: options.Index().SetUnique(true),
    })
    
    // Índice compuesto
    collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
        Keys: bson.D{{Key: "status", Value: 1}, {Key: "created_at", Value: -1}},
    })
}
```

### 2. Paginación

Implementa paginación para mejorar el rendimiento con grandes conjuntos de datos:

```go
func (r *mongoRepository) GetAll(page, limit int, filters bson.M) ([]*domain.Entity, int64, error) {
    // Configurar opciones de paginación
    opts := options.Find().
        SetSkip(int64((page - 1) * limit)).
        SetLimit(int64(limit)).
        SetSort(bson.M{"created_at": -1})
        
    // Ejecutar consulta paginada
    cursor, err := r.collection.Find(ctx, filters, opts)
    
    // Contar total para metadata de paginación
    total, err := r.collection.CountDocuments(ctx, filters)
    
    // Resto de la lógica...
}
```

### 3. Consultas Eficientes

- Usa proyecciones para seleccionar solo los campos necesarios:
  ```go
  opts := options.FindOne().SetProjection(bson.M{
      "name": 1,
      "email": 1,
      "_id": 1,
  })
  ```

- Considera la carga diferida (lazy loading) para relaciones complejas.
- Utiliza agregaciones para operaciones complejas en lugar de procesarlas en memoria.

## Seguridad

### 1. Validación de Entradas

- Usa validadores de Gin (binding tags) para la primera línea de defensa.
- Sanitiza todas las entradas para prevenir inyecciones NoSQL.
- Implementa límites de tamaño para prevenir ataques DoS.

### 2. Autenticación y Autorización

- Implementa un sistema de caducidad de tokens adecuado.
- Considera usar refresh tokens con rotación.
- Implementa bloqueo de cuentas después de intentos fallidos.
- Utiliza almacenamiento seguro de contraseñas (bcrypt).

### 3. HTTPS y Headers de Seguridad

- Configura correctamente TLS.
- Implementa headers de seguridad como:
    - Content-Security-Policy
    - X-Content-Type-Options
    - X-Frame-Options

## Despliegue

### 1. Contenedorización

Crea un Dockerfile optimizado:

```Dockerfile
# Etapa de compilación
FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o main .

# Etapa final
FROM alpine:3.15
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY .env .
EXPOSE 3000
CMD ["./main"]
```

### 2. Configuración por Entorno

- Usa archivos `.env` diferentes por entorno.
- Considera usar un servicio de configuración centralizado para entornos complejos.
- Implementa una estrategia de secretos para credenciales sensibles.

### 3. Monitorización

- Implementa health checks:
  ```go
  router.GET("/health", func(c *gin.Context) {
      if mongoClient.Ping(ctx, nil) == nil {
          c.JSON(200, gin.H{"status": "ok"})
      } else {
          c.JSON(500, gin.H{"status": "error"})
      }
  })
  ```

- Considera integrar métricas con Prometheus.
- Implementa logging estructurado para facilitar el análisis.