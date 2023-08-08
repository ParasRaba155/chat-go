// jwt utility
package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"app/config"
)

var (
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
)

// claims emebeds the jwt-go StandardClaims
type claims struct {
	jwt.StandardClaims
	Email string `json:"email"`
}

type Service struct {
	config *config.Config
}

func NewService(c *config.Config) Service {
	return Service{
		config: c,
	}
}

// GenerateJWT will generate a jwt token with given email id
func (s Service) GenerateJWT(email string) (string, error) {
	claim := claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(s.config.Server.JWTExpireMinutes)).Unix(),
		},
		Email: email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString([]byte(s.config.Server.JwtSecretKey))
}

// ExtractJWTClaimsFromToken will return the jwt.MapClaims from the token string
func (s Service) ExtractJWTClaimsFromToken(tokenStr string) (map[string]interface{}, error) {
	claim := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claim, func(t *jwt.Token) (jwtKey interface{}, err error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrInvalidSigningMethod, t.Header["alg"])
		}
		jwtKey = []byte(s.config.Server.JwtSecretKey)
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claim, nil
}
