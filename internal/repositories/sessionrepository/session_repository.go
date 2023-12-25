package sessionrepository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type SessionRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (ss *SessionRepository) Create(ctx context.Context, s *Session) error {
	const q = `
		INSERT INTO sessions (usr_id ,refresh_token, expire_at) 
			VALUES (:usr_id, :refresh_token, :expire_at)
	`
	_, err := ss.db.NamedExecContext(ctx, q, s)
	return err
}

func (ss *SessionRepository) Get(ctx context.Context, usrId uuid.UUID) (*Session, error) {
	const q = `
		SELECT usr_id, refresh_token, expire_at FROM sessions WHERE usr_id = $1
	`
	s := new(Session)
	err := ss.db.GetContext(ctx, s, q, usrId)
	return s, err
}

func (ss *SessionRepository) Update(ctx context.Context, s *Session) error {
	const q = `
		UPDATE sessions SET refresh_token=:refresh_token, expire_at=:expire_at  
			WHERE usr_id = :usr_id
	`
	_, err := ss.db.NamedExecContext(ctx, q, s)
	return err
}

func (ss *SessionRepository) Delete(ctx context.Context, usrId uuid.UUID) error {
	const q = `
		DELETE FROM sessions WHERE usr_id = $1
	`
	_, err := ss.db.ExecContext(ctx, q, usrId)
	return err
}

func (ss *SessionRepository) Validate(ctx context.Context, usrId uuid.UUID, refreshToken string) (bool, error) {
	const q = `
		SELECT expire_at FROM sessions WHERE usr_id = $1 AND refresh_token = $2
	`
	var expireAt time.Time
	err := ss.db.GetContext(ctx, expireAt, q, usrId, refreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	if time.Now().Before(expireAt) {
		return true, nil
	}

	return false, nil
}
