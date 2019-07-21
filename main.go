package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/gosuri/uilive"
	gitlab "github.com/xanzy/go-gitlab"
	"golang.org/x/crypto/ssh/terminal"
	git "gopkg.in/src-d/go-git.v4"
	git_http "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type status struct {
	path      string
	operation string
	err       error
}

type ui struct {
	isTerminal bool
	writer     *uilive.Writer
	goodCount  int
	badCount   int
	statusChan chan status
	statuses   []status
}

func newUI() ui {

	isTerminal := terminal.IsTerminal(int(os.Stdout.Fd()))

	writer := uilive.New() // TODO this is created even though its not necessarily used
	if isTerminal {
		writer.Start()
		fmt.Fprint(writer.Newline(), "getting root group... ")
	}

	return ui{
		isTerminal: isTerminal,
		writer:     writer,
		goodCount:  0,
		badCount:   0,
		statusChan: make(chan status),
		statuses:   []status{},
	}
}

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

func clone(path, url string) status {

	repo, err := git.PlainOpen(path)
	if err == git.ErrRepositoryNotExists {

		_, err := git.PlainClone(path, false, &git.CloneOptions{
			URL:      url,
			Progress: nil,
			Auth: &git_http.BasicAuth{
				Username: "token",
				Password: os.Getenv("GITLAB_TOKEN"),
			},
		})
		if err != nil {
			return status{path, "clone", err}
		}
		return status{path, "clone", nil}

	} else if err == nil {

		err = repo.Fetch(&git.FetchOptions{
			Progress: os.Stdout,
			Auth: &git_http.BasicAuth{
				Username: "token",
				Password: os.Getenv("GITLAB_TOKEN"),
			},
		})
		if err == nil || err == git.NoErrAlreadyUpToDate {
			return status{path, "fetch", nil}
		}
		return status{path, "fetch", err}

	}
	return status{path, "open", err}
}

func (ui *ui) makeUI(root string, status status) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("getting root group... %s\nprocessing projects... ", root))

	ui.statuses = append(ui.statuses, status)

	if status.err != nil {
		ui.badCount = ui.badCount + 1
	} else {
		ui.goodCount = ui.goodCount + 1
	}

	if ui.goodCount > 0 {
		sb.WriteString(fmt.Sprintf("%d \u001b[32m✔\u001b[0m", ui.goodCount))
	}
	if ui.badCount > 0 {
		sb.WriteString(fmt.Sprintf(" %d \u001b[31m✘\u001b[0m", ui.badCount))
	}

	sb.WriteString("\n")

	for _, status := range ui.statuses {
		if status.err != nil {
			sb.WriteString(fmt.Sprintf(" \u001b[31m✘\u001b[0m %s: %s\n", status.path, status.err))
		}
		// } else {
		// 	sb.WriteString(fmt.Sprintf(" \u001b[32m✔\u001b[0m %s\n", status.path))
		// }
	}

	return sb.String()
}

func (ui *ui) run(g gitlabProvider) {
	for {

		status, ok := <-ui.statusChan
		if !ok {
			break
		}

		if ui.isTerminal {
			fmt.Fprint(ui.writer.Newline(), ui.makeUI(g.root.FullPath, status))
			ui.writer.Flush() // it randomly prints multiple lines without this
		}
	}
}
func main() {

	ui := newUI()
	g := newGitlabProvider(1755573)

	g.prefix = "gitlab"

	var wg sync.WaitGroup
	wg.Add(3)

	g.projectsWG.Add(1)       // hold this open until all groups are finished processing as we don't have a 'seed' project as with groups
	g.projectsSignalWG.Add(1) // hold this open until at least one project has been found TODO need to handle if there are no projects :O

	go func() {

		for w := 0; w < 10; w++ {
			go g.recurseGroups()
		}

		g.groupsWG.Add(1)
		g.groups <- g.root.ID

		g.groupsWG.Wait()
		close(g.groups)

		g.projectsWG.Done()
		wg.Done()

	}()

	go func() {

		g.projectsSignalWG.Wait()

		for w := 0; w < 20; w++ {
			go g.processProject(ui)
		}

		g.projectsWG.Wait()
		close(g.projects)
		close(ui.statusChan)

		wg.Done()

	}()

	go func() {

		ui.run(g)
		wg.Done()

	}()

	wg.Wait()

}
