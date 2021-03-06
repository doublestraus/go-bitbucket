package bitbucket

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/fatih/structs"
)

type Projects struct {
	client *Client
}

func (c *Client) ListProjects(pagination *Pagination, filter *ProjectsFilter) ([]*Project, error) {
	body, err := c.get("projects", pagination, structs.Map(filter))
	if err != nil {
		return nil, err
	}
	var answer Answer
	err = json.Unmarshal(body, &answer)
	if err != nil {
		return nil, err
	}
	fillPagination(&answer, pagination)
	projectList := make([]*Project, 0)
	for _, v := range answer.Values {
		var project Project
		m2s(v, &project)
		projectList = append(projectList, &project)
	}
	return projectList, nil
}

func (c *Client) GetProject(projectName string) (*Project, error) {
	body, err := c.get(fmt.Sprintf("projects/%s", url.QueryEscape(projectName)), DefaultPagination(), nil)
	if err != nil {
		return nil, err
	}
	var project Project
	err = json.Unmarshal(body, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (c *Client) GetProjectRepos(projectKey string, pagination *Pagination, filter *ProjectReposFilter) ([]*Repository, error) {
	body, err := c.get(fmt.Sprintf("projects/%s/repos", url.QueryEscape(projectKey)), pagination, structs.Map(filter))
	if err != nil {
		return nil, err
	}
	var answer Answer
	err = json.Unmarshal(body, &answer)
	fillPagination(&answer, pagination)
	reposList := make([]*Repository, 0)
	for _, v := range answer.Values {
		var repo Repository
		m2s(v, &repo)
		reposList = append(reposList, &repo)
	}
	return reposList, nil
}

func (c *Client) GetProjectsReposFiles(projectKey string, repoSlug string, pagination *Pagination, filter *ProjectReposFileFilter) ([]string, error) {
	body, err := c.get(fmt.Sprintf("projects/%s/repos/%s/files", url.QueryEscape(projectKey), url.QueryEscape(repoSlug)),
		pagination, structs.Map(filter))
	if err != nil {
		return nil, err
	}
	var answer Answer
	err = json.Unmarshal(body, &answer)
	fillPagination(&answer, pagination)
	fileList := make([]string, 0)
	for _, v := range answer.Values {
		fileList = append(fileList, v.(string))
	}
	return fileList, nil
}

func (c *Client) GetProjectsReposFileRaw(projectKey string, repoSlug string, path string, filter *ProjectReposFileFilter) ([]byte, error) {
	body, err := c.get(fmt.Sprintf("projects/%s/repos/%s/raw/%s",
		url.QueryEscape(projectKey),
		url.QueryEscape(repoSlug),
		url.QueryEscape(path)),
		DefaultPagination(), structs.Map(filter))
	if err != nil {
		return nil, err
	}
	return body, nil

}

func (c *Client) GetProjectsReposCommits(projectKey string, repoSlug string, pagination *Pagination) ([]*Commit, error) {
	body, err := c.get(fmt.Sprintf("projects/%s/repos/%s/commits", url.QueryEscape(projectKey), url.QueryEscape(repoSlug)),
		pagination, nil)
	if err != nil {
		return nil, err
	}
	var answer Answer
	err = json.Unmarshal(body, &answer)
	fillPagination(&answer, pagination)
	commitList := make([]*Commit, 0)
	for _, v := range answer.Values {
		var repo Commit
		m2s(v, &repo)
		commitList = append(commitList, &repo)
	}
	return commitList, nil
}

func (c *Client) GetProjectReposCommit(projectKey, repoSlug, commitId string) (*Commit, error) {
	body, err := c.get(fmt.Sprintf("projects/%s/repos/%s/commits/%s", url.QueryEscape(projectKey), url.QueryEscape(repoSlug), url.QueryEscape(commitId)),
		DefaultPagination(), nil)
	if err != nil {
		return nil, err
	}
	var commit Commit
	err = json.Unmarshal(body, &commit)
	if err != nil {
		return nil, err
	}
	return &commit, nil

}

func (c *Client) GetProjectsReposBranches(projectKey string, repoSlug string, pagination *Pagination, filter *ProjectReposBranchesFilter) ([]*Branch, error) {
	body, err := c.get(fmt.Sprintf("projects/%s/repos/%s/branches", url.QueryEscape(projectKey), url.QueryEscape(repoSlug)),
		pagination,
		structs.Map(filter))
	if err != nil {
		return nil, err
	}
	var answer Answer
	err = json.Unmarshal(body, &answer)
	fillPagination(&answer, pagination)
	branchList := make([]*Branch, 0)
	for _, v := range answer.Values {
		var branch Branch
		m2s(v, &branch)
		branchList = append(branchList, &branch)
	}
	return branchList, nil
}

func (c *Client) GetProjectsReposBranchesDefault(projectKey string, repoSlug string, pagination *Pagination) (*Branch, error) {
	body, err := c.get(fmt.Sprintf("projects/%s/repos/%s/branches/default", url.QueryEscape(projectKey), url.QueryEscape(repoSlug)),
		pagination, nil)
	if err != nil {
		return nil, err
	}
	var branch Branch
	err = json.Unmarshal(body, &branch)
	if err != nil {
		return nil, err
	}
	return &branch, nil
}
