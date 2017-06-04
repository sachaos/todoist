package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/sachaos/todoist/lib"
)

func LoadCache(filename string, s *todoist.Store) error {
	err := ReadCache(filename, s)
	if err != nil {
		err = WriteCache(default_cache_path, s)
		if err != nil {
			return CommandFailed
		}
	}
	return nil
}

func ReadCache(filename string, s *todoist.Store) error {
	jsonString, err := ioutil.ReadFile(filename)
	if err != nil {
		return CommandFailed
	}
	err = json.Unmarshal(jsonString, &s)
	if err != nil {
		return CommandFailed
	}
	s.ConstructItemOrder()
	return nil
}

func WriteCache(filename string, s *todoist.Store) error {
	buf, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return CommandFailed
	}
	err2 := ioutil.WriteFile(filename, buf, os.ModePerm)
	if err2 != nil {
		return CommandFailed
	}
	return nil
}
