package service

import (
	"fmt"
	"gin/app/model"
	"gin/common/base"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

// MethodStat 请求方法统计
type MethodStat struct {
	Method string `json:"method"`
	Count  int64  `json:"count"`
}

// CostDist 耗时分布
type CostDist struct {
	Lt50     int64 `json:"lt50"`
	Bt50100  int64 `json:"bt50100"`
	Bt100200 int64 `json:"bt100200"`
	Bt200500 int64 `json:"bt200500"`
	Gt500    int64 `json:"gt500"`
}

// HourlyStat 小时请求量
type HourlyStat struct {
	Hour  int   `json:"hour"`
	Count int64 `json:"count"`
}

// StatusCodeStat 状态码统计
type StatusCodeStat struct {
	Code  int   `json:"code"`
	Count int64 `json:"count"`
}

// DashboardSummary 仪表盘概要
type DashboardSummary struct {
	TodayPV          int64   `json:"todayPv"`
	YesterdayPV      int64   `json:"yesterdayPv"`
	PvChange         float64 `json:"pvChange"`
	TodayUV          int64   `json:"todayUv"`
	YesterdayUV      int64   `json:"yesterdayUv"`
	UvChange         float64 `json:"uvChange"`
	AvgCost          float64 `json:"avgCost"`
	YesterdayAvgCost float64 `json:"yesterdayAvgCost"`
	CostChange       float64 `json:"costChange"`
	PeakHour         int     `json:"peakHour"`
	PeakQps          float64 `json:"peakQps"`
}

// OperatorLogStatistics 操作日志统计
type OperatorLogStatistics struct {
	MethodStats     []MethodStat     `json:"methodStats"`
	CostDist        CostDist         `json:"costDist"`
	HourlyStats     []HourlyStat     `json:"hourlyStats"`
	StatusCodeStats []StatusCodeStat `json:"statusCodeStats"`
	Summary         DashboardSummary `json:"summary"`
}

// DashboardCard 仪表盘卡片
type DashboardCard struct {
	Title          string  `json:"title"`
	Count          int64   `json:"count"`
	YesterdayCount int64   `json:"yesterdayCount"`
	Change         int64   `json:"change"`
	ChangePercent  float64 `json:"changePercent"`
	Trend          string  `json:"trend"`
}

// DashboardService 仪表盘服务
type DashboardService struct {
	base.BaseService
}

// Cards 仪表盘卡片统计
func (s *DashboardService) Cards() (cards []DashboardCard, err error) {
	cards = make([]DashboardCard, 0, 8)

	cards = append(cards, s.countCard("用户数", &model.User{}, "total"))
	cards = append(cards, s.countCard("角色数", &model.Roles{}, "total"))
	cards = append(cards, s.countCard("菜单数", &model.Menu{}, "total"))
	cards = append(cards, s.countCard("文章数", &model.Article{}, "total"))
	cards = append(cards, s.countCard("今日日志", &model.OperatorLog{}, "today"))
	cards = append(cards, s.countCard("字典项", &model.Dict{}, "total"))
	cards = append(cards, s.countCard("系统配置", &model.SystemConfig{}, "total"))
	cards = append(cards, s.countCard("导入批次", &model.ImportRecords{}, "total"))

	return cards, nil
}

// countCard 统计单个卡片 total=历史总量对比 today=今日创建数对比
func (s *DashboardService) countCard(title string, m base.Model, mode string) DashboardCard {
	var today, yesterday int64
	db := s.DB(m)

	if mode == "today" {
		db.Model(m).Where("created_at >= CURDATE()").Count(&today)
		db.Model(m).Where("created_at >= CURDATE() - INTERVAL 1 DAY AND created_at < CURDATE()").Count(&yesterday)
	} else {
		db.Model(m).Count(&today)
		db.Model(m).Where("created_at < CURDATE()").Count(&yesterday)
	}

	change := today - yesterday
	var (
		changePercent float64
		trend         string
	)

	if yesterday > 0 {
		changePercent = float64(change) / float64(yesterday) * 100
	}

	switch {
	case change > 0:
		trend = "up"
	case change < 0:
		trend = "down"
	default:
		trend = "flat"
	}

	return DashboardCard{
		Title:          title,
		Count:          today,
		YesterdayCount: yesterday,
		Change:         change,
		ChangePercent:  changePercent,
		Trend:          trend,
	}
}

// Statistics 操作日志仪表盘统计
func (s *DashboardService) Statistics() (stat OperatorLogStatistics, err error) {
	var (
		m  model.OperatorLog
		db = s.DB(&m)
	)

	// 请求方法统计
	db.Model(&m).
		Select("method, COUNT(*) as count").
		Group("method").
		Find(&stat.MethodStats)

	// 耗时分布
	db.Model(&m).
		Select(`
			SUM(CASE WHEN cost_ms < 50 THEN 1 ELSE 0 END) as lt50,
			SUM(CASE WHEN cost_ms >= 50 AND cost_ms < 100 THEN 1 ELSE 0 END) as bt50100,
			SUM(CASE WHEN cost_ms >= 100 AND cost_ms < 200 THEN 1 ELSE 0 END) as bt100200,
			SUM(CASE WHEN cost_ms >= 200 AND cost_ms < 500 THEN 1 ELSE 0 END) as bt200500,
			SUM(CASE WHEN cost_ms >= 500 THEN 1 ELSE 0 END) as gt500
		`).
		Scan(&stat.CostDist)

	// 每小时请求量
	db.Model(&m).
		Select("HOUR(created_at) as hour, COUNT(*) as count").
		Where("created_at >= CURDATE()").
		Group("HOUR(created_at)").
		Order("hour ASC").
		Find(&stat.HourlyStats)

	// 状态码统计
	db.Model(&m).
		Select("status_code as code, COUNT(*) as count").
		Group("status_code").
		Find(&stat.StatusCodeStats)

	// 今日PV
	db.Model(&m).
		Where("created_at >= CURDATE()").
		Count(&stat.Summary.TodayPV)

	// 昨日PV
	db.Model(&m).
		Where("created_at >= CURDATE() - INTERVAL 1 DAY AND created_at < CURDATE()").
		Count(&stat.Summary.YesterdayPV)

	// PV变化
	if stat.Summary.YesterdayPV > 0 {
		stat.Summary.PvChange = float64(stat.Summary.TodayPV-stat.Summary.YesterdayPV) / float64(stat.Summary.YesterdayPV) * 100
	}

	// 今日UV
	db.Model(&m).
		Select("COUNT(DISTINCT user_id)").
		Where("created_at >= CURDATE() AND user_id > 0").
		Scan(&stat.Summary.TodayUV)

	// 昨日UV
	db.Model(&m).
		Select("COUNT(DISTINCT user_id)").
		Where("created_at >= CURDATE() - INTERVAL 1 DAY AND created_at < CURDATE() AND user_id > 0").
		Scan(&stat.Summary.YesterdayUV)

	// UV变化
	if stat.Summary.YesterdayUV > 0 {
		stat.Summary.UvChange = float64(stat.Summary.TodayUV-stat.Summary.YesterdayUV) / float64(stat.Summary.YesterdayUV) * 100
	}

	// 今日平均响应时间
	type AvgRow struct {
		AvgCost float64 `gorm:"column:avg_cost"`
	}
	var avgRow AvgRow
	db.Model(&m).
		Select("COALESCE(AVG(cost_ms), 0) as avg_cost").
		Where("created_at >= CURDATE()").
		Scan(&avgRow)
	stat.Summary.AvgCost = avgRow.AvgCost

	// 昨日平均响应时间
	var yesterdayAvg AvgRow
	db.Model(&m).
		Select("COALESCE(AVG(cost_ms), 0) as avg_cost").
		Where("created_at >= CURDATE() - INTERVAL 1 DAY AND created_at < CURDATE()").
		Scan(&yesterdayAvg)
	stat.Summary.YesterdayAvgCost = yesterdayAvg.AvgCost

	if yesterdayAvg.AvgCost > 0 {
		stat.Summary.CostChange = stat.Summary.AvgCost - yesterdayAvg.AvgCost
	}

	// 峰值时段
	type PeakRow struct {
		Hour  int   `gorm:"column:hour"`
		Count int64 `gorm:"column:cnt"`
	}
	var peakRow PeakRow
	db.Model(&m).
		Select("HOUR(created_at) as hour, COUNT(*) as cnt").
		Where("created_at >= CURDATE()").
		Group("HOUR(created_at)").
		Order("cnt DESC").
		Limit(1).
		Scan(&peakRow)

	stat.Summary.PeakHour = peakRow.Hour
	if peakRow.Count > 0 {
		stat.Summary.PeakQps = float64(peakRow.Count) / 3600
	}

	return stat, nil
}

// SystemResource 系统资源
func (s *DashboardService) SystemResource() SystemResource {
	return SystemResourceInfo()
}

// CPUInfo CPU信息
type CPUInfo struct {
	Usage     float64 `json:"usage"`
	Total     float64 `json:"total"`
	Available float64 `json:"available"`
}

// MemoryInfo 内存信息
type MemoryInfo struct {
	Used      string `json:"used"`
	Total     string `json:"total"`
	Available string `json:"available"`
}

// DiskInfo 磁盘信息
type DiskInfo struct {
	Used      string `json:"used"`
	Total     string `json:"total"`
	Available string `json:"available"`
}

// NetInfo 带宽信息
type NetInfo struct {
	Used      string `json:"used"`
	Total     string `json:"total"`
	Available string `json:"available"`
}

// SystemResource 系统资源
type SystemResource struct {
	CPU    CPUInfo    `json:"cpu"`
	Memory MemoryInfo `json:"memory"`
	Disk   DiskInfo   `json:"disk"`
	Net    NetInfo    `json:"net"`
}

// SystemResourceInfo 获取系统资源(按需采集)
func SystemResourceInfo() SystemResource {
	var res SystemResource

	// CPU使用率
	if percents, err := cpu.Percent(0, false); err == nil && len(percents) > 0 {
		res.CPU.Usage = percents[0]
	}
	res.CPU.Total = 100
	res.CPU.Available = 100 - res.CPU.Usage

	// 内存
	if v, err := mem.VirtualMemory(); err == nil {
		res.Memory.Used = formatBytes(v.Used)
		res.Memory.Total = formatBytes(v.Total)
		res.Memory.Available = formatBytes(v.Available)
	}

	// 磁盘
	if v, err := disk.Usage("/"); err == nil {
		res.Disk.Used = formatBytes(v.Used)
		res.Disk.Total = formatBytes(v.Total)
		res.Disk.Available = formatBytes(v.Free)
	}

	// 带宽
	res.Net.Used = netSpeedOnce()
	res.Net.Total = "-"
	res.Net.Available = "-"

	return res
}

var (
	lastNetIn  uint64
	lastNetOut uint64
	lastNetAt  time.Time
)

// netSpeedOnce 瞬时网速采集
func netSpeedOnce() string {
	if counters, err := net.IOCounters(false); err == nil && len(counters) > 0 {
		now := time.Now()
		c := counters[0]
		if !lastNetAt.IsZero() {
			elapsed := now.Sub(lastNetAt).Seconds()
			if elapsed > 0 {
				speed := uint64(float64((c.BytesRecv-lastNetIn)+(c.BytesSent-lastNetOut)) / elapsed)
				lastNetIn = c.BytesRecv
				lastNetOut = c.BytesSent
				lastNetAt = now
				return formatBits(speed)
			}
		}
		lastNetIn = c.BytesRecv
		lastNetOut = c.BytesSent
		lastNetAt = now
	}
	return "0Mbps"
}

// formatBytes 格式化字节
func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(b)/float64(div), "KMGTPE"[exp])
}

// formatBits 格式化比特
func formatBits(b uint64) string {
	if b < 1000 {
		return fmt.Sprintf("%dbps", b)
	}
	return fmt.Sprintf("%.0fMbps", float64(b)*8/1000000)
}
