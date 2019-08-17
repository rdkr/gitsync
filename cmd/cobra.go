package cmd

import (
	"fmt"
	"gitsync/sync"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool
var cfg sync.Config

const usage = `gitsync is a tool to keep local git repos in sync with remote git servers.

It supports individual repos and git service provider groups accessed over
HTTPS and authenticated either anonymously or with a token. Groups are
recursed to find projects and projects are concurrently cloned, pulled, or
fetched as appropriate.

Supported git service providers:
 - GitLab groups and repos over HTTPS (GITLAB_TOKEN env var should be set)
 - Anonymous repos over HTTPS

A .yaml config file is expected, and will be found from:
 - $HOME/.gitsync.yaml
 - $PWD/.gitsync.yaml
 - as specified using the --config/-c flag

The format of the config file is as follows:

gitlab:
  groups:
  - group: <group-id>
    location: <local path to sync to>
  projects:
  - url: <https clone url>
    location: <local path to sync to>
anon:
  projects:
  - url: <https clone url>
	location: <local path to sync to>`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitsync",
	Short: "A tool to keep many git repos in sync with their remote origins",
	Long:  usage,
	Run: func(cmd *cobra.Command, args []string) {
		cm := sync.NewConcurrencyManager(cfg, sync.NewUI(verbose), sync.GetItemsFromCfg, sync.GitSync)
		cm.Start()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.gitsync.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose Output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {

		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1) // todo log to stderr ?
		}

		viper.SetConfigName(".gitsync")
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")

	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println(err)
	}

	err := viper.Unmarshal(&cfg)
	if err != nil {
		fmt.Println(err)
	}
}
