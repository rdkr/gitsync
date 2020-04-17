package concurrency

import (
	"context"

	"github.com/google/go-github/v30/github"
	"github.com/sirupsen/logrus"
)

type GithubManager struct {
	GroupChan   chan error
	ProjectChan chan interface{}
	manager
}

func NewGithubManager(projectAction projectActionFunc) GithubManager {
	return GithubManager{
		GroupChan:   make(chan error),
		ProjectChan: make(chan interface{}),
		manager:     newManager(projectAction),
	}
}

func (m GithubManager) projectChanSender(projectAction projectActionFunc, project Project) {
	m.ProjectChan <- projectAction(project)
}

func (m GithubManager) projectsChanCloser() {
	close(m.ProjectChan)
}

func (m GithubManager) Start(users []User, groups []Group, projects []Project) {
	m.start(users, groups, projects, m.projectChanSender, m.projectsChanCloser)
}

type GithubUser struct {
	Client   *github.Client
	Name     string
	Location string
	Token    string
}

func (u *GithubUser) GetProjects() []Project {
	var result []Project

	logrus.Debug("getting projects")

	projects, _, err := u.Client.Repositories.List(context.Background(), "", &github.RepositoryListOptions{
		Type: "owner",
	})
	if err != nil {
		logrus.Fatal(err)
	}

	for _, p := range projects {
		if *p.Archived != true {
			result = append(result, Project{*p.CloneURL, u.Location + "/" + *p.Name, u.Token})
		}
	}

	return result
}

type GithubGroup struct {
	Client       *github.Client
	Token        string
	RootFullPath string
	Location     string
	ID           int
}

func (g *GithubGroup) GetGroups() []Group {
	var result []Group

	// parent, _, err := g.Client.Groups.GetGroup(g.ID)
	// if err != nil {
	// 	logrus.Fatal(err)
	// }

	// if g.RootFullPath == "" {
	// 	g.RootFullPath = parent.FullPath
	// }

	// groups, _, err := g.Client.Groups.ListSubgroups(parent.ID, &github.ListSubgroupsOptions{
	// 	AllAvailable: github.Bool(true),
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// for _, child := range groups {
	// 	result = append(result, &GithubGroup{g.Client, g.Token, g.RootFullPath, g.Location, child.ID})
	// }

	return result
}

func (g *GithubGroup) GetProjects() []Project {
	var result []Project

	// projects, _, err := g.Client.Groups.ListGroupProjects(g.ID, &github.ListGroupProjectsOptions{
	// 	Archived: github.Bool(false),
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// for _, p := range projects {

	// 	path := strings.ReplaceAll(p.PathWithNamespace, g.RootFullPath, "")
	// 	path = strings.TrimLeft(path, "/")
	// 	path = fmt.Sprintf("%s/%s", g.Location, path)

	// 	result = append(result, Project{p.HTTPURLToRepo, path, g.Token})
	// }

	return result
}
