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
package internal

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/cli/go-gh"
)

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

type GHItems []GHItem

type GHItemList struct {
	Items      GHItems `json:"Items"`
	TotalCount int     `json:"TotalCount"`
}

func newGHItem(oldItem GHItem) GHItem {
	newItem := oldItem

	newItem.Assignees = nil
	newItem.Assignees = append(newItem.Assignees, oldItem.Assignees...)

	newItem.LinkedPullRequests = nil
	newItem.LinkedPullRequests = append(newItem.LinkedPullRequests, oldItem.LinkedPullRequests...)

	return newItem
}

func NewGHItemListFromBytes(b []byte) (*GHItemList, error) {
	var itemList GHItemList
	err := json.Unmarshal(b, &itemList)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Printf("GHItemList [len(%d) TotalCount(%d)]\n", len(itemList.Items), itemList.TotalCount)
	return &itemList, nil
}

func NewGHItemList(projNum string, owner string, limit string) (*GHItemList, error) {
	args := []string{"project", "item-list", projNum, "--owner", owner, "--limit", limit, "--format", "json"}
	stdOut, stdErr, err := gh.Exec(args...)
	if err != nil {
		log.Fatal(stdErr)
		log.Fatal(err)
		return nil, err
	}

	var itemList *GHItemList
	itemList, err = NewGHItemListFromBytes(stdOut.Bytes())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Printf("GHItemList [len(%d) TotalCount(%d)]\n", len(itemList.Items), itemList.TotalCount)

	return itemList, nil
}

func (itemList GHItemList) GetItems() (ret GHItems) {
	log.Printf("GHItemList has %d items\n", len(itemList.Items))

	for _, item := range itemList.Items {
		ret = append(ret, newGHItem(item))
	}
	return ret
}

func (items GHItems) ToJson() string {
	var b []byte
	b, err := json.MarshalIndent(items, "", "  ")
	if nil != err {
		log.Fatal(err)
		return ""
	}

	return string(b[:])
}

func (items GHItems) FilterByNotion(notion string) (ret GHItems) {
	if len(notion) == 0 {
		log.Println("Skipping FilterByNotion")
		return items
	}

	for _, item := range items {
		if strings.EqualFold(notion, item.Notion) {
			ret = append(ret, item)
		}
	}
	return ret
}

func (items GHItems) FilterByUpgrade(upgrade string) (ret GHItems) {
	if len(upgrade) == 0 {
		log.Println("Skipping FilterByUpgrade")
		return items
	}

	for _, item := range items {
		if strings.EqualFold(upgrade, item.Upgrade) {
			ret = append(ret, item)
		}
	}
	return ret
}

func (items GHItems) FilterByStatus(status string) (ret GHItems) {
	if len(status) == 0 {
		log.Println("Skipping FilterByStatus")
		return items
	}

	for _, item := range items {
		if strings.EqualFold(status, item.Status) {
			ret = append(ret, item)
		}
	}
	return ret
}
