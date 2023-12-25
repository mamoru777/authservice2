package userrepository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
func (us *UserRepository) Create(ctx context.Context, u *User) error {
	const q = `
		INSERT INTO usrs (login, email, password, isSignedUp) 
			VALUES (:login, :email, :password, false)
	`
	_, err := us.db.NamedExecContext(ctx, q, u)
	return err
}

func (us *UserRepository) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	const q = `
		SELECT id, login, email, password, isSignedUp FROM usrs WHERE id = $1
	`
	u := new(User)
	err := us.db.GetContext(ctx, u, q, id)
	return u, err
}

func (us *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	const q = `
		SELECT id, login, email, password, isSignedUp FROM usrs WHERE email = $1
	`
	u := new(User)
	err := us.db.GetContext(ctx, u, q, email)
	return u, err
}

func (us *UserRepository) GetByEmailCheck(ctx context.Context, email string) (bool, error) {
	const q = `
		SELECT id, login, email, password, isSignedUp FROM usrs WHERE email = $1
	`
	u := new(User)
	err := us.db.GetContext(ctx, u, q, email)
	if err != nil {
		return true, err
	} else {
		return false, err
	}
}

func (us *UserRepository) GetByLoginCheck(ctx context.Context, login string) (bool, error) {
	const q = `
		SELECT id, login, email, password, isSignedUp FROM usrs WHERE login = $1
	`
	u := new(User)
	err := us.db.GetContext(ctx, u, q, login)
	if err != nil {
		return true, err
	} else {
		return false, err
	}
}

func (us *UserRepository) GetByLogin(ctx context.Context, login string) (*User, error) {
	const q = `
		SELECT id, login, email, password, isSignedUp FROM usrs WHERE login = $1
	`
	u := new(User)
	err := us.db.GetContext(ctx, u, q, login)
	return u, err
}

func (us *UserRepository) GetByUserAndPassword(ctx context.Context, login string) (*User, error) {
	const q = `
		SELECT id, login, email, password, isSignedUp FROM usrs WHERE login = $1
	`
	u := new(User)
	err := us.db.GetContext(ctx, u, q, login)
	return u, err
}

func (us *UserRepository) List(ctx context.Context) ([]*User, error) {
	const q = `
		SELECT id, login, email, password, isSignedUp FROM usrs 
	`
	var list []*User
	err := us.db.SelectContext(ctx, list, q)
	return list, err
}

func (us *UserRepository) Update(ctx context.Context, u *User) error {
	const q = `
		UPDATE usrs SET login=:login, email=:email, password=:password, isSignedUp=:isSignedUp 
			WHERE id = :id
	`
	_, err := us.db.NamedExecContext(ctx, q, u)
	return err
}

func (us *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `
		DELETE FROM usrs WHERE id = $1
	`
	_, err := us.db.ExecContext(ctx, q, id)
	return err
}
