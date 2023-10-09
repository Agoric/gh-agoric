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
	"errors"
	"log"
	"os"
	"time"
)

type GHClient struct {
	CachedItemList *GHItemList `json:"CachedItemList"`

	CacheFile     string    `json:"CacheFile"`
	LastRequested time.Time `json:"LastRequested"`

	ProjNum string `json:"ProjNum"`
	Owner   string `json:"Owner"`
	Limit   string `json:"Limit"`
}

func NewGHClient(cacheFile string) (ret GHClient, err error) {
	ret.CacheFile = cacheFile

	_, err = os.Stat(ret.CacheFile)
	cacheFileExists := !errors.Is(err, os.ErrNotExist)

	if cacheFileExists {
		var b []byte
		b, err = os.ReadFile(ret.CacheFile)
		if nil != err {
			log.Fatal(err)
			return ret, err
		}

		err = json.Unmarshal(b, &ret)
		if err != nil {
			log.Fatal(err)
			return ret, err
		}
	}

	return ret, nil
}

func (client GHClient) ReqItems(projNum string, owner string, limit string) (GHItems, error) {
	var err error

	// TODO: Add time check here
	useCachedItems := (projNum == client.ProjNum) &&
		(owner == client.Owner) &&
		(limit == client.Limit)

	if useCachedItems {
		log.Println("CACHE: Found cached data")
		return client.CachedItemList.GetItems(), nil
	}

	client.CachedItemList, err = NewGHItemList(projNum, owner, limit)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	client.LastRequested = time.Now()
	client.ProjNum = projNum
	client.Owner = owner
	client.Limit = limit

	log.Println("CACHE: Saving cached data")
	err = client.save()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return client.CachedItemList.GetItems(), nil
}

func (client GHClient) save() error {
	var b []byte
	b, err := json.Marshal(client)
	if nil != err {
		log.Fatal(err)
		return err
	}

	err = os.WriteFile(client.CacheFile, b, 0666)
	if nil != err {
		log.Fatal(err)
		return err
	}

	return nil
}
