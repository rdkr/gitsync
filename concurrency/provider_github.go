package concurrency

import (
	"github.com/google/go-github/v30/github"
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
	m.start(groups, projects, m.projectChanSender, m.projectsChanCloser)
}

type GithubUser struct {
	Client *github.Client
	Name   string
}

func (g *GithubUser) GetProjects() []Project {
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
