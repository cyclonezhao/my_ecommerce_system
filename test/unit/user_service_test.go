package unit

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"my_ecommerce_system/internal/user"
	"testing"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) ExistsUserName(userName string) (bool, error) {
	args := m.Called(userName)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) AddNewUser(user *user.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestSignUp_userExists(t *testing.T) {
	// 构造参数：注册用户名为 foo 的用户
	request := user.SignUpRequest{Username: "foo", Password: "password123"}

	// Mock ExistsUserName 方法，当 userName 为 "foo" 时返回 exists=true
	mockUserService := new(MockUserRepository)
	mockUserService.On("ExistsUserName", "foo").Return(true, nil)

	// 调用正式代码
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