package authentication

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
)

const key = "asdf" // TODO: Redesign API key system

var (
	random = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	crashKey = randomString()
)

func randomString() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[rand.Int63()%int64(len(charset))]
	}
	return string(b)
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func createToken(userid string, tokenType string, audience []string) (string, *time.Time, error) {
	var err error
	var expirationTime time.Time
	if tokenType == "user" {
		expirationTime = time.Now().Add(20 * time.Minute)
	} else if tokenType == "api" {
		expirationTime = time.Now().Add(4 * (7 * (24 * time.Hour)))
	}

	if userid == "" {
		return "", nil, fmt.Errorf("createToken: missing userid")
	}

	claim := jwt.Payload{
		Issuer:         strings.Join([]string{"covey", tokenType}, "-"),
		Subject:        userid,
		Audience:       audience,
		ExpirationTime: jwt.NumericDate(expirationTime),
		IssuedAt:       jwt.NumericDate(time.Now()),
		JWTID:          randomString(),
	}

	var token []byte
	if tokenType == "user" {
		token, err = jwt.Sign(claim, jwt.NewHS256([]byte(crashKey)))
	} else if tokenType == "api" {
		token, err = jwt.Sign(claim, jwt.NewHS256([]byte(key)))
	}
	if err != nil {
		return "", nil, err
	}
	return string(token), &expirationTime, nil
}

func parseToken(tokenString string, tokenType string, audience string) (*jwt.Payload, error) {
	var claim jwt.Payload
	var err error
	iss := jwt.IssuerValidator(strings.Join([]string{"covey", tokenType}, "-"))
	exp := jwt.ExpirationTimeValidator(time.Now())
	aud := jwt.AudienceValidator(jwt.Audience{audience})
	validate := jwt.ValidatePayload(&claim, iss, exp, aud)

	if tokenType == "user" {
		_, err = jwt.Verify([]byte(tokenString), jwt.NewHS256([]byte(crashKey)), &claim, validate)
	} else if tokenType == "api" {
		_, err = jwt.Verify([]byte(tokenString), jwt.NewHS256([]byte(key)), &claim, validate)
	} else {
		return nil, fmt.Errorf("parseToken: invalid token type")
	}
	if err != nil {
		return &claim, err
	}

	return &claim, nil
}
