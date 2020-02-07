package main

import (
	"github.com/mitchellh/go-homedir"
	"github.com/rdkr/gitsync/concurrency"
	"github.com/rdkr/gitsync/sync"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

var cfgFile string
var verbose bool
var debug bool
var cfg concurrency.Config

const usage = `gitsync is a tool to keep local Git repos in sync with remote Git hosts.

It supports syncing individual Git repos and recursively syncing Git host groups.

Supported Git hosts:
 - GitLab groups and repos over HTTPS (optionally with auth token)
 - Anonymous repos over HTTPS

Groups are recursed to find projects. Projects are concurrently synced, i.e:
 - cloned, if the local repo doesn't exist
 - pulled, if the local repo exists and is on master
 - fetched, if neither of the above

A .yaml config file is expected, The format of the config file is:

gitlab:       # optional: defines GitLab resources
  token:        # required: a GitLab API token
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
	Short: "A tool to keep local Git repos in sync with remote Git servers",
	Long:  usage,
	Run: func(cmd *cobra.Command, args []string) {

		// create and run ui
		isTerminal := terminal.IsTerminal(int(os.Stdout.Fd()))
		ui := sync.NewUI(isTerminal, verbose, debug)
		go ui.Run()

		// create and run concurrency with gitsync
		cm := concurrency.NewManager(cfg, sync.GitSyncHelper)
		go cm.Start(sync.GetItemsFromCfg(cfg))

		// hook ui into cm
		for {
			status, ok := <-cm.StatusChan
			if !ok {
				break
			}
			ui.StatusChan <- status
		}

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
