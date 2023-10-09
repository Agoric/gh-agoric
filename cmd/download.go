/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	gha "github.com/agoric/gh-agoric/internal"
	"github.com/spf13/cobra"
)

type downloadOptions struct {
	notion  string
	status  string
	upgrade string
}

var downloadOpts downloadOptions

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads all issues.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		client, err := gha.NewGHClient(options.cacheFile)
		if err != nil {
			log.Fatal(err)
			return
		}

		var items gha.GHItems
		items, err = client.ReqItems(options.projNum, options.owner, options.limit)
		if err != nil {
			log.Fatal(err)
			return
		}

		items = items.FilterByNotion(downloadOpts.notion)
		items = items.FilterByStatus(downloadOpts.status)
		items = items.FilterByUpgrade(downloadOpts.upgrade)

		log.Println(items.ToJson())
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	downloadCmd.Flags().StringVar(&downloadOpts.notion, "notion", "", "Filter by value of notion")
	downloadCmd.Flags().StringVar(&downloadOpts.status, "status", "", "Filter by value of status")
	downloadCmd.Flags().StringVar(&downloadOpts.upgrade, "upgrade", "", "Filter by value of upgrade")
}
