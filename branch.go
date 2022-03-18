package bitbucket

type Branch struct {
	Id              string `json:"id"`
	DisplayID       string `json:"displayId"`
	Type            string `json:"type"`
	LatestCommit    string `json:"latestCommit"`
	LatestChangeset string `json:"latestCahngeset`
	IsDefault       bool   `json:"isDefault"`
}
