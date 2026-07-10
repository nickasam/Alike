package storage

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/internal/middleware"
	"github.com/Alike/backend/pkg/response"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestClassifyImage(t *testing.T) {
	s := &Store{maxImageSize: 5 << 20, maxDocSize: 10 << 20}
	cat, ext, max, err := s.classify("image/png")
	if err != nil {
		t.Fatalf("classify image/png err: %v", err)
	}
	if cat != CategoryImage || ext != ".png" || max != 5<<20 {
		t.Fatalf("classify image/png = (%s,%s,%d), want image/.png/5MB", cat, ext, max)
	}
}

func TestClassifyDocument(t *testing.T) {
	s := &Store{maxImageSize: 5 << 20, maxDocSize: 10 << 20}
	cat, ext, max, err := s.classify("application/pdf")
	if err != nil {
		t.Fatalf("classify pdf err: %v", err)
	}
	if cat != CategoryDocument || ext != ".pdf" || max != 10<<20 {
		t.Fatalf("classify pdf = (%s,%s,%d), want document/.pdf/10MB", cat, ext, max)
	}
}

func TestClassifyRejectsUnsupported(t *testing.T) {
	s := &Store{maxImageSize: 5 << 20, maxDocSize: 10 << 20}
	for _, mime := range []string{"application/x-msdownload", "video/mp4", "application/octet-stream", ""} {
		if _, _, _, err := s.classify(mime); err != ErrUnsupportedType {
			t.Errorf("classify(%q) err = %v, want ErrUnsupportedType", mime, err)
		}
	}
}

func TestBaseMIME(t *testing.T) {
	cases := map[string]string{
		"text/plain; charset=utf-8": "text/plain",
		"image/png":                 "image/png",
		"text/html;charset=utf-8":   "text/html",
	}
	for in, want := range cases {
		if got := baseMIME(in); got != want {
			t.Errorf("baseMIME(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestBuildObjectNameStructure(t *testing.T) {
	name := buildObjectName("42", "image", ".png")
	if !strings.HasPrefix(name, "image/42/") {
		t.Errorf("object name %q missing category/prefix", name)
	}
	if !strings.HasSuffix(name, ".png") {
		t.Errorf("object name %q missing extension", name)
	}
	parts := strings.Split(name, "/")
	if len(parts) != 4 {
		t.Fatalf("object name %q should have 4 path segments, got %d", name, len(parts))
	}
}

func TestBuildObjectNameEmptyPrefixFallsBackToAnon(t *testing.T) {
	name := buildObjectName("", "document", ".pdf")
	if !strings.HasPrefix(name, "document/anon/") {
		t.Errorf("empty prefix should fall back to 'anon', got %q", name)
	}
	name2 := buildObjectName("///", "image", ".jpg")
	if !strings.HasPrefix(name2, "image/anon/") {
		t.Errorf("slash-only prefix should fall back to 'anon', got %q", name2)
	}
}

func TestRandomHexUniqueAndLength(t *testing.T) {
	a := randomHex(16)
	b := randomHex(16)
	if len(a) != 32 {
		t.Errorf("randomHex(16) len = %d, want 32", len(a))
	}
	if a == b {
		t.Errorf("randomHex produced identical values %q", a)
	}
}

// --- handler-level guards ---

func invoke(h gin.HandlerFunc, prep func(*gin.Context)) (int, response.Body) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/upload", strings.NewReader(""))
	if prep != nil {
		prep(c)
	}
	h(c)
	var parsed response.Body
	_ = json.Unmarshal(w.Body.Bytes(), &parsed)
	return w.Code, parsed
}

func TestUploadRequiresAuth(t *testing.T) {
	h := NewHandler(nil)
	_, body := invoke(h.Upload, nil)
	if body.Code != response.CodeUnauthorized {
		t.Fatalf("code=%d, want %d (upload without auth)", body.Code, response.CodeUnauthorized)
	}
}

func TestUploadStorageDisabled(t *testing.T) {
	h := NewHandler(nil) // store nil => disabled
	_, body := invoke(h.Upload, func(c *gin.Context) {
		c.Set(middleware.ContextUserID, int64(7))
	})
	if body.Code != response.CodeInternalError {
		t.Fatalf("code=%d, want %d (storage disabled)", body.Code, response.CodeInternalError)
	}
}
