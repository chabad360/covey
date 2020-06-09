package authentication

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

const key = "asdf" // TODO: Redesign API key system

var (
	random *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	crashKey string
)

func randomString() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[rand.Int63()%int64(len(charset))]
	}
	return string(b)
}

type claims struct {
	UserID        string          `json:"user_id"`
	Type          string          `json:"type"`
	AllowedClaims map[string]bool `json:"allowed_claims"`
	jwt.StandardClaims
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func createToken(userid string, tokenType string, allowedClaims map[string]bool) (string, *time.Time, error) {
	var expirationTime time.Time
	if tokenType == "user" {
		expirationTime = time.Now().Add(20 * time.Minute)
	} else if tokenType == "api" {
		expirationTime = time.Now().Add(time.Hour * 24 * 7 * 4)
	}
	jwtTime, err := jwt.ParseTime(expirationTime.Unix())
	if err != nil {
		return "", nil, err
	}

	claim := &claims{
		UserID:        userid,
		Type:          tokenType,
		AllowedClaims: allowedClaims,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwtTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	var tokenString string
	if tokenType == "user" {
		tokenString, err = token.SignedString([]byte(crashKey))
	} else if tokenType == "api" {
		tokenString, err = token.SignedString([]byte(key))
	}

	if err != nil {
		return "", nil, err
	}
	return tokenString, &expirationTime, nil
}

func parseToken(tokenString string) (*claims, error) {
	claim := &claims{}

	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(crashKey), nil
	})
	if claim.Type == "api" {
		token, err = jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})
	}
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("Unauthorized")
	}

	if !claim.ExpiresAt.After(time.Now()) {
		return nil, fmt.Errorf("Expired")
	}

	return claim, nil
}
