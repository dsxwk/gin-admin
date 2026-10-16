package enum

import (
	"fmt"
	"gin/app/enum"
	enumx "gin/pkg/enum"
	"testing"

	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, []enumx.Item[int]{
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
		assert.Equal(t, []enumx.Item[int64]{
			{Value: 1, Desc: "启用"},
			{Value: 2, Desc: "停用"},
		}, list)

		// 测试Desc方法
		assert.Equal(t, "启用", status.Desc(1))
		assert.Equal(t, "停用", status.Desc(2))
		assert.Equal(t, "", status.Desc(0))

		// 测试Value方法
		assert.Equal(t, int64(1), status.Value("启用"))
		assert.Equal(t, int64(2), status.Value("停用"))
		assert.Equal(t, int64(0), status.Value("未知"))

		// 测试Map方法
		m := status.Map()
		assert.Equal(t, map[int64]string{
			1: "启用",
			2: "停用",
		}, m)

		// 测试ContainsValue
		assert.True(t, status.ContainsValue(1))
		assert.True(t, status.ContainsValue(2))
		assert.False(t, status.ContainsValue(0))

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
		assert.Equal(t, int64(1), enum.UserStatusEnabled)
		assert.Equal(t, int64(2), enum.UserStatusDisabled)
	})
}

// TestBaseEnum 测试基础枚举功能
func TestBaseEnum(t *testing.T) {
	// 创建测试枚举
	testEnum := enumx.NewEnum(
		enumx.Item[int]{Value: 1, Desc: "一"},
		enumx.Item[int]{Value: 2, Desc: "二"},
		enumx.Item[int]{Value: 3, Desc: "三"},
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
			_ = status.Desc(1)
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
	desc := userEnum.Status().Desc(1)
	fmt.Printf("1的描述: %s\n", desc) // 去掉空格

	// 获取状态值
	value := userEnum.Status().Value("启用")
	fmt.Printf("启用的值: %d\n", value) // 去掉空格

	// Output:
	// 状态列表: [{Value:1 Desc:启用} {Value:2 Desc:停用}]
	// 1的描述: 启用
	// 启用的值: 1
}

// TestBaseEnumGetReturnsClone 测试获取枚举列表返回副本
func TestBaseEnumGetReturnsClone(t *testing.T) {
	enum := enumx.NewEnum(
		enumx.Item[int]{Value: 1, Desc: "一"},
		enumx.Item[int]{Value: 2, Desc: "二"},
	)

	items := enum.Get()
	items[0].Desc = "已修改"

	if enum.Get()[0].Desc != "一" {
		t.Fatal("枚举列表不应被调用方修改")
	}
}

// TestBaseEnumLookup 测试枚举值存在性
func TestBaseEnumLookup(t *testing.T) {
	enum := enumx.NewEnum(
		enumx.Item[int]{Value: 0, Desc: "未知"},
		enumx.Item[int]{Value: 1, Desc: "有效"},
	)

	value, ok := enum.ValueOK("未知")
	if !ok || value != 0 {
		t.Fatalf("描述查询失败,value=%d,ok=%v", value, ok)
	}
	if _, ok = enum.ValueOK("不存在"); ok {
		t.Fatal("不存在的描述不应返回有效标记")
	}

	desc, ok := enum.DescOK(0)
	if !ok || desc != "未知" {
		t.Fatalf("值查询失败,desc=%s,ok=%v", desc, ok)
	}
	if _, ok = enum.DescOK(99); ok {
		t.Fatal("不存在的值不应返回有效标记")
	}
}
