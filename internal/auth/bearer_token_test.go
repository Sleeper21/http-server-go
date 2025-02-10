package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {

	type headerType map[string]string

	validHeader := headerType{
		"Authorization": "Bearer 1234Abc",
	}

	invalidPrefix := headerType{
		"Authorization": "token 1234Abc",
	}

	emptyPrefix := headerType{
		"Authorization": "1234Abc",
	}

	emptyHeaderValue := headerType{
		"Authorization": "",
	}

	wrongHeaderKey := headerType{
		"auth": "Bearer 1234Abc",
	}

	emptyToken := headerType{
		"Authorization": "Bearer ",
	}

	tests := []struct {
		name                string
		header              headerType
		expectedBearerToken string
		expectedError       bool
	}{
		{
			name:                "Valid header",
			header:              validHeader,
			expectedBearerToken: "1234Abc",
			expectedError:       false,
		},
		{
			name:                "Invalid prefix",
			header:              invalidPrefix,
			expectedBearerToken: "",
			expectedError:       true,
		},
		{
			name:                "Empty prefix",
			header:              emptyPrefix,
			expectedBearerToken: "",
			expectedError:       true,
		},
		{
			name:                "Empty header value",
			header:              emptyHeaderValue,
			expectedBearerToken: "",
			expectedError:       true,
		},
		{
			name:                "Wrong header key",
			header:              wrongHeaderKey,
			expectedBearerToken: "",
			expectedError:       true,
		},
		{
			name:                "Empty token",
			header:              emptyToken,
			expectedBearerToken: "",
			expectedError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			for key, value := range tt.header {
				header := http.Header{
					key: []string{value},
				}
				bearerToken, err := GetBearerToken(header)
				if err != nil && !tt.expectedError {
					t.Errorf("got error: %s - expected Error: %v", err, tt.expectedError)
				}

				if bearerToken != tt.expectedBearerToken {
					t.Errorf("got token: %s - expected token: %s", bearerToken, tt.expectedBearerToken)
				}
			}

		})
	}
}
