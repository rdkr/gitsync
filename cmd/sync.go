/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"sync"

	"github.com/spf13/cobra"
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
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
