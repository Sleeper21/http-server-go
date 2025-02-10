package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		testName    string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		expectError bool
	}{
		{
			testName:    "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			expectError: false,
		},
		{
			testName:    "Invalid token",
			tokenString: "invalid.Token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			expectError: true,
		},
		{
			testName:    "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			expectError: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.testName, func(t *testing.T) {
			gotUserID, err := ValidateJWT(testCase.tokenString, testCase.tokenSecret)
			if (err != nil) != testCase.expectError {
				t.Errorf("ValidateJWT() error = %v, expectedError %v", err, testCase.expectError)
				return
			}
			if gotUserID != testCase.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, expected UserID %v", gotUserID, testCase.wantUserID)
			}
		})
	}
}
