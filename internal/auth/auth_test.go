package auth

import "testing"

func TestCheckPasswordHash(t *testing.T) {

	// store some hashed passwords to test
	password1 := "someTestingPassword123!"
	password2 := "anotherPassword_123"

	hashed1, _ := HashPassword(password1)
	hashed2, _ := HashPassword(password2)

	cases := []struct {
		name        string
		password    string
		hash        string
		expectError bool
	}{
		{
			name:        "test a correct Password",
			password:    password1,
			hash:        hashed1,
			expectError: false,
		},
		{
			name:        "test an incorrect Password",
			password:    "myPassword",
			hash:        hashed1,
			expectError: true,
		},
		{
			name:        "test password does not match other hashes",
			password:    password1,
			hash:        hashed2,
			expectError: true,
		},
		{
			name:        "test empty Password",
			password:    "   ",
			hash:        hashed1,
			expectError: true,
		},
		{
			name:        "test an incorrect Hash",
			password:    password2,
			hash:        "someMadeUpHashoiwehfgiowehf3290f",
			expectError: true,
		},
		{
			name:        "test another correct password",
			password:    password2,
			hash:        hashed2,
			expectError: false,
		},
	}

	for _, testCase := range cases {
		err := CheckPasswordHash(testCase.hash, testCase.password)
		if (err != nil) != testCase.expectError {
			t.Errorf(`---------------------------
			test: %s FAILED 
			expected error: %v 
			got error: %v	
			error: %s\n
			`, testCase.name, testCase.expectError, err != nil, err)
		}
	}
}
