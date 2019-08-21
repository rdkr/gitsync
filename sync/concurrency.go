package sync

import (
	"sync"
)

type ProviderProcessor interface {
	GetGroups() []ProviderProcessor
	GetProjects() []Project
}

type ConcurrencyManager struct {
	cfg Config

	groups   chan ProviderProcessor
	projects chan Project

	groupsWG, groupsSignalWG, projectsWG, projectsSignalWG *sync.WaitGroup
	groupsSignalOnce, projectsSignalOnce                   *sync.Once

	getItemsFromCfg ConfigParser
	gitSync         GitSyncer
	ui              UI
}

func NewConcurrencyManager(cfg Config, ui UI, configParser ConfigParser, gitSync GitSyncer) ConcurrencyManager {
	return ConcurrencyManager{
		cfg:                cfg,
		groups:             make(chan ProviderProcessor),
		projects:           make(chan Project),
		groupsWG:           new(sync.WaitGroup),
		groupsSignalWG:     new(sync.WaitGroup),
		groupsSignalOnce:   new(sync.Once),
		projectsWG:         new(sync.WaitGroup),
		projectsSignalWG:   new(sync.WaitGroup),
		projectsSignalOnce: new(sync.Once),
		getItemsFromCfg:    configParser,
		gitSync:            gitSync,
		ui:                 ui,
	}
}

func (cm ConcurrencyManager) Start() {

	var wg sync.WaitGroup
	wg.Add(3)

	cm.groupsWG.Add(1)
	cm.groupsSignalWG.Add(1)

	cm.projectsWG.Add(2)
	cm.projectsSignalWG.Add(1)

	// groups manager goroutine
	go func() {

		// Start some groups processors
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

		// Start some projects processors
		for w := 0; w < 20; w++ {
			go cm.processProject()
		}

		// wait for a signal that we have a Project to process
		cm.projectsSignalWG.Wait()

		// wait to finish processing all projects
		cm.projectsWG.Wait()
		close(cm.projects)

		// ensure we have processed all projects before stopping the UI
		close(cm.ui.StatusChan)

		// stop the projects manager goroutine
		wg.Done()

	}()

	// UI manager goroutine
	go func() {

		cm.ui.Run()
		wg.Done()

	}()

	groups, projects := cm.getItemsFromCfg(cm.cfg)

	for _, g := range groups {
		cm.groupsWG.Add(1)
		cm.groups <- g
	}

	for _, p := range projects {
		cm.projectsWG.Add(1)
		cm.projectsSignalOnce.Do(func() { cm.projectsSignalWG.Done() })

		go func(project Project) {
			cm.projects <- project
		}(p)
	}

	cm.groupsWG.Done()
	cm.groupsSignalOnce.Do(func() { cm.groupsSignalWG.Done() })

	cm.projectsWG.Done()
	cm.projectsSignalOnce.Do(func() { cm.projectsSignalWG.Done() })

	wg.Wait()
}

func (cm ConcurrencyManager) processGroups() {
	for {

		parent, ok := <-cm.groups
		if !ok {
			break
		}

		childGroups := parent.GetGroups()

		for _, child := range childGroups {
			cm.groupsWG.Add(1)
			cm.groupsSignalOnce.Do(func() { cm.groupsSignalWG.Done() })
			go func(group ProviderProcessor) {
				cm.groups <- group
			}(child)
		}

		childProjects := parent.GetProjects()

		for _, child := range childProjects {
			cm.projectsWG.Add(1)
			cm.projectsSignalOnce.Do(func() { cm.projectsSignalWG.Done() })
			go func(project Project) {
				cm.projects <- project
			}(child)
		}

		cm.groupsWG.Done()
	}
}

func (cm ConcurrencyManager) processProject() {
	for {

		project, ok := <-cm.projects
		if !ok {
			break
		}

		cm.ui.StatusChan <- cm.gitSync(project, project.Location)
		cm.projectsWG.Done()
	}
}
