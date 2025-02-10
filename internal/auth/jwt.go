package auth

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	// TokenType is the token type for access tokens
	TokenTypeAccess TokenType = "chirpy"
)

// this function creates a jwt token for a user and signs it with the secret
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	signingKey := []byte(tokenSecret)

	// Create new token
	regClaims := jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		Subject:   userID.String(), // converts the uuid to string
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, regClaims)

	// Sign the token with the secret
	tokenStr, err := token.SignedString(signingKey)
	if err != nil {
		log.Printf("error signing token: %s", err)
		return "", err
	}
	return tokenStr, nil
}

// this function validates a received jwt token and extracts the user id
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// Use the ParseWithClaims functions to validate the signature of the JWT and extract the claims
	claimsStruct := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claimsStruct, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return id, nil
}
