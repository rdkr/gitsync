package main

import (
	"errors"
	"os"
	"sort"
	"sync"

	"github.com/mitchellh/go-homedir"
	"github.com/rdkr/gitsync/concurrency"
	gitsync "github.com/rdkr/gitsync/sync"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

var cfgFile string
var verbose bool
var debug bool
var unmanaged bool
var cfg concurrency.Config

const usage = `gitsync is a tool to keep many local repos in sync with their remote hosts.

It supports recursively syncing GitHub orgs, teams, and users; GitLab groups; and individual
repos. Repos are synced over HTTPS, optionally / where required using auth tokens.

              Orgs'    Groups' / Teams'    Users'    Repos'
    GitHub      x              x             x
    GitLab                     x                       x
    HTTPS                                              x

Orgs / groups / user profiles are enumerated / recursed to find projects. All projects
are then concurrently synced, i.e:
 - cloned, if the local repo doesn't exist
 - pulled, if the local repo exists and is on main
 - fetched, if neither of the above

A .yaml config file is expected, The format of the config file is:

github:       # optional: defines GitHub resources
- baseurl:      # optional: a custom GitHub API URL
  token:        # required: a GitHub API token
  users:        # optional: defines GitHub users
  - name:         # required: GitHub username
    location:     # required: local path to sync to
  orgs:         # optional: defines GitHub organisations
  - name:         # required: GitHub org name
    location:     # required: local path to sync to
  teams:        # optional: defines GitHub teams
  - org:          # required: GitHub org name
    name:         # required: GitHub team name
    location:     # required: local path to sync to
gitlab:       # optional: defines GitLab resources
- baseurl:      # optional: a custom GitLab API URL
  token:        # optional: a GitLab API token
  groups:       # optional: defines GitLab groups
  - group:        # required: group ID number
    location:     # required: local path to sync to
  projects:     # optional: defines GitLab projects
  - url:          # required: https clone url
    location:     # required: local path to sync to
anon:         # optional: defines any other resources
  projects:     # optional: defines any HTTPS projects
  - url:          # required: https clone url
    location:     # required: local path to sync to
    token:        # optional: HTTPS token to use

The config file will will be found, by order of precedence, from:
 - $HOME/.gitsync.yaml
 - $PWD/.gitsync.yaml
 - as specified using the --config/-c flag

Treat this file with care, as it may contain secrets.`

var rootCmd = &cobra.Command{
	Use:   "gitsync",
	Short: "gitsync is a tool to keep many local repos in sync with their remote hosts.",
	Long:  usage,
	Run: func(cmd *cobra.Command, args []string) {

		// empty lists to store root directories to sync to and all paths that
		// are synced to in order to find 'unmanaged' git directories
		var syncLocations []string
		var syncedPaths []string

		// create ui
		isTerminal := terminal.IsTerminal(int(os.Stdout.Fd()))
		ui := gitsync.NewUI(isTerminal, verbose, debug)

		// create concurrency managers
		gl := concurrency.NewManager(gitsync.GitSyncHelper)
		gh := concurrency.NewManager(gitsync.GitSyncHelper)

		// create status merger channel
		statuses := make(chan interface{})

		// create wait group to manage the above
		wg := sync.WaitGroup{}
		wg.Add(4)

		// start concurrency manager
		go func() {
			groups, projects, gitlabLocations := concurrency.GetGitlabItemsFromCfg(cfg)
			gl.Start(groups, projects)
			syncLocations = append(syncLocations, gitlabLocations...)

			groups, projects, githubLocations := concurrency.GetGithubItemsFromCfg(cfg)
			gh.Start(groups, projects)
			syncLocations = append(syncLocations, githubLocations...)

			wg.Done()
		}()

		// start ui
		go func() {
			ui.Run()
			wg.Done()
		}()

		// start status merger
		go func() {
			concurrency.ChannelMerger(statuses, gl.ProjectChan, gh.ProjectChan)
			wg.Done()
		}()

		// connect cm and ui
		go func() {
			for {
				status, ok := <-statuses
				if !ok {
					break
				}
				s, _ := status.(gitsync.Status)
				syncedPaths = append(syncedPaths, s.Path)
				ui.StatusChan <- s
			}

			if unmanaged {
				var unmanagedPaths []string
				for _, path := range gitsync.Unmanaged(syncLocations, syncedPaths) {
					unmanagedPaths = append(unmanagedPaths, path)
				}
				sort.Strings(unmanagedPaths)
				for _, path := range unmanagedPaths {
					err := errors.New("not in upstream parent")
					ui.StatusChan <- gitsync.Status{path, gitsync.StatusUnmanaged, "", err}
				}
			}

			close(ui.StatusChan)
			wg.Done()
		}()

		// wait until all the above are done before exiting
		wg.Wait()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file location")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output instead of pretty output")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug output (implies verbose)")
	rootCmd.PersistentFlags().BoolVarP(&unmanaged, "unmanaged", "u", false, "shows local repos not in upstream parent")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {

		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			logrus.Fatal(err)
		}

		viper.SetConfigName(".gitsync")
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")

	}

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		logrus.Fatal(err)
	}
}
