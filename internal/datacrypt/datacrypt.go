package datacrypt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pwd string) (string, error) {
	hp, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hp), nil
}

func CheckPassword(pwd string, hpwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hpwd), []byte(pwd))
}

const TokenExp = time.Hour * 2
const MaxAge = 3600 * 2

type Claims struct {
	jwt.RegisteredClaims
	userID int
}

func BuildUserJWT(id int, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		userID: id,
	})

	return token.SignedString([]byte(key))
}

func GetUserID(unparsed string, key string) (int, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(unparsed, claims,
		func(t *jwt.Token) (any, error) {
			return []byte(key), nil
		})

	if err != nil {
		return -1, err
	}

	if !token.Valid {
		return -1, errors.New("token is not valid")
	}

	return claims.userID, nil
}
