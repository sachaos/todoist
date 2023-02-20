package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/sachaos/todoist/lib"
)

func LoadCache(filename string, s *todoist.Store) error {
	err := ReadCache(filename, s)
	if err != nil {
		err = WriteCache(cachePath, s)
		if err != nil {
			return err
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
	s.ConstructItemTree()
	return nil
}

func WriteCache(filename string, s *todoist.Store) error {
	buf, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return CommandFailed
	}
	err = AssureExists(filename)
	if err != nil {
		return err
	}
	err2 := ioutil.WriteFile(filename, buf, os.ModePerm)
	if err2 != nil {
		return errors.New("Couldn't write to the cache file")
	}
	return nil
}
