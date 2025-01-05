package models

type QueryParams struct {
	Offset int
	Limit  int
}

func NewQueryParams() *QueryParams {
	return &QueryParams{
		Offset: 0,
		Limit:  5,
	}
}
