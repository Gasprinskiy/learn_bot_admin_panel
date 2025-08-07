package global

type CommotListSearchResponse[T any] struct {
	Data       []T `json:"data"`
	TotalCount int `json:"total"`
	LeftCount  int `json:"left"`
}

func NewCommotListSearchResponse[T any](
	data []T,
	totalCount, limit, pageCount int,
) CommotListSearchResponse[T] {
	loadedCount := limit * pageCount
	left := totalCount - loadedCount

	if left < 0 {
		left = 0
	}

	return CommotListSearchResponse[T]{
		Data:       data,
		TotalCount: totalCount,
		LeftCount:  left,
	}
}
