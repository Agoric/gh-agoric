/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"log"

	gh "github.com/cli/go-gh"
	"github.com/spf13/cobra"
)

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
		itemList, err := NewGHItemList(options.projNum, options.owner, options.limit)
		if err != nil {
			log.Fatal(err)
			return
		}

		var b []byte
		b, err = json.MarshalIndent(&itemList, "", "  ")
		if nil != err {
			log.Fatal(err)
			return
		}

		log.Print(string(b[:]))

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
}

type GHContent struct {
	Type       string `json:"type"`
	Body       string `json:"body"`
	Title      string `json:"title"`
	Number     int    `json:"number"`
	Repository string `json:"repository"`
	URL        string `json:"url"`
}

type GHItem struct {
	Assignees          []string  `json:"assignees"`
	LinkedPullRequests []string  `json:"linked pull requests"`
	Content            GHContent `json:"content"`
	Repository         string    `json:"repository"`
	Status             string    `json:"status"`
	Id                 string    `json:"id"`
	DueDate            string    `json:"due Date"`
	StartDate          string    `json:"start Date"`
	Notion             string    `json:"notion"`
	Upgrade            string    `json:"upgrade"`
}

type GHItemList struct {
	Items      []GHItem
	TotalCount int
}

func NewGHItemList(projNum string, owner string, limit string) (*GHItemList, error) {
	args := []string{"project", "item-list", options.projNum, "--owner", options.owner, "--limit", options.limit, "--format", "json"}
	stdOut, stdErr, err := gh.Exec(args...)
	if err != nil {
		log.Fatal(stdErr)
		log.Fatal(err)
		return nil, err
	}

	var itemList GHItemList
	err = json.Unmarshal(stdOut.Bytes(), &itemList)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &itemList, nil
}
