package security

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "rahasia123"

	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if hash == "" {
		t.Fatal("expected hash, got empty string")
	}

	if hash == password {
		t.Fatal("password must not be stored as plaintext")
	}

	if !strings.HasPrefix(hash, "$argon2id$") {
		t.Fatalf(
			"expected argon2id hash, got: %s",
			hash,
		)
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "rahasia123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !VerifyPassword(password, hash) {
		t.Fatal("expected password to be valid")
	}

	if VerifyPassword("password-salah", hash) {
		t.Fatal("expected password to be invalid")
	}
}
