package utils

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FilterDefinition define los parámetros para un campo filtrable
type FilterDefinition struct {
	AllowedValues []string                 // Valores permitidos para el campo (si es una lista de opciones)
	Validator     func(string) bool        // Función de validación personalizada (si no es una lista simple)
	Transformer   func(string) interface{} // Función para transformar el valor antes de usarlo (ej: convertir a ObjectID)
}

// FilterConfig define qué campos pueden ser filtrados y cómo
type FilterConfig map[string]FilterDefinition

// BuildMongoFilter construye un filtro seguro para MongoDB basado en parámetros de consulta
func BuildMongoFilter(queryParams map[string]string, config FilterConfig) bson.M {
	filter := bson.M{}

	for param, value := range queryParams {
		// Verificar si este parámetro está permitido para filtrado
		if definition, exists := config[param]; exists {
			// Ignorar valores vacíos
			if value == "" {
				continue
			}

			// Verificar si es un valor permitido (si hay lista de valores permitidos)
			if len(definition.AllowedValues) > 0 {
				isValidValue := false
				for _, allowed := range definition.AllowedValues {
					if value == allowed {
						isValidValue = true
						break
					}
				}

				if !isValidValue {
					continue // Ignorar valores no permitidos
				}
			}

			// Aplicar validador personalizado (si está definido)
			if definition.Validator != nil {
				if !definition.Validator(value) {
					continue // Ignorar valores que no pasan la validación
				}
			}

			// Aplicar transformador (si está definido)
			var finalValue interface{} = value
			if definition.Transformer != nil {
				finalValue = definition.Transformer(value)
			}

			// Añadir al filtro
			filter[param] = finalValue
		}
	}

	return filter
}

// Validadores y transformadores comunes

// IsValidObjectID verifica si un string puede convertirse en un ObjectID válido
func IsValidObjectID(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}

// TransformToObjectID convierte un string a ObjectID
func TransformToObjectID(id string) interface{} {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil
	}
	return objID
}

// IsValidDate verifica si un string puede convertirse en una fecha válida (formato ISO)
func IsValidDate(date string) bool {
	_, err := time.Parse(time.RFC3339, date)
	return err == nil
}

// TransformToDate convierte un string a Date
func TransformToDate(date string) interface{} {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return nil
	}
	return t
}

// TransformToRegex convierte un string a una regex case-insensitive para búsquedas parciales
func TransformToRegex(text string) interface{} {
	// Escapar caracteres especiales de regex
	text = strings.ReplaceAll(text, "\\", "\\\\")
	text = strings.ReplaceAll(text, ".", "\\.")
	text = strings.ReplaceAll(text, "+", "\\+")
	text = strings.ReplaceAll(text, "*", "\\*")
	text = strings.ReplaceAll(text, "?", "\\?")
	text = strings.ReplaceAll(text, "|", "\\|")
	text = strings.ReplaceAll(text, "(", "\\(")
	text = strings.ReplaceAll(text, ")", "\\)")
	text = strings.ReplaceAll(text, "[", "\\[")
	text = strings.ReplaceAll(text, "]", "\\]")
	text = strings.ReplaceAll(text, "{", "\\{")
	text = strings.ReplaceAll(text, "}", "\\}")
	text = strings.ReplaceAll(text, "$", "\\$")
	text = strings.ReplaceAll(text, "^", "\\^")

	return primitive.Regex{
		Pattern: text,
		Options: "i", // case-insensitive
	}
}

// Función auxiliar para crear un rango de fechas
func DateRangeFilter(fromDate, toDate string) bson.M {
	filter := bson.M{}

	if fromDate != "" {
		from, err := time.Parse(time.RFC3339, fromDate)
		if err == nil {
			filter["$gte"] = from
		}
	}

	if toDate != "" {
		to, err := time.Parse(time.RFC3339, toDate)
		if err == nil {
			filter["$lte"] = to
		}
	}

	if len(filter) == 0 {
		return nil
	}

	return filter
}

// Constantes para filtros
const (
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusArchived = "archived"

	RoleAdmin     = "admin"
	RoleUser      = "user"
	RoleModerator = "moderator"
)

// Define filtros comunes para reutilizar
var (
	CommonUserFilterConfig = FilterConfig{
		"status": FilterDefinition{
			AllowedValues: []string{StatusActive, StatusInactive, StatusArchived},
		},
		"role": FilterDefinition{
			AllowedValues: []string{RoleAdmin, RoleUser, RoleModerator},
		},
		"name": FilterDefinition{
			Validator:   func(s string) bool { return len(s) > 0 && len(s) <= 100 },
			Transformer: TransformToRegex,
		},
		"email": FilterDefinition{
			Validator:   func(s string) bool { return strings.Contains(s, "@") && len(s) <= 100 },
			Transformer: TransformToRegex,
		},
		"_id": FilterDefinition{
			Validator:   IsValidObjectID,
			Transformer: TransformToObjectID,
		},
		"created_at": FilterDefinition{
			Validator: IsValidDate,
		},
		"archived_at": FilterDefinition{
			Validator: IsValidDate,
		},
	}
)
