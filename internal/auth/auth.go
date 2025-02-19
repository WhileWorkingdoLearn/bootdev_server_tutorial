package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const cost = 13
const pw_length = 1

func HashPassword(password string) (string, error) {
	if len(password) < pw_length {
		return "", fmt.Errorf("password to short")
	}
	data, errGPW := bcrypt.GenerateFromPassword([]byte(password), cost)
	if errGPW != nil {
		return "", errGPW
	}

	return string(data), nil
}

func CheckPasswordHash(password string, hash string) error {
	if len(password) < pw_length {
		return fmt.Errorf("password to short")
	}

	errPW := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if errPW != nil {
		return errPW
	}
	return nil
}

func MakeJWT(userId uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	newtoken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chripy",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn).UTC()),
			Subject:   userId.String(),
		})

	token, errSign := newtoken.SignedString([]byte(tokenSecret))
	if errSign != nil {
		return "", errSign
	}
	return token, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	token, errToken := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if errToken != nil {
		return uuid.UUID{}, errToken
	}

	subject, errSub := token.Claims.GetSubject()
	if errSub != nil {
		return uuid.UUID{}, errSub
	}
	uuId, errParse := uuid.Parse(subject)
	if errParse != nil {
		return uuid.UUID{}, errParse
	}

	return uuId, nil
}

func GetToken(header http.Header, keyType string) (string, error) {
	result := strings.Split(header.Get("Authorization"), " ")

	if len(result) != 2 {
		return "", fmt.Errorf("wrong authroization header format")
	}

	if strings.Trim(result[0], " ") != keyType {
		return "", fmt.Errorf("missing Bearer")
	}

	return strings.Trim(result[1], " "), nil
}

func MakeRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
