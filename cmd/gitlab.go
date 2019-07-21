package cmd

import (
	"fmt"
	"os"
	"strings"
	"sync"

	gitlab "github.com/xanzy/go-gitlab"
)

type gitlabProvider struct {
	c        *gitlab.Client
	root     *gitlab.Group
	prefix   string
	groups   chan int
	projects chan *gitlab.Project

	groupsWG, projectsWG, projectsSignalWG *sync.WaitGroup
	projectSignalOnce                      *sync.Once
}

func newGitlabProvider(group int) gitlabProvider {

	c := gitlab.NewClient(nil, os.Getenv("GITLAB_TOKEN"))

	root, _, err := c.Groups.GetGroup(group)
	if err != nil {
		panic("o no 123")
	}

	return gitlabProvider{
		c:                 c,
		root:              root,
		groups:            make(chan int),
		projects:          make(chan *gitlab.Project),
		groupsWG:          new(sync.WaitGroup),
		projectsWG:        new(sync.WaitGroup),
		projectsSignalWG:  new(sync.WaitGroup),
		projectSignalOnce: new(sync.Once),
	}
}

func (g gitlabProvider) recurseGroups() {
	// fmt.Println("recurseGroups...")
	for {

		groupID, ok := <-g.groups
		if !ok {
			break
		}

		subgroups, _, err := g.c.Groups.ListSubgroups(groupID, &gitlab.ListSubgroupsOptions{
			AllAvailable: gitlab.Bool(true),
		})
		if err != nil {
			panic(err)
		}

		for _, group := range subgroups {
			// fmt.Println("group: " + group.FullName)
			g.groupsWG.Add(1)
			go func(id int) {
				g.groups <- id
			}(group.ID)
		}

		groupProjects, _, err := g.c.Groups.ListGroupProjects(groupID, &gitlab.ListGroupProjectsOptions{Archived: gitlab.Bool(false)})
		if err != nil {
			panic(err)
		}

		for _, project := range groupProjects {
			g.projectsWG.Add(1)
			g.projectSignalOnce.Do(func() { g.projectsSignalWG.Done() })
			go func(p *gitlab.Project) {
				g.projects <- p
			}(project)
		}

		g.groupsWG.Done()

	}
}

func (g gitlabProvider) processProject(ui ui) {
	// fmt.Println("processProject...")
	for {

		project, ok := <-g.projects
		if !ok {
			break
		}

		path := strings.ReplaceAll(project.PathWithNamespace, g.root.FullPath, "")
		path = strings.TrimLeft(path, "/")
		path = fmt.Sprintf("%s/%s", g.prefix, path)

		ui.statusChan <- clone(path, project.HTTPURLToRepo)

		g.projectsWG.Done()
	}
}
