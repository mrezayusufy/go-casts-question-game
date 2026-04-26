package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenRepository interface {
	Generate(userID uint, phone_number string) (string, int64, error)
	Validate(tokenString string) (uint, error)
}

type Token struct {
	secretKey      []byte
	expirationTime time.Duration
}

func NewToken(secretKey []byte) *Token {
	expirationTime := 1000
	return &Token{
		secretKey:      []byte(secretKey),
		expirationTime: time.Duration(expirationTime) * time.Hour,
	}
}

// methods
// generate
func (s *Token) Generate(userID uint) (string, error) {
	expiresAt := time.Now().Add(s.expirationTime)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}
	// pass the claims to jwt.NewWithClaims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

// validate
func (s *Token) Validate(tokenString string) (uint, error) {
	// parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})
	// error handling
	if err != nil {
		return 0, err
	}
	// claim and get user id
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(uint))
		return userID, nil
	}
	// return user id
	return 0, errors.New("Invalid Token")
}
