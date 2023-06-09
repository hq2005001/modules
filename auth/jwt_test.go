package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"testing"
	"time"
)

type CustomClaims struct {
	Name string `json:"name,omitempty"`
	jwt.RegisteredClaims
}

func TestToken(t *testing.T) {
	conf := JwtConf{
		Key:     "api",
		Refresh: "10h",
		Expired: "10m",
	}
	t.Log(Builder(&conf).SetClaims(CustomClaims{
		Name: "test",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "111",
		},
	}).Build())
}
