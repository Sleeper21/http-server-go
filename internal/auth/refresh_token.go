package auth

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

// Just generate a random 256-bit (32-byte) hex-encoded string

func MakeRefreshToken() (string, error) {
	// Create the array of 32 bytes (256 bits)
	bytesArr := make([]byte, 32)

	// Fill the array with random values
	_, err := rand.Read(bytesArr)
	if err != nil {
		log.Printf("error creating the refresh token: %s", err)
		return "", err
	}

	// Encode to a hex string
	refreshToken := hex.EncodeToString(bytesArr)

	return refreshToken, nil
}
