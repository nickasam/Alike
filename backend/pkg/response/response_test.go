package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func performRequest(handler gin.HandlerFunc) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	handler(c)
	return w
}

func TestSuccess(t *testing.T) {
	w := performRequest(func(c *gin.Context) {
		Success(c, gin.H{"hello": "world"})
	})
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var body Body
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body.Code != CodeSuccess {
		t.Errorf("code = %d, want %d", body.Code, CodeSuccess)
	}
	if body.Message != "success" {
		t.Errorf("message = %q, want success", body.Message)
	}
	if body.Data == nil {
		t.Error("data should not be nil")
	}
}

func TestPage(t *testing.T) {
	w := performRequest(func(c *gin.Context) {
		Page(c, []int{1, 2, 3}, 100, 2, 20)
	})
	var raw struct {
		Code int `json:"code"`
		Data struct {
			Total    int64 `json:"total"`
			Page     int   `json:"page"`
			PageSize int   `json:"page_size"`
			List     []int `json:"list"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if raw.Data.Total != 100 || raw.Data.Page != 2 || raw.Data.PageSize != 20 {
		t.Errorf("page data = %+v", raw.Data)
	}
	if len(raw.Data.List) != 3 {
		t.Errorf("list len = %d, want 3", len(raw.Data.List))
	}
}

func TestErrorStatusMapping(t *testing.T) {
	cases := []struct {
		code       int
		wantStatus int
	}{
		{CodeBadRequest, http.StatusBadRequest},
		{CodeUnauthorized, http.StatusUnauthorized},
		{CodeForbidden, http.StatusForbidden},
		{CodeNotFound, http.StatusNotFound},
		{CodeConflict, http.StatusConflict},
		{CodeValidationError, http.StatusUnprocessableEntity},
		{CodeTooManyRequests, http.StatusTooManyRequests},
		{CodeInternalError, http.StatusInternalServerError},
	}
	for _, tc := range cases {
		w := performRequest(func(c *gin.Context) {
			Fail(c, tc.code)
		})
		if w.Code != tc.wantStatus {
			t.Errorf("code %d -> http %d, want %d", tc.code, w.Code, tc.wantStatus)
		}
		var body Body
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if body.Code != tc.code {
			t.Errorf("body code = %d, want %d", body.Code, tc.code)
		}
		if body.Message == "" {
			t.Errorf("code %d: default message should not be empty", tc.code)
		}
	}
}

func TestErrorCustomMessage(t *testing.T) {
	w := performRequest(func(c *gin.Context) {
		Error(c, CodeConflict, "邮箱已注册")
	})
	var body Body
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body.Message != "邮箱已注册" {
		t.Errorf("message = %q, want 邮箱已注册", body.Message)
	}
}
