package service

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

type JWTService interface {
	GenerateToken(name string, admin bool) (string, int64)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

// jwtCustomClaims are custom claims extending default ones.
type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService() JWTService {
	return &jwtService{
		secretKey: getSecretKey(),
		issuer:    "sms.vcs.com",
	}
}

func getSecretKey() string {
	err := godotenv.Load()
	if err != nil {
		godotenv.Load("./../.env")
		// panic("Failed to load env file")
	}
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		secret = "secret"
	}
	return secret
}

func (jwtSrv *jwtService) GenerateToken(username string, admin bool) (string, int64) {

	expirationTime := time.Now().Add(time.Hour * 24).Unix()

	// Set custom and standard claims
	claims := &jwtCustomClaims{
		username,
		admin,
		jwt.StandardClaims{
			ExpiresAt: expirationTime,
			Issuer:    jwtSrv.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token using the secret signing key
	encoded_token, err := token.SignedString([]byte(jwtSrv.secretKey))
	if err != nil {
		panic(err)
	}
	return encoded_token, expirationTime
}

func (jwtSrv *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Signing method validation
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret signing key
		return []byte(jwtSrv.secretKey), nil
	})
}
