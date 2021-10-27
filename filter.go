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
