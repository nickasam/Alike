package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/pkg/jwt"
	"github.com/Alike/backend/pkg/response"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestHealthEndpoint 在依赖全部为 nil/未配置时，应返回统一格式的 200、
// 各依赖标记 unavailable、整体 status=ok（unavailable 不判定为故障）。
func TestHealthEndpoint(t *testing.T) {
	engine, _, _ := New(&Deps{
		Cfg: nil, // 无配置 → MinIO store 为 nil → unavailable
		JWT: jwt.NewManager("s", 0, 0),
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	var body response.Body
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body.Code != response.CodeSuccess {
		t.Errorf("code = %d, want 0", body.Code)
	}

	data, ok := body.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("data type = %T, want object", body.Data)
	}
	if data["status"] != "ok" {
		t.Errorf("status field = %v, want ok", data["status"])
	}
	if data["database"] != "unavailable" {
		t.Errorf("database = %v, want unavailable (nil db)", data["database"])
	}
	if data["redis"] != "unavailable" {
		t.Errorf("redis = %v, want unavailable (nil redis)", data["redis"])
	}
	if data["minio"] != "unavailable" {
		t.Errorf("minio = %v, want unavailable (no cfg)", data["minio"])
	}
}
