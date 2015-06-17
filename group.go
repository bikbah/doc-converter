package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

type Group struct {
	Name       string `json:"name,omitempty"`
	Text       string `json:"text,omitempty"`
	Link       string `json:"link,omitempty"`
	Dependency string `json:"dependency,omitempty"`
	Id         string `json:"-"`
}

var groups []Group

func addGroup(g Group) {
	groups = append(groups, g)
}

func saveGroups() {

	res := ""

	for _, group := range groups {
		groupStr, err := json.Marshal(group)
		if err != nil {
			log.Printf("Group marshalling to JSON error: %v\n", err)
		}

		res += `"` + group.Id + `":` + string(groupStr[:]) + `,`
	}

	res = `{` + res[:len(res)-1] + `}`

	res = strings.Replace(res, "and", " && ", -1)
	res = strings.Replace(res, "or", " || ", -1)

	if err := ioutil.WriteFile("metadata.json", []byte(res), 0644); err != nil {
		panic("Write metadata.json file error..")
	}
}

func init() {
	groups = make([]Group, 0)
}
