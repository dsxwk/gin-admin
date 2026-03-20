package tests

import (
	"gin/common/errcode"
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
	e := errcode.NewError(1, "msg").
		WithPrefix(1000).
		WithMsg("new msg").
		WithData("data")

	if e.Code != 10001 {
		t.Errorf("unexpected code: %d", e.Code)
	}

	if e.Msg != "new msg" {
		t.Errorf("unexpected msg: %s", e.Msg)
	}

	if e.Data != "data" {
		t.Errorf("unexpected data: %v", e.Data)
	}

	e1 := errcode.SystemError().WithMsg("hello").WithPrefix(2000)

	if e1.Code != 2000500 {
		t.Errorf("unexpected code: %d", e1.Code)
	}

	if e1.Msg != "hello" {
		t.Errorf("unexpected msg: %s", e1.Msg)
	}
}

// Error()方法测试
func TestErrorCode_Error(t *testing.T) {
	e := errcode.NewError(123, "hello")

	expected := "123=>hello"
	if e.Error() != expected {
		t.Errorf("expected %s, got %s", expected, e.Error())
	}
}

// 内置方法测试
func TestBuiltInErrors(t *testing.T) {
	if errcode.Success().Code != 0 {
		t.Error("Success code should be 0")
	}

	if errcode.SystemError().Code != 500 {
		t.Error("SystemError code should be 500")
	}

	if errcode.ArgsError().Code != 400 {
		t.Error("ArgsError code should be 400")
	}
}
