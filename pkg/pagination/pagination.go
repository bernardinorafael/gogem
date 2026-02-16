package pagination

type Pagination struct {
	TotalItems      int  `json:"totalItems"`
	CurrentPage     int  `json:"currentPage"`
	ItemsPerPage    int  `json:"itemsPerPage"`
	TotalPages      int  `json:"totalPages"`
	HasPreviousPage bool `json:"hasPreviousPage"`
	HasNextPage     bool `json:"hasNextPage"`
	IsFirstPage     bool `json:"isFirstPage"`
	IsLastPage      bool `json:"isLastPage"`
}

type Paginated[T any] struct {
	Pagination Pagination `json:"pagination"`
	Items      []T        `json:"items"`
}

func New[T any](data []T, totalItems, currentPage, itemsPerPage int) Paginated[T] {
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
	hasPreviousPage := currentPage > 1
	hasNextPage := currentPage < totalPages
	isFirstPage := currentPage == 1
	isLastPage := currentPage == totalPages

	return Paginated[T]{
		Items: data,
		Pagination: Pagination{
			TotalItems:      totalItems,
			CurrentPage:     currentPage,
			ItemsPerPage:    itemsPerPage,
			TotalPages:      totalPages,
			HasPreviousPage: hasPreviousPage,
			HasNextPage:     hasNextPage,
			IsFirstPage:     isFirstPage,
			IsLastPage:      isLastPage,
		},
	}
}
