package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

// Config almacena toda la configuración de la aplicación
type Config struct {
	// Servidor
	Port string
	Env  string

	// MongoDB
	MongoURI     string
	MongoDB      string
	MongoTimeout time.Duration

	// OAuth
	JWTSecret  string
	TokenExp   time.Duration
	RefreshExp time.Duration

	// Cliente OAuth (solo si tu aplicación es también un cliente)
	OAuthClientID     string
	OAuthClientSecret string
}

// LoadConfig carga la configuración desde variables de entorno
func LoadConfig() (*Config, error) {
	// Cargar variables desde .env
	godotenv.Load()

	// Configuración predeterminada
	config := &Config{
		Port:         getEnv("PORT", "3000"),
		Env:          getEnv("ENV", "development"),
		MongoURI:     getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:      getEnv("MONGO_DB", "my_database"),
		MongoTimeout: time.Duration(getEnvAsInt("MONGO_TIMEOUT", 10)) * time.Second,
		JWTSecret:    getEnv("JWT_SECRET", "mi_secret_super_seguro"),
		TokenExp:     time.Duration(getEnvAsInt("TOKEN_EXP", 2)) * time.Hour,
		RefreshExp:   time.Duration(getEnvAsInt("REFRESH_EXP", 7*24)) * time.Hour, // 7 días

	}

	return config, nil
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt obtiene una variable de entorno como entero o retorna un valor por defecto
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool obtiene una variable de entorno como booleano o retorna un valor por defecto
func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// IsDevelopment verifica si estamos en entorno de desarrollo
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// IsProduction verifica si estamos en entorno de producción
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}
