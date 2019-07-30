package cmd

import (
	"os"
	"sync"

	"github.com/spf13/cobra"
	gitlab "github.com/xanzy/go-gitlab"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		ui := newUI()
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
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	// syncCmd.Flags().IntVarP(&source, "source", "s", 0, "source group to read from")
	// syncCmd.MarkFlagRequired("source")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
