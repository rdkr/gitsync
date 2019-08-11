package cmd

import (
	"github.com/xanzy/go-gitlab"
	"os"
	"sync"
)

type groupProcessor interface {
	getGroups() []groupProcessor
	getProjects() []project
	rootFullPath() string
	rootLocation() string
}

type concurrencyManager struct {
	groups   chan groupProcessor
	projects chan project

	groupsWG, groupsSignalWG, projectsWG, projectsSignalWG *sync.WaitGroup
	groupsSignalOnce, projectsSignalOnce                   *sync.Once

	ui               ui
	projectProcessor func(Git, string) Status
}

func newConcurrencyManager(ui ui, projectProcessor func(Git, string) Status) concurrencyManager {
	return concurrencyManager{
		groups:             make(chan groupProcessor),
		projects:           make(chan project),
		groupsWG:           new(sync.WaitGroup),
		groupsSignalWG:     new(sync.WaitGroup),
		groupsSignalOnce:   new(sync.Once),
		projectsWG:         new(sync.WaitGroup),
		projectsSignalWG:   new(sync.WaitGroup),
		projectsSignalOnce: new(sync.Once),
		ui:                 ui,
		projectProcessor:   projectProcessor,
	}
}

func (cm concurrencyManager) start() {

	var wg sync.WaitGroup
	wg.Add(3)

	cm.groupsWG.Add(1)
	cm.groupsSignalWG.Add(1)

	cm.projectsWG.Add(2)
	cm.projectsSignalWG.Add(1)

	// groups manager goroutine
	go func() {

		// start some groups processors
		for w := 0; w < 10; w++ {
			go cm.processGroups()
		}

		// wait for a signal that we have a group to process
		cm.groupsSignalWG.Wait()

		// wait to finish processing all groups
		cm.groupsWG.Wait()
		close(cm.groups)

		// ensure we have processed all groups before stopping on projects
		cm.projectsWG.Done()

		// stop the groups manager goroutine
		wg.Done()

	}()

	// projects manager goroutine
	go func() {

		// start some projects processors
		for w := 0; w < 20; w++ {
			go cm.processProject()
		}

		// wait for a signal that we have a project to process
		cm.projectsSignalWG.Wait()

		// wait to finish processing all projects
		cm.projectsWG.Wait()
		close(cm.projects)

		// ensure we have processed all projects before stopping the ui
		close(cm.ui.statusChan)

		// stop the projects manager goroutine
		wg.Done()

	}()

	// ui manager goroutine
	go func() {

		cm.ui.run()
		wg.Done()

	}()

	cm.addItemsFromCfg()

	cm.groupsWG.Done()
	cm.groupsSignalOnce.Do(func() { cm.groupsSignalWG.Done() })

	cm.projectsWG.Done()
	cm.projectsSignalOnce.Do(func() { cm.projectsSignalWG.Done() })

	wg.Wait()
}

func (cm concurrencyManager) addItemsFromCfg() {

	var rootGroups []groupProcessor

	if len(cfg.Gitlab.Groups) > 0 || len(cfg.Gitlab.Projects) > 0 {

		// TODO improve the handline of no / bad token

		token := os.Getenv("GITLAB_TOKEN")
		if len(token) == 0 {
			panic("bad token?")
		}

		c := gitlab.NewClient(nil, token)

		for _, item := range cfg.Gitlab.Groups {
			root, _, err := c.Groups.GetGroup(item.Group)
			if err != nil {
				panic("bad token?")
			}
			rootGroups = append(rootGroups, gitlabGroupProvider{c, token, root.FullPath, item.Location, root})
		}

		for _, p := range cfg.Gitlab.Projects {
			cm.projectsWG.Add(1)
			cm.projectsSignalOnce.Do(func() { cm.projectsSignalWG.Done() })

			if p.Token == "" {
				p.Token = token
			}
			go func(project project) {
				cm.projects <- project
			}(p)
		}

	}

	for _, group := range rootGroups {
		cm.groupsWG.Add(1)
		cm.groups <- group
	}

	for _, p := range cfg.Anon.Projects {
		cm.projectsWG.Add(1)
		cm.projectsSignalOnce.Do(func() { cm.projectsSignalWG.Done() })

		go func(project project) {
			cm.projects <- project
		}(p)
	}
}

func (cm concurrencyManager) processGroups() {
	for {

		parent, ok := <-cm.groups
		if !ok {
			break
		}

		cm.ui.currentParent = parent.rootFullPath()

		childGroups := parent.getGroups()

		for _, child := range childGroups {
			cm.groupsWG.Add(1)
			cm.groupsSignalOnce.Do(func() { cm.groupsSignalWG.Done() })
			go func(group groupProcessor) {
				cm.groups <- group
			}(child)
		}

		childProjects := parent.getProjects()

		for _, child := range childProjects {
			cm.projectsWG.Add(1)
			cm.projectsSignalOnce.Do(func() { cm.projectsSignalWG.Done() })
			go func(project project) {
				cm.projects <- project
			}(child)
		}

		cm.groupsWG.Done()
	}
}

func (cm concurrencyManager) processProject() {
	for {

		project, ok := <-cm.projects
		if !ok {
			break
		}

		cm.ui.statusChan <- cm.projectProcessor(project, project.Location)
		cm.projectsWG.Done()
	}
}
