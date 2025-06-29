package auth

import "testing"

func TestSamePassword(t *testing.T) {
	const password string = "mypassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("Error, trying to hash password: %v", err)
	}

	if err := CheckPasswordHash(password, hash); err != nil {
		t.Errorf("Passwords do not match: %v", err)
	}
}

func TestDifferentPassword(t *testing.T) {
	const password1 string = "mypassword"
	const password2 string = "mypassword2"
	hash1, err := HashPassword(password1)
	if err != nil {
		t.Errorf("Error, trying to hash password: %v", err)
	}

	if err := CheckPasswordHash(password2, hash1); err == nil {
		t.Errorf("H(%s) == H(%s)", password1, password2)
	}
}

func TestEmptyPassword(t *testing.T) {
	_, err := HashPassword("")
	if err != nil {
		t.Errorf("Error, hashing empty password: %v", err)
	}
}
