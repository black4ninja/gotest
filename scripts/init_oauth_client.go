package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/black4ninja/mi-proyecto/internal/oauth/domain"
	"github.com/black4ninja/mi-proyecto/pkg/utils"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("Archivo .env no encontrado, usando variables de entorno del sistema")
	}

	// Configurar MongoDB
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	mongoDBName := getEnv("MONGO_DB", "my_database")

	// Conectar a MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error al conectar a MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Crear ID de cliente y secreto aleatorios
	clientID, err := utils.GenerateRandomString(16)
	if err != nil {
		log.Fatalf("Error al generar client_id: %v", err)
	}

	clientSecret, err := utils.GenerateRandomString(32)
	if err != nil {
		log.Fatalf("Error al generar client_secret: %v", err)
	}

	// Crear cliente OAuth
	now := time.Now()
	oauthClient := domain.Client{
		ID:           primitive.NewObjectID(),
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Name:         "Cliente de prueba",
		RedirectURIs: []string{"http://localhost:3000/callback"},
		GrantTypes:   []string{domain.GrantTypePassword, domain.GrantTypeRefreshToken, domain.GrantTypeClientCredentials},
		Scopes:       []string{"read", "write", "admin"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Guardar en la base de datos
	collection := client.Database(mongoDBName).Collection("oauth_clients")
	_, err = collection.InsertOne(ctx, oauthClient)
	if err != nil {
		log.Fatalf("Error al insertar cliente OAuth: %v", err)
	}

	// Mostrar información
	log.Println("Cliente OAuth creado correctamente:")
	log.Printf("ClientID: %s", clientID)
	log.Printf("ClientSecret: %s", clientSecret)
	log.Printf("Scopes disponibles: %v", oauthClient.Scopes)
	log.Printf("Tipos de concesión: %v", oauthClient.GrantTypes)
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
