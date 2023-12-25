package sessionrepository

import (
	"context"
	"github.com/google/uuid"
)

//go:generate mockgen -destination mock_session_repository.go -package sessionrepository . ISessionRepository

type ISessionRepository interface {
	Create(ctx context.Context, s *Session) error
	Get(ctx context.Context, usrId uuid.UUID) (*Session, error)
	Update(ctx context.Context, s *Session) error
	Delete(ctx context.Context, usrId uuid.UUID) error
	Validate(ctx context.Context, usrId uuid.UUID, refreshToken string) (bool, error)
}
