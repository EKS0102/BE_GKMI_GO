package token

import (
	"testing"
)

func TestGenerateRefreshToken(t *testing.T) {
	token, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}
	if token == "" {
		t.Fatal("expected refresh token, got empty string")
	}
	if len(token) < 40 {
		t.Fatalf("refresh token length = %d, expected at least 40", len(token))
	}
}

func TestGenerateRefreshTokenProducesDifferentTokens(t *testing.T) {
	first, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("first GenerateRefreshToken() error = %v", err)
	}
	second, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("second GenerateRefreshToken() error = %v", err)
	}
	if first == second {
		t.Fatal("expected generated refresh tokens to be different")
	}
}

func TestHashRefreshToken(t *testing.T) {
	token := "test-refresh-token"
	hash := HashRefreshToken(token)
	if hash == "" {
		t.Fatal("expected hash, got empty string")
	}
	if len(hash) != 64 {
		t.Fatalf("hash length = %d, want 64", len(hash))
	}
}

func TestHashRefreshTokenIsDeterministic(t *testing.T) {
	token := "test-refresh-token"
	first := HashRefreshToken(token)
	second := HashRefreshToken(token)
	if first != second {
		t.Fatal("expected same token to produce same hash")
	}
}

func TestHashRefreshTokenProducesDifferentHash(t *testing.T) {
	first := HashRefreshToken("refresh-token-one")
	second := HashRefreshToken("refresh-token-two")
	if first == second {
		t.Fatal("expected different tokens to produce different hashes")
	}
}
