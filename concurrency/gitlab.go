package concurrency

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
	"strings"
)

type GitlabGroupProvider struct {
	Client       *gitlab.Client
	Token        string
	RootFullPath string
	Location     string
	ID           int
}

type GitlabManager struct {
	GroupChan   chan Status
	ProjectChan chan Status
	manager
}

func NewGitlabManager(projectAction projectActionFunc) GitlabManager {
	return GitlabManager{
		GroupChan:   make(chan Status),
		ProjectChan: make(chan Status),
		manager:     newManager(projectAction),
	}
}

func (m GitlabManager) projectChanSender(projectAction projectActionFunc, project Project) {
	m.ProjectChan <- projectAction(project)
}

func (m GitlabManager) projectsChanCloser(){
	close(m.ProjectChan)
}

func (m GitlabManager) Start(groups []ProviderProcessor, projects []Project) {
	m.start(groups, projects, m.projectChanSender, m.projectsChanCloser)
}

func (g *GitlabGroupProvider) GetGroups() []ProviderProcessor {
	var result []ProviderProcessor

	parent, _, err := g.Client.Groups.GetGroup(g.ID)
	if err != nil {
		logrus.Fatal(err)
	}

	if g.RootFullPath == "" {
		g.RootFullPath = parent.FullPath
	}

	groups, _, err := g.Client.Groups.ListSubgroups(parent.ID, &gitlab.ListSubgroupsOptions{
		AllAvailable: gitlab.Bool(true),
	})
	if err != nil {
		panic(err)
	}

	for _, child := range groups {
		result = append(result, &GitlabGroupProvider{g.Client, g.Token, g.RootFullPath, g.Location, child.ID})
	}

	return result
}

func (g *GitlabGroupProvider) GetProjects() []Project {
	var result []Project

	projects, _, err := g.Client.Groups.ListGroupProjects(g.ID, &gitlab.ListGroupProjectsOptions{
		Archived: gitlab.Bool(false),
	})
	if err != nil {
		panic(err)
	}

	for _, p := range projects {

		path := strings.ReplaceAll(p.PathWithNamespace, g.RootFullPath, "")
		path = strings.TrimLeft(path, "/")
		path = fmt.Sprintf("%s/%s", g.Location, path)

		result = append(result, Project{p.HTTPURLToRepo, path, g.Token})
	}

	return result
}
