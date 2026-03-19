package tests

import (
	_t "gin/pkg/time"
	"testing"
	"time"
)

func init() {
	_ = _t.SetLocation("Asia/Shanghai")
}

// 时间格式化测试
func TestFormat(t *testing.T) {
	now := _t.Now()

	result := now.Format("Y-m-d H:i:s")

	if result == "" {
		t.Fatal("Format failed")
	}

	t.Log("Format:", result)
}

// 时间解析测试
func TestParse(t *testing.T) {
	str := "2024-01-02 03:04:05"
	tt, err := _t.Parse(str, "Y-m-d H:i:s")
	if err != nil {
		t.Fatal(err)
	}

	if tt.Format("Y-m-d H:i:s") != str {
		t.Fatal("Parse mismatch")
	}
}

// 时间测试
func TestFromTime(t *testing.T) {
	raw := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	tt := _t.FromTime(raw)
	if tt.Format("Y-m-d H:i:s") == "" {
		t.Fatal("FromTime failed")
	}
}

// 添加时间测试
func TestAdd(t *testing.T) {
	base, _ := _t.Parse("2024-01-01 00:00:00", "Y-m-d H:i:s")
	tt := base.AddDays(1)
	if tt.Format("Y-m-d") != "2024-01-02" {
		t.Fatal("AddDays failed")
	}

	tt = base.AddMonths(1)
	if tt.Format("Y-m") != "2024-02" {
		t.Fatal("AddMonths failed")
	}

	tt = base.AddYears(1)
	if tt.Format("Y") != "2025" {
		t.Fatal("AddYears failed")
	}
}

// 时间差测试
func TestDiff(t *testing.T) {
	t1, _ := _t.Parse("2024-01-02 00:00:00", "Y-m-d H:i:s")
	t2, _ := _t.Parse("2024-01-01 00:00:00", "Y-m-d H:i:s")
	if t1.DiffDays(t2) != 1 {
		t.Fatal("DiffDays failed")
	}

	if t1.DiffHours(t2) != 24 {
		t.Fatal("DiffHours failed")
	}
}

// 时间戳测试
func TestTimestamp(t *testing.T) {
	now := _t.Now()
	if now.Timestamp() <= 0 {
		t.Fatal("Timestamp failed")
	}

	if now.TimestampMilli() <= 0 {
		t.Fatal("TimestampMilli failed")
	}
}

// 开始结束时间测试
func TestStartEnd(t *testing.T) {
	base, _ := _t.Parse("2024-01-02 15:04:05", "Y-m-d H:i:s")
	start := base.StartOfDay()
	end := base.EndOfDay()

	if start.Format("H:i:s") != "00:00:00" {
		t.Fatal("StartOfDay failed")
	}

	if end.Format("H:i:s") != "23:59:59" {
		t.Fatal("EndOfDay failed")
	}
}

// 月份测试
func TestMonth(t *testing.T) {
	base, _ := _t.Parse("2024-02-15 12:00:00", "Y-m-d H:i:s")
	start := base.StartOfMonth()
	end := base.EndOfMonth()

	if start.Format("Y-m-d") != "2024-02-01" {
		t.Fatal("StartOfMonth failed")
	}

	if end.Format("Y-m-d") != "2024-02-29" {
		// 闰年测试
		t.Fatal("EndOfMonth failed")
	}
}

// 时区测试
func TestTimezone(t *testing.T) {
	_ = _t.SetLocation("UTC")
	nowUTC := _t.Now().Format("Y-m-d H:i:s")
	_ = _t.SetLocation("Asia/Shanghai")

	nowCN := _t.Now().Format("Y-m-d H:i:s")
	if nowUTC == nowCN {
		t.Fatal("Timezone not working")
	}
}
