package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type configItem struct {
	Group    int    `yaml:"group"`
	Location string `yaml:"location"`
}

type config struct {
	Items []configItem `yaml:"items"`
}

var cfgFile string
var cfg config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitsync",
	Short: "A tool to keep many git repos in sync with their remote origins",
	Long: `gitsync is a tool to help keep many local instances of git repos in sync
  with their remote origins. It supports individual repos and git service
  provider groups.`,
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gitsync.yaml)")
	// rootCmd.PersistentFlags().StringVar(&author, "author", "YOUR NAME", "Author name for copyright attribution")
	// viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile != "" {

		// Use config file from the flag.
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

	viper.AutomaticEnv() // read in environment variables that match

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
