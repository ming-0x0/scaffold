package jwt

import (
	"maps"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func Decode(JWTToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(JWTToken, func(token *jwt.Token) (any, error) {
		return nil, nil
	})
	if strings.Contains(err.Error(), jwt.ErrInvalidKeyType.Error()) {
		return token, nil
	}
	return nil, err
}

func GenerateHS256JWT(payload map[string]any) (string, error) {
	claims := jwt.MapClaims{}
	maps.Copy(claims, payload)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	return signedToken, err
}

func VerifyJWT(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return false
	}

	return token.Valid
}
