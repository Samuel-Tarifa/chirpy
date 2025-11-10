package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now.UTC()),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn).UTC()),
		Subject:   userID.String(),
	})

	signed, err := token.SignedString([]byte(tokenSecret))

	if err != nil {
		return "", err
	}

	return signed, nil

}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(tokenSecret), nil

	})

	if err!=nil{
		return uuid.Nil,err
	}

	if !token.Valid{
		return uuid.Nil,errors.New("invalid token")
	}

	if claims.Subject==""{
		return uuid.Nil,errors.New("token missing subject")
	}

	id,err:=uuid.Parse(claims.Subject)

	if err!=nil{
		return uuid.Nil,err
	}

	return id,nil

}
