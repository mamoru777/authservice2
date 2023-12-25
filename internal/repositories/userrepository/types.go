package userrepository

import "github.com/google/uuid"

type User struct {
	Id         uuid.UUID `db:"id"`
	Login      string    `db:"login"`
	Email      string    `db:"email"`
	Password   []byte    `db:"password"`
	IsSignedUp bool      `db:"isSignedUp"`
}
