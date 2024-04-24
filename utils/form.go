package utils

const (
	defaultPage     = 1
	defaultPageSize = 10
)

func Pagination(page, pageSize int) (int, int) {
	if page == 0 {
		page = defaultPage
	}
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	offset := pageSize * (page - 1)
	return offset, pageSize
}
