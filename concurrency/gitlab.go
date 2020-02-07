package concurrency

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/xanzy/go-gitlab"
)

type GitlabGroupProvider struct {
	Client       *gitlab.Client
	Token        string
	RootFullPath string
	Location     string
	ID           int
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
