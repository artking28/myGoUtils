package jwt

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CreateToken creates a JWT token using a secret key
func CreateToken(claims map[string]interface{}, customKey, subject string, expirationTimeInMillis int64) (string, error) {

	key := sha256.Sum256([]byte(customKey))
	tokenClaims := jwt.MapClaims{
		"sub": subject,
		"iat": time.Now().Unix(),
	}
	for k, v := range claims {
		tokenClaims[k] = v
	}

	if expirationTimeInMillis > 0 {
		tokenClaims["exp"] = time.Now().Add(time.Millisecond * time.Duration(expirationTimeInMillis)).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	return token.SignedString(key[:])
}

// IsValidJwt validates JWT tokens using a secret key
func IsValidJwt(tokenString, customKey string) bool {

	key := sha256.Sum256([]byte(customKey))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return key[:], nil
	})
	if err != nil {
		return false
	}

	return token.Valid
}

// GetValuesFrom extracts the values from JSW using a secret key
func GetValuesFrom(tokenString, customKey string) (map[string]interface{}, error) {

	key := sha256.Sum256([]byte(customKey))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return key[:], nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
