package cmd

import (
	"os"
	"sync"

	gitlab "github.com/xanzy/go-gitlab"
)

type syncer struct {
	groups   chan group
	projects chan project

	groupsWG, projectsWG, projectsSignalWG *sync.WaitGroup
	projectSignalOnce                      *sync.Once

	ui ui
}

func newSyncer(ui ui) syncer {
	return syncer{
		groups:            make(chan group),
		projects:          make(chan project),
		groupsWG:          new(sync.WaitGroup),
		projectsWG:        new(sync.WaitGroup),
		projectsSignalWG:  new(sync.WaitGroup),
		projectSignalOnce: new(sync.Once),
		ui:                ui,
	}
}

func (s syncer) recurseGroups() {
	for {

		parent, ok := <-s.groups
		if !ok {
			break
		}

		s.ui.currentParent = parent.rootFullPath()

		childGroups := parent.getGroups()

		for _, child := range childGroups {
			s.groupsWG.Add(1)
			go func(group group) {
				s.groups <- group
			}(child)
		}

		childProjects := parent.getProjects()

		for _, child := range childProjects {
			s.projectsWG.Add(1)
			s.projectSignalOnce.Do(func() { s.projectsSignalWG.Done() })
			go func(project project) {
				s.projects <- project
			}(child)
		}

		s.groupsWG.Done()
	}
}

func (s syncer) processProject() {
	for {

		project, ok := <-s.projects
		if !ok {
			break
		}

		s.ui.statusChan <- status{project.Location, "", "", nil}
		s.ui.statusChan <- project.clone()
		s.projectsWG.Done()
	}
}

func syncronise() {
	ui := newUI(verbose)
	s := newSyncer(ui)

	token := os.Getenv("GITLAB_TOKEN")
	c := gitlab.NewClient(nil, token)

	var rootGroups []gitlabGroupProvider

	for _, item := range cfg.Gitlab.Groups {
		root, _, err := c.Groups.GetGroup(item.Group)
		if err != nil {
			panic("bad token?")
		}
		rootGroups = append(rootGroups, gitlabGroupProvider{c, token, root.FullPath, item.Location, root})
	}

	var wg sync.WaitGroup
	wg.Add(3)

	s.projectsWG.Add(1)       // hold this open until all groups are finished processing as we don't have a 'seed' project as with groups
	s.projectsSignalWG.Add(1) // hold this open until at least one project has been found TODO need to handle if there are no projects :O

	go func() {

		for w := 0; w < 10; w++ {
			go s.recurseGroups()
		}

		for _, group := range rootGroups {
			s.groupsWG.Add(1)
			s.groups <- group
		}

		s.groupsWG.Wait()
		close(s.groups)

		s.projectsWG.Done()
		wg.Done()

	}()

	go func() {

		for _, p := range cfg.Gitlab.Projects {
			s.projectsWG.Add(1)
			s.projectSignalOnce.Do(func() { s.projectsSignalWG.Done() })

			if p.Token == "" {
				p.Token = token
			}
			go func(project project) {
				s.projects <- project
			}(p)
		}

		for _, p := range cfg.Anon.Projects {
			s.projectsWG.Add(1)
			s.projectSignalOnce.Do(func() { s.projectsSignalWG.Done() })

			go func(project project) {
				s.projects <- project
			}(p)
		}

		s.projectsSignalWG.Wait()

		for w := 0; w < 20; w++ {
			go s.processProject()
		}

		s.projectsWG.Wait()
		close(s.projects)
		close(ui.statusChan)

		wg.Done()

	}()

	go func() {

		ui.run()
		wg.Done()

	}()

	wg.Wait()
}
