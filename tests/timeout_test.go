package tests

import (
	"gin/app/middleware"
	"gin/common/errcode"
	"gin/common/response"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTimeoutMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var timeoutMiddleware middleware.Timeout
	r := gin.New()
	r.Use(timeoutMiddleware.Handle(2 * time.Second))

	// 慢接口(3秒)
	r.GET("/slow", func(c *gin.Context) {
		time.Sleep(3 * time.Second)
		errCode := errcode.Success()
		response.Success(c, &errCode)
	})

	// 快接口(1秒)
	r.GET("/fast", func(c *gin.Context) {
		time.Sleep(1 * time.Second)
		errCode := errcode.Success()
		response.Success(c, &errCode)
	})

	t.Run("timeout case", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/slow", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)

		if resp.Code != 504 {
			t.Fatalf("expected 504, got %d", resp.Code)
		}
	})

	t.Run("normal case", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/fast", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		var resp response.Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)

		if resp.Code != 0 {
			t.Fatalf("expected 0, got %d", resp.Code)
		}
	})
}
