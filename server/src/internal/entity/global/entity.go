package global

type CommotListSearchResponse[T any] struct {
	Data       []T `json:"data"`
	TotalCount int `json:"total_count"`
}

func NewCommotListSearchResponse[T any](
	data []T,
	totalCount int,
) CommotListSearchResponse[T] {
	return CommotListSearchResponse[T]{
		Data:       data,
		TotalCount: totalCount,
	}
}
