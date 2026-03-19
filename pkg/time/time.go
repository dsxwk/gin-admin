package time

import (
	"strings"
	"time"
)

var defaultLocation = time.Local

// SetLocation 设置默认时区
func SetLocation(name string) error {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return err
	}
	defaultLocation = loc
	return nil
}

// Location 获取当前时区
func Location() *time.Location {
	return defaultLocation
}

type Time struct {
	t time.Time
}

// Now 当前时间(带时区)
func Now() Time {
	return Time{t: time.Now().In(defaultLocation)}
}

// FromTime 从time.Time创建
func FromTime(t time.Time) Time {
	return Time{t: t.In(defaultLocation)}
}

// Parse 解析时间
func Parse(value string, format string) (Time, error) {
	layout := convertFormat(format)
	t, err := time.ParseInLocation(layout, value, defaultLocation)
	if err != nil {
		return Time{}, err
	}
	return Time{t: t}, nil
}

// Format 格式化
func (c Time) Format(format string) string {
	return c.t.In(defaultLocation).Format(convertFormat(format))
}

// AddDays 增加天数
func (c Time) AddDays(days int) Time {
	return Time{t: c.t.AddDate(0, 0, days)}
}

// AddMonths 增加月数
func (c Time) AddMonths(months int) Time {
	return Time{t: c.t.AddDate(0, months, 0)}
}

// AddYears 增加年数
func (c Time) AddYears(years int) Time {
	return Time{t: c.t.AddDate(years, 0, 0)}
}

// AddHours 增加小时数
func (c Time) AddHours(hours int) Time {
	return Time{t: c.t.Add(time.Duration(hours) * time.Hour)}
}

// AddMinutes 增加分钟数
func (c Time) AddMinutes(min int) Time {
	return Time{t: c.t.Add(time.Duration(min) * time.Minute)}
}

// AddSeconds 增加秒数
func (c Time) AddSeconds(sec int) Time {
	return Time{t: c.t.Add(time.Duration(sec) * time.Second)}
}

// DiffSeconds 秒数差值
func (c Time) DiffSeconds(other Time) int64 {
	return int64(c.t.Sub(other.t).Seconds())
}

// DiffMinutes 分钟数差值
func (c Time) DiffMinutes(other Time) int64 {
	return int64(c.t.Sub(other.t).Minutes())
}

// DiffHours 小时数差值
func (c Time) DiffHours(other Time) int64 {
	return int64(c.t.Sub(other.t).Hours())
}

// DiffDays 天数差值
func (c Time) DiffDays(other Time) int64 {
	return int64(c.t.Sub(other.t).Hours() / 24)
}

// Timestamp 时间戳
func (c Time) Timestamp() int64 {
	return c.t.Unix()
}

// TimestampMilli 毫秒时间戳
func (c Time) TimestampMilli() int64 {
	return c.t.UnixMilli()
}

// StartOfDay 当日开始时间
func (c Time) StartOfDay() Time {
	y, m, d := c.t.In(defaultLocation).Date()
	return Time{t: time.Date(y, m, d, 0, 0, 0, 0, defaultLocation)}
}

// EndOfDay 当日结束时间
func (c Time) EndOfDay() Time {
	y, m, d := c.t.In(defaultLocation).Date()
	return Time{t: time.Date(y, m, d, 23, 59, 59, 999999999, defaultLocation)}
}

// StartOfMonth 当月开始时间
func (c Time) StartOfMonth() Time {
	y, m, _ := c.t.In(defaultLocation).Date()
	return Time{t: time.Date(y, m, 1, 0, 0, 0, 0, defaultLocation)}
}

// EndOfMonth 当月结束时间
func (c Time) EndOfMonth() Time {
	y, m, _ := c.t.In(defaultLocation).Date()
	firstNext := time.Date(y, m+1, 1, 0, 0, 0, 0, defaultLocation)
	return Time{t: firstNext.Add(-time.Nanosecond)}
}

// 格式转换
func convertFormat(format string) string {
	replacer := strings.NewReplacer(
		"Y", "2006",
		"y", "06",
		"m", "01",
		"n", "1",
		"d", "02",
		"j", "2",
		"H", "15",
		"G", "15",
		"h", "03",
		"g", "3",
		"i", "04",
		"s", "05",
	)
	return replacer.Replace(format)
}
