package jwttokens

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/mamoru777/authservice2/internal/mylogger"

	"time"
)

const (
	AccessTokenDuration  = time.Hour
	RefreshTokenDuration = 30 * 24 * time.Hour
	SecretKey            = "mamoru" // Замените на ваш секретный ключ
)

func (t *Tokens) CreateAccessToken(userId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(AccessTokenDuration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}
func (t *Tokens) CreateRefreshToken(userId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(RefreshTokenDuration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}

func (t *Tokens) VerifyToken(tokenString string, logger *mylogger.Logger) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			expirationTime, expOk := claims["exp"].(float64)
			if !expOk {
				logger.Logger.Println("Неверный токен: пропущено 'exp' claim")
				err = errors.New("Неверный токен: пропущено 'exp' claim")
				return nil, err
			}
			expiration := time.Unix(int64(expirationTime), 0)
			if time.Now().After(expiration) {
				logger.Logger.Println("У токена вышел срок годности")
				err = errors.New("У токена вышел срок годности")
				return nil, err
			}
			return claims, nil
		}
		return claims, nil
	}
	logger.Logger.Println("Неверный токен")
	err = errors.New("Неверный токен")
	return nil, err
}
