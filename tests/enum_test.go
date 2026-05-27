package tests

import (
	"fmt"
	"gin/app/enum"
	"gin/common/base"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestUserEnum 测试用户枚举
func TestUserEnum(t *testing.T) {
	userEnum := enum.UserEnum{}

	// 性别枚举测试(int类型)
	t.Run("Gender Enum", func(t *testing.T) {
		gender := userEnum.Gender()

		// 测试Get方法
		list := gender.Get()
		assert.Len(t, list, 3)
		assert.Equal(t, []base.Item[int]{
			{Value: 0, Desc: "保密"},
			{Value: 1, Desc: "男"},
			{Value: 2, Desc: "女"},
		}, list)

		// 测试Desc方法
		assert.Equal(t, "保密", gender.Desc(0))
		assert.Equal(t, "男", gender.Desc(1))
		assert.Equal(t, "女", gender.Desc(2))
		assert.Equal(t, "", gender.Desc(99)) // 不存在的值

		// 测试Value方法
		assert.Equal(t, 0, gender.Value("保密"))
		assert.Equal(t, 1, gender.Value("男"))
		assert.Equal(t, 2, gender.Value("女"))
		assert.Equal(t, 0, gender.Value("未知")) // 不存在的描述返回零值

		// 测试Map方法
		m := gender.Map()
		assert.Equal(t, map[int]string{
			0: "保密",
			1: "男",
			2: "女",
		}, m)

		// 测试ContainsValue
		assert.True(t, gender.ContainsValue(0))
		assert.True(t, gender.ContainsValue(1))
		assert.True(t, gender.ContainsValue(2))
		assert.False(t, gender.ContainsValue(99))

		// 测试ContainsDesc
		assert.True(t, gender.ContainsDesc("保密"))
		assert.True(t, gender.ContainsDesc("男"))
		assert.True(t, gender.ContainsDesc("女"))
		assert.False(t, gender.ContainsDesc("未知"))

		// 测试Len
		assert.Equal(t, 3, gender.Len())
	})

	// 状态枚举测试(string类型)
	t.Run("Status Enum", func(t *testing.T) {
		status := userEnum.Status()

		// 测试Get方法
		list := status.Get()
		assert.Len(t, list, 2)
		assert.Equal(t, []base.Item[string]{
			{Value: "enable", Desc: "启用"},
			{Value: "disable", Desc: "停用"},
		}, list)

		// 测试Desc方法
		assert.Equal(t, "启用", status.Desc("enable"))
		assert.Equal(t, "停用", status.Desc("disable"))
		assert.Equal(t, "", status.Desc("unknown"))

		// 测试Value方法
		assert.Equal(t, "enable", status.Value("启用"))
		assert.Equal(t, "disable", status.Value("停用"))
		assert.Equal(t, "", status.Value("未知"))

		// 测试Map方法
		m := status.Map()
		assert.Equal(t, map[string]string{
			"enable":  "启用",
			"disable": "停用",
		}, m)

		// 测试ContainsValue
		assert.True(t, status.ContainsValue("enable"))
		assert.True(t, status.ContainsValue("disable"))
		assert.False(t, status.ContainsValue("unknown"))

		// 测试ContainsDesc
		assert.True(t, status.ContainsDesc("启用"))
		assert.True(t, status.ContainsDesc("停用"))
		assert.False(t, status.ContainsDesc("未知"))

		// 测试Len
		assert.Equal(t, 2, status.Len())
	})

	t.Run("Constants", func(t *testing.T) {
		// 测试性别常量
		assert.Equal(t, 0, enum.UserGenderSecret)
		assert.Equal(t, 1, enum.UserGenderMale)
		assert.Equal(t, 2, enum.UserGenderFemale)

		// 测试状态常量
		assert.Equal(t, "enable", enum.UserStatusEnabled)
		assert.Equal(t, "disable", enum.UserStatusDisabled)
	})
}

// TestBaseEnum 测试基础枚举功能
func TestBaseEnum(t *testing.T) {
	// 创建测试枚举
	testEnum := base.NewEnum(
		base.Item[int]{Value: 1, Desc: "一"},
		base.Item[int]{Value: 2, Desc: "二"},
		base.Item[int]{Value: 3, Desc: "三"},
	)

	t.Run("Get", func(t *testing.T) {
		list := testEnum.Get()
		assert.Len(t, list, 3)
		assert.Equal(t, 1, list[0].Value)
		assert.Equal(t, "一", list[0].Desc)
	})

	t.Run("Desc", func(t *testing.T) {
		assert.Equal(t, "一", testEnum.Desc(1))
		assert.Equal(t, "二", testEnum.Desc(2))
		assert.Equal(t, "三", testEnum.Desc(3))
		assert.Equal(t, "", testEnum.Desc(99))
	})

	t.Run("Value", func(t *testing.T) {
		assert.Equal(t, 1, testEnum.Value("一"))
		assert.Equal(t, 2, testEnum.Value("二"))
		assert.Equal(t, 3, testEnum.Value("三"))
		assert.Equal(t, 0, testEnum.Value("未知"))
	})

	t.Run("Map", func(t *testing.T) {
		m := testEnum.Map()
		assert.Equal(t, map[int]string{
			1: "一",
			2: "二",
			3: "三",
		}, m)
	})

	t.Run("ContainsValue", func(t *testing.T) {
		assert.True(t, testEnum.ContainsValue(1))
		assert.True(t, testEnum.ContainsValue(2))
		assert.True(t, testEnum.ContainsValue(3))
		assert.False(t, testEnum.ContainsValue(99))
	})

	t.Run("ContainsDesc", func(t *testing.T) {
		assert.True(t, testEnum.ContainsDesc("一"))
		assert.True(t, testEnum.ContainsDesc("二"))
		assert.True(t, testEnum.ContainsDesc("三"))
		assert.False(t, testEnum.ContainsDesc("未知"))
	})

	t.Run("Len", func(t *testing.T) {
		assert.Equal(t, 3, testEnum.Len())
	})
}

// BenchmarkEnum 基准测试
func BenchmarkEnum(b *testing.B) {
	userEnum := enum.UserEnum{}
	status := userEnum.Status()

	b.Run("Desc", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = status.Desc("enable")
		}
	})

	b.Run("Value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = status.Value("启用")
		}
	})

	b.Run("Map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = status.Map()
		}
	})
}

// ExampleUserEnum 示例
func ExampleUserEnum() {
	userEnum := enum.UserEnum{}

	// 获取状态列表
	statusList := userEnum.Status().Get()
	fmt.Printf("状态列表: %+v\n", statusList)

	// 获取状态描述
	desc := userEnum.Status().Desc("enable")
	fmt.Printf("enable的描述: %s\n", desc) // 去掉空格

	// 获取状态值
	value := userEnum.Status().Value("启用")
	fmt.Printf("启用的值: %s\n", value) // 去掉空格

	// Output:
	// 状态列表: [{Value:enable Desc:启用} {Value:disable Desc:停用}]
	// enable的描述: 启用
	// 启用的值: enable
}
