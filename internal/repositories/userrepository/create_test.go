package userrepository

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	var pass []byte
	pass = []byte("$2a$10$eu9gKHpYI36Q7OVZKGCgzuWiimSOUwsVuoXgL5JpuYZIOGg0n2Qvi")
	type mockBehavior func(u *User)
	testTable := []struct {
		name         string
		user         *User
		mockBehavior mockBehavior
		expectedErr  error
	}{
		{
			name: "FullInfo",
			user: &User{
				Login:    "testlogin",
				Email:    "test@example.com",
				Password: pass,
			},
			mockBehavior: func(u *User) {
				mock.ExpectExec("INSERT INTO usrs").
					WithArgs(u.Login, u.Email, u.Password).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			userRepo := New(sqlxDB)
			testCase.mockBehavior(testCase.user)
			err = userRepo.Create(context.Background(), testCase.user)
			assert.Equal(t, testCase.expectedErr, err)
			err = mock.ExpectationsWereMet()
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}
