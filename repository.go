package bitbucket

type Repository struct {
	Slug          string       `json:"slug"`
	Id            int          `json:"id"`
	Name          string       `json:"name"`
	ScmId         string       `json:"scmId"`
	State         string       `json:"state"`
	StatusMessage string       `json:"statusMessage"`
	Forkable      bool         `json:"forkable"`
	Project       Project      `json:"project"`
	Public        bool         `json:"public"`
	Links         ProjectLinks `json:"links"`
}
