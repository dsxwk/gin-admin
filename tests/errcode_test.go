package tests

import (
	"gin/common/errcode"
	"net/http"
	"testing"
)

// 错误码添加前缀测试
func TestErrorCode_WithPrefix(t *testing.T) {
	tests := []struct {
		name     string
		prefix   int64
		code     int64
		expected int64
	}{
		{"normal", 1000, 1, 10001},
		{"zero code", 1000, 0, 10000},
		{"large code", 1001, 999, 1001999},
		{"multiple digits", 10, 123, 10123},
		{"large prefix", 9999, 999, 9999999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := errcode.NewError(tt.code, "test")
			e2 := e.WithPrefix(tt.prefix)

			if e2.Code != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, e2.Code)
			}

			// 原对象不应被修改
			if e.Code != tt.code {
				t.Errorf("original object modified, expected %d, got %d", tt.code, e.Code)
			}
		})
	}
}

// 链式调用测试
func TestErrorCode_Chaining(t *testing.T) {
	t.Run("new error chain", func(t *testing.T) {
		e := errcode.NewError(1, "msg").
			WithPrefix(1000).
			WithMsg("new msg").
			WithData("data")

		if e.Code != 10001 {
			t.Errorf("unexpected code: %d, expected 10001", e.Code)
		}
		if e.Msg != "new msg" {
			t.Errorf("unexpected msg: %s", e.Msg)
		}
		if e.Data != "data" {
			t.Errorf("unexpected data: %v", e.Data)
		}
	})

	t.Run("system error chain", func(t *testing.T) {
		e1 := errcode.SystemError().WithMsg("hello").WithPrefix(2000)

		if e1.Code != 2000500 {
			t.Errorf("unexpected code: %d, expected 2000500", e1.Code)
		}
		if e1.Msg != "hello" {
			t.Errorf("unexpected msg: %s", e1.Msg)
		}
		if e1.HttpCode != http.StatusInternalServerError {
			t.Errorf("unexpected httpCode: %d", e1.HttpCode)
		}
	})

	t.Run("full chain with data", func(t *testing.T) {
		e := errcode.NotFound().
			WithMsg("user not found").
			WithData(map[string]interface{}{"user_id": 123}).
			WithHttpCode(http.StatusNotFound)

		if e.Code != 404 {
			t.Errorf("expected code 404, got %d", e.Code)
		}
		if e.Msg != "user not found" {
			t.Errorf("expected msg 'user not found', got %s", e.Msg)
		}
		if e.HttpCode != http.StatusNotFound {
			t.Errorf("expected httpCode 404, got %d", e.HttpCode)
		}
	})
}

// Error()方法测试
func TestErrorCode_Error(t *testing.T) {
	tests := []struct {
		name     string
		code     int64
		msg      string
		expected string
	}{
		{"normal", 123, "hello", "123=>hello"},
		{"zero code", 0, "success", "0=>success"},
		{"large code", 999999, "error", "999999=>error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := errcode.NewError(tt.code, tt.msg)
			if e.Error() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, e.Error())
			}
		})
	}
}

// NewError 测试
func TestNewError(t *testing.T) {
	t.Run("with httpCode", func(t *testing.T) {
		e := errcode.NewError(400, "Bad Request", 400)
		if e.Code != 400 {
			t.Errorf("expected code 400, got %d", e.Code)
		}
		if e.Msg != "Bad Request" {
			t.Errorf("expected msg 'Bad Request', got %s", e.Msg)
		}
		if e.HttpCode != 400 {
			t.Errorf("expected httpCode 400, got %d", e.HttpCode)
		}
	})

	t.Run("without httpCode", func(t *testing.T) {
		e := errcode.NewError(400, "Bad Request")
		if e.Code != 400 {
			t.Errorf("expected code 400, got %d", e.Code)
		}
		if e.Msg != "Bad Request" {
			t.Errorf("expected msg 'Bad Request', got %s", e.Msg)
		}
		if e.HttpCode != http.StatusOK {
			t.Errorf("expected httpCode %d, got %d", http.StatusOK, e.HttpCode)
		}
	})

	t.Run("with zero httpCode", func(t *testing.T) {
		e := errcode.NewError(400, "Bad Request", 0)
		if e.HttpCode != http.StatusOK {
			t.Errorf("expected httpCode %d, got %d", http.StatusOK, e.HttpCode)
		}
	})
}

// 内置错误码测试
func TestBuiltInErrors(t *testing.T) {
	tests := []struct {
		name         string
		err          errcode.ErrorCode
		expectedCode int64
		expectedMsg  string
		expectedHttp int
	}{
		{"Success", errcode.Success(), 0, "Success", http.StatusOK},
		{"Redirect", errcode.Redirect(), 301, "Redirect", http.StatusMovedPermanently},
		{"ArgsError", errcode.ArgsError(), 400, "Invalid arguments", http.StatusBadRequest},
		{"Unauthorized", errcode.Unauthorized(), 401, "Unauthorized", http.StatusUnauthorized},
		{"NotFound", errcode.NotFound(), 404, "Resource not found", http.StatusNotFound},
		{"RateLimitError", errcode.RateLimitError(), 429, "Rate limit exceeded", http.StatusTooManyRequests},
		{"SystemError", errcode.SystemError(), 500, "Internal server error", http.StatusInternalServerError},
		{"TimeoutError", errcode.TimeoutError(), 504, "Request Timeout", http.StatusGatewayTimeout},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.expectedCode {
				t.Errorf("Code: expected %d, got %d", tt.expectedCode, tt.err.Code)
			}
			if tt.err.Msg != tt.expectedMsg {
				t.Errorf("Msg: expected %s, got %s", tt.expectedMsg, tt.err.Msg)
			}
			if tt.err.HttpCode != tt.expectedHttp {
				t.Errorf("HttpCode: expected %d, got %d", tt.expectedHttp, tt.err.HttpCode)
			}
		})
	}
}

// WithCode 测试
func TestErrorCode_WithCode(t *testing.T) {
	e := errcode.Success().WithCode(100)
	if e.Code != 100 {
		t.Errorf("expected code 100, got %d", e.Code)
	}
	if e.Msg != "Success" {
		t.Errorf("expected msg 'Success', got %s", e.Msg)
	}
}

// WithMsg 测试
func TestErrorCode_WithMsg(t *testing.T) {
	e := errcode.SystemError().WithMsg("custom error message")
	if e.Msg != "custom error message" {
		t.Errorf("expected msg 'custom error message', got %s", e.Msg)
	}
	if e.Code != 500 {
		t.Errorf("expected code 500, got %d", e.Code)
	}
}

// WithData 测试
func TestErrorCode_WithData(t *testing.T) {
	data := map[string]interface{}{
		"field1": "value1",
		"field2": 123,
	}
	e := errcode.Success().WithData(data)
	if e.Data == nil {
		t.Error("Data should not be nil")
	}
	if e.Data.(map[string]interface{})["field1"] != "value1" {
		t.Error("Data content mismatch")
	}
}

// WithHttpCode 测试
func TestErrorCode_WithHttpCode(t *testing.T) {
	e := errcode.Success().WithHttpCode(http.StatusCreated)
	if e.HttpCode != http.StatusCreated {
		t.Errorf("expected httpCode %d, got %d", http.StatusCreated, e.HttpCode)
	}
	if e.Code != 0 {
		t.Errorf("expected code 0, got %d", e.Code)
	}
}

// 完整流程测试
func TestErrorCode_Integration(t *testing.T) {
	// 模拟 API 错误处理流程
	err := errcode.NotFound().
		WithMsg("user_not_found").
		WithData(map[string]interface{}{"user_id": 12345}).
		WithHttpCode(http.StatusNotFound)

	if err.Code != 404 {
		t.Errorf("unexpected code: %d", err.Code)
	}
	if err.HttpCode != http.StatusNotFound {
		t.Errorf("unexpected httpCode: %d", err.HttpCode)
	}
	if err.Msg != "user_not_found" {
		t.Errorf("unexpected msg: %s", err.Msg)
	}

	// 验证 Error 接口实现
	errMsg := err.Error()
	if errMsg == "" {
		t.Error("Error() should not return empty string")
	}
}

// 并发测试
func TestErrorCode_Concurrency(t *testing.T) {
	t.Run("concurrent WithPrefix", func(t *testing.T) {
		e := errcode.NewError(1, "test")
		done := make(chan bool)

		for i := 0; i < 100; i++ {
			go func() {
				_ = e.WithPrefix(1000)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}

		// 原对象不应被修改
		if e.Code != 1 {
			t.Errorf("original object modified, expected 1, got %d", e.Code)
		}
	})
}
