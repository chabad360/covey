package authentication

import (
	"fmt"
	"strings"
	"time"

	"github.com/chabad360/covey/common"
	"github.com/gbrlsnchs/jwt/v3"
)

const key = "asdf" // TODO: Redesign API key system

var (
	crashKey string
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func createToken(userid string, tokenType string, audience []string) (string, *time.Time, error) {
	refreshKey()
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
		JWTID:          common.RandomString(),
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
	refreshKey()
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

func refreshKey() {
	if crashKey == "" {
		// TODO: Don't release with this
		crashKey = fmt.Sprintf("12345")
		// crashKey = commmon.RandomString()
	}
}
