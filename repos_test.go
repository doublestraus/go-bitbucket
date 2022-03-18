package bitbucket

import (
	"os"
	"testing"
	"time"
)

func createClient() *Client {
	token := os.Getenv("BB_TOKEN")
	url := os.Getenv("BB_URL")
	return New(token, url)
}

func createMaxConnClient(maxConns int) *Client {
	token := os.Getenv("BB_TOKEN")
	url := os.Getenv("BB_URL")
	return New(token, url, WithMaxConnections(maxConns))
}

func createMaxWaitConnClient(maxWaitTimout time.Duration) *Client {
	token := os.Getenv("BB_TOKEN")
	url := os.Getenv("BB_URL")
	return New(token, url, WithMaxTimeoutWait(maxWaitTimout))
}

func ceateMaxConnMaxWaitClient(maxConns int, maxWaitTimeout time.Duration) *Client {
	token := os.Getenv("BB_TOKEN")
	url := os.Getenv("BB_URL")
	return New(token, url, WithMaxConnections(maxConns), WithMaxTimeoutWait(maxWaitTimeout))
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
	client := createMaxConnClient(10)
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
	client := createMaxWaitConnClient(5 * time.Minute)
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
	client := ceateMaxConnMaxWaitClient(15, 5*time.Minute)
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
