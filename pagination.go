package bitbucket

type Pagination struct {
	Size          int  `json:"size"`
	Limit         int  `json:"limit"`
	IsLastPage    bool `json:"isLastPage"`
	Start         int  `json:"start"`
	NextPageStart int  `json:"nextPageStart"`
}

func DefaultPagination() *Pagination {
	return &Pagination{Limit: 25, Start: 0}
}
