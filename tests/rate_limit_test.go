package tests

import (
	"context"
	"gin/app/facade"
	"gin/app/middleware"
	"gin/common/ctxkey"
	"gin/pkg/container"
	"gin/pkg/errcode"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/ratelimit"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"golang.org/x/time/rate"
)

// useRateLimiter 替换测试使用的限流管理器
func useRateLimiter(t *testing.T, manager *ratelimit.Manager) {
	t.Helper()

	previous := facade.RateLimiter()
	container.Default().Set(serviceprovider.ServiceRateLimit, manager)

	t.Cleanup(func() {
		manager.Close()
		if previous == nil {
			container.Default().Delete(serviceprovider.ServiceRateLimit)
			return
		}
		container.Default().Set(serviceprovider.ServiceRateLimit, previous)
	})
}

// decodeRateLimitResponse 解析测试响应
func decodeRateLimitResponse(t *testing.T, recorder *httptest.ResponseRecorder) errcode.Response {
	t.Helper()

	var response errcode.Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	return response
}

// performRequest 执行测试请求
func performRequest(t *testing.T, engine *gin.Engine, method, path string) errcode.Response {
	t.Helper()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(method, path, nil)
	engine.ServeHTTP(recorder, request)
	return decodeRateLimitResponse(t, recorder)
}

func TestRateLimiterFacade(t *testing.T) {
	if facade.RateLimiter() == nil {
		t.Fatal("限流门面未返回管理器")
	}
}

func TestRateLimitManager(t *testing.T) {
	var nilManager *ratelimit.Manager
	if !nilManager.AllowGlobal() {
		t.Fatal("空管理器应允许请求")
	}
	if !nilManager.AllowIP("127.0.0.1", "/ping", 1, 1) {
		t.Fatal("空管理器应允许IP请求")
	}
	if err := nilManager.WaitUser(context.Background(), "1", "/user", 1, 1); err != nil {
		t.Fatalf("空管理器应允许用户请求: %v", err)
	}
	nilManager.Close()

	manager := ratelimit.NewManager(time.Minute, 1, 1)
	if !manager.AllowGlobal() {
		t.Fatal("首次全局请求应通过")
	}
	if manager.AllowGlobal() {
		t.Fatal("第二次全局请求应被限流")
	}

	if !manager.AllowIP("127.0.0.1", "/ping", 1, 1) {
		t.Fatal("首次IP请求应通过")
	}
	if manager.AllowIP("127.0.0.1", "/ping", 1, 1) {
		t.Fatal("第二次IP请求应被限流")
	}

	if err := manager.WaitUser(context.Background(), "1", "/user", 1, 1); err != nil {
		t.Fatalf("首次用户请求应通过: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if err := manager.WaitUser(ctx, "1", "/user", 1, 1); err == nil {
		t.Fatal("第二次用户请求应因等待超时被限流")
	}

	manager.Close()
	manager.Close()
}

func TestRateLimitGlobalMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	useRateLimiter(t, ratelimit.NewManager(time.Minute, 1, 1))

	engine := gin.New()
	var rateLimit middleware.RateLimit
	engine.Use(rateLimit.Handle())
	engine.GET("/global", func(c *gin.Context) {
		facade.Response().Success(c, errcode.Success())
	})

	if response := performRequest(t, engine, http.MethodGet, "/global"); response.Code != 0 {
		t.Fatalf("首次请求应通过,实际错误码: %d", response.Code)
	}
	if response := performRequest(t, engine, http.MethodGet, "/global"); response.Code != 429 {
		t.Fatalf("第二次请求应被限流,实际错误码: %d", response.Code)
	}
}

func TestRateLimitIPMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	useRateLimiter(t, ratelimit.NewManager(time.Minute, rate.Inf, 1))

	engine := gin.New()
	var rateLimit middleware.RateLimit
	engine.Use(rateLimit.IpRateLimit(1, 1))
	engine.GET("/ip", func(c *gin.Context) {
		facade.Response().Success(c, errcode.Success())
	})

	if response := performRequest(t, engine, http.MethodGet, "/ip"); response.Code != 0 {
		t.Fatalf("首次请求应通过,实际错误码: %d", response.Code)
	}
	if response := performRequest(t, engine, http.MethodGet, "/ip"); response.Code != 429 {
		t.Fatalf("第二次请求应被限流,实际错误码: %d", response.Code)
	}
}

func TestRateLimitUserMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	useRateLimiter(t, ratelimit.NewManager(time.Minute, rate.Inf, 1))

	engine := gin.New()
	engine.Use(func(c *gin.Context) {
		c.Set(ctxkey.UserIdKey, int64(1))
		c.Next()
	})
	var rateLimit middleware.RateLimit
	engine.Use(rateLimit.UserRateLimit(1, 1))
	engine.GET("/user", func(c *gin.Context) {
		facade.Response().Success(c, errcode.Success())
	})

	if response := performRequest(t, engine, http.MethodGet, "/user"); response.Code != 0 {
		t.Fatalf("首次请求应通过,实际错误码: %d", response.Code)
	}
	if response := performRequest(t, engine, http.MethodGet, "/user"); response.Code != 429 {
		t.Fatalf("第二次请求应被限流,实际错误码: %d", response.Code)
	}
}
