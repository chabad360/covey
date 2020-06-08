package authentication

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/chabad360/covey/common"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gorilla/mux"
)

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
	UserID uint32 `json:"user_id"`
	jwt.StandardClaims
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func createToken(userid uint32) (string, error) {
	expirationTime := time.Now().Add(20 * time.Minute)
	eTime := expirationTime.Unix()
	jwtTime, err := jwt.ParseTime(eTime)
	if err != nil {
		return "", err
	}

	claim := &claims{
		UserID: userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwtTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(crashKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func parseToken(tokenString string) (uint32, error) {
	claim := &claims{}

	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(crashKey), nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, fmt.Errorf("Unauthorized")
	}

	return claim.UserID, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/auth/token" {
			next.ServeHTTP(w, r)
			return
		}
		var tokenString string

		c, err := r.Cookie("token")
		if err == nil {
			tokenString = c.Value
		} else {
			if err != http.ErrNoCookie {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		header := r.Header.Get("Authorization")
		if header != "" {
			splitToken := strings.Split(header, "Bearer ")
			tokenString = splitToken[1]
		}

		if tokenString == "" {
			common.ErrorWriterCustom(w, fmt.Errorf("Unauthorized"), http.StatusUnauthorized)
			return
		}

		_, err = parseToken(tokenString)
		if err != nil {
			common.ErrorWriter(w, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func tokenGet(w http.ResponseWriter, r *http.Request) {
	user := &credentials{}
	var ok bool
	user.Username, user.Password, ok = r.BasicAuth()
	if !ok {
		common.ErrorWriterCustom(w, fmt.Errorf("Unauthorized"), http.StatusUnauthorized)
		return
	}

	id, err := GetUser(*user)
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}

	token, err := createToken(id)
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{ Token string }{Token: token})
}

func RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/token", tokenGet).Methods("GET")

	crashKey = randomString()
}
