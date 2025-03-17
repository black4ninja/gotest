package delivery_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/black4ninja/mi-proyecto/internal/user/delivery"
	"github.com/black4ninja/mi-proyecto/internal/user/domain"
)

// Caso de uso simulado (mock) para pruebas
type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) GetUser(id string) (*domain.UserResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserResponse), args.Error(1)
}

func (m *MockUserUseCase) GetUserByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserUseCase) GetAllUsers(filters bson.M) ([]*domain.UserResponse, error) {
	args := m.Called(filters)
	return args.Get(0).([]*domain.UserResponse), args.Error(1)
}

func (m *MockUserUseCase) CreateUser(req *domain.CreateUserRequest) (*domain.UserResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserResponse), args.Error(1)
}

func (m *MockUserUseCase) UpdateUser(id string, req *domain.UpdateUserRequest) (*domain.UserResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserResponse), args.Error(1)
}

func (m *MockUserUseCase) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserUseCase) ArchiveUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserUseCase) ChangePassword(userID string, req *domain.ChangePasswordRequest) error {
	args := m.Called(userID, req)
	return args.Error(0)
}

func (m *MockUserUseCase) ValidateCredentials(email string, password string) (*domain.User, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserUseCase) UpdateRefreshToken(userID string, refreshToken string) error {
	args := m.Called(userID, refreshToken)
	return args.Error(0)
}

func (m *MockUserUseCase) GetUserByRefreshToken(refreshToken string) (*domain.User, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// Configuración para pruebas HTTP
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

// Pruebas para el Handler HTTP
func TestCreateUserHandler(t *testing.T) {
	// Configurar el mock
	mockUseCase := new(MockUserUseCase)

	// Configurar router
	r := setupRouter()
	userGroup := r.Group("/api/users")

	// Registrar handler
	delivery.NewUserHandler(userGroup, mockUseCase)

	// Datos de prueba
	createUserReq := domain.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "user",
	}

	userID := primitive.NewObjectID()
	mockResponse := &domain.UserResponse{
		ID:     userID.Hex(),
		Name:   "Test User",
		Email:  "test@example.com",
		Role:   "user",
		Status: domain.UserStatusActive,
	}

	// Configurar comportamiento esperado del mock
	mockUseCase.On("CreateUser", mock.AnythingOfType("*domain.CreateUserRequest")).Return(mockResponse, nil)

	// Crear solicitud HTTP
	jsonValue, _ := json.Marshal(createUserReq)
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	// Ejecutar solicitud
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verificaciones
	assert.Equal(t, http.StatusCreated, w.Code)

	// Parsear respuesta
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verificar estructura de respuesta
	assert.Equal(t, "success", response["status"])

	// Verificar datos de usuario
	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, userID.Hex(), data["id"])
	assert.Equal(t, "test@example.com", data["email"])

	// Verificar que se llamó al caso de uso como esperamos
	mockUseCase.AssertExpectations(t)
}

func TestGetUserHandler(t *testing.T) {
	// Configurar el mock
	mockUseCase := new(MockUserUseCase)

	// Configurar router
	r := setupRouter()
	userGroup := r.Group("/api/users")

	// Registrar handler
	delivery.NewUserHandler(userGroup, mockUseCase)

	// Datos de prueba
	userID := primitive.NewObjectID()
	mockResponse := &domain.UserResponse{
		ID:     userID.Hex(),
		Name:   "Test User",
		Email:  "test@example.com",
		Role:   "user",
		Status: domain.UserStatusActive,
	}

	// Configurar comportamiento esperado del mock
	mockUseCase.On("GetUser", userID.Hex()).Return(mockResponse, nil)

	// Crear solicitud HTTP
	req, _ := http.NewRequest("GET", "/api/users/"+userID.Hex(), nil)

	// Ejecutar solicitud
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verificaciones
	assert.Equal(t, http.StatusOK, w.Code)

	// Parsear respuesta
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verificar estructura de respuesta
	assert.Equal(t, "success", response["status"])

	// Verificar datos de usuario
	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, userID.Hex(), data["id"])
	assert.Equal(t, "test@example.com", data["email"])

	// Verificar que se llamó al caso de uso como esperamos
	mockUseCase.AssertExpectations(t)
}

func TestGetUserHandlerNotFound(t *testing.T) {
	// Configurar el mock
	mockUseCase := new(MockUserUseCase)

	// Configurar router
	r := setupRouter()
	userGroup := r.Group("/api/users")

	// Registrar handler
	delivery.NewUserHandler(userGroup, mockUseCase)

	// Datos de prueba
	userID := primitive.NewObjectID().Hex()

	// Configurar comportamiento esperado del mock
	mockUseCase.On("GetUser", userID).Return(nil, errors.New("usuario no encontrado"))

	// Crear solicitud HTTP
	req, _ := http.NewRequest("GET", "/api/users/"+userID, nil)

	// Ejecutar solicitud
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verificaciones
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Parsear respuesta
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verificar estructura de respuesta
	assert.Equal(t, "error", response["status"])
	assert.Contains(t, response["error"], "no encontrado")

	// Verificar que se llamó al caso de uso como esperamos
	mockUseCase.AssertExpectations(t)
}
