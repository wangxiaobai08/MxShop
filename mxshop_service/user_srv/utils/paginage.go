package utils

// Paginate
// @Description: 分页
// @param page
// @param pageSize
// @return func(db *gorm.DB) *gorm.DB
func Paginate(page int, pageSize int) int {

	if page == 0 {
		page = 1
	}
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return offset
}
