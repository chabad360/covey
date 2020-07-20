package authentication

import (
	"fmt"
	"github.com/chabad360/covey/models"
	"net/http"
	"strings"
	"time"

	"github.com/chabad360/covey/common"
	"github.com/gbrlsnchs/jwt/v3"
)

var (
	crashKey string //nolint:gochecknoglobals
)

const (
	user = "user"
	api  = "api"
	key  = "asdfg" // TODO: Redesign API key system
)

func createToken(userID string, tokenType string, audience []string) (string, *time.Time, error) {
	refreshKey()

	var err error
	var expirationTime time.Time

	if userID == "" {
		return "", nil, fmt.Errorf("createToken: missing userID")
	}

	switch tokenType {
	case user:
		// TODO: implement refresh token
		expirationTime = time.Now().Add(12 * time.Hour)
	case api:
		expirationTime = time.Now().Add(4 * (7 * (24 * time.Hour)))
	default:
		return "", nil, fmt.Errorf("createToken: invalid token type")
	}

	claim := jwt.Payload{
		Issuer:         strings.Join([]string{"covey", tokenType}, "-"),
		Subject:        userID,
		Audience:       audience,
		ExpirationTime: jwt.NumericDate(expirationTime),
		IssuedAt:       jwt.NumericDate(time.Now()),
		JWTID:          common.RandomString(),
	}

	var token []byte
	switch tokenType {
	case user:
		token, err = jwt.Sign(claim, jwt.NewHS256([]byte(crashKey)))
	case api:
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

	switch tokenType {
	case user:
		_, err = jwt.Verify([]byte(tokenString), jwt.NewHS256([]byte(crashKey)), &claim, validate)
	case api:
		_, err = jwt.Verify([]byte(tokenString), jwt.NewHS256([]byte(key)), &claim, validate)
	default:
		return nil, fmt.Errorf("parseToken: invalid token type")
	}
	if err != nil {
		return &claim, err
	}

	return &claim, nil
}

func tokenCookie(user *models.User) (*http.Cookie, error) {
	id, err := GetUser(*user)
	if err != nil {
		return nil, err
	}

	token, eTime, err := createToken(id, "user", []string{"all"})
	if err != nil {
		return nil, fmt.Errorf("internal error")
	}

	return &http.Cookie{
		Name: "token",
		// Domain:   r.Host, // ?
		Value:   token,
		Expires: *eTime,
		MaxAge:  int(time.Until(*eTime).Seconds()),
		Path:    "/",
	}, nil
}

func refreshKey() {
	if crashKey == "" {
		crashKey = common.RandomString()
	}
}
