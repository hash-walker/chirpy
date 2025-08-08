package auth

import (
	"log"
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"net/http"
	"strings"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"crypto/rand"
	"encoding/hex"
)

func HashPassword(password string) (string, error){

	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Printf("Error generating hashed password: %v", err)
		return "", err
	}

	return string(hashed_password), nil
}

func CheckPasswordHash(hash string, password string) error{
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		return err
	}

	return nil

}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){

	claims := jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims) 

	signedToken, err := token.SignedString([]byte(tokenSecret))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error){

	parsedClaims := &jwt.RegisteredClaims{}
	
	token, err := jwt.ParseWithClaims(tokenString, parsedClaims, func (token *jwt.Token) (interface{}, error){
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token: %w", err)
	}

	claims := token.Claims.(*jwt.RegisteredClaims)

	userID, err := uuid.Parse(claims.Subject)

	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil

}

func GetBearerToken(headers http.Header) (string, error){
	authorization := headers.Get("Authorization")

	if authorization == ""{
		return "", fmt.Errorf("error getting the authorization header")
	}


	authorizationSplit := strings.Split(authorization, " ")

	return authorizationSplit[1], nil 
}

func MakeRefreshToken() (string, error){
	key := make([]byte, 32)
	rand.Read(key)
	refreshToken := hex.EncodeToString(key)

	return refreshToken, nil
}

func GetApiKey(headers http.Header) (string, error){
	authorization := headers.Get("Authorization")

	if authorization == ""{
		return "", fmt.Errorf("error getting the authorization header")
	}


	authorizationSplit := strings.Split(authorization, " ")

	return authorizationSplit[1], nil 
}