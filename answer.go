package bitbucket

type Answer struct {
	Size          int           `json:"size"`
	Limit         int           `json:"limit"`
	IsLastPage    bool          `json:"isLastPage"`
	Start         int           `json:"start"`
	NextPageStart int           `json:"nextPageStart"`
	Values        []interface{} `json:"values"`
}

func fillPagination(answer *Answer, pagination *Pagination) {
	pagination.Size = answer.Size
	pagination.Limit = answer.Limit
	pagination.IsLastPage = answer.IsLastPage
	pagination.Start = answer.Start
	pagination.NextPageStart = answer.NextPageStart
}
