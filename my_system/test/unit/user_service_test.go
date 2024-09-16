package unit

import (
	"my_system/internal/user"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
	NewUserPassword string
}

func (m *MockUserRepository) ExistsUserName(userName string) (bool, error) {
	args := m.Called(userName)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) AddNewUser(user *user.User) error {
	m.NewUserPassword = user.Password

	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetTokenExpirationTime() time.Duration {
	return time.Duration(900) * time.Second
}

func (m *MockUserRepository) SetTokenIntoRedis(cacheKey string, tokenString string, expirationTime time.Duration) error {
	args := m.Called(cacheKey, tokenString, expirationTime)
	return args.Error(0)
}

func TestSignUp_userExists(t *testing.T) {
	// Mock ExistsUserName 方法，当 userName 为 "foo" 时返回 exists=true
	mockUserService := new(MockUserRepository)
	mockUserService.On("ExistsUserName", "foo").Return(true, nil)

	// 调用正式代码
	request := user.SignUpRequest{Username: "foo", Password: "password123"}
	token, err := user.SignUpService(request, mockUserService)

	// 断言返回结果
	assert.Error(t, err)
	// 由于exists=true，所以预期会返回这个错误
	assert.Contains(t, err.Error(), "用户名[foo]已存在")
	// 由于exists=true，所以预期返回的token为空
	assert.Equal(t, "", token)
	// 将断言设置回传给 testify
	mockUserService.AssertExpectations(t)
}

func TestSignUp_AddNewUser(t *testing.T) {
	mockUserService := new(MockUserRepository)

	// Mock ExistsUserName 方法，当 userName 为 "bar" 时返回 false
	mockUserService.On("ExistsUserName", "bar").Return(false, nil)

	// Mock AddNewUser 方法，判断密码是否被哈希化
	mockUserService.On("AddNewUser", mock.Anything).Return(nil)
	// SetTokenIntoRedis 方法什么都不做，直接返回 nil
	mockUserService.On("SetTokenIntoRedis", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// 调用正式代码
	request := user.SignUpRequest{Username: "bar", Password: "password123"}
	token, err := user.SignUpService(request, mockUserService)

	// 断言
	assert.NoError(t, err)
	// 密码是被hash的，和原始的不一样
	assert.NotEqual(t, "password123", mockUserService.NewUserPassword)
	assert.NotEmpty(t, token)

	mockUserService.AssertExpectations(t)
}
