package jwt

import (
	"crypto/rand"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int    `json:"user_id"`
	NIP    string `json:"nip"`
	Role   string `json:"role"`
	IsTemp bool   `json:"is_temp,omitempty"`
	jwt.RegisteredClaims
}

var secretKey = loadSecretKey()

// loadSecretKey reads JWT_SECRET. Without it a random per-process key is used, which keeps tokens
// unforgeable but logs everyone out on restart, so production must set JWT_SECRET.
func loadSecretKey() []byte {
	if s := strings.TrimSpace(os.Getenv("JWT_SECRET")); s != "" {
		return []byte(s)
	}
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic("jwt: JWT_SECRET is not set and no random key could be generated: " + err.Error())
	}
	log.Println("[JWT Warning] JWT_SECRET tidak diset; memakai kunci acak sementara (semua sesi hilang saat restart).")
	return key
}

func GenerateToken(userID int, nip string, role string, expiration time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		NIP:    nip,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func GenerateTempToken(userID int, expiration time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		IsTemp: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
