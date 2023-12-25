package jwttokens

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/mamoru777/authservice2/internal/mylogger"
)

//go:generate mockgen -destination mock_jwt_tokens.go -package jwttokens . ITokens

type ITokens interface {
	CreateAccessToken(userId uuid.UUID) (string, error)
	CreateRefreshToken(userId uuid.UUID) (string, error)
	VerifyToken(tokenString string, logger *mylogger.Logger) (jwt.MapClaims, error)
}
