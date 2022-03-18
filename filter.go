package bitbucket

type ProjectsFilter struct {
	Name       string
	Permission string
}

type ProjectReposFilter struct {
	ProjectKey string `structs:"projectKey"`
}

type ProjectReposFileFilter struct {
	At string `structs:"at"`
}

type ProjectReposBranchesFilter struct {
	Base    string `structs:"base"`
	Details bool   `structs:"details"`
	Text    string `structs:"filterText`
	OrderBy string `structs:"orderBy"`
}
