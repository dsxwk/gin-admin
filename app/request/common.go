package request

import "github.com/dsxwk/validatex"

type PageListValidate struct {
	Page     int            `json:"page" form:"page" label:"页码"`
	PageSize int            `json:"pageSize" form:"pageSize" label:"每页数量"`
	NotPage  bool           `json:"notPage" form:"notPage" label:"不分页"`
	Search   map[string]any `json:"__search" form:"__search"`
	Sort     map[string]any `json:"__sort" form:"__sort"`
}

// Validate 校验分页参数
func (s PageListValidate) Validate(fields validatex.Fields) validatex.Errors {
	return validatex.When(!s.NotPage, func() validatex.Errors {
		return validatex.Check(
			fields.Int("PageListValidate.Page", s.Page).Required().Gt(0).Errors(),
			fields.Int("PageListValidate.PageSize", s.PageSize).Required().Gt(0).Errors(),
		)
	})
}

// commonMessages 返回请求公共错误消息
func commonMessages() validatex.Messages {
	return validatex.Messages{
		"required": "字段 {field} 必填",
		"gt":       "字段 {field} 需大于 0",
		"min":      "{field} 长度不能少于 {min} 个字符",
		"max":      "{field} 长度不能超过 {max} 个字符",
		"email":    "{field} 邮箱格式错误",
	}
}

// requiredValue 校验任意值必填
func requiredValue(fields validatex.Fields, path string, value any) validatex.Errors {
	return fields.Bool(path, value != nil).Rule("required", "", func(valid bool) bool {
		return valid
	}).Errors()
}

// PageData 公共分页数据
type PageData struct {
	//  总条数
	Total int64 `json:"total"`
	// 不分页
	NotPage bool `json:"notPage" example:"false"`
	// 当前页
	Page int `json:"page"`
	// 每页条数
	PageSize int `json:"pageSize"`
	// 数据列表
	List any `json:"list"`
}

// Pagination 计算分页
func Pagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	return offset, pageSize
}
