package request

type Search struct {
	Search map[string]interface{} `form:"__search" json:"__search"`
	Sort   map[string]interface{} `form:"__sort" json:"__sort"`
}

type PageListValidate struct {
	Page     int  `form:"page" validate:"required|int|gt:0" label:"页码"`
	PageSize int  `form:"pageSize" validate:"required|int|gt:0" label:"每页数量"`
	NotPage  bool `form:"notPage" validate:"bool" label:"不分页"`
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
	List interface{} `json:"list"`
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
