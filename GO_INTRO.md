# Guía de Go para Desarrolladores de Node.js

Esta guía está diseñada para ayudar a desarrolladores familiarizados con Node.js a entender la estructura, sintaxis y patrones de nuestra arquitectura en Go. Nos enfocaremos en las diferencias clave y proporcionar ejemplos comparativos para facilitar la transición.

## Tabla de Contenidos

1. [Conceptos Básicos de Go vs. Node.js](#conceptos-básicos-de-go-vs-nodejs)
2. [Estructura del Proyecto](#estructura-del-proyecto)
3. [Patrones Comunes en Go](#patrones-comunes-en-go)
4. [Trabajando con MongoDB](#trabajando-con-mongodb)
5. [Manejo de Rutas HTTP en Gin vs. Express](#manejo-de-rutas-http-en-gin-vs-express)
6. [Control de Errores](#control-de-errores)
7. [Patrones de Arquitectura Limpia](#patrones-de-arquitectura-limpia)
8. [Glosario de Términos](#glosario-de-términos)

## Conceptos Básicos de Go vs. Node.js

### Declaración de Variables

**Node.js (JavaScript):**
```javascript
// Variables
let counter = 0;
const name = "John";
var oldStyle = true;

// Objetos
const user = {
    id: 1,
    name: "John",
    roles: ["admin", "editor"]
};
```

**Go:**
```go
// Variables
var counter int = 0     // Declaración explícita de tipo
counter := 0            // Inferencia de tipo con :=
name := "John"          // String
oldStyle := true        // Boolean

// Structs (similar a objetos)
type User struct {
  ID    int      `json:"id"`     // Las etiquetas definen cómo se serializa
  Name  string   `json:"name"`
  Roles []string `json:"roles"`  // Array/slice
}

user := User{
  ID:    1,
  Name:  "John",
  Roles: []string{"admin", "editor"},
}
```

### Funciones

**Node.js:**
```javascript
// Función básica
function add(a, b) {
  return a + b;
}

// Arrow function
const multiply = (a, b) => a * b;

// Método de objeto
const calculator = {
  add: function(a, b) {
    return a + b;
  }
};

// Async/await
async function fetchData() {
  try {
    const response = await fetch('https://api.example.com/data');
    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Error:', error);
    throw error;
  }
}
```

**Go:**
```go
// Función básica
func add(a, b int) int {
  return a + b
}

// Función con múltiples retornos (común en Go)
func divide(a, b float64) (float64, error) {
  if b == 0 {
    return 0, errors.New("división por cero")
  }
  return a / b, nil
}

// Método (función asociada a un tipo)
type Calculator struct{}

func (c Calculator) Add(a, b int) int {
  return a + b
}

// Manejo asíncrono con goroutines (diferente al async/await)
func fetchData() ([]byte, error) {
  resp, err := http.Get("https://api.example.com/data")
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()  // Se ejecuta al final de la función
  
  return ioutil.ReadAll(resp.Body)
}
```

### Tipos de Datos

| JavaScript | Go | Notas |
|------------|----|----|
| Number | int, int8, int16, int32, int64, float32, float64 | Go distingue entre varios tipos numéricos |
| String | string | Similar |
| Boolean | bool | Similar pero con nombre diferente |
| Object | struct | Los structs tienen campos tipados y definidos previamente |
| Array, Array.map | slice, for loops | Go no tiene métodos funcionales como map o filter |
| null/undefined | nil | nil es similar pero se usa solo con tipos específicos |
| Promise | Channels, goroutines | Concurrencia diferente a promesas |

## Estructura del Proyecto

Nuestra aplicación sigue una arquitectura limpia (Clean Architecture) organizada por dominio:

```
proyecto-raiz/
├── internal/          # Código específico de la aplicación
│   ├── user/          # Módulo de usuario (dominio)
│   │   ├── domain/    # Definiciones y contratos
│   │   ├── repository/# Implementación de persistencia
│   │   ├── usecase/   # Lógica de negocio
│   │   └── delivery/  # Controladores HTTP
│   └── product/       # Otro módulo de dominio
├── pkg/               # Código compartido y utilidades
├── cmd/               # Entradas de la aplicación
└── main.go            # Punto de entrada
```

**Comparación con Node.js:**

En Node.js (Express) puedes tener:
```
proyecto-node/
├── routes/            # Rutas HTTP
├── controllers/       # Controladores
├── models/            # Modelos (Mongoose)
├── services/          # Servicios/lógica
└── app.js             # Punto de entrada
```

La diferencia principal es que en nuestra estructura Go, organizamos por módulos/dominios (como "user", "product") en lugar de por capas técnicas (como "routes", "controllers").

## Patrones Comunes en Go

### Interfaces

Las interfaces en Go son diferentes a las de TypeScript o a las clases en JavaScript. En Go, las interfaces se implementan implícitamente:

```go
// Definición de una interfaz
type Repository interface {
  GetByID(id string) (*Entity, error)
  Save(entity *Entity) error
}

// Implementación implícita (no se declara explícitamente)
type MongoRepository struct {
  collection *mongo.Collection
}

// Si MongoRepository implementa todos los métodos de Repository,
// automáticamente satisface la interfaz
func (r *MongoRepository) GetByID(id string) (*Entity, error) {
  // implementación...
}

func (r *MongoRepository) Save(entity *Entity) error {
  // implementación...
}
```

Esto es diferente de TypeScript/JavaScript donde las interfaces o clases se implementan explícitamente:

```typescript
interface Repository {
  getById(id: string): Promise<Entity>;
  save(entity: Entity): Promise<void>;
}

// Implementación explícita
class MongoRepository implements Repository {
  constructor(private collection: Collection) {}
  
  async getById(id: string): Promise<Entity> {
    // implementación...
  }
  
  async save(entity: Entity): Promise<void> {
    // implementación...
  }
}
```

### Inyección de Dependencias

En Go no tenemos decoradores ni inyección automática. La inyección se hace manualmente en los constructores:

```go
// Constructor de repositorio
func NewMongoRepository(collection *mongo.Collection) domain.Repository {
  return &mongoRepository{
    collection: collection,
  }
}

// Constructor de caso de uso
func NewUseCase(repo domain.Repository) domain.UseCase {
  return &useCase{
    repo: repo,
  }
}

// En main.go - Ensamblaje manual
collection := mongoClient.Database("mydb").Collection("entities")
repository := repository.NewMongoRepository(collection)
useCase := usecase.NewUseCase(repository)
handler := delivery.NewHandler(useCase)
```

Esto contrasta con Node.js donde podrías usar decoradores (con NestJS) o un contenedor DI:

```typescript
// Con NestJS
@Injectable()
export class EntityService {
  constructor(
    @InjectRepository(Entity)
    private entityRepository: Repository<Entity>,
  ) {}
}
```

## Trabajando con MongoDB

### Inicialización y Conexión

**Node.js con Mongoose:**
```javascript
const mongoose = require('mongoose');

mongoose.connect('mongodb://localhost:27017/mydatabase', {
  useNewUrlParser: true,
  useUnifiedTopology: true
})
.then(() => console.log('Connected to MongoDB'))
.catch(err => console.error('Error connecting to MongoDB', err));

// Definición de modelo
const userSchema = new mongoose.Schema({
  name: String,
  email: { type: String, required: true, unique: true },
  age: Number,
  createdAt: { type: Date, default: Date.now }
});

const User = mongoose.model('User', userSchema);
```

**Go con driver oficial:**
```go
import (
  "context"
  "log"
  "time"
  
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

func connectToMongoDB() (*mongo.Client, error) {
  // Crear contexto con timeout
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()
  
  // Conectar al servidor
  client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
  if err != nil {
    return nil, err
  }
  
  // Verificar la conexión
  err = client.Ping(ctx, nil)
  if err != nil {
    return nil, err
  }
  
  log.Println("Conectado a MongoDB")
  return client, nil
}

// No hay esquemas explícitos como en Mongoose
// Usamos structs para modelar documentos
type User struct {
  ID        primitive.ObjectID `bson:"_id,omitempty"`
  Name      string             `bson:"name"`
  Email     string             `bson:"email"`
  Age       int                `bson:"age"`
  CreatedAt time.Time          `bson:"created_at"`
}
```

### Operaciones Comunes con MongoDB

#### Insertar Documentos

**Node.js:**
```javascript
// Insertar un documento
const newUser = new User({
  name: 'John Doe',
  email: 'john@example.com',
  age: 30
});

await newUser.save();

// Insertar múltiples documentos
await User.insertMany([
  { name: 'Alice', email: 'alice@example.com', age: 25 },
  { name: 'Bob', email: 'bob@example.com', age: 28 }
]);
```

**Go:**
```go
// Insertar un documento
user := User{
  ID:        primitive.NewObjectID(),
  Name:      "John Doe",
  Email:     "john@example.com",
  Age:       30,
  CreatedAt: time.Now(),
}

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

collection := client.Database("mydatabase").Collection("users")
_, err := collection.InsertOne(ctx, user)
if err != nil {
  log.Fatal(err)
}

// Insertar múltiples documentos
users := []interface{}{
  User{
    ID:        primitive.NewObjectID(),
    Name:      "Alice",
    Email:     "alice@example.com",
    Age:       25,
    CreatedAt: time.Now(),
  },
  User{
    ID:        primitive.NewObjectID(),
    Name:      "Bob",
    Email:     "bob@example.com",
    Age:       28,
    CreatedAt: time.Now(),
  },
}

_, err = collection.InsertMany(ctx, users)
```

#### Consultar Documentos

**Node.js:**
```javascript
// Encontrar por ID
const user = await User.findById('60d5ec9af682fbd12a8924c4');

// Encontrar uno
const admin = await User.findOne({ role: 'admin' });

// Encontrar muchos
const adultUsers = await User.find({ age: { $gte: 18 } }).sort({ name: 1 }).limit(10);

// Consulta con proyección
const userEmails = await User.find({}, 'name email').exec();
```

**Go:**
```go
// Encontrar por ID
id, _ := primitive.ObjectIDFromHex("60d5ec9af682fbd12a8924c4")
var user User
err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)

// Encontrar uno
var admin User
err = collection.FindOne(ctx, bson.M{"role": "admin"}).Decode(&admin)

// Encontrar muchos
opts := options.Find().SetSort(bson.M{"name": 1}).SetLimit(10)
cursor, err := collection.Find(ctx, bson.M{"age": bson.M{"$gte": 18}}, opts)
if err != nil {
  log.Fatal(err)
}
defer cursor.Close(ctx)

var adultUsers []User
if err = cursor.All(ctx, &adultUsers); err != nil {
  log.Fatal(err)
}

// Consulta con proyección
opts = options.Find().SetProjection(bson.M{"name": 1, "email": 1})
cursor, err = collection.Find(ctx, bson.M{}, opts)
// ...procesar cursor
```

#### Actualizar Documentos

**Node.js:**
```javascript
// Actualizar uno
await User.updateOne(
  { _id: '60d5ec9af682fbd12a8924c4' },
  { $set: { name: 'John Smith', updatedAt: new Date() } }
);

// Actualizar muchos
await User.updateMany(
  { age: { $lt: 18 } },
  { $set: { status: 'minor' } }
);

// Encontrar y actualizar
const updatedUser = await User.findByIdAndUpdate(
  '60d5ec9af682fbd12a8924c4',
  { $set: { status: 'active' } },
  { new: true } // Devuelve el documento actualizado
);
```

**Go:**
```go
// Actualizar uno
id, _ := primitive.ObjectIDFromHex("60d5ec9af682fbd12a8924c4")
update := bson.M{
  "$set": bson.M{
    "name":       "John Smith",
    "updated_at": time.Now(),
  },
}
_, err = collection.UpdateOne(ctx, bson.M{"_id": id}, update)

// Actualizar muchos
_, err = collection.UpdateMany(
  ctx,
  bson.M{"age": bson.M{"$lt": 18}},
  bson.M{"$set": bson.M{"status": "minor"}},
)

// Encontrar y actualizar (no hay método directo como en Mongoose)
opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
var updatedUser User
err = collection.FindOneAndUpdate(
  ctx,
  bson.M{"_id": id},
  bson.M{"$set": bson.M{"status": "active"}},
  opts,
).Decode(&updatedUser)
```

## Manejo de Rutas HTTP en Gin vs. Express

### Definición de Rutas

**Express (Node.js):**
```javascript
const express = require('express');
const router = express.Router();

// Middleware
const authMiddleware = require('./middleware/auth');

// Rutas
router.get('/users', authMiddleware, async (req, res) => {
  try {
    const users = await UserModel.find({});
    res.json(users);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

router.post('/users', authMiddleware, async (req, res) => {
  try {
    const newUser = new UserModel(req.body);
    const savedUser = await newUser.save();
    res.status(201).json(savedUser);
  } catch (err) {
    res.status(400).json({ error: err.message });
  }
});

// Parámetros en ruta
router.get('/users/:id', authMiddleware, async (req, res) => {
  try {
    const user = await UserModel.findById(req.params.id);
    if (!user) return res.status(404).json({ error: 'User not found' });
    res.json(user);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Exportar router
module.exports = router;
```

**Gin (Go):**
```go
package delivery

import (
  "net/http"

  "github.com/gin-gonic/gin"
  
  "myproject/internal/user/domain"
  "myproject/pkg/middleware"
  "myproject/pkg/utils"
)

type UserHandler struct {
  userUseCase domain.UserUseCase
}

// Constructor - equivalente a exportar el router
func NewUserHandler(router *gin.RouterGroup, useCase domain.UserUseCase) {
  handler := &UserHandler{
    userUseCase: useCase,
  }
  
  // Aplicar middleware
  userRoutes := router.Group("/users")
  userRoutes.Use(middleware.Auth())
  
  // Definir rutas
  userRoutes.GET("/", handler.GetAllUsers)
  userRoutes.POST("/", handler.CreateUser)
  userRoutes.GET("/:id", handler.GetUser)
}

// Métodos de handler
func (h *UserHandler) GetAllUsers(c *gin.Context) {
  users, err := h.userUseCase.GetAllUsers()
  if err != nil {
    utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
    return
  }
  
  utils.SuccessResponse(c, http.StatusOK, "Usuarios obtenidos con éxito", users)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
  var req domain.CreateUserRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    utils.ValidationErrorResponse(c, err.Error())
    return
  }
  
  user, err := h.userUseCase.CreateUser(&req)
  if err != nil {
    utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
    return
  }
  
  utils.SuccessResponse(c, http.StatusCreated, "Usuario creado con éxito", user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
  id := c.Param("id")
  
  user, err := h.userUseCase.GetUser(id)
  if err != nil {
    utils.ErrorResponse(c, http.StatusNotFound, err.Error())
    return
  }
  
  utils.SuccessResponse(c, http.StatusOK, "Usuario obtenido con éxito", user)
}
```

### Diferencias Clave

1. **Estructura:** En Express, las rutas suelen definirse como funciones anónimas. En Gin, definimos métodos de un struct.

2. **Middleware:** En Express, los middleware se pasan como argumentos a las rutas. En Gin, se aplican con el método `.Use()`.

3. **Manejadores:** En Express, los manejadores reciben `req` y `res`. En Gin, reciben un contexto `c`.

4. **Parámetros:** En Express, se accede con `req.params.id`. En Gin, con `c.Param("id")`.

5. **Body Parsing:** En Express, usamos `req.body` (ya procesado por middleware como `express.json()`). En Gin, usamos `c.ShouldBindJSON(&req)`.

6. **Respuestas:** En Express, usamos `res.json()`, `res.status()`. En Gin, usamos `c.JSON()` o nuestras funciones de utilidad como `utils.SuccessResponse()`.

## Control de Errores

### En Node.js (Express):

```javascript
// Manejo de errores en controladores
async function getUser(req, res) {
  try {
    const user = await UserModel.findById(req.params.id);
    if (!user) {
      return res.status(404).json({ error: 'User not found' });
    }
    res.json(user);
  } catch (err) {
    console.error('Error fetching user:', err);
    res.status(500).json({ error: 'Internal server error' });
  }
}

// Middleware de error global
app.use((err, req, res, next) => {
  console.error(err.stack);
  res.status(500).json({ error: 'Something broke!' });
});

// Custom error class
class AppError extends Error {
  constructor(message, statusCode) {
    super(message);
    this.statusCode = statusCode;
    this.isOperational = true;
  }
}
```

### En Go (Gin):

```go
// Devolver errores en repositories
func (r *mongoRepository) GetByID(id string) (*domain.User, error) {
  objID, err := primitive.ObjectIDFromHex(id)
  if err != nil {
    return nil, err  // Devuelve el error
  }
  
  var user domain.User
  err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
  if err != nil {
    if err == mongo.ErrNoDocuments {
      return nil, errors.New("usuario no encontrado")  // Error específico
    }
    return nil, err  // Error genérico
  }
  
  return &user, nil  // Éxito - sin error
}

// Manejo de errores en usecase
func (u *userUseCase) GetUser(id string) (*domain.UserResponse, error) {
  user, err := u.userRepo.GetByID(id)
  if err != nil {
    return nil, err  // Propaga el error
  }
  
  // Mapear a la respuesta
  return &domain.UserResponse{
    ID:   user.ID.Hex(),
    Name: user.Name,
    // ...
  }, nil
}

// Manejo de errores en delivery (handlers HTTP)
func (h *UserHandler) GetUser(c *gin.Context) {
  id := c.Param("id")
  
  user, err := h.userUseCase.GetUser(id)
  if err != nil {
    // Determinar código de error apropiado
    if strings.Contains(err.Error(), "no encontrado") {
      utils.ErrorResponse(c, http.StatusNotFound, err.Error())
      return
    }
    utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
    return
  }
  
  utils.SuccessResponse(c, http.StatusOK, "Usuario obtenido con éxito", user)
}

// Función de utilidad para respuestas de error
func ErrorResponse(c *gin.Context, status int, message string) {
  c.JSON(status, gin.H{
    "status": "error",
    "error":  message,
  })
}
```

### Diferencias Clave

1. **Manejo múltiple de retorno**: En Go, las funciones comúnmente devuelven un resultado y un error: `result, err := someFunction()`.

2. **Comprobación de errores explícita**: Go no tiene excepciones, así que debes comprobar explícitamente `if err != nil { ... }`.

3. **Propagación de errores**: Los errores generalmente se propagan hacia arriba en la pila de llamadas.

4. **Creación de errores**: Usa `errors.New("mensaje")` o `fmt.Errorf("error: %v", err)` para crear errores.

5. **No hay try/catch**: Go utiliza flujo de control condicional para manejar errores, no bloques try/catch.

# Guía de Go para Desarrolladores de Node.js

Esta guía está diseñada para ayudar a desarrolladores familiarizados con Node.js a entender la estructura, sintaxis y patrones de nuestra arquitectura en Go. Nos enfocaremos en las diferencias clave y proporcionar ejemplos comparativos para facilitar la transición.

## Tabla de Contenidos

1. [Conceptos Básicos de Go vs. Node.js](#conceptos-básicos-de-go-vs-nodejs)
2. [Estructura del Proyecto](#estructura-del-proyecto)
3. [Patrones Comunes en Go](#patrones-comunes-en-go)
4. [Trabajando con MongoDB](#trabajando-con-mongodb)
5. [Manejo de Rutas HTTP en Gin vs. Express](#manejo-de-rutas-http-en-gin-vs-express)
6. [Control de Errores](#control-de-errores)
7. [Patrones de Arquitectura Limpia](#patrones-de-arquitectura-limpia)
8. [Glosario de Términos](#glosario-de-términos)

## Conceptos Básicos de Go vs. Node.js

### Declaración de Variables

**Node.js (JavaScript):**
```javascript
// Variables
let counter = 0;
const name = "John";
var oldStyle = true;

// Objetos
const user = {
  id: 1,
  name: "John",
  roles: ["admin", "editor"]
};
```

**Go:**
```go
// Variables
var counter int = 0     // Declaración explícita de tipo
counter := 0            // Inferencia de tipo con :=
name := "John"          // String
oldStyle := true        // Boolean

// Structs (similar a objetos)
type User struct {
  ID    int      `json:"id"`     // Las etiquetas definen cómo se serializa
  Name  string   `json:"name"`
  Roles []string `json:"roles"`  // Array/slice
}

user := User{
  ID:    1,
  Name:  "John",
  Roles: []string{"admin", "editor"},
}
```

### Funciones

**Node.js:**
```javascript
// Función básica
function add(a, b) {
  return a + b;
}

// Arrow function
const multiply = (a, b) => a * b;

// Método de objeto
const calculator = {
  add: function(a, b) {
    return a + b;
  }
};

// Async/await
async function fetchData() {
  try {
    const response = await fetch('https://api.example.com/data');
    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Error:', error);
    throw error;
  }
}
```

**Go:**
```go
// Función básica
func add(a, b int) int {
  return a + b
}

// Función con múltiples retornos (común en Go)
func divide(a, b float64) (float64, error) {
  if b == 0 {
    return 0, errors.New("división por cero")
  }
  return a / b, nil
}

// Método (función asociada a un tipo)
type Calculator struct{}

func (c Calculator) Add(a, b int) int {
  return a + b
}

// Manejo asíncrono con goroutines (diferente al async/await)
func fetchData() ([]byte, error) {
  resp, err := http.Get("https://api.example.com/data")
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()  // Se ejecuta al final de la función
  
  return ioutil.ReadAll(resp.Body)
}
```

### Tipos de Datos

| JavaScript | Go | Notas |
|------------|----|----|
| Number | int, int8, int16, int32, int64, float32, float64 | Go distingue entre varios tipos numéricos |
| String | string | Similar |
| Boolean | bool | Similar pero con nombre diferente |
| Object | struct | Los structs tienen campos tipados y definidos previamente |
| Array, Array.map | slice, for loops | Go no tiene métodos funcionales como map o filter |
| null/undefined | nil | nil es similar pero se usa solo con tipos específicos |
| Promise | Channels, goroutines | Concurrencia diferente a promesas |

## Estructura del Proyecto

Nuestra aplicación sigue una arquitectura limpia (Clean Architecture) organizada por dominio:

```
proyecto-raiz/
├── internal/          # Código específico de la aplicación
│   ├── user/          # Módulo de usuario (dominio)
│   │   ├── domain/    # Definiciones y contratos
│   │   ├── repository/# Implementación de persistencia
│   │   ├── usecase/   # Lógica de negocio
│   │   └── delivery/  # Controladores HTTP
│   └── product/       # Otro módulo de dominio
├── pkg/               # Código compartido y utilidades
├── cmd/               # Entradas de la aplicación
└── main.go            # Punto de entrada
```

**Comparación con Node.js:**

En Node.js (Express) puedes tener:
```
proyecto-node/
├── routes/            # Rutas HTTP
├── controllers/       # Controladores
├── models/            # Modelos (Mongoose)
├── services/          # Servicios/lógica
└── app.js             # Punto de entrada
```

La diferencia principal es que en nuestra estructura Go, organizamos por módulos/dominios (como "user", "product") en lugar de por capas técnicas (como "routes", "controllers").

## Patrones Comunes en Go

### Interfaces

Las interfaces en Go son diferentes a las de TypeScript o a las clases en JavaScript. En Go, las interfaces se implementan implícitamente:

```go
// Definición de una interfaz
type Repository interface {
  GetByID(id string) (*Entity, error)
  Save(entity *Entity) error
}

// Implementación implícita (no se declara explícitamente)
type MongoRepository struct {
  collection *mongo.Collection
}

// Si MongoRepository implementa todos los métodos de Repository,
// automáticamente satisface la interfaz
func (r *MongoRepository) GetByID(id string) (*Entity, error) {
  // implementación...
}

func (r *MongoRepository) Save(entity *Entity) error {
  // implementación...
}
```

Esto es diferente de TypeScript/JavaScript donde las interfaces o clases se implementan explícitamente:

```typescript
interface Repository {
  getById(id: string): Promise<Entity>;
  save(entity: Entity): Promise<void>;
}

// Implementación explícita
class MongoRepository implements Repository {
  constructor(private collection: Collection) {}
  
  async getById(id: string): Promise<Entity> {
    // implementación...
  }
  
  async save(entity: Entity): Promise<void> {
    // implementación...
  }
}
```

### Inyección de Dependencias

En Go no tenemos decoradores ni inyección automática. La inyección se hace manualmente en los constructores:

```go
// Constructor de repositorio
func NewMongoRepository(collection *mongo.Collection) domain.Repository {
  return &mongoRepository{
    collection: collection,
  }
}

// Constructor de caso de uso
func NewUseCase(repo domain.Repository) domain.UseCase {
  return &useCase{
    repo: repo,
  }
}

// En main.go - Ensamblaje manual
collection := mongoClient.Database("mydb").Collection("entities")
repository := repository.NewMongoRepository(collection)
useCase := usecase.NewUseCase(repository)
handler := delivery.NewHandler(useCase)
```

Esto contrasta con Node.js donde podrías usar decoradores (con NestJS) o un contenedor DI:

```typescript
// Con NestJS
@Injectable()
export class EntityService {
  constructor(
    @InjectRepository(Entity)
    private entityRepository: Repository<Entity>,
  ) {}
}
```

## Trabajando con MongoDB

### Inicialización y Conexión

**Node.js con Mongoose:**
```javascript
const mongoose = require('mongoose');

mongoose.connect('mongodb://localhost:27017/mydatabase', {
  useNewUrlParser: true,
  useUnifiedTopology: true
})
.then(() => console.log('Connected to MongoDB'))
.catch(err => console.error('Error connecting to MongoDB', err));

// Definición de modelo
const userSchema = new mongoose.Schema({
  name: String,
  email: { type: String, required: true, unique: true },
  age: Number,
  createdAt: { type: Date, default: Date.now }
});

const User = mongoose.model('User', userSchema);
```

**Go con driver oficial:**
```go
import (
  "context"
  "log"
  "time"
  
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

func connectToMongoDB() (*mongo.Client, error) {
  // Crear contexto con timeout
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()
  
  // Conectar al servidor
  client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
  if err != nil {
    return nil, err
  }
  
  // Verificar la conexión
  err = client.Ping(ctx, nil)
  if err != nil {
    return nil, err
  }
  
  log.Println("Conectado a MongoDB")
  return client, nil
}

// No hay esquemas explícitos como en Mongoose
// Usamos structs para modelar documentos
type User struct {
  ID        primitive.ObjectID `bson:"_id,omitempty"`
  Name      string             `bson:"name"`
  Email     string             `bson:"email"`
  Age       int                `bson:"age"`
  CreatedAt time.Time          `bson:"created_at"`
}
```

### Operaciones Comunes con MongoDB

#### Insertar Documentos

**Node.js:**
```javascript
// Insertar un documento
const newUser = new User({
  name: 'John Doe',
  email: 'john@example.com',
  age: 30
});

await newUser.save();

// Insertar múltiples documentos
await User.insertMany([
  { name: 'Alice', email: 'alice@example.com', age: 25 },
  { name: 'Bob', email: 'bob@example.com', age: 28 }
]);
```

**Go:**
```go
// Insertar un documento
user := User{
  ID:        primitive.NewObjectID(),
  Name:      "John Doe",
  Email:     "john@example.com",
  Age:       30,
  CreatedAt: time.Now(),
}

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

collection := client.Database("mydatabase").Collection("users")
_, err := collection.InsertOne(ctx, user)
if err != nil {
  log.Fatal(err)
}

// Insertar múltiples documentos
users := []interface{}{
  User{
    ID:        primitive.NewObjectID(),
    Name:      "Alice",
    Email:     "alice@example.com",
    Age:       25,
    CreatedAt: time.Now(),
  },
  User{
    ID:        primitive.NewObjectID(),
    Name:      "Bob",
    Email:     "bob@example.com",
    Age:       28,
    CreatedAt: time.Now(),
  },
}

_, err = collection.InsertMany(ctx, users)
```

#### Consultar Documentos

**Node.js:**
```javascript
// Encontrar por ID
const user = await User.findById('60d5ec9af682fbd12a8924c4');

// Encontrar uno
const admin = await User.findOne({ role: 'admin' });

// Encontrar muchos
const adultUsers = await User.find({ age: { $gte: 18 } }).sort({ name: 1 }).limit(10);

// Consulta con proyección
const userEmails = await User.find({}, 'name email').exec();
```

**Go:**
```go
// Encontrar por ID
id, _ := primitive.ObjectIDFromHex("60d5ec9af682fbd12a8924c4")
var user User
err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)

// Encontrar uno
var admin User
err = collection.FindOne(ctx, bson.M{"role": "admin"}).Decode(&admin)

// Encontrar muchos
opts := options.Find().SetSort(bson.M{"name": 1}).SetLimit(10)
cursor, err := collection.Find(ctx, bson.M{"age": bson.M{"$gte": 18}}, opts)
if err != nil {
  log.Fatal(err)
}
defer cursor.Close(ctx)

var adultUsers []User
if err = cursor.All(ctx, &adultUsers); err != nil {
  log.Fatal(err)
}

// Consulta con proyección
opts = options.Find().SetProjection(bson.M{"name": 1, "email": 1})
cursor, err = collection.Find(ctx, bson.M{}, opts)
// ...procesar cursor
```

#### Actualizar Documentos

**Node.js:**
```javascript
// Actualizar uno
await User.updateOne(
  { _id: '60d5ec9af682fbd12a8924c4' },
  { $set: { name: 'John Smith', updatedAt: new Date() } }
);

// Actualizar muchos
await User.updateMany(
  { age: { $lt: 18 } },
  { $set: { status: 'minor' } }
);

// Encontrar y actualizar
const updatedUser = await User.findByIdAndUpdate(
  '60d5ec9af682fbd12a8924c4',
  { $set: { status: 'active' } },
  { new: true } // Devuelve el documento actualizado
);
```

**Go:**
```go
// Actualizar uno
id, _ := primitive.ObjectIDFromHex("60d5ec9af682fbd12a8924c4")
update := bson.M{
  "$set": bson.M{
    "name":       "John Smith",
    "updated_at": time.Now(),
  },
}
_, err = collection.UpdateOne(ctx, bson.M{"_id": id}, update)

// Actualizar muchos
_, err = collection.UpdateMany(
  ctx,
  bson.M{"age": bson.M{"$lt": 18}},
  bson.M{"$set": bson.M{"status": "minor"}},
)

// Encontrar y actualizar (no hay método directo como en Mongoose)
opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
var updatedUser User
err = collection.FindOneAndUpdate(
  ctx,
  bson.M{"_id": id},
  bson.M{"$set": bson.M{"status": "active"}},
  opts,
).Decode(&updatedUser)
```

## Manejo de Rutas HTTP en Gin vs. Express

### Definición de Rutas

**Express (Node.js):**
```javascript
const express = require('express');
const router = express.Router();

// Middleware
const authMiddleware = require('./middleware/auth');

// Rutas
router.get('/users', authMiddleware, async (req, res) => {
  try {
    const users = await UserModel.find({});
    res.json(users);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

router.post('/users', authMiddleware, async (req, res) => {
  try {
    const newUser = new UserModel(req.body);
    const savedUser = await newUser.save();
    res.status(201).json(savedUser);
  } catch (err) {
    res.status(400).json({ error: err.message });
  }
});

// Parámetros en ruta
router.get('/users/:id', authMiddleware, async (req, res) => {
  try {
    const user = await UserModel.findById(req.params.id);
    if (!user) return res.status(404).json({ error: 'User not found' });
    res.json(user);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Exportar router
module.exports = router;
```

**Gin (Go):**
```go
package delivery

import (
  "net/http"

  "github.com/gin-gonic/gin"
  
  "myproject/internal/user/domain"
  "myproject/pkg/middleware"
  "myproject/pkg/utils"
)

type UserHandler struct {
  userUseCase domain.UserUseCase
}

// Constructor - equivalente a exportar el router
func NewUserHandler(router *gin.RouterGroup, useCase domain.UserUseCase) {
  handler := &UserHandler{
    userUseCase: useCase,
  }
  
  // Aplicar middleware
  userRoutes := router.Group("/users")
  userRoutes.Use(middleware.Auth())
  
  // Definir rutas
  userRoutes.GET("/", handler.GetAllUsers)
  userRoutes.POST("/", handler.CreateUser)
  userRoutes.GET("/:id", handler.GetUser)
}

// Métodos de handler
func (h *UserHandler) GetAllUsers(c *gin.Context) {
  users, err := h.userUseCase.GetAllUsers()
  if err != nil {
    utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
    return
  }
  
  utils.SuccessResponse(c, http.StatusOK, "Usuarios obtenidos con éxito", users)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
  var req domain.CreateUserRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    utils.ValidationErrorResponse(c, err.Error())
    return
  }
  
  user, err := h.userUseCase.CreateUser(&req)
  if err != nil {
    utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
    return
  }
  
  utils.SuccessResponse(c, http.StatusCreated, "Usuario creado con éxito", user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
  id := c.Param("id")
  
  user, err := h.userUseCase.GetUser(id)
  if err != nil {
    utils.ErrorResponse(c, http.StatusNotFound, err.Error())
    return
  }
  
  utils.SuccessResponse(c, http.StatusOK, "Usuario obtenido con éxito", user)
}
```

### Diferencias Clave

1. **Estructura:** En Express, las rutas suelen definirse como funciones anónimas. En Gin, definimos métodos de un struct.

2. **Middleware:** En Express, los middleware se pasan como argumentos a las rutas. En Gin, se aplican con el método `.Use()`.

3. **Manejadores:** En Express, los manejadores reciben `req` y `res`. En Gin, reciben un contexto `c`.

4. **Parámetros:** En Express, se accede con `req.params.id`. En Gin, con `c.Param("id")`.

5. **Body Parsing:** En Express, usamos `req.body` (ya procesado por middleware como `express.json()`). En Gin, usamos `c.ShouldBindJSON(&req)`.

6. **Respuestas:** En Express, usamos `res.json()`, `res.status()`. En Gin, usamos `c.JSON()` o nuestras funciones de utilidad como `utils.SuccessResponse()`.

## Control de Errores

### En Node.js (Express):

```javascript
// Manejo de errores en controladores
async function getUser(req, res) {
  try {
    const user = await UserModel.findById(req.params.id);
    if (!user) {
      return res.status(404).json({ error: 'User not found' });
    }
    res.json(user);
  } catch (err) {
    console.error('Error fetching user:', err);
    res.status(500).json({ error: 'Internal server error' });
  }
}

// Middleware de error global
app.use((err, req, res, next) => {
  console.error(err.stack);
  res.status(500).json({ error: 'Something broke!' });
});

// Custom error class
class AppError extends Error {
  constructor(message, statusCode) {
    super(message);
    this.statusCode = statusCode;
    this.isOperational = true;
  }
}
```

### En Go (Gin):

```go
// Devolver errores en repositories
func (r *mongoRepository) GetByID(id string) (*domain.User, error) {
  objID, err := primitive.ObjectIDFromHex(id)
  if err != nil {
    return nil, err  // Devuelve el error
  }
  
  var user domain.User
  err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
  if err != nil {
    if err == mongo.ErrNoDocuments {
      return nil, errors.New("usuario no encontrado")  // Error específico
    }
    return nil, err  // Error genérico
  }
  
  return &user, nil  // Éxito - sin error
}

// Manejo de errores en usecase
func (u *userUseCase) GetUser(id string) (*domain.UserResponse, error) {
  user, err := u.userRepo.GetByID(id)
  if err != nil {
    return nil, err  // Propaga el error
  }
  
  // Mapear a la respuesta
  return &domain.UserResponse{
    ID:   user.ID.Hex(),
    Name: user.Name,
    // ...
  }, nil
}

// Manejo de errores en delivery (handlers HTTP)
func (h *UserHandler) GetUser(c *gin.Context) {
  id := c.Param("id")
  
  user, err := h.userUseCase.GetUser(id)
  if err != nil {
    // Determinar código de error apropiado
    if strings.Contains(err.Error(), "no encontrado") {
      utils.ErrorResponse(c, http.StatusNotFound, err.Error())
      return
    }
    utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
    return
  }
  
  utils.SuccessResponse(c, http.StatusOK, "Usuario obtenido con éxito", user)
}

// Función de utilidad para respuestas de error
func ErrorResponse(c *gin.Context, status int, message string) {
  c.JSON(status, gin.H{
    "status": "error",
    "error":  message,
  })
}
```

### Diferencias Clave

1. **Manejo múltiple de retorno**: En Go, las funciones comúnmente devuelven un resultado y un error: `result, err := someFunction()`.

2. **Comprobación de errores explícita**: Go no tiene excepciones, así que debes comprobar explícitamente `if err != nil { ... }`.

3. **Propagación de errores**: Los errores generalmente se propagan hacia arriba en la pila de llamadas.

4. **Creación de errores**: Usa `errors.New("mensaje")` o `fmt.Errorf("error: %v", err)` para crear errores.

5. **No hay try/catch**: Go utiliza flujo de control condicional para manejar errores, no bloques try/catch.

## Patrones de Arquitectura Limpia

Nuestra aplicación utiliza Clean Architecture con estas capas:

### 1. Entities/Domain Layer (Capa de Dominio)

En Node.js podrías tener:
```javascript
// models/user.js
class User {
  constructor(id, name, email) {
    this.id = id;
    this.name = name;
    this.email = email;
  }
}

module.exports = User;
```

En nuestra arquitectura Go:
```go
// internal/user/domain/user.domain.go
package domain

import (
  "time"
  "go.mongodb.org/mongo-driver/bson/primitive"
)

// User representa la entidad de usuario
type User struct {
  ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
  Name      string             `json:"name" bson:"name"`
  Email     string             `json:"email" bson:"email"`
  Password  string             `json:"-" bson:"password"`  // "-" oculta este campo en JSON
  CreatedAt time.Time          `json:"created_at" bson:"created_at"`
  UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// UserRepository define el contrato para la capa de persistencia
type UserRepository interface {
  GetByID(id string) (*User, error)
  // ...otros métodos
}

// UserUseCase define el contrato para la capa de casos de uso
type UserUseCase interface {
  GetUser(id string) (*UserResponse, error)
  // ...otros métodos
}
```

### 2. Use Cases / Application Layer

En Node.js podrías tener:
```javascript
// services/userService.js
class UserService {
  constructor(userRepository) {
    this.userRepository = userRepository;
  }

  async getUser(id) {
    return this.userRepository.findById(id);
  }
  
  // ...otros métodos
}

module.exports = UserService;
```

En nuestra arquitectura Go:
```go
// internal/user/usecase/user.usecase.go
package usecase

import (
  "github.com/myproject/internal/user/domain"
)

type userUseCase struct {
  userRepo domain.UserRepository
}

// NewUserUseCase - constructor
func NewUserUseCase(userRepo domain.UserRepository) domain.UserUseCase {
  return &userUseCase{
    userRepo: userRepo,
  }
}

// GetUser implementa el método de la interfaz
func (u *userUseCase) GetUser(id string) (*domain.UserResponse, error) {
  user, err := u.userRepo.GetByID(id)
  if err != nil {
    return nil, err
  }
  
  return &domain.UserResponse{
    ID:    user.ID.Hex(),
    Name:  user.Name,
    Email: user.Email,
    // ...mapear a una respuesta sin campos sensibles
  }, nil
}
```

### 3. Interface Adapters Layer

En Node.js podrías tener:
```javascript
// repositories/userRepository.js
class UserRepository {
  constructor(database) {
    this.database = database;
  }

  async findById(id) {
    return this.database.collection('users').findOne({ _id: ObjectId(id) });
  }
  
  // ...otros métodos
}

// controllers/userController.js
class UserController {
  constructor(userService) {
    this.userService = userService;
  }

  async getUser(req, res) {
    try {
      const user = await this.userService.getUser(req.params.id);
      if (!user) return res.status(404).json({ error: 'User not found' });
      res.json(user);
    } catch (err) {
      res.status(500).json({ error: err.message });
    }
  }
  
  // ...otros métodos
}
```

En nuestra arquitectura Go:
```go
// Repository
// internal/user/repository/mongo.user.repository.go
package repository

import (
  "context"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  
  "github.com/myproject/internal/user/domain"
)

type mongoUserRepository struct {
  collection *mongo.Collection
}

func NewMongoUserRepository(collection *mongo.Collection) domain.UserRepository {
  return &mongoUserRepository{
    collection: collection,
  }
}

func (r *mongoUserRepository) GetByID(id string) (*domain.User, error) {
  // implementación...
}

// Delivery (Controllers)
// internal/user/delivery/user.delivery.go
package delivery

import (
  "github.com/gin-gonic/gin"
  
  "github.com/myproject/internal/user/domain"
)

type UserHandler struct {
  userUseCase domain.UserUseCase
}

func NewUserHandler(router *gin.RouterGroup, useCase domain.UserUseCase) {
  handler := &UserHandler{
    userUseCase: useCase,
  }
  
  router.GET("/:id", handler.GetUser)
  // ...otras rutas
}

func (h *UserHandler) GetUser(c *gin.Context) {
  // implementación...
}
```

### 4. Frameworks & Drivers Layer (framework y drivers)

En Node.js:
```javascript
// app.js
const express = require('express');
const { MongoClient } = require('mongodb');
const app = express();

// Inicializar
async function initApp() {
  // Conectar a la base de datos
  const client = await MongoClient.connect('mongodb://localhost:27017');
  const db = client.db('mydatabase');
  
  // Crear dependencias
  const userRepository = new UserRepository(db);
  const userService = new UserService(userRepository);
  const userController = new UserController(userService);
  
  // Configurar rutas
  app.get('/users/:id', userController.getUser.bind(userController));
  
  app.listen(3000, () => console.log('Server running on port 3000'));
}

initApp().catch(console.error);
```

En nuestra arquitectura Go:
```go
// main.go
package main

import (
  "context"
  "log"
  
  "github.com/gin-gonic/gin"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  
  "github.com/myproject/internal/user/delivery"
  "github.com/myproject/internal/user/repository"
  "github.com/myproject/internal/user/usecase"
)

func main() {
  // Conectar a MongoDB
  client, err := connectMongoDB()
  if err != nil {
    log.Fatal(err)
  }
  defer client.Disconnect(context.Background())
  
  // Inicializar router
  router := gin.Default()
  
  // Configurar módulo de usuario
  userCollection := client.Database("mydatabase").Collection("users")
  userRepository := repository.NewMongoUserRepository(userCollection)
  userUseCase := usecase.NewUserUseCase(userRepository)
  
  // Configurar rutas
  userGroup := router.Group("/api/users")
  delivery.NewUserHandler(userGroup, userUseCase)
  
  // Iniciar servidor
  if err := router.Run(":3000"); err != nil {
    log.Fatal(err)
  }
}

func connectMongoDB() (*mongo.Client, error) {
  // Implementación de conexión...
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()
  
  client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
  if err != nil {
    return nil, err
  }
  
  err = client.Ping(ctx, nil)
  if err != nil {
    return nil, err
  }
  
  return client, nil
}
```

La principal diferencia es la forma en que construimos nuestra aplicación en Go:
1. Usamos inyección de dependencias manual en vez de clases y constructores
2. No tenemos promesas ni async/await, usamos retornos múltiples (valor, error)
3. Organizamos el código por dominio en lugar de por capas técnicas

## Comparación Conceptual Node.js vs Go

| Concepto | Node.js | Go | 
|----------|---------|-------|
| Async    | Promise, async/await | Goroutines, canales |
| Errores  | try/catch, Promise reject | if err != nil {} |
| Clases   | class User {} | type User struct {} |
| Métodos  | user.save() | func (u *User) Save() {} |
| Interfaces | interface UserRepo {} | type UserRepo interface {} |
| Módulos  | import/export, require | package, import |
| Variables | let, const, var | var, := |
| Nulls    | null, undefined | nil |
| Formato  | Flexible, ESLint | go fmt (estándar) |
| Generics | Sí (tipado dinámico) | Sí (desde Go 1.18) |

## Recursos Adicionales

- [Tour de Go](https://go.dev/tour) - Tutorial interactivo oficial
- [Go by Example](https://gobyexample.com/) - Ejemplos prácticos de conceptos de Go
- [Effective Go](https://go.dev/doc/effective_go) - Guía de buenas prácticas
- [Error Handling in Go](https://go.dev/blog/error-handling-and-go)
- [Gin Framework](https://gin-gonic.com/docs/) - Documentación oficial de Gin
- [MongoDB Go Driver](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo)

## Consejos para la Transición

1. **Flujo de errores**: Acostumbrate a revisar errores explícitamente en vez de usar try/catch.
2. **Tipado estático**: Define bien tus estructuras y piensa en los tipos desde el principio.
3. **Concurrencia**: Entiende cómo funcionan goroutines y canales, son diferentes de async/await.
4. **Zero values**: Go inicializa variables con valores por defecto (0, "", false, nil) según el tipo.
5. **Convenciones de nombres**:
    - PascalCase para exportación (públicos)
    - camelCase para no exportación (privados)
6. **Estructura del proyecto**: Organiza por dominio, no por capa técnica.
7. **Testing**: Go tiene soporte para testing incorporado, úsalo desde el principio.
8. **Punteros vs valores**: Entiende cuándo usar punteros (*Type) y cuándo usar valores directos.
9. **Context**: Familiarízate con Context para manejar deadlines, cancelaciones y valores de petición.
10. **Interfaces implícitas**: Recuerda que en Go las interfaces son satisfechas implícitamente.

## Diferencias Principales en el Flujo de Datos

### Express (Node.js)
```
Cliente → Routers → Middleware → Controllers → Services → Models → DB
          ↑                                                    |
          └────────────────── Respuesta ───────────────────────┘
```

### Gin (Go con Arquitectura Limpia)
```
Cliente → Routers → Middleware → Handlers → UseCases → Repository → DB
          ↑                          ↓           ↑           ↓
          └─────────── Respuesta ────┴───── Entidades ───────┘
```

## Patrones Comunes en la Arquitectura

### 1. Constructor para Inyección de Dependencias

```go
// Crear e inyectar dependencias
func NewService(repository Repository) Service {
  return &serviceImpl{
    repo: repository,
  }
}
```

### 2. Interfaces como Contratos

```go
// Definir interfaces que representan comportamiento
type Repository interface {
  Get(id string) (*Entity, error)
  Save(entity *Entity) error
}
```

### 3. Patrón de Retorno con Error

```go
// Retornar valor y error
func (r *repository) Get(id string) (*Entity, error) {
  if id == "" {
    return nil, errors.New("id no puede estar vacío")
  }
  
  // Lógica...
  return entity, nil
}
```

### 4. Handlers HTTP que Llaman a Casos de Uso

```go
func (h *handler) GetEntity(c *gin.Context) {
  id := c.Param("id")
  entity, err := h.useCase.GetEntity(id)
  
  if err != nil {
    // Manejar error
    return
  }
  
  // Responder con éxito
}
```

## Conclusión

La transición de Node.js a Go implica algunos cambios de mentalidad importantes, especialmente en cuanto al manejo de errores, tipado estático y modelo de concurrencia. Sin embargo, muchos conceptos de arquitectura limpia son transferibles entre ambos lenguajes.

La mayor ventaja de nuestra arquitectura actual es que proporciona una estructura clara y modular que facilita la evolución del código, promueve la reutilización y simplifica las pruebas, independientemente del lenguaje utilizado.

Recuerda que esta guía es solo un punto de partida. A medida que te familiarices más con Go y con la arquitectura del proyecto, desarrollarás tus propios patrones y preferencias adaptados a tus necesidades específicas.
