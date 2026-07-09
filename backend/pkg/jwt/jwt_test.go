package jwt

import (
	"testing"
	"time"
)

func newTestManager() *Manager {
	return NewManager("test-secret", time.Hour, 24*time.Hour)
}

func TestGenerateAndParseAccess(t *testing.T) {
	m := newTestManager()
	token, err := m.GenerateAccess(42)
	if err != nil {
		t.Fatalf("generate access: %v", err)
	}

	claims, err := m.Parse(token)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if claims.UserID != 42 {
		t.Errorf("user id = %d, want 42", claims.UserID)
	}
	if claims.Type != AccessToken {
		t.Errorf("type = %q, want %q", claims.Type, AccessToken)
	}
}

func TestGenerateRefreshType(t *testing.T) {
	m := newTestManager()
	token, err := m.GenerateRefresh(7)
	if err != nil {
		t.Fatalf("generate refresh: %v", err)
	}
	claims, err := m.Parse(token)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if claims.Type != RefreshToken {
		t.Errorf("type = %q, want %q", claims.Type, RefreshToken)
	}
}

func TestParseInvalidToken(t *testing.T) {
	m := newTestManager()
	if _, err := m.Parse("not-a-jwt"); err == nil {
		t.Error("expected error for malformed token, got nil")
	}
}

func TestParseWrongSecret(t *testing.T) {
	m := newTestManager()
	token, _ := m.GenerateAccess(1)

	other := NewManager("different-secret", time.Hour, 24*time.Hour)
	if _, err := other.Parse(token); err == nil {
		t.Error("expected error when parsing with wrong secret, got nil")
	}
}

func TestParseExpiredToken(t *testing.T) {
	m := NewManager("test-secret", -time.Hour, -time.Hour) // 立即过期
	token, err := m.GenerateAccess(1)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if _, err := m.Parse(token); err == nil {
		t.Error("expected error for expired token, got nil")
	}
}
