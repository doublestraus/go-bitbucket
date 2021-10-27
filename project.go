package bitbucket

type Project struct {
	Key         string       `json:"key"`
	Id          int          `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Public      bool         `json:"public"`
	Type        string       `json:"type"`
	Links       ProjectLinks `json:"links,omitempty"`
}

type Links struct {
	Href string `json:"href,omitempty"`
	Name string `json:"name,omitempty"`
}

type ProjectLinks struct {
	Self  []Links `json:"self,omitempty"`
	Clone []Links `json:"clone,omitempty"`
}
