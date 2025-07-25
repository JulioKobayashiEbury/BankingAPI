package service

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"BankingAPI/internal/model"
	"BankingAPI/internal/model/user"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserCreate(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	userMockDB := user.NewMockUserRepository()
	userService := NewUserService(userMockDB, nil)

	type table []struct {
		name          string
		userRequest   *user.User
		expectedUser  *user.User
		expectedError *echo.HTTPError
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
			expectedError: model.ErrMissingCredentials,
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
			expectedError: model.ErrMissingCredentials,
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
			expectedError: model.ErrMissingCredentials,
		},
	}

	for _, test := range testCases {
		got, err := userService.Create(ctx, test.userRequest)
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
	ctx := context.Background()
	defer ctx.Done()

	userMockDB := user.NewMockUserRepository()
	userService := NewUserService(userMockDB, nil)

	if _, err := userService.Create(ctx, &user.User{
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
		expectedError *echo.HTTPError
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
			expectedError: model.ErrIDnotFound,
		},
	}
	for _, test := range testCases {
		if err := userService.Delete(ctx, &test.id); err != test.expectedError {
			t.Errorf("Test %s failed: expected error %v, got %v", test.name, test.expectedError, err)
			continue
		}
	}
}

func TestUserGet(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	userMockDB := user.NewMockUserRepository()
	userService := NewUserService(userMockDB, nil)

	testUser := &user.User{
		User_id:  "1",
		Name:     "Edson Cesar",
		Document: "12345678901234",
		Password: "edsonPass",
	}

	if _, err := userService.Create(ctx, testUser); err != nil {
		t.Fatalf("Failed to create user for get test: %v", err)
		return
	}
	type table []struct {
		name          string
		id            string
		expectedUser  *user.User
		expectedError *echo.HTTPError
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
			expectedError: model.ErrIDnotFound,
		},
	}
	for _, test := range testCases {
		got, err := userService.Get(ctx, &test.id)
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
	ctx := context.Background()
	defer ctx.Done()

	userMockDB := user.NewMockUserRepository()
	userService := NewUserService(userMockDB, nil)
	type table []struct {
		name          string
		userRequest   *user.User
		expectedUser  *user.User
		expectedError *echo.HTTPError
	}
	testUser := &user.User{
		User_id:  "1",
		Name:     "Edson Cesar",
		Document: "12345678901234",
		Password: "edsonPass",
	}
	if _, err := userService.Create(ctx, testUser); err != nil {
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
			},
			expectedUser: &user.User{
				User_id:  "1",
				Name:     "Edson Cesar Updated",
				Document: "12345678901234updated",
				Password: "edsonPassUpdated",
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
			},
			expectedError: nil,
		},
	}
	for _, test := range testCases {
		got, err := userService.Update(ctx, test.userRequest)
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
	ctx := context.Background()
	defer ctx.Done()

	userMockDB := user.NewMockUserRepository()
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
		if _, err := userService.Create(ctx, testUser); err != nil {
			t.Fatalf("Failed to create user for update test: %v", err)
			return
		} else {
			users = append(users, *testUser)
		}
	}
	type table []struct {
		name          string
		expectedUsers *[]user.User
		expectedError *echo.HTTPError
	}
	testCases := table{
		{
			name:          "Get All Users Successfully",
			expectedUsers: &users,
			expectedError: nil,
		},
		{
			name:          "No users found",
			expectedUsers: nil,
			expectedError: model.ErrIDnotFound,
		},
	}
	for _, test := range testCases {
		got, err := userService.GetAll(ctx)
		if err != test.expectedError {
			t.Errorf("Test %s failed: expected error %v, got %v", test.name, test.expectedError, err)
			continue
		}
		if !assert.Equal(t, test.expectedUsers, got) {
			t.Errorf("Test %s failed: expected users %v, got %v", test.name, test.expectedUsers, *got)
			continue
		}
		if test.name != "No users found" {
			for _, user := range users {
				if err := userService.Delete(ctx, &user.User_id); err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}
}
