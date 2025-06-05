package service

import (
	"BankingAPI/internal/model"
	"BankingAPI/internal/model/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	createdUser := user.User{
		User_id:       "notIdTestUser",
		Name:          "NotRegularUser",
		Document:      "43210987654321",
		Password:      "wrongPass",
		Status:        false,
		Register_date: "2023-10-01T00:00:01Z",
	}
	type table []struct {
		name string
		arg  string
		want user.User
		err  *model.Erro
	}
	tests := table{
		{
			name: "Test Not Pass",
			arg:  "idTestUser",
			want: user.User{
				User_id:       "idTestUser",
				Name:          "RegularUser",
				Document:      "12345678901234",
				Password:      "testPass",
				Status:        true,
				Register_date: "2023-10-01T00:00:00Z",
			},
			err: model.IDnotFound,
		},
		{
			name: "Test Pass",
			arg:  "notIdTestUser",
			want: createdUser,
			err:  nil,
		},
		{
			name: "Not Found",
			arg:  "notFoundId",
			want: user.User{},
			err:  model.IDnotFound,
		},
	}

	userDatabaseMock := user.NewMockUserRepository()
	serviceGet := NewGetService(nil, nil, userDatabaseMock)

	serviceCreate := NewCreateService(nil, nil, userDatabaseMock, serviceGet)
	if _, err := serviceCreate.CreateUser(&createdUser); err != nil {
		assert.Fail(t, "Failed to create user in mock repository: %v", err.Err)
		return
	}

	for _, test := range tests {
		got, err := serviceGet.User(test.arg)
		if got == nil {
			if !assert.Error(t, test.err.Err, err.Err, test.name) {
				assert.Fail(t, "Expected error: %v, got: %v", test.err.Err, err.Err, test.name)
			}
			continue
		}
		if err != nil {
			assert.Error(t, test.err.Err, err.Err, test.name)
		}
		if !assert.Equal(t, test.want, (*got), test.name) {
			assert.Fail(t, "Expected: %v, got: %v", test.want, *got, test.name)
		}
	}
}
