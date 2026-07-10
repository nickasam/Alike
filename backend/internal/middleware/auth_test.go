package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/pkg/jwt"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newManager() *jwt.Manager {
	return jwt.NewManager("test-secret", time.Hour, 24*time.Hour)
}

// run executes the Auth middleware with the given Authorization header and
// returns the HTTP status plus whether the terminal handler ran.
func run(mgr *jwt.Manager, authHeader string) (int, bool) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/protected", nil)
	if authHeader != "" {
		c.Request.Header.Set("Authorization", authHeader)
	}

	reached := false
	handlers := gin.HandlersChain{
		Auth(mgr),
		func(c *gin.Context) { reached = true },
	}
	c.Set("__test", struct{}{})
	// Manually drive the chain: gin's CreateTestContext has no handlers.
	for _, h := range handlers {
		if c.IsAborted() {
			break
		}
		h(c)
	}
	return w.Code, reached
}

func TestAuthMissingHeader(t *testing.T) {
	code, reached := run(newManager(), "")
	if reached {
		t.Fatal("handler should not run without Authorization header")
	}
	if code != http.StatusUnauthorized {
		t.Fatalf("status=%d, want 401", code)
	}
}

func TestAuthMalformedHeader(t *testing.T) {
	for _, h := range []string{"Bearer", "Basic abc", "Bearertoken", "token-only"} {
		code, reached := run(newManager(), h)
		if reached {
			t.Errorf("handler should not run for header %q", h)
		}
		if code != http.StatusUnauthorized {
			t.Errorf("header %q status=%d, want 401", h, code)
		}
	}
}

func TestAuthRejectsInvalidToken(t *testing.T) {
	code, reached := run(newManager(), "Bearer not.a.jwt")
	if reached {
		t.Fatal("handler should not run for garbage token")
	}
	if code != http.StatusUnauthorized {
		t.Fatalf("status=%d, want 401", code)
	}
}

func TestAuthRejectsRefreshTokenAsAccess(t *testing.T) {
	mgr := newManager()
	refresh, err := mgr.GenerateRefresh(99)
	if err != nil {
		t.Fatalf("generate refresh: %v", err)
	}
	code, reached := run(mgr, "Bearer "+refresh)
	if reached {
		t.Fatal("refresh token must not authenticate access-protected routes")
	}
	if code != http.StatusUnauthorized {
		t.Fatalf("status=%d, want 401", code)
	}
}

func TestAuthRejectsTokenFromDifferentSecret(t *testing.T) {
	attacker := jwt.NewManager("other-secret", time.Hour, 24*time.Hour)
	forged, _ := attacker.GenerateAccess(1)
	code, reached := run(newManager(), "Bearer "+forged)
	if reached {
		t.Fatal("token signed with a different secret must be rejected")
	}
	if code != http.StatusUnauthorized {
		t.Fatalf("status=%d, want 401", code)
	}
}

func TestAuthAcceptsValidAccessTokenAndSetsUserID(t *testing.T) {
	mgr := newManager()
	access, err := mgr.GenerateAccess(1234)
	if err != nil {
		t.Fatalf("generate access: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer "+access)

	Auth(mgr)(c)
	if c.IsAborted() {
		t.Fatal("valid access token should not abort")
	}
	uid, ok := CurrentUserID(c)
	if !ok || uid != 1234 {
		t.Fatalf("CurrentUserID = (%d,%v), want (1234,true)", uid, ok)
	}
}

func TestAuthLowercaseBearerScheme(t *testing.T) {
	mgr := newManager()
	access, _ := mgr.GenerateAccess(5)
	code, reached := run(mgr, "bearer "+access)
	if !reached {
		t.Fatal("scheme match must be case-insensitive per EqualFold")
	}
	if code != http.StatusOK && code != 200 {
		// terminal handler wrote nothing, recorder defaults to 200
	}
}

func TestCurrentUserIDMissing(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	if _, ok := CurrentUserID(c); ok {
		t.Fatal("CurrentUserID should return false when unset")
	}
}

func TestCurrentUserIDWrongType(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(ContextUserID, "not-an-int64")
	if _, ok := CurrentUserID(c); ok {
		t.Fatal("CurrentUserID should return false for non-int64 value")
	}
}

// runOptional drives OptionalAuth and reports whether it aborted and the
// user id it set (0/false when none).
func runOptional(mgr *jwt.Manager, authHeader string) (aborted bool, uid int64, hasUID bool) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/diaries/1", nil)
	if authHeader != "" {
		c.Request.Header.Set("Authorization", authHeader)
	}
	OptionalAuth(mgr)(c)
	uid, hasUID = CurrentUserID(c)
	return c.IsAborted(), uid, hasUID
}

func TestOptionalAuthNoHeaderPassesThrough(t *testing.T) {
	aborted, _, hasUID := runOptional(newManager(), "")
	if aborted {
		t.Fatal("OptionalAuth must not abort when no header present")
	}
	if hasUID {
		t.Fatal("OptionalAuth must not set user_id without a token")
	}
}

func TestOptionalAuthInvalidTokenPassesThroughAnonymously(t *testing.T) {
	aborted, _, hasUID := runOptional(newManager(), "Bearer garbage")
	if aborted {
		t.Fatal("OptionalAuth must not abort on invalid token")
	}
	if hasUID {
		t.Fatal("OptionalAuth must not set user_id for invalid token")
	}
}

func TestOptionalAuthRefreshTokenNotAccepted(t *testing.T) {
	mgr := newManager()
	refresh, _ := mgr.GenerateRefresh(9)
	_, _, hasUID := runOptional(mgr, "Bearer "+refresh)
	if hasUID {
		t.Fatal("OptionalAuth must not accept refresh token as identity")
	}
}

func TestOptionalAuthValidTokenSetsUserID(t *testing.T) {
	mgr := newManager()
	access, _ := mgr.GenerateAccess(777)
	aborted, uid, hasUID := runOptional(mgr, "Bearer "+access)
	if aborted {
		t.Fatal("OptionalAuth must not abort on valid token")
	}
	if !hasUID || uid != 777 {
		t.Fatalf("OptionalAuth user_id = (%d,%v), want (777,true)", uid, hasUID)
	}
}
