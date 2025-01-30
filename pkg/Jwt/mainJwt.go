package Jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtDate struct {
	Email string
}

type JwtSecret struct {
	Secret string
}

func NewJWT(secret string) *JwtSecret {
	return &JwtSecret{
		Secret: secret,
	}
}

func (j *JwtSecret) Create(date JwtDate) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": date.Email,
	})
	s, err := t.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return s, nil
}

func (j *JwtSecret) Parse(token string) (bool, *JwtDate) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		return false, nil
	}
	email := t.Claims.(jwt.MapClaims)["email"].(string)
	return t.Valid, &JwtDate{
		Email: email,
	}
}
