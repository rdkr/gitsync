package concurrency

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

type GitlabGroup struct {
	Client       *gitlab.Client
	Token        string
	RootFullPath string
	Location     string
	ID           int
}

func (g *GitlabGroup) GetGroups() []Group {
	var result []Group

	logrus.WithField("user", g.ID).Debug("getting gitlab groups by group")

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
		// Response code is 404 if missing permissions for ListSubgroups
		if errVal := err.(*gitlab.ErrorResponse); errVal.Response.StatusCode == 404 {
			logrus.Warn(err)
		} else {
			logrus.Fatal(err)
		}
	}

	for _, child := range groups {
		result = append(result, &GitlabGroup{g.Client, g.Token, g.RootFullPath, g.Location, child.ID})
	}

	return result
}

func (g *GitlabGroup) GetProjects() []Project {
	var result []Project

	logrus.WithField("user", g.ID).Debug("getting gitlab projects by group")

	projects, _, err := g.Client.Groups.ListGroupProjects(g.ID, &gitlab.ListGroupProjectsOptions{
		Archived: gitlab.Bool(false),
	})
	if err != nil {
		logrus.Fatal(err)
	}

	for _, p := range projects {

		path := strings.ReplaceAll(p.PathWithNamespace, g.RootFullPath, "")
		path = strings.TrimLeft(path, "/")
		path = fmt.Sprintf("%s/%s", g.Location, path)

		result = append(result, Project{p.HTTPURLToRepo, path, g.Token})
	}

	return result
}
