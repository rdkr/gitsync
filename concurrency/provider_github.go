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

func (g *GithubUserGroup) GetGroups() []Group {
	return []Group{}
}

func (g *GithubUserGroup) GetProjects() []Project {
	var result []Project

	logrus.Debug("getting projects by user")

	var allRepos []*github.Repository
	opt := &github.RepositoryListOptions{
		Type: "owner",
	}
	for {
		repos, resp, err := g.Client.Repositories.List(context.Background(), "", opt)
		if err != nil {
			logrus.Fatal(err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, r := range allRepos {
		if !*r.Archived {
			result = append(result, Project{*r.CloneURL, g.Location + "/" + *r.Name, g.Token})
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

func (g *GithubOrgGroup) GetGroups() []Group {
	return []Group{}
}

func (g *GithubOrgGroup) GetProjects() []Project {
	var result []Project

	logrus.Debug("getting projects by org")

	var allRepos []*github.Repository
	opt := &github.RepositoryListByOrgOptions{}
	for {
		repos, resp, err := g.Client.Repositories.ListByOrg(context.Background(), g.Name, opt)
		if err != nil {
			logrus.Fatal(err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, r := range allRepos {
		if !*r.Archived {
			result = append(result, Project{*r.CloneURL, g.Location + "/" + *r.Name, g.Token})
		}
	}

	return result
}

type GithubTeamGroup struct {
	Client   *github.Client
	Org      string
	Name     string
	Location string
	Token    string
}

func (g *GithubTeamGroup) GetGroups() []Group {
	var result []Group

	logrus.Debug("getting groups by team")

	var allTeams []*github.Team
	opt := &github.ListOptions{}
	for {
		teams, resp, err := g.Client.Teams.ListChildTeamsByParentSlug(context.Background(), g.Org, g.Name, opt)
		if err != nil {
			logrus.Fatal(err)
		}
		allTeams = append(allTeams, teams...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, child := range allTeams {
		result = append(result, &GithubTeamGroup{g.Client, g.Org, *child.Slug, g.Location + "/" + *child.Slug, g.Token})
	}

	return result

}

func (g *GithubTeamGroup) GetProjects() []Project {
	var result []Project

	logrus.Debug("getting projects by team")

	var allRepos []*github.Repository
	opt := &github.ListOptions{}
	for {
		repos, resp, err := g.Client.Teams.ListTeamReposBySlug(context.Background(), g.Org, g.Name, opt)
		if err != nil {
			logrus.Fatal(err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, r := range allRepos {
		if !*r.Archived {
			result = append(result, Project{*r.CloneURL, g.Location + "/" + *r.Name, g.Token})
		}
	}

	return result
}
