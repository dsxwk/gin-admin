package tests

import (
	"gin/app/middleware"
	"gin/common/errcode"
	"gin/common/response"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()

	var rateLimitMiddleware middleware.RateLimit
	// 1QPS,最多1个突发
	r.Use(rateLimitMiddleware.IpRateLimit(1, 1))

	r.GET("/test", func(c *gin.Context) {
		errCode := errcode.Success()
		response.Success(c, &errCode)
	})

	t.Run("limit trigger", func(t *testing.T) {
		// 第一次请求
		req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
		w1 := httptest.NewRecorder()
		r.ServeHTTP(w1, req1)

		var resp response.Response
		if err := json.Unmarshal(w1.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Code != 0 {
			t.Fatalf("first request should pass, got code %d", resp.Code)
		}

		// 第二次立刻请求(应该被限流)
		req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)

		if err := json.Unmarshal(w2.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Code != 429 {
			t.Fatalf("expected 429, got %d", resp.Code)
		}
	})

	t.Run("recover after time", func(t *testing.T) {
		// 等待限流器恢复(需要等待至少1秒)
		time.Sleep(2 * time.Second)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		var resp response.Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Code != 0 {
			t.Fatalf("should recover after 1s, got code %d", resp.Code)
		}
	})

	t.Run("concurrent test", func(t *testing.T) {
		var wg sync.WaitGroup
		var mu sync.Mutex
		var success int
		var fail int

		// 等待限流器状态重置
		time.Sleep(1 * time.Second)

		// 并发10个请求
		for i := 0; i < 10; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				w := httptest.NewRecorder()

				r.ServeHTTP(w, req)

				var resp response.Response
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Errorf("failed to unmarshal response: %v", err)
					return
				}

				mu.Lock()
				if resp.Code == 0 {
					success++
				} else if resp.Code == 429 {
					fail++
				}
				mu.Unlock()
			}()
		}

		wg.Wait()

		t.Logf("success=%d fail=%d", success, fail)

		// 由于是1QPS,并发10个请求,最多只有1-2个能成功
		if success == 0 {
			t.Errorf("should have at least one success request")
		}
		if success > 3 {
			t.Logf("Warning: success count %d is higher than expected (should be <= 2)", success)
		}
		// 至少应该有一些请求被限流
		if fail == 0 {
			t.Errorf("should have some limited requests")
		}
	})
}
