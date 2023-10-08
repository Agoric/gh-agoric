/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type rootOptions struct {
	cfgFile   string
	projNum   string
	owner     string
	cacheFile string
	limit     string
}

var options rootOptions

const (
	DEFAULT_PROJNUM   = "19"
	DEFAULT_OWNER     = "Agoric"
	DEFAULT_CACHEFILE = "./.gh-agoric-cache.json"
	DEFAULT_LIMIT     = "4"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh-agoric",
	Short: "A gh extension that helps with Agoric project management.",
	Long:  `A gh extension that helps with Agoric project management.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&options.cfgFile, "config", "", "config file (default is $HOME/.gh-agoric.yaml)")

	rootCmd.PersistentFlags().StringVar(&options.projNum, "projnum", DEFAULT_PROJNUM, "GitHub Project number")
	rootCmd.PersistentFlags().StringVar(&options.owner, "owner", DEFAULT_OWNER, "GitHub Project owner")
	rootCmd.PersistentFlags().StringVar(&options.cacheFile, "cacheFile", DEFAULT_CACHEFILE, "gh-agoric cache file")
	rootCmd.PersistentFlags().StringVar(&options.limit, "limit", DEFAULT_LIMIT, "gh command limit")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if options.cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(options.cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".gh-agoric" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gh-agoric")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
