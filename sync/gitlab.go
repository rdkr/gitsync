package sync

import (
	"fmt"
	"strings"

	"github.com/xanzy/go-gitlab"
)

type gitlabGroupProvider struct {
	client       *gitlab.Client
	token        string
	rootFullPath string // TODO rename group root path
	location     string
	ID           int
}

func (g *gitlabGroupProvider) GetGroups() []ProviderProcessor {
	var result []ProviderProcessor

	parent, _, err := g.client.Groups.GetGroup(g.ID)
	if err != nil {
		panic("bad token?") // TODO fix this
	}

	if g.rootFullPath == "" {
		g.rootFullPath = parent.FullPath
	}

	groups, _, err := g.client.Groups.ListSubgroups(parent.ID, &gitlab.ListSubgroupsOptions{
		AllAvailable: gitlab.Bool(true),
	})
	if err != nil {
		panic(err)
	}

	for _, child := range groups {
		result = append(result, &gitlabGroupProvider{g.client, g.token, g.rootFullPath, g.location, child.ID})
	}

	return result
}

func (g *gitlabGroupProvider) GetProjects() []Project {
	var result []Project

	projects, _, err := g.client.Groups.ListGroupProjects(g.ID, &gitlab.ListGroupProjectsOptions{
		Archived: gitlab.Bool(false),
	})
	if err != nil {
		panic(err)
	}

	for _, p := range projects {

		path := strings.ReplaceAll(p.PathWithNamespace, g.rootFullPath, "")
		path = strings.TrimLeft(path, "/")
		path = fmt.Sprintf("%s/%s", g.location, path)

		result = append(result, Project{p.HTTPURLToRepo, path, g.token})
	}

	return result
}
