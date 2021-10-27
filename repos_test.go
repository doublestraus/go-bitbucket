package bitbucket

import (
	"os"
	"testing"
)

func createClient() *Client {
	token := os.Getenv("BB_TOKEN")
	url := os.Getenv("BB_URL")
	return New(token, url)
}

func TestClient_Repos(t *testing.T) {
	client := createClient()
	repPag := DefaultPagination()
	repCounter := 0
	flt := &ProjectsFilter{}
	for {
		reps, err := client.ListProjects(repPag, flt)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		repCounter += len(reps)
		if repPag.IsLastPage {
			break
		}
		repPag.Start = repPag.NextPageStart
	}
	t.Logf("Total reps: %d", repCounter)
}

func TestClient_GetNonExistentProject(t *testing.T) {
	projectName := os.Getenv("BB_PROJECTNAME")
	client := createClient()
	_, err := client.GetProject(projectName)
	if err, ok := err.(Errors); ok {
		if err.StatusCode != 404 {
			t.Fail()
			return
		}
	}
}

func TestClient_GetExistentProject(t *testing.T) {
	client := createClient()
	projectName := os.Getenv("BB_PROJECTNAME")
	project, err := client.GetProject(projectName)
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	t.Log(project.Name)
}

func TestClient_GetProjectRepos(t *testing.T) {
	client := createClient()
	projectName := os.Getenv("BB_PROJECTNAME")
	pagination := DefaultPagination()
	filter := &ProjectReposFilter{}
	_, err := client.GetProjectRepos(projectName, pagination, filter)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestClient_GetProjectsReposFiles(t *testing.T) {
	client := createClient()
	pagination := DefaultPagination()
	filter := &ProjectReposFileFilter{}
	projectName := os.Getenv("BB_PROJECTNAME")
	repoSlug := os.Getenv("BB_REPOSLUG")
	for {
		_, err := client.GetProjectsReposFiles(projectName, repoSlug, pagination, filter)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		if pagination.IsLastPage {
			break
		}
		pagination.Start = pagination.NextPageStart
	}
}

func TestClient_GetProjectsReposFileRaw(t *testing.T) {
	client := createClient()
	projectName := os.Getenv("BB_PROJECTNAME")
	repoSlug := os.Getenv("BB_REPOSLUG")
	fileName := os.Getenv("BB_FILENAME")
	content, err := client.GetProjectsReposFileRaw(projectName, repoSlug, fileName)
	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		t.Logf("Sneak peek: %35s", content)
	}
}

func TestClient_GetProjectsReposCommits(t *testing.T) {
	client := createClient()
	pagination := DefaultPagination()
	projectName := os.Getenv("BB_PROJECTNAME")
	repoSlug := os.Getenv("BB_REPOSLUG")
	commits := make([]*Commit, 0)
	cms, err := client.GetProjectsReposCommits(projectName, repoSlug, pagination)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	commits = append(commits, cms...)

}
