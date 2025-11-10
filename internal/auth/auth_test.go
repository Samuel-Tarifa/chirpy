package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT_Success(t *testing.T) {
	secret := "test-secret"
	id := uuid.New()

	token, err := MakeJWT(id, secret, time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT error: %v", err)
	}

	gotID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT error: %v", err)
	}
	if gotID != id {
		t.Fatalf("expected %s, got %s", id, gotID)
	}
}

func TestValidateJWT_Expired(t *testing.T) {
	secret := "test-secret"
	id := uuid.New()

	token, err := MakeJWT(id, secret, -time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT error: %v", err)
	}

	if _, err := ValidateJWT(token, secret); err == nil {
		t.Fatalf("expected error for expired token, got nil")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	id := uuid.New()

	token, err := MakeJWT(id, "correct", time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT error: %v", err)
	}

	if _, err := ValidateJWT(token, "wrong"); err == nil {
		t.Fatalf("expected error for wrong secret, got nil")
	}
}