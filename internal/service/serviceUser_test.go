package service

import (
	"strconv"
	"testing"

	"BankingAPI/internal/model"
	"BankingAPI/internal/model/user"

	"github.com/stretchr/testify/assert"
)

var userMockDB = user.NewMockUserRepository()

func TestUserCreate(t *testing.T) {
	userService := NewUserService(userMockDB, nil)

	type table []struct {
		name          string
		userRequest   *user.User
		expectedUser  *user.User
		expectedError *model.Erro
	}
	testCases := table{
		{
			name: "Create User Succesfully",
			userRequest: &user.User{
				User_id:  "1",
				Name:     "Edson Cesar",
				Document: "12345678901234", // 14 digits
				Password: "edsonPass",
			},
			expectedUser: &user.User{
				User_id:  "1",
				Name:     "Edson Cesar",
				Document: "12345678901234", // 14 digits
				Password: "edsonPass",
			},
			expectedError: nil,
		},
		{
			name: "Failed to create user with empty name",
			userRequest: &user.User{
				User_id:  "1",
				Name:     "",
				Document: "12345678901234", // 14 digits
				Password: "edsonPass",
			},
			expectedUser:  nil,
			expectedError: ErrorMissingCredentials,
		},
		{
			name: "Failed to create user with empty document",
			userRequest: &user.User{
				User_id:  "1",
				Name:     "Edson Cesar",
				Document: "", // 14 digits
				Password: "edsonPass",
			},
			expectedUser:  nil,
			expectedError: ErrorMissingCredentials,
		},
		{
			name: "Failed to create user with empty password",
			userRequest: &user.User{
				User_id:  "1",
				Name:     "Edson Cesar",
				Document: "12345678901234", // 14 digits
				Password: "",
			},
			expectedUser:  nil,
			expectedError: ErrorMissingCredentials,
		},
	}

	for _, test := range testCases {
		got, err := userService.Create(test.userRequest)
		if err != test.expectedError {
			t.Errorf("Test %s failed: expected error %v, got %v", test.name, test.expectedError, err)
			continue
		}
		if !assert.Equal(t, test.expectedUser, got) {
			t.Fail()
		}
	}
}

func TestUserDelete(t *testing.T) {
	userService := NewUserService(userMockDB, nil)

	if _, err := userService.Create(&user.User{
		User_id:  "1",
		Name:     "Edson Cesar",
		Document: "12345678901234",
		Password: "edsonPass",
	}); err != nil {
		t.Fatalf("Failed to create user for delete test: %v", err)
		return
	}
	type table []struct {
		name          string
		id            string
		expectedError *model.Erro
	}
	testCases := table{
		{
			name:          "User deleted successfully",
			id:            "1",
			expectedError: nil,
		},
		{
			name:          "User not found",
			id:            "2",
			expectedError: model.IDnotFound,
		},
	}
	for _, test := range testCases {
		if err := userService.Delete(&test.id); err != test.expectedError {
			t.Errorf("Test %s failed: expected error %v, got %v", test.name, test.expectedError, err)
			continue
		}
	}
}

func TestUserGet(t *testing.T) {
	userService := NewUserService(userMockDB, nil)

	testUser := &user.User{
		User_id:  "1",
		Name:     "Edson Cesar",
		Document: "12345678901234",
		Password: "edsonPass",
	}

	if _, err := userService.Create(testUser); err != nil {
		t.Fatalf("Failed to create user for get test: %v", err)
		return
	}
	type table []struct {
		name          string
		id            string
		expectedUser  *user.User
		expectedError *model.Erro
	}
	testCases := table{
		{
			name:         "Get User Successfully",
			id:           "1",
			expectedUser: testUser,
		},
		{
			name:          "User not found",
			id:            "2",
			expectedUser:  nil,
			expectedError: model.IDnotFound,
		},
	}
	for _, test := range testCases {
		got, err := userService.Get(&test.id)
		if err != test.expectedError {
			t.Errorf("Test %s failed: expected error %v, got %v", test.name, test.expectedError, err)
			continue
		}
		if !assert.Equal(t, test.expectedUser, got) {
			t.Errorf("Test %s failed: expected user %v, got %v", test.name, test.expectedUser, got)
			continue
		}
	}
}

func TestUserUpdate(t *testing.T) {
	userService := NewUserService(userMockDB, nil)
	type table []struct {
		name          string
		userRequest   *user.User
		expectedUser  *user.User
		expectedError *model.Erro
	}
	testUser := &user.User{
		User_id:  "1",
		Name:     "Edson Cesar",
		Document: "12345678901234",
		Password: "edsonPass",
	}
	if _, err := userService.Create(testUser); err != nil {
		t.Fatalf("Failed to create user for update test: %v", err)
		return
	}
	testCases := table{
		{
			name: "Update User Successfully",
			userRequest: &user.User{
				User_id:  "1",
				Name:     "Edson Cesar Updated",
				Document: "12345678901234updated",
				Password: "edsonPassUpdated",
				Status:   false,
			},
			expectedUser: &user.User{
				User_id:  "1",
				Name:     "Edson Cesar Updated",
				Document: "12345678901234updated",
				Password: "edsonPassUpdated",
				Status:   false,
			},
			expectedError: nil,
		},
		{
			name: "Update User Only Name",
			userRequest: &user.User{
				User_id: "1",
				Name:    "Edson Cesar",
			},
			expectedUser: &user.User{
				User_id:  "1",
				Name:     "Edson Cesar",
				Document: "12345678901234updated",
				Password: "edsonPassUpdated",
				Status:   false,
			},
			expectedError: nil,
		},
		{
			name: "Update User Only DOcument",
			userRequest: &user.User{
				User_id:  "1",
				Document: "12345678901234",
			},
			expectedUser: &user.User{
				User_id:  "1",
				Name:     "Edson Cesar",
				Document: "12345678901234",
				Password: "edsonPassUpdated",
				Status:   false,
			},
			expectedError: nil,
		},
	}
	for _, test := range testCases {
		got, err := userService.Update(test.userRequest)
		if err != test.expectedError {
			t.Errorf("Test %s failed: expected error %v, got %v", test.name, test.expectedError, err)
			continue
		}
		if !assert.Equal(t, test.expectedUser, got) {
			t.Errorf("Test %s failed: expected user %v, got %v", test.name, test.userRequest, got)
			continue
		}
	}
}

func TestUserGetAll(t *testing.T) {
	userService := NewUserService(userMockDB, nil)
	var usersQTD int = 10
	users := make([]user.User, 0, usersQTD)
	for i := 0; i < usersQTD; i++ {
		testUser := &user.User{
			User_id:  strconv.Itoa(i + 1),
			Name:     "Edson Cesar " + strconv.Itoa(i+1),
			Document: "12345678901234",
			Password: "edsonPass",
		}
		if _, err := userService.Create(testUser); err != nil {
			t.Fatalf("Failed to create user for update test: %v", err)
			return
		} else {
			users = append(users, *testUser)
		}
	}
	type table []struct {
		name          string
		expectedUsers []user.User
		expectedError *model.Erro
	}
	/* testCases := table{
		{
			name:          "Get All Users Successfully",
			expectedUsers: users,
			expectedError: nil,
		},
		{
			name:          "No users found",
			expectedUsers: nil,
			expectedError: model.IDnotFound,
		},
		{
			name: "",
		},
	}
	*/
}
