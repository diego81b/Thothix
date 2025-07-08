package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"thothix-backend/internal/shared/dto"
	usersDto "thothix-backend/internal/users/dto"
)

// MockUserService is a mock implementation of the UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserByID(userID string) *usersDto.GetUserResponse {
	args := m.Called(userID)
	return args.Get(0).(*usersDto.GetUserResponse)
}

func (m *MockUserService) GetUserByClerkID(clerkID string) *usersDto.GetUserResponse {
	args := m.Called(clerkID)
	return args.Get(0).(*usersDto.GetUserResponse)
}

func (m *MockUserService) GetUsers(req *dto.PaginationRequest) *usersDto.GetUsersResponse {
	args := m.Called(req)
	return args.Get(0).(*usersDto.GetUsersResponse)
}

func (m *MockUserService) CreateUser(req *usersDto.CreateUserRequest) *usersDto.CreateUserResponse {
	args := m.Called(req)
	return args.Get(0).(*usersDto.CreateUserResponse)
}

func (m *MockUserService) UpdateUser(userID string, req *usersDto.UpdateUserRequest) *usersDto.UpdateUserResponse {
	args := m.Called(userID, req)
	return args.Get(0).(*usersDto.UpdateUserResponse)
}

func (m *MockUserService) DeleteUser(userID string) *usersDto.DeleteUserResponse {
	args := m.Called(userID)
	return args.Get(0).(*usersDto.DeleteUserResponse)
}

func (m *MockUserService) SyncUserFromClerk(req *usersDto.ClerkUserSyncRequest) *usersDto.CreateUserResponse {
	args := m.Called(req)
	return args.Get(0).(*usersDto.CreateUserResponse)
}

func (m *MockUserService) ProcessClerkWebhook(userData interface{}) *usersDto.ClerkSyncUserResponse {
	args := m.Called(userData)
	return args.Get(0).(*usersDto.ClerkSyncUserResponse)
}

type UserHandlerTestSuite struct {
	suite.Suite
	handler     *UserHandler
	mockService *MockUserService
	router      *gin.Engine
}

func (suite *UserHandlerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (suite *UserHandlerTestSuite) SetupTest() {
	suite.mockService = new(MockUserService)
	suite.handler = NewUserHandler(suite.mockService)

	suite.router = gin.New()
	suite.router.GET("/users/:id", suite.handler.GetUserByID)
	suite.router.POST("/users", suite.handler.CreateUser)
	suite.router.GET("/users", suite.handler.GetUsers)
	suite.router.PUT("/users/:id", suite.handler.UpdateUser)
	suite.router.DELETE("/users/:id", suite.handler.DeleteUser)
}

func (suite *UserHandlerTestSuite) TestNewUserHandler() {
	// Act
	handler := NewUserHandler(suite.mockService)

	// Assert
	assert.NotNil(suite.T(), handler)
	assert.Equal(suite.T(), suite.mockService, handler.userService)
}

func (suite *UserHandlerTestSuite) TestGetUserByID_Success() {
	// Arrange
	userResponse := &usersDto.UserResponse{
		ID:    "test-id",
		Email: "test@example.com",
		Name:  "Test User",
	}

	mockResponse := usersDto.NewGetUserResponse(func() dto.Validation[*usersDto.UserResponse] {
		return dto.Success(userResponse)
	})

	suite.mockService.On("GetUserByID", "test-id").Return(mockResponse)

	// Act
	req, _ := http.NewRequest("GET", "/users/test-id", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *UserHandlerTestSuite) TestGetUserByID_NotFound() {
	// Arrange
	mockResponse := usersDto.NewGetUserResponse(func() dto.Validation[*usersDto.UserResponse] {
		return dto.Invalid[*usersDto.UserResponse](dto.NewError("USER_NOT_FOUND", "User not found", nil))
	})

	suite.mockService.On("GetUserByID", "nonexistent").Return(mockResponse)

	// Act
	req, _ := http.NewRequest("GET", "/users/nonexistent", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *UserHandlerTestSuite) TestCreateUser_Success() {
	// Arrange
	createReq := &usersDto.CreateUserRequest{
		Email: "test@example.com",
		Name:  "Test User",
	}

	userResponse := &usersDto.UserResponse{
		ID:    "new-id",
		Email: "test@example.com",
		Name:  "Test User",
	}

	mockResponse := usersDto.NewCreateUserResponse(func() dto.Validation[*usersDto.UserResponse] {
		return dto.Success(userResponse)
	})

	suite.mockService.On("CreateUser", mock.MatchedBy(func(req *usersDto.CreateUserRequest) bool {
		return req.Email == "test@example.com" && req.Name == "Test User"
	})).Return(mockResponse)

	// Act
	reqBody, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
}

func (suite *UserHandlerTestSuite) TestCreateUser_ValidationError() {
	// Arrange
	createReq := &usersDto.CreateUserRequest{
		Email: "", // Invalid
		Name:  "Test User",
	}

	mockResponse := usersDto.NewCreateUserResponse(func() dto.Validation[*usersDto.UserResponse] {
		return dto.Invalid[*usersDto.UserResponse](dto.NewError("VALIDATION_ERROR", "Email is required", nil))
	})

	suite.mockService.On("CreateUser", mock.AnythingOfType("*dto.CreateUserRequest")).Return(mockResponse)

	// Act
	reqBody, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *UserHandlerTestSuite) TestUpdateUser_Success() {
	// Arrange
	newEmail := "updated@example.com"
	updateReq := &usersDto.UpdateUserRequest{
		Email: &newEmail,
	}

	userResponse := &usersDto.UserResponse{
		ID:    "test-id",
		Email: "updated@example.com",
		Name:  "Test User",
	}

	mockResponse := usersDto.NewUpdateUserResponse(func() dto.Validation[*usersDto.UserResponse] {
		return dto.Success(userResponse)
	})

	suite.mockService.On("UpdateUser", "test-id", mock.AnythingOfType("*dto.UpdateUserRequest")).Return(mockResponse)

	// Act
	reqBody, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PUT", "/users/test-id", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_Success() {
	// Arrange
	mockResponse := usersDto.NewDeleteUserResponse(func() dto.Validation[string] {
		return dto.Success("User deleted successfully")
	})

	suite.mockService.On("DeleteUser", "test-id").Return(mockResponse)

	// Act
	req, _ := http.NewRequest("DELETE", "/users/test-id", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *UserHandlerTestSuite) TestGetUsers_Success() {
	// Arrange
	userListResponse := &usersDto.UserListResponse{}

	mockResponse := usersDto.NewGetUsersResponse(func() dto.Validation[*usersDto.UserListResponse] {
		return dto.Success(userListResponse)
	})

	suite.mockService.On("GetUsers", mock.AnythingOfType("*dto.PaginationRequest")).Return(mockResponse)

	// Act
	req, _ := http.NewRequest("GET", "/users?page=1&per_page=10", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
