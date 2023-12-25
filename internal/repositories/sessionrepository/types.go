package sessionrepository

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	UsrId        uuid.UUID `db:"usr_id"`
	RefreshToken string    `db:"refresh_token"`
	ExpireAt     time.Time `db:"expire_at"`
}
