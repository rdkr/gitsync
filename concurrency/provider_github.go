package concurrency

import (
	"context"

	"github.com/google/go-github/v30/github"
	"github.com/sirupsen/logrus"
)

type GithubUserGroup struct {
	Client   *github.Client
	Name     string
	Location string
	Token    string
}

func (u *GithubUserGroup) GetGroups() []Group {
	return []Group{}
}

func (u *GithubUserGroup) GetProjects() []Project {
	var result []Project

	logrus.Debug("getting projects")

	projects, _, err := u.Client.Repositories.List(context.Background(), "", &github.RepositoryListOptions{
		Type: "owner",
	})
	if err != nil {
		logrus.Fatal(err)
	}

	for _, p := range projects {
		if !*p.Archived {
			result = append(result, Project{*p.CloneURL, u.Location + "/" + *p.Name, u.Token})
		}
	}

	return result
}

type GithubOrgGroup struct {
	Client   *github.Client
	Name     string
	Location string
	Token    string
}

func (u *GithubOrgGroup) GetGroups() []Group {
	return []Group{}
}

func (o *GithubOrgGroup) GetProjects() []Project {
	var result []Project

	logrus.Debug("getting projects by org")

	if o.Name != "" {
		projectsByOrg, _, err := o.Client.Repositories.ListByOrg(context.Background(), o.Name, &github.RepositoryListByOrgOptions{})

		if err != nil {
			logrus.Fatal(err)
		}

		for _, p := range projectsByOrg {
			if !*p.Archived {
				result = append(result, Project{*p.CloneURL, o.Location + "/" + *p.Name, o.Token})
			}
		}
	}

	return result
}
