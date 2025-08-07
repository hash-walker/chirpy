package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestHashPasswordSuccess checks if the password is hashed correctly and is not equal to the plain text.
func TestHashPasswordSuccess(t *testing.T) {
	password := "01234"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Logf("Hashed password: %q", hashed)

	if hashed == password {
		t.Errorf("hashed password should not equal the original password")
	}

	// Validate the hash actually matches the password
	if err := CheckPasswordHash(hashed, password); err != nil {
		t.Errorf("expected password to match hash, but got error: %v", err)
	}
}

func TestJwtMake(t *testing.T){

	
	id := "10c18434-f020-47f3-ad07-61b7b225957f"
	userID, err := uuid.Parse(id)

	if err != nil {
		fmt.Printf("Cannot parse user id: %v", err)
	}
	

	tokenSecret := "secretToken"
	var expiresIn time.Duration
	expiresIn = 10000000000

	signedToken, err := MakeJWT(userID, tokenSecret, expiresIn)

	if err != nil {
		fmt.Printf("Cannot create JWT token: %v", err)
	}

	fmt.Printf("The signed token is: %v", signedToken)

	userID, err = ValidateJWT(signedToken, tokenSecret)

	if err != nil {
		fmt.Printf("It can't be validated: %v", err)
	}

	fmt.Printf("The user id is: %v", userID)
}
