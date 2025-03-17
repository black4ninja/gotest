package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoConfig contiene la configuración para MongoDB
type MongoConfig struct {
	URI      string
	Database string
	Timeout  time.Duration
}

// NewMongoClient crea un nuevo cliente de MongoDB
func NewMongoClient(config MongoConfig) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.URI)
	client, err := mongo.Connect(ctx, clientOptions)
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

// GetCollection obtiene una colección de MongoDB
func GetCollection(client *mongo.Client, dbName, colName string) *mongo.Collection {
	return client.Database(dbName).Collection(colName)
}
